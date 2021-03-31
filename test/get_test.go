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
		"str":               "TestString",
		"int32":             999,
		"int64":             998,
		"float":             18.1,
		"double":            18.4,
		"bool":              true,
		"bytes":             "TestBytes",
		"timestamp":         defaultTimestamp,
		"duration_std":      "1h",
		"duration_custom":   "1m",
		"timestamp_n":       defaultTimestamp,
		"string_a":          []interface{}{"TestString1", "TestString2"},
		"bool_a":            []interface{}{false, true, false},
		"bytes_a":           []interface{}{"TestBytes1", "TestBytes2"},
		"timestamp_a":       []interface{}{defaultTimestamp},
		"duration_custom_a": []interface{}{"1m"},

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
		},
	}
)

// buildSubjectGet builds Test struct from test fixture data
func buildSubjectGet(t *testing.T) (*Test, error) {
	subject := &Test{}
	data := schema.TestResourceDataRaw(t, SchemaTest(), fixture)
	err := GetTestFromResourceData(data, subject)
	return subject, err
}

// TestElementariesGet ensures decoding of elementary types
func TestElementariesGet(t *testing.T) {
	subject, err := buildSubjectGet(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, subject.Str, "TestString", "Test.Str")
	assert.Equal(t, subject.Int32, int32(999), "Test.Int32")
	assert.Equal(t, subject.Int64, int64(998), "Test.Int64")
	assert.Equal(t, subject.Float, float32(18.1), "Test.Float")
	assert.Equal(t, subject.Double, float64(18.4), "Test.Dobule")
	assert.Equal(t, subject.Bool, true, "Test.Bool")
	assert.Equal(t, subject.Bytes, []byte("TestBytes"), "Test.Bytes")
}

// TestTimesGet ensures decoding of time and duration fields
func TestTimesGet(t *testing.T) {
	subject, err := buildSubjectGet(t)
	require.NoError(t, err, "failed to unmarshal test data")

	timestamp, err := time.Parse(time.RFC3339Nano, defaultTimestamp)
	require.NoError(t, err, "failed to parse example timestamp")

	durationStd, err := time.ParseDuration("1h")
	require.NoError(t, err, "failed to parse example duration")

	durationCustom, err := time.ParseDuration("1m")
	require.NoError(t, err, "failed to parse example duration")

	assert.Equal(t, subject.Timestamp, timestamp, "Test.Timestamp")
	assert.Equal(t, subject.DurationStd, durationStd, "Test.DurationStd")
	assert.Equal(t, subject.DurationCustom, Duration(durationCustom), "Test.DurationCustom")
	assert.Equal(t, *(subject.TimestampN), timestamp, "Test.TimestampN")
}

// TestArraysGet ensures decoding of arrays
func TestArraysGet(t *testing.T) {
	subject, err := buildSubjectGet(t)
	require.NoError(t, err, "failed to unmarshal test data")

	timestamp, err := time.Parse(time.RFC3339Nano, defaultTimestamp)
	require.NoError(t, err, "failed to parse example timestamp")

	duration, err := time.ParseDuration("1m")
	require.NoError(t, err, "failed to parse example duration")

	assert.Equal(t, subject.StringA, []string{"TestString1", "TestString2"})
	assert.Equal(t, subject.BoolA, []BoolCustom{false, true, false})
	assert.Equal(t, subject.BytesA, [][]byte{[]byte("TestBytes1"), []byte("TestBytes2")})
	assert.Equal(t, subject.TimestampA, []*time.Time{&timestamp})
	assert.Equal(t, subject.DurationCustomA, []Duration{Duration(duration)})
}

// TestNestedMessageGet ensures decoding of nested messages
func TestNestedMessageGet(t *testing.T) {
	subject, err := buildSubjectGet(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, subject.Nested.Str, "TestString", "Test.Nested.Str")
}

// TestNestedMessageArrayGet ensures decoding of array of messages
func TestNestedMessageArrayGet(t *testing.T) {
	subject, err := buildSubjectGet(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, subject.Nested.Nested[0].Str, "TestString1")
	assert.Equal(t, subject.Nested.Nested[1].Str, "TestString2")
}

// TestMapGet ensures decoding of a maps
func TestMapGet(t *testing.T) {
	subject, err := buildSubjectGet(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, subject.NestedM["k1"], "v1")
	assert.Equal(t, subject.NestedM["k2"], "v2")
	assert.Equal(t, subject.Nested.NestedM["kn1"], "vn1")
	assert.Equal(t, subject.Nested.NestedM["kn1"], "vn1")
}

// TestObjectMapGet ensures decoding of maps of messages
func TestObjectMapGet(t *testing.T) {
	subject, err := buildSubjectGet(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, subject.NestedMObj["obj1"].Str, "TestString1")
}
