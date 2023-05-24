/*
Copyright 2023 Gravitational, Inc.

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

package main

import (
	"testing"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/stretchr/testify/require"
)

func TestFieldBuildContextGetName(t *testing.T) {
	tests := []struct {
		name      string
		embedded  bool
		fieldName string
		goType    string
		expected  string
	}{
		{
			name:      "regular name",
			fieldName: "Name",
			expected:  "Name",
		},
		{
			name:      "embedded name",
			embedded:  true,
			fieldName: "Name",
			goType:    "EmbeddedName",
			expected:  "EmbeddedName",
		},
		{
			name:      "embedded name with pointer",
			embedded:  true,
			fieldName: "Name",
			goType:    "*EmbeddedName",
			expected:  "EmbeddedName",
		},
		{
			name:      "embedded name in another package",
			embedded:  true,
			fieldName: "Name",
			goType:    "*someotherpackage.EmbeddedName",
			expected:  "EmbeddedName",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			msg := &descriptor.FieldDescriptorProto{
				Name: &test.fieldName,
			}

			fbc := FieldBuildContext{
				field:   &FieldDescriptorProtoExt{msg},
				goType:  test.goType,
				isEmbed: test.embedded,
			}
			require.Equal(t, test.expected, fbc.GetName())
		})
	}
}
