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

// Package test contains protoc-gen-terraform tests
package test

import (
	"reflect"
	"testing"
	"time"

	"github.com/gravitational/trace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// defaultTimestamp predefined timestamp used in tests
	defaultTimestamp = "2022-10-12T07:20:50.5Z"
)

var (
	// fixture raw source data for schema.ResourceData
	fixture map[string]interface{} = map[string]interface{}{
		"str":                              "TestString",
		"int32":                            999,
		"int64":                            998,
		"float":                            18.1,
		"double":                           18.4,
		"bool":                             true,
		"bytes":                            "TestBytes",
		"timestamp":                        defaultTimestamp,
		"timestamp_nullable":               defaultTimestamp,
		"timestamp_nullable_with_no_value": nil,
		"duration_standard":                "1h",
		"duration_custom":                  "1m",
		"string_a":                         []interface{}{"TestString1", "TestString2"},
		"bool_a":                           []interface{}{false, true, false},
		"bytes_a":                          []interface{}{"TestBytes1", "TestBytes2"},
		"timestamp_a":                      []interface{}{defaultTimestamp},
		"duration_custom_a":                []interface{}{"1m"},

		"nested": []interface{}{
			map[string]interface{}{
				"str": "TestString",
				"nested": []interface{}{
					map[string]interface{}{
						"str": "TestString1",
					},
					map[string]interface{}{
						"str": "TestString2",
					},
				},
				"nested_m": map[string]interface{}{
					"kn1": "vn1",
					"kn2": "vn2",
				},
			},
		},

		"nested_m": map[string]interface{}{
			"k1": "v1",
			"k2": "v2",
		},

		"nested_m_obj": []interface{}{
			map[string]interface{}{
				"key": "obj1",
				"value": []interface{}{
					map[string]interface{}{
						"str": "TestString1",
					},
				},
			},
			map[string]interface{}{
				"key": "obj2",
				"value": []interface{}{
					map[string]interface{}{
						"str": "TestString2",
					},
				},
			},
		},
	}
)

type SchemaMeta struct {
	name       string
	schemaName string
	isTime     bool
	isDuration bool
	nested     map[string]*SchemaMeta
}

// NOTE: make compound struct, values instead of methods
func SchemaTestMeta() map[string]*SchemaMeta {
	return map[string]*SchemaMeta{
		"str": {
			name:       "Str",
			schemaName: "str",
		},
		"int32": {
			name:       "Int32",
			schemaName: "int32",
		},
		"int64": {
			name:       "Int64",
			schemaName: "int64",
		},
		"float": {
			name:       "Float",
			schemaName: "float",
		},
		"double": {
			name:       "Double",
			schemaName: "double",
		},
		"bool": {
			name:       "Bool",
			schemaName: "bool",
		},
		"bytes": {
			name:       "Bytes",
			schemaName: "bytes",
		},
		"timestamp": {
			name:       "Timestamp",
			schemaName: "timestamp",
			isTime:     true,
		},
		"timestamp_nullable": {
			name:       "TimestampNullable",
			schemaName: "timestamp_nullable",
			isTime:     true,
		},
		"timestamp_nullable_with_nil_value": {
			name:       "TimestampNullableWithNilValue",
			schemaName: "timestamp_nullable_with_nil_value",
			isTime:     true,
		},
		"duration_standard": {
			name:       "DurationStandard",
			schemaName: "duration_standard",
			isDuration: true,
		},
		"duration_custom": {
			name:       "DurationCustom",
			schemaName: "duration_custom",
			isDuration: true,
		},
	}
}

func GetFromResourceData(
	sch map[string]*schema.Schema,
	data *schema.ResourceData,
	meta map[string]*SchemaMeta,
	obj interface{},
) error {
	o := reflect.ValueOf(obj)

	if o.IsNil() {
		return trace.Errorf("obj must not be nil")
	}

	if o.Kind() != reflect.Ptr {
		return trace.Errorf("pass a reference to an object you need to set")
	}

	o = reflect.Indirect(o)

	for k, sch := range sch {
		m, ok := meta[k]
		if !ok {
			continue
			return trace.Errorf("field not found in corresponding meta " + k)
		}

		v := o.FieldByName(m.name)

		switch {
		//case me.object_map:
		//case me.custom
		case sch.Type == schema.TypeInt ||
			sch.Type == schema.TypeFloat ||
			sch.Type == schema.TypeBool ||
			sch.Type == schema.TypeString:

			err := setAtomic(&v, m, sch, data)
			if err != nil {
				return trace.Wrap(err)
			}
		// case sch.Type == schema.TypeList:
		// case sch.Type == schema.TypeMap:
		// case sch.Type == schema.TypeSet:
		default:
			return trace.Errorf("unknown type %v", sch.Type)
		}
	}

	return nil
}

// setAtomic sets atomic value (scalar, string, time, duration)
func setAtomic(target *reflect.Value, meta *SchemaMeta, sch *schema.Schema, data *schema.ResourceData) error {
	s, ok := data.GetOk(meta.schemaName)
	if !ok {
		target.Set(reflect.Zero(target.Type()))

		return nil
	}

	switch {
	case meta.isTime:
		err := assignTime(s, target)
		if err != nil {
			return trace.Wrap(err)
		}
	case meta.isDuration:
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

// func setList()
// func setMap()
// func setSet()
// func setCustom()

// buildSubjectGet builds Test struct from test fixture data
func buildSubjectGet(t *testing.T, subject *Test) (*Test, error) {
	data := schema.TestResourceDataRaw(t, SchemaTest(), fixture)
	err := GetFromResourceData(SchemaTest(), data, SchemaTestMeta(), subject)
	return subject, err
}

// TestElementariesGet ensures decoding of elementary types
func TestElementariesGet(t *testing.T) {
	subject, err := buildSubjectGet(t, &Test{})
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, int32(999), subject.Int32, "Test.Int32")
	assert.Equal(t, int64(998), subject.Int64, "Test.Int64")
	assert.Equal(t, "TestString", subject.Str, "Test.Str")
	assert.Equal(t, float32(18.1), subject.Float, "Test.Float")
	assert.Equal(t, float64(18.4), subject.Double, "Test.Dobule")
	assert.Equal(t, true, subject.Bool, "Test.Bool")
	assert.Equal(t, []byte("TestBytes"), subject.Bytes, "Test.Bytes")
}

// TestTimesGet ensures decoding of time and duration fields
func TestTimesGet(t *testing.T) {
	now := time.Now()

	// Ensure nullify
	subject, err := buildSubjectGet(t, &Test{TimestampNullableWithNilValue: &now})
	require.NoError(t, err, "failed to unmarshal test data")

	timestamp, err := time.Parse(time.RFC3339Nano, defaultTimestamp)
	require.NoError(t, err, "failed to parse example timestamp")

	durationStd, err := time.ParseDuration("1h")
	require.NoError(t, err, "failed to parse example duration")

	durationCustom, err := time.ParseDuration("1m")
	require.NoError(t, err, "failed to parse example duration")

	assert.Equal(t, timestamp, subject.Timestamp, "Test.Timestamp")
	assert.Equal(t, timestamp, *subject.TimestampNullable, "Test.TimestampNullable")
	assert.Nil(t, subject.TimestampNullableWithNilValue, "Test.Timestamp.TimestampNullableWithNilValue")

	assert.Equal(t, durationStd, subject.DurationStandard, "Test.DurationStandard")
	assert.Equal(t, Duration(durationCustom), subject.DurationCustom, "Test.DurationCustom")
}

// // TestArraysGet ensures decoding of arrays
// func TestArraysGet(t *testing.T) {
// 	subject, err := buildSubjectGet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	timestamp, err := time.Parse(time.RFC3339Nano, defaultTimestamp)
// 	require.NoError(t, err, "failed to parse example timestamp")

// 	duration, err := time.ParseDuration("1m")
// 	require.NoError(t, err, "failed to parse example duration")

// 	assert.Equal(t, subject.StringA, []string{"TestString1", "TestString2"})
// 	assert.Equal(t, subject.BoolA, []BoolCustom{false, true, false})
// 	assert.Equal(t, subject.BytesA, [][]byte{[]byte("TestBytes1"), []byte("TestBytes2")})
// 	assert.Equal(t, subject.TimestampA, []*time.Time{&timestamp})
// 	assert.Equal(t, subject.DurationCustomA, []Duration{Duration(duration)})
// }

// // TestNestedMessageGet ensures decoding of nested messages
// func TestNestedMessageGet(t *testing.T) {
// 	subject, err := buildSubjectGet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, subject.Nested.Str, "TestString", "Test.Nested.Str")
// }

// // TestNestedMessageArrayGet ensures decoding of array of messages
// func TestNestedMessageArrayGet(t *testing.T) {
// 	subject, err := buildSubjectGet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, subject.Nested.Nested[0].Str, "TestString1")
// 	assert.Equal(t, subject.Nested.Nested[1].Str, "TestString2")
// }

// // TestMapGet ensures decoding of a maps
// func TestMapGet(t *testing.T) {
// 	subject, err := buildSubjectGet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, subject.NestedM["k1"], "v1")
// 	assert.Equal(t, subject.NestedM["k2"], "v2")
// 	assert.Equal(t, subject.Nested.NestedM["kn1"], "vn1")
// 	assert.Equal(t, subject.Nested.NestedM["kn1"], "vn1")
// }

// // TestObjectMapGet ensures decoding of maps of messages
// func TestObjectMapGet(t *testing.T) {
// 	subject, err := buildSubjectGet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, subject.NestedMObj["obj1"].Str, "TestString1")
// 	assert.Equal(t, subject.NestedMObj["obj2"].Str, "TestString2")
// }
