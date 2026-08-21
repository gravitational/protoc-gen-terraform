package resource

import (
	"testing"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/stretchr/testify/require"
)

func TestFieldBuildContextCustomTypeFromSchemaType(t *testing.T) {
	c := &FieldBuildContext{
		MessageBuildContext: MessageBuildContext{
			config: &Config{
				SchemaTypes: map[string]SchemaType{
					"Test.StringOverride": {
						CustomType: "StringCustom",
					},
				},
			},
		},
		field: &FieldDescriptorProtoExt{
			FieldDescriptorProto: &descriptor.FieldDescriptorProto{},
		},
		path:     "Test.StringOverride",
		typeName: "Test.StringOverride",
	}

	require.True(t, c.IsCustomType())
	require.Equal(t, "StringCustom", c.GetCustomType())
}
