/*
Copyright 2015-2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package accessors

import (
	"fmt"
	"reflect"
	"time"

	"github.com/gravitational/trace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Get(
	sch map[string]*schema.Schema,
	data *schema.ResourceData,
	meta map[string]*SchemaMeta,
	obj interface{},
) error {
	if obj == nil {
		return trace.Errorf("obj must not be nil")
	}

	t := reflect.Indirect(reflect.ValueOf(obj))
	return getFragment("", &t, meta, sch, data)
}

// getFragment iterates over schema fragment and calls appropriate getters for each field
func getFragment(
	path string,
	target *reflect.Value,
	meta map[string]*SchemaMeta,
	sch map[string]*schema.Schema,
	data *schema.ResourceData,
) error {
	for k, m := range meta {
		s, ok := sch[k]
		if !ok {
			return trace.Errorf("field %v.%v not found in corresponding schema", path, k)
		}

		v := target.FieldByName(m.Name)

		switch {
		case s.Type == schema.TypeInt ||
			s.Type == schema.TypeFloat ||
			s.Type == schema.TypeBool ||
			s.Type == schema.TypeString:

			err := setAtomic(path+k, &v, m, s, data)
			if err != nil {
				return trace.Wrap(err)
			}
		case s.Type == schema.TypeList:
			err := setList(path+k, &v, m, s, data)
			if err != nil {
				return trace.Wrap(err)
			}

		case s.Type == schema.TypeMap:
			err := setMap(path+k, &v, m, s, data)
			if err != nil {
				return trace.Wrap(err)
			}
		case s.Type == schema.TypeSet:
			err := setSet(path+k, &v, m, s, data)
			if err != nil {
				return trace.Wrap(err)
			}

		default:
			return trace.Errorf("unknown type %v", s.Type.String())
		}
	}

	return nil
}

// setAtomic sets atomic value (scalar, string, time, duration)
func setAtomic(path string, target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
	s, ok := data.GetOk(path)
	if !ok {
		target.Set(reflect.Zero(target.Type()))

		return nil
	}

	switch {
	case meta.IsTime:
		err := assignTime(s, target)
		if err != nil {
			return trace.Wrap(err)
		}
	case meta.IsDuration:
		err := assignDuration(s, target)
		if err != nil {
			return trace.Wrap(err)
		}
	default:
		err := assignAtomic(s, target)
		if err != nil {
			return trace.Wrap(err)
		}
	}

	return nil
}

// TODO: change assignAtomic
// assignAtomic reads atomic value form
func assignAtomic(source interface{}, target *reflect.Value) error {
	v := reflect.ValueOf(source)
	t := target.Type()

	// If target type is at the pointer reference use underlying type
	if target.Type().Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Convert value to target type
	if reflect.TypeOf(source) != t {
		if !v.Type().ConvertibleTo(target.Type()) {
			return trace.Errorf("can not convert %v to %v", v.Type().Name(), t.Name())
		}

		v = v.Convert(t)
	}

	if !v.Type().AssignableTo(t) {
		return trace.Errorf("can not assign %s to %s", v.Type().Name(), t.Name())
	}

	// If target type is a reference, create new pointer to this reference and assign
	if target.Type().Kind() == reflect.Ptr {
		e := reflect.New(v.Type())
		e.Elem().Set(v)
		v = e
	}

	target.Set(v)

	return nil
}

// assignTime assigns time value from a string
func assignTime(source interface{}, target *reflect.Value) error {
	s, ok := source.(string)
	if !ok {
		return trace.Errorf("can not convert %T to string", source)
	}

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return trace.Errorf("can not parse time: %w", err)
	}

	return assignAtomic(t, target)
}

// assignTime assigns duration value from a string
func assignDuration(source interface{}, target *reflect.Value) error {
	s, ok := source.(string)
	if !ok {
		return trace.Errorf("can not convert %T to string", source)
	}

	t, err := time.ParseDuration(s)
	if err != nil {
		return trace.Errorf("can not parse duration: %w", err)
	}

	return assignAtomic(t, target)
}

// setList sets atomic value (scalar, string, time, duration)
func setList(path string, target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
	// Get list count variable
	n, okn := data.GetOk(path + ".#")
	len, okc := n.(int)

	if !okc {
		return trace.Errorf("failed to convert list count to number")
	}

	// If list is empty, set target list to empty value
	if !okn || len == 0 {
		target.Set(reflect.Zero(target.Type()))

		return nil
	}

	if target.Type().Kind() == reflect.Slice {
		r := reflect.MakeSlice(target.Type(), len, len)

		for i := 0; i < len; i++ {
			el := r.Index(i)
			p := fmt.Sprintf("%v.%v", path, i)

			switch s := sch.Elem.(type) {
			case *schema.Schema:
				err := setAtomic(p, &el, meta, s, data)
				if err != nil {
					return trace.Wrap(err)
				}
			case *schema.Resource:
				if el.Kind() == reflect.Ptr {
					el.Set(reflect.New(el.Type().Elem()))
					el = reflect.Indirect(el)
				}

				err := getFragment(p+".", &el, meta.Nested, s.Schema, data)
				if err != nil {
					return trace.Wrap(err)
				}
			default:
				return trace.Errorf("unknown Elem type")
			}
		}

		target.Set(r)
	} else {
		s, ok := sch.Elem.(*schema.Resource)
		if !ok {
			return trace.Errorf("failed to convert %T to *schema.Resource", sch.Elem)
		}

		// Construct blank object
		t := target.Type()
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		r := reflect.Indirect(reflect.New(t))

		// Fill blank object in
		err := getFragment(path+".0.", &r, meta.Nested, s.Schema, data)
		if err != nil {
			return trace.Wrap(err)
		}

		// Assign blank object
		err = assignAtomic(r.Interface(), target)
		if err != nil {
			return trace.Wrap(err)
		}
	}

	return nil
}

// setMap sets map of atomic values (scalar, string, time, duration)
func setMap(path string, target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
	md, ok := data.GetOk(path)
	if !ok {
		return nil
	}

	m, ok := md.(map[string]interface{})
	if !ok {
		return trace.Errorf("failed to convert %T to map[string]interface{}", md)
	}

	// If map is empty, set target empty map
	if len(m) == 0 {
		target.Set(reflect.Zero(target.Type()))

		return nil
	}

	if target.Type().Kind() != reflect.Map {
		return trace.Errorf("target time is not a map")
	}

	r := reflect.MakeMap(target.Type())

	for k := range m {
		kv := reflect.ValueOf(k)

		el := reflect.Indirect(reflect.New(target.Type().Elem()))

		switch s := sch.Elem.(type) {
		case *schema.Schema:
			err := setAtomic(path+"."+k, &el, meta, s, data)
			if err != nil {
				return trace.Wrap(err)
			}
		default:
			return trace.Errorf("unknown Elem type: map key must be *schema.Schmea")
		}

		r.SetMapIndex(kv, el)
	}

	target.Set(r)

	return nil
}

// setMessageMap sets map of atomic values (scalar, string, time, duration)
func setSet(path string, target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
	// md, ok := data.GetOk(path)
	// if !ok {
	// 	return nil
	// }

	// m, ok := md.(map[string]interface{})
	// if !ok {
	// 	return trace.Errorf("failed to convert %T to map[string]interface{}", md)
	// }

	// // If map is empty, set target empty map
	// if len(m) == 0 {
	// 	target.Set(reflect.Zero(target.Type()))

	// 	return nil
	// }

	// if target.Type().Kind() != reflect.Map {
	// 	return trace.Errorf("target time is not a map")
	// }

	// r := reflect.MakeMap(target.Type())

	// for k := range m {
	// 	kv := reflect.ValueOf(k)

	// 	el := reflect.Indirect(reflect.New(target.Type().Elem()))

	// 	switch s := sch.Elem.(type) {
	// 	case *schema.Schema:
	// 		err := setAtomic(path+"."+k, &el, meta, s, data)
	// 		if err != nil {
	// 			return trace.Wrap(err)
	// 		}
	// 	default:
	// 		return trace.Errorf("unknown Elem type: map key must be *schema.Schmea")
	// 	}

	// 	r.SetMapIndex(kv, el)
	// }

	// target.Set(r)

	return nil
}
