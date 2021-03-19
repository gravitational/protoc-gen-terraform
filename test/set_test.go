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

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	test = Test{
		Str:    "TestString",
		Int32:  2,
		Int64:  3,
		Float:  18.99,
		Double: 19.21,
		Bool:   true,
		Bytes:  []byte("TestBytes"),
	}
)

// buildSubjectSet builds Test struct from test fixture data
func buildSubjectSet(t *testing.T) (*schema.ResourceData, error) {
	subject := schema.TestResourceDataRaw(t, SchemaTest(), map[string]interface{}{})
	err := SetTestToResourceData(subject, &test)
	return subject, err
}

// TestElementariesSet ensures decoding of elementary types
func TestElementariesSet(t *testing.T) {
	subject, err := buildSubjectSet(t)
	require.NoError(t, err, "failed to marshal test data")

	str := subject.Get("str")
	assert.Equal(t, str, "TestString", "schema.ResourceData['str']")

	i32 := subject.Get("int32")
	assert.Equal(t, i32, 2, "schema.ResourceData['int32']")

	i64 := subject.Get("int64")
	assert.Equal(t, i64, 3, "schema.ResourceData['int32']")

}
