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

// Get reads object data from schema.ResourceData
func Get(
	sch map[string]*schema.Schema,
	data *schema.ResourceData,
	meta map[string]*SchemaMeta,
	obj interface{},
) error {
	if obj == nil {
		return trace.Errorf("obj must not be nil")
	}

	// Obj must be reference
	t := reflect.Indirect(reflect.ValueOf(obj))
	return getFragment("", &t, meta, sch, data)
}

// getFragment iterates over a schema fragment and calls appropriate getters for each field
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
	len, err := getLen(path, data)
	if err != nil {
		return trace.Wrap(err)
	}

	if len == 0 {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	// Target is a slice of elementary values or objects
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
		// Target is a singular object
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
			return trace.Errorf("unknown Elem type: map key must be *schema.Schema")
		}

		r.SetMapIndex(kv, el)
	}

	target.Set(r)

	return nil
}

// setSet reads set from resource data
func setSet(path string, target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
	len, err := getLen(path, data)
	if err != nil {
		return trace.Wrap(err)
	}

	if len == 0 {
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	switch target.Kind() {
	case reflect.Slice:
		// This set must be converted to normal slice
		return nil
	case reflect.Map:
		// This set must be read into a map, so, it contains artificial key and value arguments
		r := reflect.MakeMap(target.Type())

		ds, ok := data.GetOk(path)
		if !ok {
			return fmt.Errorf("can not read key " + path)
		}

		s, ok := ds.(*schema.Set)
		if !ok {
			return fmt.Errorf("can not convert %T to *schema.Set", ds)
		}

		for _, i := range s.List() {
			m := i.(map[string]interface{})
			k := m["key"]

			t := target.Type().Elem()
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			v := reflect.Indirect(reflect.New(t))

			e := sch.Elem.(*schema.Resource).Schema["value"].Elem.(*schema.Resource)

			p := fmt.Sprintf("%v.%v.value.0.", path, s.F(i))

			err := getFragment(p, &v, meta.Nested, e.Schema, data)
			if err != nil {
				return err
			}

			if target.Type().Elem().Kind() == reflect.Ptr {
				r.SetMapIndex(reflect.ValueOf(k), v.Addr())
			} else {
				r.SetMapIndex(reflect.ValueOf(k), v)
			}
		}

		target.Set(r)

		return nil
	default:
		return fmt.Errorf("unknown set target type")
	}

	// md, ok := data.GetOk(path)
	// if !ok {
	// 	return nil
	// }

	// m, ok := md.(map[string]interface{})
	// if !ok {
	// 	return trace.Errorf("failed to convert %T to map[string]interface{}", md)
	// }

	// target is list: this is list of objects
	// target is map: this is map of objects

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

func getLen(path string, data *schema.ResourceData) (int, error) {
	n, ok := data.GetOk(path + ".#")
	if !ok || n == nil {
		return 0, nil
	}

	len, ok := n.(int)
	if !ok {
		return 0, trace.Errorf("failed to convert list count to number")
	}

	return len, nil
}
