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

	base, err := readFragment(reflect.Indirect(reflect.ValueOf(obj)), meta, sch)
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

func readFragment(
	source reflect.Value,
	meta map[string]*SchemaMeta,
	sch map[string]*schema.Schema,
) (map[string]interface{}, error) {
	target := make(map[string]interface{})

	for k, m := range meta {
		s, ok := sch[k]
		if !ok {
			return nil, trace.Errorf("field %v not found in corresponding schema", k)
		}

		v := source.FieldByName(m.Name)

		switch {
		case m.Setter != nil:
			r, err := m.Setter(v, m, s)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			err = setConvertedKey(target, k, r, m)
			if err != nil {
				return nil, err
			}
		case s.Type == schema.TypeInt ||
			s.Type == schema.TypeFloat ||
			s.Type == schema.TypeBool ||
			s.Type == schema.TypeString:

			r, err := readAtomic(v, m, s)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			err = setConvertedKey(target, k, r, m)
			if err != nil {
				return nil, trace.Wrap(err)
			}

		case s.Type == schema.TypeList:
			r, err := readList(v, m, s)
			if err != nil {
				return nil, trace.Wrap(err)
			}
			target[k] = r

		case s.Type == schema.TypeMap:
			r, err := readMap(v, m, s)
			if err != nil {
				return nil, trace.Wrap(err)
			}
			target[k] = r

		case s.Type == schema.TypeSet:
			r, err := readSet(v, m, s)
			if err != nil {
				return nil, trace.Wrap(err)
			}
			target[k] = r

		default:
			return nil, trace.Errorf("unknown type %v", s.Type.String())
		}
	}

	return target, nil
}

// readAtomic gets atomic value (scalar, string, time, duration)
func readAtomic(source reflect.Value, meta *SchemaMeta, sch *schema.Schema) (interface{}, error) {
	if source.Kind() == reflect.Ptr && source.IsNil() {
		return nil, nil
	}

	switch {
	case meta.IsTime:
		t, err := readTime(source)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		return t, nil
	case meta.IsDuration:
		d, err := readDuration(source)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		return d, nil
	default:
		return source.Interface(), nil
	}
}

// readList converts source value to list
func readList(source reflect.Value, meta *SchemaMeta, sch *schema.Schema) (interface{}, error) {
	if source.Type().Kind() == reflect.Slice {
		t := make([]interface{}, source.Len())

		for i := 0; i < source.Len(); i++ {
			v := source.Index(i)
			el, err := readEnumerableElement(v, meta, sch)
			if err != nil {
				return nil, err
			}

			n, err := convert(reflect.ValueOf(el), meta)
			if err != nil {
				return nil, err
			}

			t[i] = n
		}

		return t, nil
	}

	t := make([]interface{}, 1)

	item, err := readEnumerableElement(reflect.Indirect(source), meta, sch)
	if err != nil {
		return nil, err
	}

	if item != nil {
		t[0] = item
		return t, nil
	}

	return nil, nil
}

// readMap converts source value to map
func readMap(source reflect.Value, meta *SchemaMeta, sch *schema.Schema) (interface{}, error) {
	return nil, nil
}

// readSet converts source value to set
func readSet(source reflect.Value, meta *SchemaMeta, sch *schema.Schema) (interface{}, error) {
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
	t := reflect.Indirect(source)
	d := reflect.TypeOf((*time.Duration)(nil)).Elem()

	if !t.Type().ConvertibleTo(d) {
		return nil, trace.Errorf("can not convert %T to time.Duration", t)
	}

	s, ok := t.Convert(d).Interface().(time.Duration)
	if !ok {
		return nil, trace.Errorf("can not convert %T to time.Duration", t)
	}

	return s.String(), nil
}

// convert converts source value to schema tyoe given in meta
func convert(source reflect.Value, meta *SchemaMeta) (interface{}, error) {
	t := reflect.Indirect(source)

	if !t.Type().ConvertibleTo(meta.SchemaValueType) {
		return nil, trace.Errorf("can not convert %T to %T", t.Type(), meta.SchemaValueType)
	}

	return t.Convert(meta.SchemaValueType).Interface(), nil
}

// setConvertedKey converts value to target schema type and sets it to resulting map if not nil
func setConvertedKey(target map[string]interface{}, key string, source interface{}, meta *SchemaMeta) error {
	if source != nil {
		f, err := convert(reflect.ValueOf(source), meta)
		if err != nil {
			return trace.Wrap(err)
		}
		target[key] = f
	}

	return nil
}

// readEnumerableElement gets singular slice element from a resource data. If enumerable element is empty, it assigns
// an empty value to the target.
func readEnumerableElement(
	source reflect.Value,
	meta *SchemaMeta,
	sch *schema.Schema,
) (interface{}, error) {
	switch s := sch.Elem.(type) {
	case *schema.Schema:
		return readAtomic(source, meta, s)
	case *schema.Resource:
		return nil, nil
		// v := newEmptyValue(target.Type())

		// _, ok := data.GetOk(path)
		// if ok {
		// 	err := getFragment(path+".", v, meta.Nested, s.Schema, data)
		// 	if err != nil {
		// 		return err
		// 	}
		// }

		// return assign(v, target)
	default:
		return nil, trace.Errorf("unknown Elem type")
	}
}
