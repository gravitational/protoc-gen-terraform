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

// Get reads object data from schema.ResourceData to object
//
// Example:
//   user := UserV2{}
//   Get(&user, data, SchemaUserV2, MetaUserV2)
func Get(
	obj interface{},
	data *schema.ResourceData,
	sch map[string]*schema.Schema,
	meta map[string]*SchemaMeta,
) error {
	if obj == nil {
		return trace.Errorf("obj must not be nil")
	}

	// Obj must be reference
	t := reflect.Indirect(reflect.ValueOf(obj))
	return getFragment("", &t, meta, sch, data)
}

// getFragment iterates over a schema fragment and calls appropriate getters for a fields of passed target.
// Target must be struct.
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
		p := path + k

		switch {
		case s.Type == schema.TypeInt ||
			s.Type == schema.TypeFloat ||
			s.Type == schema.TypeBool ||
			s.Type == schema.TypeString:

			err := getAtomic(p, &v, m, s, data)
			if err != nil {
				return trace.Wrap(err)
			}
		case s.Type == schema.TypeList:
			err := getList(p, &v, m, s, data)
			if err != nil {
				return trace.Wrap(err)
			}

		case s.Type == schema.TypeMap:
			err := getMap(p, &v, m, s, data)
			if err != nil {
				return trace.Wrap(err)
			}

		case s.Type == schema.TypeSet:
			err := getSet(p, &v, m, s, data)
			if err != nil {
				return trace.Wrap(err)
			}

		default:
			return trace.Errorf("unknown type %v for %s", s.Type.String(), p)
		}
	}

	return nil
}

// getEnumerableElement gets singular slice element from a resource data
func getEnumerableElement(
	path string,
	target *reflect.Value,
	sch *schema.Schema,
	meta *SchemaMeta,
	data *schema.ResourceData,
) error {
	switch s := sch.Elem.(type) {
	case *schema.Schema:
		err := getAtomic(path, target, meta, s, data)
		if err != nil {
			return trace.Wrap(err)
		}
	case *schema.Resource:
		t := target

		if target.Kind() == reflect.Ptr {
			target.Set(reflect.New(target.Type().Elem()))
			i := reflect.Indirect(*target)
			t = &i
		}

		err := getFragment(path+".", t, meta.Nested, s.Schema, data)
		if err != nil {
			return trace.Wrap(err)
		}
	default:
		return trace.Errorf("unknown Elem type")
	}

	return nil
}

// getAtomic gets atomic value (scalar, string, time, duration)
func getAtomic(path string, target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
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
		v := reflect.ValueOf(s)
		err := assign(&v, target)
		if err != nil {
			return trace.Wrap(err)
		}
	}

	return nil
}

// assignTime parses time value from a string and assigns it to target
func assignTime(source interface{}, target *reflect.Value) error {
	s, ok := source.(string)
	if !ok {
		return trace.Errorf("can not convert %T to string", source)
	}

	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return trace.Errorf("can not parse time: %w", err)
	}

	v := reflect.ValueOf(t)
	return assign(&v, target)
}

// assignTime parses duration value from a string and assigns it to target
func assignDuration(source interface{}, target *reflect.Value) error {
	s, ok := source.(string)
	if !ok {
		return trace.Errorf("can not convert %T to string", source)
	}

	t, err := time.ParseDuration(s)
	if err != nil {
		return trace.Errorf("can not parse duration: %w", err)
	}

	v := reflect.ValueOf(t)
	return assign(&v, target)
}

// setList gets list from ResourceData
func getList(path string, target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
	len, err := getLen(path, data)
	if err != nil {
		return trace.Wrap(err)
	}

	// Empty list: do nothing, but set target field to empty value
	if len == 0 {
		assignZeroValue(target)
		return nil
	}

	// Target is a slice of elementary values or objects
	if target.Type().Kind() == reflect.Slice {
		r := reflect.MakeSlice(target.Type(), len, len)

		for i := 0; i < len; i++ {
			el := r.Index(i)
			p := fmt.Sprintf("%v.%v", path, i)

			err := getEnumerableElement(p, &el, sch, meta, data)
			if err != nil {
				return err
			}
		}

		return assign(&r, target)
	} else {
		// Target is an object represented by a single element list
		err := getEnumerableElement(path+".0", target, sch, meta, data)
		if err != nil {
			return err
		}

		return nil
	}
}

// setMap sets map of atomic values (scalar, string, time, duration)
func getMap(path string, target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
	raw, ok := data.GetOk(path)
	if !ok {
		return nil
	}

	m, ok := raw.(map[string]interface{})
	if !ok {
		return trace.Errorf("failed to convert %T to map[string]interface{}", raw)
	}

	// If map is empty, set target empty map
	if len(m) == 0 {
		target.Set(reflect.Zero(target.Type()))

		return nil
	}

	if target.Type().Kind() != reflect.Map {
		return trace.Errorf("target time is not a map")
	}

	t := reflect.MakeMap(target.Type())

	// Iterate over map keys
	for k := range m {
		// Construct new map value
		// Map values can be only elementary so no need to handle reflect.Ptr here
		v := reflect.Indirect(reflect.New(target.Type().Elem()))

		err := getEnumerableElement(path+"."+k, &v, sch, meta, data)
		if err != nil {
			return err
		}

		t.SetMapIndex(reflect.ValueOf(k), v)
	}

	return assign(&t, target)
}

// setSet reads set from resource data
func getSet(path string, target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
	len, err := getLen(path, data)
	if err != nil {
		return trace.Wrap(err)
	}

	if len == 0 {
		assignZeroValue(target)
		return nil
	}

	switch target.Kind() {
	case reflect.Slice:
		// This set must be converted to normal slice
		return trace.NotImplemented("set acting as list on target is not implemented yet")
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
}

// getLen returns length of a list or set
func getLen(path string, data *schema.ResourceData) (int, error) {
	n, ok := data.GetOk(path + ".#")
	if !ok || n == nil {
		return 0, nil
	}

	len, ok := n.(int)
	if !ok {
		return 0, trace.Errorf("failed to convert list count to number %s", path)
	}

	return len, nil
}
