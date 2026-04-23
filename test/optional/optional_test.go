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

	// Populated map
	m := o.Attrs["optional_map"].(types.Map)
	require.False(t, m.Null)
	require.Len(t, m.Elems, 2)
	require.Equal(t, "val1", m.Elems["key1"].(types.String).Value)
	require.Equal(t, "val2", m.Elems["key2"].(types.String).Value)

	// Populated inner message
	inner := o.Attrs["optional_inner_message"].(types.Object)
	require.False(t, inner.Null)
	require.True(t, inner.Attrs["inner_bool"].(types.Bool).Value)

	// Populated list
	l := o.Attrs["string_list"].(types.List)
	require.False(t, l.Null)
	require.Len(t, l.Elems, 2)
	require.Equal(t, "test1", l.Elems[0].(types.String).Value)
	require.Equal(t, "test2", l.Elems[1].(types.String).Value)
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

	diags := CopyOptionalTestToTerraform(context.Background(), obj, &o)
	require.False(t, diags.HasError())

	// Optional fields with nil should be null
	require.True(t, o.Attrs["optional_str"].(types.String).Null)
	require.True(t, o.Attrs["optional_int64"].(types.Int64).Null)
	require.True(t, o.Attrs["optional_bool"].(types.Bool).Null)
	require.True(t, o.Attrs["optional_inner_message"].(types.Object).Null)

	// Nil map and slice set with Null true on the Terraform side.
	require.True(t, o.Attrs["optional_map"].(types.Map).Null)
	require.True(t, o.Attrs["string_list"].(types.List).Null)
}

func TestCopyFromOptionalFields(t *testing.T) {
	s, d := GenSchemaOptionalTest(context.Background())
	require.False(t, d.HasError())

	typ := s.AttributeType()
	obj, ok := typ.(types.ObjectType)
	require.True(t, ok)

	innerType := obj.AttrTypes["optional_inner_message"].(types.ObjectType)
	mapType := obj.AttrTypes["optional_map"].(types.MapType)
	listType := obj.AttrTypes["string_list"].(types.ListType)

	tf := types.Object{
		Null:    false,
		Unknown: false,
		Attrs: map[string]attr.Value{
			"optional_str":   types.String{Value: "test"},
			"optional_int64": types.Int64{Value: 42},
			"optional_bool":  types.Bool{Value: true},
			"choice_a":       types.String{Null: true},
			"choice_b":       types.String{Value: "picked_b"},
			"optional_map": types.Map{
				ElemType: mapType.ElemType,
				Elems: map[string]attr.Value{
					"key1": types.String{Value: "val1"},
					"key2": types.String{Value: "val2"},
				},
			},
			"optional_inner_message": types.Object{
				AttrTypes: innerType.AttrTypes,
				Attrs: map[string]attr.Value{
					"inner_bool": types.Bool{Value: true},
				},
			},
			"string_list": types.List{
				ElemType: listType.ElemType,
				Elems: []attr.Value{
					types.String{Value: "test1"},
					types.String{Value: "test2"},
				},
			},
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

	typ := s.AttributeType()
	objType, ok := typ.(types.ObjectType)
	require.True(t, ok)

	tf := types.Object{
		Null:    false,
		Unknown: false,
		Attrs: map[string]attr.Value{
			"optional_str":           types.String{Null: true},
			"optional_int64":         types.Int64{Null: true},
			"optional_bool":          types.Bool{Null: true},
			"choice_a":               types.String{Null: true},
			"choice_b":               types.String{Null: true},
			"optional_map":           types.Map{Null: true},
			"optional_inner_message": types.Object{Null: true},
			"string_list":            types.List{Null: true},
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
	require.Nil(t, obj.OptionalInnerMessage)

	// Null map and list decode to empty containers
	require.Empty(t, obj.OptionalMap)
	require.Empty(t, obj.StringList)
}
