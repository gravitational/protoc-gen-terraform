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
	"testing"
	"time"

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
		// Elementary values, timestamps, durations
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
		"string_list":                      []interface{}{"TestString1", "TestString2"},
		"bool_custom_list":                 []interface{}{false, true, false},
		"bytes_list":                       []interface{}{"TestBytes1", "TestBytes2"},
		"timestamp_list":                   []interface{}{defaultTimestamp},
		"duration_custom_list":             []interface{}{"1m"},

		// Messages
		"nested": []interface{}{
			map[string]interface{}{
				"str": "TestString",
				"nested_list": []interface{}{
					map[string]interface{}{
						"str": "TestString1",
					},
					map[string]interface{}{
						"str": "TestString2",
					},
				},
				"map": map[string]interface{}{
					"kn1": "vn1",
					"kn2": "vn2",
				},
				"map_object_nested": []interface{}{
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
			},
		},
		"nested_nullable": []interface{}{
			map[string]interface{}{
				"str": "TestString",
			},
		},

		// List of messages
		"nested_list": []interface{}{
			map[string]interface{}{
				"str": "TestString1",
			},
			map[string]interface{}{
				"str": "TestString2",
			},
		},
		"nested_list_nullable": []interface{}{
			map[string]interface{}{
				"str": "TestString1",
			},
			map[string]interface{}{
				"str": "TestString2",
			},
		},

		// Map
		"map": map[string]interface{}{
			"k1": "v1",
			"k2": "v2",
		},

		"map_object": []interface{}{
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

		"map_object_nullable": []interface{}{
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

// buildSubjectGet builds Test struct from test fixture data
func buildSubjectGet(t *testing.T, subject *Test) (*Test, error) {
	data := schema.TestResourceDataRaw(t, SchemaTest, fixture)
	err := GetTest(subject, data)
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
func TestArraysGet(t *testing.T) {
	subject, err := buildSubjectGet(t, &Test{StringListEmpty: []string{"a"}})
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, []string(nil), subject.StringListEmpty, "Test.StringListEmpty")
	assert.Equal(t, []string{"TestString1", "TestString2"}, subject.StringList)

	timestamp, err := time.Parse(time.RFC3339Nano, defaultTimestamp)
	require.NoError(t, err, "failed to parse example timestamp")

	duration, err := time.ParseDuration("1m")
	require.NoError(t, err, "failed to parse example duration")

	assert.Equal(t, subject.BoolCustomList, []BoolCustom{false, true, false})
	assert.Equal(t, [][]byte{[]byte("TestBytes1"), []byte("TestBytes2")}, subject.BytesList, "Test.BytesList")
	assert.Equal(t, []*time.Time{&timestamp}, subject.TimestampList, "Test.TimestampList")
	assert.Equal(t, []Duration{Duration(duration)}, subject.DurationCustomList, "Test.DurationCustomList")
}

// TestNestedMessageGet ensures decoding of nested messages
func TestNestedMessageGet(t *testing.T) {
	subject, err := buildSubjectGet(t, &Test{NestedNullableWithNilValue: &Nested{Str: "5"}})
	require.NoError(t, err, "failed to unmarshal test data")

	var x *Nested = nil

	assert.Equal(t, x, subject.NestedNullableWithNilValue, "Test.NestedNullableWithNilValue")
	assert.Equal(t, "TestString", subject.Nested.Str, "Test.Nested.Str")
	assert.Equal(t, "TestString", subject.NestedNullable.Str, "Test.NestedNullable.Str")
}

// TestNestedMessageArrayGet ensures decoding of array of messages
func TestNestedMessageArrayGet(t *testing.T) {
	subject, err := buildSubjectGet(t, &Test{NestedNullableWithNilValue: &Nested{Str: "5"}})
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, 2, len(subject.NestedList), "len(NestedList)")
	assert.Equal(t, 2, len(subject.NestedListNullable), "len(NestedListNullable)")

	assert.Equal(t, "TestString1", subject.NestedList[0].Str, "NestedList[0].Str")
	assert.Equal(t, "TestString2", subject.NestedList[1].Str, "NestedList[1].Str")
	assert.Equal(t, "TestString1", subject.NestedListNullable[0].Str, "NestedListNullable[0].Str")
	assert.Equal(t, "TestString2", subject.NestedListNullable[1].Str, "NestedListNullable[1].Str")
}

// TestMapGet ensures decoding of a maps
func TestMapGet(t *testing.T) {
	subject, err := buildSubjectGet(t, &Test{NestedNullableWithNilValue: &Nested{Str: "5"}})
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, "v1", subject.Map["k1"], "Test.Map['k1']")
	assert.Equal(t, "v2", subject.Map["k2"], "Test.Map['k2']")
	assert.Equal(t, "vn1", subject.Nested.Map["kn1"], "Test.Nested.Map['kn1']")
	assert.Equal(t, "vn2", subject.Nested.Map["kn2"], "Test.Nested.Map['kn2']")
}

// TestObjectMapGet ensures decoding of maps of messages
func TestObjectMapGet(t *testing.T) {
	subject, err := buildSubjectGet(t, &Test{NestedNullableWithNilValue: &Nested{Str: "5"}})
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, "TestString1", subject.MapObject["obj1"].Str, "MapObject['obj1'].Str")
	assert.Equal(t, "TestString2", subject.MapObject["obj2"].Str, "MapObject['obj2'].Str")

	assert.Equal(t, "TestString1", subject.MapObjectNullable["obj1"].Str)
	assert.Equal(t, "TestString2", subject.MapObjectNullable["obj2"].Str)

	assert.Equal(t, "TestString1", subject.Nested.MapObjectNested["obj1"].Str)
	assert.Equal(t, "TestString2", subject.Nested.MapObjectNested["obj2"].Str)
}
