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

// Package accessors contains Get and Set methods for ResourceData
package accessors

import (
	"reflect"
	"time"

	"github.com/gravitational/trace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Set assigns object data from object to schema.ResourceData
//
// Example:
//   user := UserV2{Name: "example"}
//   Set(&user, data, SchemaUserV2, MetaUserV2)
func Set(
	obj interface{},
	data *schema.ResourceData,
	sch map[string]*schema.Schema,
	meta map[string]*SchemaMeta,
) error {
	if obj == nil {
		return trace.Errorf("obj must not be nil")
	}

	base, err := setFragment(reflect.Indirect(reflect.ValueOf(obj)), meta, sch, data)
	if err != nil {
		return trace.Wrap(err)
	}

	for k, v := range base {
		if v != nil {
			err := data.Set(k, v)
			if err != nil {
				return trace.Wrap(err)
			}
		}
	}

	return nil
}

func setFragment(
	source reflect.Value,
	meta map[string]*SchemaMeta,
	sch map[string]*schema.Schema,
	data *schema.ResourceData,
) (map[string]interface{}, error) {
	target := make(map[string]interface{})

	for k, m := range meta {
		s, ok := sch[k]
		if !ok {
			return nil, trace.Errorf("field %v not found in corresponding schema", k)
		}

		v := source.FieldByName(m.Name)

		switch {
		// case m.Setter != nil:
		// err := m.Setter(target, v, m, s, data)
		// if err != nil {
		// 	return err
		// }
		case s.Type == schema.TypeInt ||
			s.Type == schema.TypeFloat ||
			s.Type == schema.TypeBool ||
			s.Type == schema.TypeString:

			r, err := setAtomic(v, m, s, data)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			target[k] = r
			// case s.Type == schema.TypeList:
			// 	err := getList(p, &v, m, s, data)
			// 	if err != nil {
			// 		return trace.Wrap(err)
			// 	}

			// case s.Type == schema.TypeMap:
			// 	err := getMap(p, &v, m, s, data)
			// 	if err != nil {
			// 		return trace.Wrap(err)
			// 	}

			// case s.Type == schema.TypeSet:
			// 	err := getSet(p, &v, m, s, data)
			// 	if err != nil {
			// 		return trace.Wrap(err)
			// 	}

			// default:
			// 	return trace.Errorf("unknown type %v for %s", s.Type.String(), p)
		}
	}

	return target, nil
	// // Fragment value must be nil in case it's empty
	// if len(target) > 0 {

	// }

	// return nil
}

// getAtomic gets atomic value (scalar, string, time, duration)
func setAtomic(source reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) (interface{}, error) {
	if source.Kind() == reflect.Ptr && source.IsNil() {
		return nil, nil
	}

	switch {
	case meta.IsTime:
		return readTime(source)
	case meta.IsDuration:
		return readDuration(source)
	default:
		return source.Interface(), nil
	}

	return nil, nil
}

// readTime returns value as time
func readTime(source reflect.Value) (interface{}, error) {
	t := reflect.Indirect(source).Interface()
	v, ok := t.(time.Time)
	if !ok {
		return nil, trace.Errorf("can not convert %T to time.Time", t)
	}
	return v.Format(time.RFC3339Nano), nil
}

// readDuration returns value as duration
func readDuration(source reflect.Value) (interface{}, error) {
	var _d time.Duration

	t := reflect.Indirect(source)
	d := reflect.TypeOf(_d)

	if !t.Type().ConvertibleTo(d) {
		return nil, trace.Errorf("can not convert %T to time.Duration", t)
	}

	return t.Convert(d).MethodByName("String").Call([]reflect.Value{}), nil
}
