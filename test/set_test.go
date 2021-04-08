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
	time "time"

	"github.com/gravitational/protoc-gen-terraform/accessors"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	test = Test{
		Str:    "TestString",
		Int32:  2,
		Int64:  3,
		Float:  18.5,
		Double: 19.21,
		Bool:   true,
		Bytes:  []byte("TestBytes"),
		// StringA: []string{"TestString1", "TestString2"},
		// BoolA:   []BoolCustom{false, true, false},
		// BytesA:  [][]byte{[]byte("TestBytes1"), []byte("TestBytes2")},
		// Nested: &Nested{
		// 	Str:    "TestStringA",
		// 	Nested: []*NestedLevel2{{Str: "NestedString1"}, {Str: "NestedString2"}},
		// 	NestedM: map[string]string{
		// 		"kn1": "vn1",
		// 		"kn2": "vn2",
		// 	},
		// },
		// NestedM: map[string]string{
		// 	"k1": "v1",
		// 	"k2": "v2",
		// },
		// NestedMObj: map[string]*Nested{
		// 	"n1": {
		// 		Str: "NestedObjString1",
		// 	},
		// 	"n2": {
		// 		Str: "NestedObjString2",
		// 	},
		// 	"n3": {
		// 		Str: "NestedObjString3",
		// 	},
		// },
	}
)

// // fillTimestamps parses time and duration from predefined strings and fills in correspoding fields in test structure
func fillTimestamps(t *Test) error {
	ti, err := time.Parse(time.RFC3339Nano, defaultTimestamp)
	if err != nil {
		return err
	}

	d, err := time.ParseDuration("1h")
	if err != nil {
		return err
	}

	t.Timestamp = ti
	t.DurationStandard = d
	t.DurationCustom = Duration(d)
	// t.TimestampNullable = &ti
	// t.TimestampA = []*time.Time{&ti, &ti}
	// t.DurationCustomA = []Duration{Duration(d), Duration(d)}

	return nil
}

// buildSubjectSet builds Test struct from test fixture data
func buildSubjectSet(t *testing.T) (*schema.ResourceData, error) {
	subject := schema.TestResourceDataRaw(t, SchemaTest, map[string]interface{}{})
	err := fillTimestamps(&test)
	if err != nil {
		return nil, err
	}

	err = accessors.Set(&test, subject, SchemaTest, SchemaMetaTest)
	return subject, err
}

// TestElementariesSet ensures decoding of elementary types
func TestElementariesSet(t *testing.T) {
	subject, err := buildSubjectSet(t)
	require.NoError(t, err, "failed to marshal test data")

	assert.Equal(t, subject.Get("str"), "TestString", "schema.ResourceData['str']")
	assert.Equal(t, subject.Get("int32"), 2, "schema.ResourceData['int32']")
	assert.Equal(t, subject.Get("int64"), 3, "schema.ResourceData['int64']")
	assert.Equal(t, subject.Get("float"), 18.5, "schema.ResourceData['float']")
	assert.Equal(t, subject.Get("double"), 19.21, "schema.ResourceData['d']")
	assert.Equal(t, subject.Get("bool"), true, "schema.ResourceData['bool']")
	assert.Equal(t, subject.Get("bytes"), "TestBytes", "schema.ResourceData['bytes']")
}

// // TestTimesSet ensures decoding of time and duration fields
// func TestTimesSet(t *testing.T) {
// 	subject, err := buildSubjectSet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, test.Timestamp.Format(time.RFC3339Nano), subject.Get("timestamp"), "Test.Timestamp")
// 	assert.Equal(t, test.DurationStd.String(), subject.Get("duration_std"), "Test.DurationStd")
// 	assert.Equal(t, time.Duration(test.DurationCustom).String(), subject.Get("duration_custom"), "Test.DurationCustom")
// 	assert.Equal(t, test.TimestampN.Format(time.RFC3339Nano), subject.Get("timestamp_n"), "Test.TimestampN")
// }

// // TestArraysSet ensures decoding of arrays
// func TestArraysSet(t *testing.T) {
// 	subject, err := buildSubjectSet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, []interface{}{"TestString1", "TestString2"}, subject.Get("string_a"), "Test.StringA")
// 	assert.Equal(t, []interface{}{false, true, false}, subject.Get("bool_a"), "Test.BoolA")
// 	assert.Equal(t, []interface{}{"TestBytes1", "TestBytes2"}, subject.Get("bytes_a"), "Test.BytesA")
// 	assert.Equal(t, []interface{}{"1h0m0s", "1h0m0s"}, subject.Get("duration_custom_a"))

// 	raw := subject.Get("timestamp_a")
// 	a, ok := raw.([]interface{})
// 	if !ok {
// 		assert.Fail(t, "can not convert %T to []interface{}", raw)
// 	}

// 	for n, v := range a {
// 		assert.Equal(t, v, test.TimestampA[n].Format(time.RFC3339Nano), "Test.TimestampA[]")
// 	}
// }

// // TestNestedMessageSet ensures decoding of nested messages
// func TestNestedMessageSet(t *testing.T) {
// 	subject, err := buildSubjectSet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, test.Nested.Str, subject.Get("nested.0.str"), "Test.Nested.Str")
// }

// // TestNestedMessageArraySet ensures decoding of array of messages
// func TestNestedMessageArraySet(t *testing.T) {
// 	subject, err := buildSubjectSet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, test.Nested.Nested[0].Str, subject.Get("nested.0.nested.0.str"))
// 	assert.Equal(t, test.Nested.Nested[1].Str, subject.Get("nested.0.nested.1.str"))
// }

// // TestMapGet ensures decoding of a maps
// func TestMapSet(t *testing.T) {
// 	subject, err := buildSubjectSet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, test.NestedM["k1"], subject.Get("nested_m.k1"))
// 	assert.Equal(t, test.NestedM["k2"], subject.Get("nested_m.k2"))
// 	assert.Equal(t, test.Nested.NestedM["kn1"], "vn1")
// 	assert.Equal(t, test.Nested.NestedM["kn1"], "vn1")
// }

// // TestObjectMapGet ensures decoding of maps of messages
// func TestObjectMapSet(t *testing.T) {
// 	subject, err := buildSubjectSet(t)
// 	require.NoError(t, err, "failed to unmarshal test data")

// 	assert.Equal(t, test.NestedMObj["n1"].Str, subject.Get("nested_m_obj.0.value.0.str"))
// 	assert.Equal(t, test.NestedMObj["n2"].Str, subject.Get("nested_m_obj.1.value.0.str"))
// 	assert.Equal(t, test.NestedMObj["n3"].Str, subject.Get("nested_m_obj.2.value.0.str"))
// }
