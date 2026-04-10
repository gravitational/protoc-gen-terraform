/*
Copyright 2026 Gravitational, Inc.

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
package optional

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func schemaObject(t *testing.T) types.Object {
	t.Helper()
	s, d := GenSchemaOptionalTest(context.Background())
	require.False(t, d.HasError())
	typ := s.AttributeType()
	obj, ok := typ.(types.ObjectType)
	require.True(t, ok)
	return types.Object{
		Null:      false,
		Unknown:   false,
		Attrs:     make(map[string]attr.Value),
		AttrTypes: obj.AttrTypes,
	}
}

func TestSchemaHasOptionalFields(t *testing.T) {
	s, d := GenSchemaOptionalTest(context.Background())
	require.False(t, d.HasError())

	_, ok := s.Attributes["optional_str"]
	require.True(t, ok)

	_, ok = s.Attributes["optional_int64"]
	require.True(t, ok)

	_, ok = s.Attributes["optional_bool"]
	require.True(t, ok)

	_, ok = s.Attributes["choice_a"]
	require.True(t, ok)

	_, ok = s.Attributes["choice_b"]
	require.True(t, ok)
}

func TestCopyToOptionalFieldsSet(t *testing.T) {
	o := schemaObject(t)
	obj := &OptionalTest{
		OptionalStr:   proto.String("world"),
		OptionalInt64: proto.Int64(42),
		OptionalBool:  proto.Bool(true),
		RealOneOf:     &OptionalTest_ChoiceB{ChoiceB: "picked_b"},
	}

	diags := CopyOptionalTestToTerraform(context.Background(), obj, &o)
	require.False(t, diags.HasError())

	// Optional fields with values set
	require.Equal(t, "world", o.Attrs["optional_str"].(types.String).Value)
	require.False(t, o.Attrs["optional_str"].(types.String).Null)

	require.Equal(t, int64(42), o.Attrs["optional_int64"].(types.Int64).Value)
	require.False(t, o.Attrs["optional_int64"].(types.Int64).Null)

	require.True(t, o.Attrs["optional_bool"].(types.Bool).Value)
	require.False(t, o.Attrs["optional_bool"].(types.Bool).Null)

	// Real oneof
	require.Equal(t, "picked_b", o.Attrs["choice_b"].(types.String).Value)
}

func TestCopyToOptionalFieldsNil(t *testing.T) {
	o := schemaObject(t)
	obj := &OptionalTest{
		OptionalStr:   nil,
		OptionalInt64: nil,
		OptionalBool:  nil,
	}

	diags := CopyOptionalTestToTerraform(context.Background(), obj, &o)
	require.False(t, diags.HasError())

	// Optional fields with nil should be null
	require.True(t, o.Attrs["optional_str"].(types.String).Null)
	require.True(t, o.Attrs["optional_int64"].(types.Int64).Null)
	require.True(t, o.Attrs["optional_bool"].(types.Bool).Null)
}

func TestCopyFromOptionalFields(t *testing.T) {
	s, d := GenSchemaOptionalTest(context.Background())
	require.False(t, d.HasError())

	typ := s.AttributeType()
	obj, ok := typ.(types.ObjectType)
	require.True(t, ok)

	tf := types.Object{
		Null:    false,
		Unknown: false,
		Attrs: map[string]attr.Value{
			"optional_str":   types.String{Value: "test"},
			"optional_int64": types.Int64{Value: 42},
			"optional_bool":  types.Bool{Value: true},
			"choice_a":       types.String{Null: true},
			"choice_b":       types.String{Value: "picked_b"},
		},
		AttrTypes: obj.AttrTypes,
	}

	optionalTest := &OptionalTest{}

	diags := CopyOptionalTestFromTerraform(context.Background(), tf, optionalTest)
	require.False(t, diags.HasError())

	// Optional fields should have pointer values set
	require.NotNil(t, optionalTest.OptionalStr)
	require.Equal(t, "test", *optionalTest.OptionalStr)

	require.NotNil(t, optionalTest.OptionalInt64)
	require.Equal(t, int64(42), *optionalTest.OptionalInt64)

	require.NotNil(t, optionalTest.OptionalBool)
	require.True(t, *optionalTest.OptionalBool)

	// Real oneof
	choiceB, ok := optionalTest.RealOneOf.(*OptionalTest_ChoiceB)
	require.True(t, ok)
	require.Equal(t, "picked_b", choiceB.ChoiceB)

	_, ok = optionalTest.RealOneOf.(*OptionalTest_ChoiceA)
	require.False(t, ok)
}

func TestCopyFromOptionalFieldsNull(t *testing.T) {
	s, d := GenSchemaOptionalTest(context.Background())
	require.False(t, d.HasError())

	typ := s.AttributeType()
	objType, ok := typ.(types.ObjectType)
	require.True(t, ok)

	tf := types.Object{
		Null:    false,
		Unknown: false,
		Attrs: map[string]attr.Value{
			"optional_str":   types.String{Null: true},
			"optional_int64": types.Int64{Null: true},
			"optional_bool":  types.Bool{Null: true},
			"choice_a":       types.String{Null: true},
			"choice_b":       types.String{Null: true},
		},
		AttrTypes: objType.AttrTypes,
	}

	obj := &OptionalTest{}
	diags := CopyOptionalTestFromTerraform(context.Background(), tf, obj)
	require.False(t, diags.HasError())

	// Null optional fields should remain nil pointers
	require.Nil(t, obj.OptionalStr)
	require.Nil(t, obj.OptionalInt64)
	require.Nil(t, obj.OptionalBool)

	require.Nil(t, obj.RealOneOf)
}
