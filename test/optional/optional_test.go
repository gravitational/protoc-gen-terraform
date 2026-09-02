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
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func schemaObject(t *testing.T) types.Object {
	s, d := GenSchemaOptionalTest(context.Background())
	require.False(t, d.HasError())

	obj, ok := s.Type().(types.ObjectType)
	require.True(t, ok)

	attrTypes := obj.AttributeTypes()

	return types.ObjectValueMust(attrTypes, nullAttrs(attrTypes))
}

func nullAttrs(attrTypes map[string]attr.Type) map[string]attr.Value {
	attrs := make(map[string]attr.Value, len(attrTypes))

	for name, attrType := range attrTypes {
		attrs[name] = nullValue(attrType)
	}

	return attrs
}

func nullValue(attrType attr.Type) attr.Value {
	if attrType.Equal(types.StringType) {
		return types.StringNull()
	}
	if attrType.Equal(types.Int64Type) {
		return types.Int64Null()
	}
	if attrType.Equal(types.Float64Type) {
		return types.Float64Null()
	}
	if attrType.Equal(types.BoolType) {
		return types.BoolNull()
	}
	switch t := attrType.(type) {
	case types.ListType:
		return types.ListNull(t.ElementType())
	case types.MapType:
		return types.MapNull(t.ElementType())
	case types.ObjectType:
		return types.ObjectNull(t.AttributeTypes())
	default:
		panic(fmt.Sprintf("unsupported fixture attr type %T", attrType))
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

	_, ok = s.Attributes["optional_map"]
	require.True(t, ok)

	_, ok = s.Attributes["optional_inner_message"]
	require.True(t, ok)

	_, ok = s.Attributes["string_list"]
	require.True(t, ok)
}

func TestCopyToOptionalFieldsSet(t *testing.T) {
	o := schemaObject(t)
	obj := &OptionalTest{
		OptionalStr:          proto.String("world"),
		OptionalInt64:        proto.Int64(42),
		OptionalBool:         proto.Bool(true),
		RealOneOf:            &OptionalTest_ChoiceB{ChoiceB: "picked_b"},
		OptionalMap:          map[string]string{"key1": "val1", "key2": "val2"},
		OptionalInnerMessage: &InnerMessage{InnerBool: proto.Bool(true)},
		StringList:           []string{"test1", "test2"},
	}

	o, diags := CopyOptionalTestToTerraform(context.Background(), obj, &o)
	require.False(t, diags.HasError())

	// Optional fields with values set
	require.Equal(t, "world", o.Attributes()["optional_str"].(types.String).ValueString())
	require.False(t, o.Attributes()["optional_str"].(types.String).IsNull())

	require.Equal(t, int64(42), o.Attributes()["optional_int64"].(types.Int64).ValueInt64())
	require.False(t, o.Attributes()["optional_int64"].(types.Int64).IsNull())

	require.True(t, o.Attributes()["optional_bool"].(types.Bool).ValueBool())
	require.False(t, o.Attributes()["optional_bool"].(types.Bool).IsNull())

	// Real oneof
	require.Equal(t, "picked_b", o.Attributes()["choice_b"].(types.String).ValueString())

	// Populated map
	m := o.Attributes()["optional_map"].(types.Map)
	require.False(t, m.IsNull())
	require.Len(t, m.Elements(), 2)
	require.Equal(t, "val1", m.Elements()["key1"].(types.String).ValueString())
	require.Equal(t, "val2", m.Elements()["key2"].(types.String).ValueString())

	// Populated inner message
	inner := o.Attributes()["optional_inner_message"].(types.Object)
	require.False(t, inner.IsNull())
	require.True(t, inner.Attributes()["inner_bool"].(types.Bool).ValueBool())

	// Populated list
	l := o.Attributes()["string_list"].(types.List)
	require.False(t, l.IsNull())
	require.Len(t, l.Elements(), 2)
	require.Equal(t, "test1", l.Elements()[0].(types.String).ValueString())
	require.Equal(t, "test2", l.Elements()[1].(types.String).ValueString())
}

func TestCopyToOptionalFieldsNil(t *testing.T) {
	o := schemaObject(t)
	obj := &OptionalTest{
		OptionalStr:          nil,
		OptionalInt64:        nil,
		OptionalBool:         nil,
		OptionalMap:          nil,
		OptionalInnerMessage: nil,
		StringList:           nil,
	}

	o, diags := CopyOptionalTestToTerraform(context.Background(), obj, &o)
	require.False(t, diags.HasError())

	// Optional fields with nil should be null
	require.True(t, o.Attributes()["optional_str"].(types.String).IsNull())
	require.True(t, o.Attributes()["optional_int64"].(types.Int64).IsNull())
	require.True(t, o.Attributes()["optional_bool"].(types.Bool).IsNull())
	require.True(t, o.Attributes()["optional_inner_message"].(types.Object).IsNull())

	// Nil map and slice are normalized to the empty value on the Terraform side.
	require.False(t, o.Attributes()["optional_map"].(types.Map).IsNull())
	require.False(t, o.Attributes()["string_list"].(types.List).IsNull())
}

func TestCopyFromOptionalFields(t *testing.T) {
	s, d := GenSchemaOptionalTest(context.Background())
	require.False(t, d.HasError())

	typ := s.Type()
	obj, ok := typ.(types.ObjectType)
	require.True(t, ok)

	innerType := obj.AttributeTypes()["optional_inner_message"].(types.ObjectType)
	mapType := obj.AttributeTypes()["optional_map"].(types.MapType)
	listType := obj.AttributeTypes()["string_list"].(types.ListType)

	tf := types.ObjectValueMust(
		obj.AttributeTypes(),
		map[string]attr.Value{
			"optional_str":   types.StringValue("test"),
			"optional_int64": types.Int64Value(42),
			"optional_bool":  types.BoolValue(true),
			"choice_a":       types.StringNull(),
			"choice_b":       types.StringValue("picked_b"),
			"optional_map": types.MapValueMust(
				mapType.ElementType(),
				map[string]attr.Value{
					"key1": types.StringValue("val1"),
					"key2": types.StringValue("val2"),
				},
			),
			"optional_inner_message": types.ObjectValueMust(
				innerType.AttributeTypes(),
				map[string]attr.Value{
					"inner_bool": types.BoolValue(true),
				},
			),
			"string_list": types.ListValueMust(
				listType.ElementType(),
				[]attr.Value{
					types.StringValue("test1"),
					types.StringValue("test2"),
				},
			),
		},
	)

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

	// Populated map
	require.Equal(t, map[string]string{"key1": "val1", "key2": "val2"}, optionalTest.OptionalMap)

	// Populated inner message
	require.NotNil(t, optionalTest.OptionalInnerMessage)
	require.True(t, *optionalTest.OptionalInnerMessage.InnerBool)

	// Populated list
	require.Equal(t, []string{"test1", "test2"}, optionalTest.StringList)
}

func TestCopyFromOptionalFieldsNull(t *testing.T) {
	s, d := GenSchemaOptionalTest(context.Background())
	require.False(t, d.HasError())

	typ := s.Type()
	objType, ok := typ.(types.ObjectType)
	require.True(t, ok)

	innerType := objType.AttributeTypes()["optional_inner_message"].(types.ObjectType)
	mapType := objType.AttributeTypes()["optional_map"].(types.MapType)
	listType := objType.AttributeTypes()["string_list"].(types.ListType)

	tf := types.ObjectValueMust(
		objType.AttributeTypes(),
		map[string]attr.Value{
			"optional_str":           types.StringNull(),
			"optional_int64":         types.Int64Null(),
			"optional_bool":          types.BoolNull(),
			"choice_a":               types.StringNull(),
			"choice_b":               types.StringNull(),
			"optional_map":           types.MapNull(mapType.ElementType()),
			"optional_inner_message": types.ObjectNull(innerType.AttributeTypes()),
			"string_list":            types.ListNull(listType.ElementType()),
		},
	)

	obj := &OptionalTest{}
	diags := CopyOptionalTestFromTerraform(context.Background(), tf, obj)
	require.False(t, diags.HasError())

	// Null optional fields should remain nil pointers
	require.Nil(t, obj.OptionalStr)
	require.Nil(t, obj.OptionalInt64)
	require.Nil(t, obj.OptionalBool)

	require.Nil(t, obj.RealOneOf)
	require.Nil(t, obj.OptionalInnerMessage)

	// Null map and list decode to empty containers
	require.Empty(t, obj.OptionalMap)
	require.Empty(t, obj.StringList)
}
