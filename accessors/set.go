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

// readFragment returns map[string]interface{} of a block
func readFragment(
	source reflect.Value,
	meta map[string]*SchemaMeta,
	sch map[string]*schema.Schema,
) (map[string]interface{}, error) {
	target := make(map[string]interface{})

	if !source.IsValid() {
		return nil, nil
	}

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

			if r != nil {
				target[k] = r
			}
		case s.Type == schema.TypeInt ||
			s.Type == schema.TypeFloat ||
			s.Type == schema.TypeBool ||
			s.Type == schema.TypeString:

			r, err := readAtomic(v, m, s)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			err = setConvertedKey(target, k, r, s)
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

			t[i] = el
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
	if source.Len() == 0 {
		return nil, nil
	}

	m := make(map[string]interface{})

	for _, k := range source.MapKeys() {
		i := source.MapIndex(k)

		v, err := readEnumerableElement(i, meta, sch)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		reflect.ValueOf(m).SetMapIndex(k, reflect.ValueOf(v))
	}

	return m, nil
}

// readSet converts source value to set
func readSet(source reflect.Value, meta *SchemaMeta, sch *schema.Schema) (interface{}, error) {
	if source.Len() == 0 {
		return nil, nil
	}

	s, ok := sch.ZeroValue().(*schema.Set)
	if !ok {
		return nil, trace.Errorf("zero value for schema set element is not *schema.Set")
	}

	switch source.Kind() {
	case reflect.Slice:
		// TODO: This case is not important for now
		return nil, trace.NotImplemented("set acting as list on target is not implemented yet")
	case reflect.Map:
		for _, k := range source.MapKeys() {
			i := source.MapIndex(k)

			vsch := sch.Elem.(*schema.Resource).Schema["value"]

			v, err := readEnumerableElement(i, meta, vsch)
			if err != nil {
				return nil, trace.Wrap(err)
			}

			t := map[string]interface{}{
				"key":   k.Interface(),
				"value": []interface{}{v},
			}

			s.Add(t)
		}

		return s, nil
	default:
		return nil, trace.Errorf("unknown set source type")
	}
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

// convert converts source value to schema type given in meta
func convert(source reflect.Value, sch *schema.Schema) (interface{}, error) {
	t := reflect.Indirect(source)
	s, err := schemaValueType(sch)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	if !t.Type().ConvertibleTo(s) {
		return nil, trace.Errorf("can not convert %T to %T", t.Type(), s)
	}

	return t.Convert(s).Interface(), nil
}

// setConvertedKey converts value to target schema type and sets it to resulting map if not nil
func setConvertedKey(target map[string]interface{}, key string, source interface{}, sch *schema.Schema) error {
	if source != nil {
		f, err := convert(reflect.ValueOf(source), sch)
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
		a, err := readAtomic(source, meta, s)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		n, err := convert(reflect.ValueOf(a), s)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		return n, nil

	case *schema.Resource:
		m, err := readFragment(reflect.Indirect(source), meta.Nested, s.Schema)
		if err != nil {
			return nil, err
		}

		return m, nil
	default:
		return nil, trace.Errorf("unknown Elem type")
	}
}

// schemaValueType returns type to convert value to
func schemaValueType(sch *schema.Schema) (reflect.Type, error) {
	switch sch.Type {
	case schema.TypeFloat:
		return reflect.TypeOf((*float64)(nil)).Elem(), nil
	case schema.TypeInt:
		return reflect.TypeOf((*int)(nil)).Elem(), nil
	case schema.TypeBool:
		return reflect.TypeOf((*bool)(nil)).Elem(), nil
	case schema.TypeString:
		return reflect.TypeOf((*string)(nil)).Elem(), nil
	default:
		return nil, trace.Errorf("unknown schema type: %v", sch.Type.String())
	}

}
