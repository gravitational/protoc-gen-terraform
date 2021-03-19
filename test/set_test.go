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
	fmt "fmt"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	test = Test{
		Str:     "TestString",
		Int32:   2,
		Int64:   3,
		Float:   18.5,
		Double:  19.21,
		Bool:    true,
		Bytes:   []byte("TestBytes"),
		StringA: []string{"TestString1", "TestString2"},
	}
)

// toSlice is a helper method which takes a slice of any type and converts it to []interface{}
func toSlice(s interface{}) ([]interface{}, error) {
	v := reflect.ValueOf(s)
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		result := make([]interface{}, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result, nil
	default:
		return nil, fmt.Errorf("failed to convert %T to []interface{}", s)
	}
}

// fillTimestamps parses time and duration from predefined strings and fills in correspoding fields in test structure
func fillTimestamps(t *Test) error {
	ti, err := time.Parse(time.RFC3339, defaultTimestamp)
	if err != nil {
		return err
	}

	d, err := time.ParseDuration("1h")
	if err != nil {
		return err
	}

	t.Timestamp = ti
	t.DurationStd = d
	t.DurationCustom = Duration(d)
	t.TimestampN = &ti

	return nil
}

// buildSubjectSet builds Test struct from test fixture data
func buildSubjectSet(t *testing.T) (*schema.ResourceData, error) {
	subject := schema.TestResourceDataRaw(t, SchemaTest(), map[string]interface{}{})
	err := fillTimestamps(&test)
	if err != nil {
		return nil, err
	}
	err = SetTestToResourceData(subject, &test)
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

// TestTimesSet ensures decoding of time and duration fields
func TestTimesSet(t *testing.T) {
	subject, err := buildSubjectSet(t)
	require.NoError(t, err, "failed to unmarshal test data")

	assert.Equal(t, test.Timestamp.Format(time.RFC3339), subject.Get("timestamp"), "Test.Timestamp")
	assert.Equal(t, test.DurationStd.String(), subject.Get("duration_std"), "Test.DurationStd")
	assert.Equal(t, test.DurationCustom.String(), subject.Get("duration_custom"), "Test.DurationCustom")
	assert.Equal(t, test.TimestampN.Format(time.RFC3339), subject.Get("timestamp_n"), "Test.TimestampN")
}

// TestArraysSet ensures decoding of arrays
func TestArraysSet(t *testing.T) {
	subject, err := buildSubjectSet(t)
	require.NoError(t, err, "failed to unmarshal test data")

	stringA, err := toSlice(test.StringA)
	require.NoError(t, err, "failed to convert string_a to []interface{}")

	assert.Equal(t, stringA, subject.Get("string_a"), "Test.StringA")
	// assert.Equal(t, subject.BoolA, []BoolCustom{false, true, false})
	// assert.Equal(t, subject.BytesA, [][]byte{[]byte("TestBytes1"), []byte("TestBytes2")})
	// assert.Equal(t, subject.TimestampA, []*time.Time{&timestamp})
	// assert.Equal(t, subject.DurationCustomA, []Duration{Duration(duration)})
}