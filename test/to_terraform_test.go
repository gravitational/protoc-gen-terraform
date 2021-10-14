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
	"sort"
	"testing"
	time "time"

	"github.com/gravitational/trace"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	test = Test{
		Str:            "TestString",
		Int32:          2,
		Int64:          3,
		Float:          18.5,
		Double:         19.21,
		Bool:           true,
		Bytes:          []byte("TestBytes"),
		StringList:     []string{"TestString1", "TestString2"},
		BoolCustomList: []BoolCustom{false, true, false},
		BytesList:      [][]byte{[]byte("TestBytes1"), []byte("TestBytes2")},
		Nested: Nested{
			Str:        "TestStringA",
			NestedList: []*OtherNested{{Str: "NestedString1"}, {Str: "NestedString2"}},
			Map: map[string]string{
				"kn1": "vn1",
				"kn2": "vn2",
			},
		},
		NestedNullable: &Nested{
			Str:        "TestStringA",
			NestedList: []*OtherNested{{Str: "NestedString1"}, {Str: "NestedString2"}},
			Map: map[string]string{
				"kn1": "vn1",
				"kn2": "vn2",
			},
		},
		Map: map[string]string{
			"k1": "v1",
			"k2": "v2",
		},

		MapObject: map[string]Nested{
			"n1": {
				Str: "NestedObjString1",
			},
			"n2": {
				Str: "NestedObjString2",
			},
			"n3": {
				Str: "NestedObjString3",
			},
		},

		MapObjectNullable: map[string]*Nested{
			"n1": {
				Str: "NestedObjString1",
			},
			"n2": {
				Str: "NestedObjString2",
			},
			"n3": {
				Str: "NestedObjString3",
			},
		},
	}
)

// fillTimestamps parses time and duration from predefined strings and fills in correspoding fields in test structure
func fillTimestamps(t *Test) error {
	ti, err := time.Parse(time.RFC3339Nano, defaultTimestamp)
	if err != nil {
		return trace.Wrap(err)
	}

	d, err := time.ParseDuration("1h")
	if err != nil {
		return trace.Wrap(err)
	}

	t.Timestamp = ti
	t.TimestampNullable = &ti
	t.DurationStandard = d
	t.DurationCustom = Duration(d)
	t.TimestampList = []*time.Time{&ti, &ti}
	t.DurationCustomList = []Duration{Duration(d), Duration(d)}

	return nil
}

// buildSubjectToTerraform builds Test struct from test fixture data
func buildSubjectToTerraform(t *testing.T) (*schema.ResourceData, error) {
	subject := schema.TestResourceDataRaw(t, SchemaTest, map[string]interface{}{
		"string_list":          []interface{}{"nil", "nil"},
		"bytes_list":           []interface{}{"bytes", "bytes"},
		"timestamp_list":       []interface{}{"nil"},
		"duration_custom_list": []interface{}{"nil"},
		"nested": []interface{}{
			map[string]interface{}{
				"nested_list": []interface{}{
					map[string]interface{}{},
					map[string]interface{}{},
				},
				"map": map[string]interface{}{
					"nil1": "nil1",
					"nil2": "nil2",
				},
				"map_object_nested": []interface{}{
					map[string]interface{}{
						"key": "obj1",
						"value": []interface{}{
							map[string]interface{}{},
						},
					},
					map[string]interface{}{
						"key": "obj2",
						"value": []interface{}{
							map[string]interface{}{},
						},
					},
				},
			},
		},
		"nested_nullable": []interface{}{
			map[string]interface{}{},
		},

		// List of messages
		"nested_list": []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		},
		"nested_list_nullable": []interface{}{
			map[string]interface{}{},
			map[string]interface{}{},
		},

		"map": map[string]interface{}{
			"k1": "nil1",
			"k2": "nil2",
		},

		"map_object": []interface{}{
			map[string]interface{}{
				"key": "obj1",
				"value": []interface{}{
					map[string]interface{}{},
				},
			},
			map[string]interface{}{
				"key": "obj2",
				"value": []interface{}{
					map[string]interface{}{},
				},
			},
		},

		"map_object_nullable": []interface{}{
			map[string]interface{}{
				"key": "obj1",
				"value": []interface{}{
					map[string]interface{}{},
				},
			},
			map[string]interface{}{
				"key": "obj2",
				"value": []interface{}{
					map[string]interface{}{},
				},
			},
		},
	})

	err := fillTimestamps(&test)
	if err != nil {
		return nil, err
	}

	err = ToTerraformTest(&test, subject)
	return subject, err
}

// TestElementariesToTerraform ensures decoding of elementary types
func TestElementariesToTerraform(t *testing.T) {
	subject, err := buildSubjectToTerraform(t)
	require.NoError(t, err, "failed to marshal test data")

	assert.Equal(t, subject.Get("str"), "TestString", "schema.ResourceData['str']")
	assert.Equal(t, subject.Get("int32"), 2, "schema.ResourceData['int32']")
	assert.Equal(t, subject.Get("int64"), 3, "schema.ResourceData['int64']")
	assert.Equal(t, subject.Get("float"), 18.5, "schema.ResourceData['float']")
	assert.Equal(t, subject.Get("double"), 19.21, "schema.ResourceData['d']")
	assert.Equal(t, subject.Get("bool"), true, "schema.ResourceData['bool']")
	assert.Equal(t, subject.Get("bytes"), "TestBytes", "schema.ResourceData['bytes']")
}

// TestTimesToTerraform ensures decoding of time and duration fields
func TestTimesToTerraform(t *testing.T) {
	subject, err := buildSubjectToTerraform(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, test.Timestamp.Format(time.RFC3339Nano), subject.Get("timestamp"), "Test.Timestamp")
	assert.Equal(t, test.DurationStandard.String(), subject.Get("duration_standard"), "Test.DurationStandard")
	assert.Equal(t, time.Duration(test.DurationCustom).String(), subject.Get("duration_custom"), "Test.DurationCustom")
	assert.Equal(t, test.TimestampNullable.Format(time.RFC3339Nano), subject.Get("timestamp_nullable"), "Test.TimestampNullable")
}

// TestArraysToTerraform ensures decoding of arrays
func TestArraysToTerraform(t *testing.T) {
	subject, err := buildSubjectToTerraform(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, []interface{}{"TestString1", "TestString2"}, subject.Get("string_list"), "Test.StringList")
	assert.Equal(t, []interface{}{false, true, false}, subject.Get("bool_custom_list"), "Test.BoolCustomList")
	assert.Equal(t, []interface{}{"TestBytes1", "TestBytes2"}, subject.Get("bytes_list"), "Test.BytesList")
	assert.Equal(t, []interface{}{"1h0m0s", "1h0m0s"}, subject.Get("duration_custom_list"))

	raw := subject.Get("timestamp_list")
	a, ok := raw.([]interface{})
	if !ok {
		assert.Fail(t, "can not convert %T to []interface{}", raw)
	}

	for n, v := range a {
		assert.Equal(t, v, test.TimestampList[n].Format(time.RFC3339Nano), "Test.TimestampA[]")
	}
}

// TestNestedMessageToTerraform ensures decoding of nested messages
func TestNestedMessageToTerraform(t *testing.T) {
	subject, err := buildSubjectToTerraform(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, test.Nested.Str, subject.Get("nested.0.str"), "Test.Nested.Str")
	assert.Equal(t, test.NestedNullable.Str, subject.Get("nested_nullable.0.str"), "Test.NestedNullable.Str")
}

// TestNestedMessageArrayToTerraform ensures decoding of array of messages
func TestNestedMessageArrayToTerraform(t *testing.T) {
	subject, err := buildSubjectToTerraform(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, test.Nested.NestedList[0].Str, subject.Get("nested.0.nested_list.0.str"))
	assert.Equal(t, test.Nested.NestedList[1].Str, subject.Get("nested.0.nested_list.1.str"))
}

// TestMapToTerraform ensures decoding of a maps
func TestMapToTerraform(t *testing.T) {
	subject, err := buildSubjectToTerraform(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, test.Map["k1"], subject.Get("map.k1"))
	assert.Equal(t, test.Map["k2"], subject.Get("map.k2"))
	assert.Equal(t, test.Nested.Map["kn1"], "vn1")
	assert.Equal(t, test.Nested.Map["kn1"], "vn1")
}

// TestObjectMapToTerraform ensures decoding of maps of messages
func TestObjectMapToTerraform(t *testing.T) {
	subject, err := buildSubjectToTerraform(t)
	require.NoError(t, err, "failed to unmarshal test data")

	m := subject.Get("map_object")
	s, ok := m.(*schema.Set)
	if !ok {
		assert.Fail(t, "can not convert %T to *schema.Set", s)
	}

	keys := make([]string, s.Len())

	for i, n := range s.List() {
		v, ok := n.(map[string]interface{})
		if !ok {
			assert.Fail(t, "can not convert %T to map[string]interface{}", s)
		}

		s, ok := v["key"].(string)
		if !ok {
			assert.Fail(t, "can not convert %T to string", s)
		}

		keys[i] = s
	}

	sort.Strings(keys)

	assert.Equal(t, keys, []string{"n1", "n2", "n3"})
}
