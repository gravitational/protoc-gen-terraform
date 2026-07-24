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

package test

import (
	"context"
	"testing"
	time "time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCopyToTerraformPrimitives(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, "TestString", o.Attributes()["str"].(types.String).ValueString())
	require.False(t, o.Attributes()["str"].(types.String).IsUnknown())
	require.False(t, o.Attributes()["str"].(types.String).IsNull())

	require.Equal(t, int64(888), o.Attributes()["int32"].(types.Int64).ValueInt64())
	require.False(t, o.Attributes()["int32"].(types.Int64).IsUnknown())
	require.False(t, o.Attributes()["int32"].(types.Int64).IsNull())

	require.Equal(t, int64(999), o.Attributes()["int64"].(types.Int64).ValueInt64())
	require.False(t, o.Attributes()["int64"].(types.Int64).IsUnknown())
	require.False(t, o.Attributes()["int64"].(types.Int64).IsNull())

	require.Equal(t, float64(88.5), o.Attributes()["float"].(types.Float64).ValueFloat64())
	require.False(t, o.Attributes()["float"].(types.Float64).IsUnknown())
	require.False(t, o.Attributes()["float"].(types.Float64).IsNull())

	require.Equal(t, float64(99.5), o.Attributes()["double"].(types.Float64).ValueFloat64())
	require.False(t, o.Attributes()["double"].(types.Float64).IsUnknown())
	require.False(t, o.Attributes()["double"].(types.Float64).IsNull())

	require.True(t, o.Attributes()["bool"].(types.Bool).ValueBool())
	require.False(t, o.Attributes()["bool"].(types.Bool).IsUnknown())
	require.False(t, o.Attributes()["bool"].(types.Bool).IsNull())

	require.Equal(t, "TestBytes", o.Attributes()["bytes"].(types.String).ValueString())
	require.False(t, o.Attributes()["bytes"].(types.String).IsUnknown())
	require.False(t, o.Attributes()["bytes"].(types.String).IsNull())
}

func TestCopyToTime(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, timestamp, o.Attributes()["timestamp"].(TimeValue).ValueTime())
	require.False(t, o.Attributes()["timestamp"].(TimeValue).IsUnknown())
	require.False(t, o.Attributes()["timestamp"].(TimeValue).IsNull())

	require.Equal(t, time.Time{}, o.Attributes()["timestamp_missing"].(TimeValue).ValueTime())
	require.False(t, o.Attributes()["timestamp_missing"].(TimeValue).IsUnknown())
	// Handle empty time value
	// require.True(t, o.Attributes()["timestamp_missing"].(TimeValue).IsNull())
}

func TestCopyToDuration(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, duration, o.Attributes()["duration_standard"].(DurationValue).ValueDuration())
	require.False(t, o.Attributes()["duration_standard"].(DurationValue).IsUnknown())
	require.False(t, o.Attributes()["duration_standard"].(DurationValue).IsNull())

	require.Equal(t, duration, o.Attributes()["duration_custom"].(DurationValue).ValueDuration())
	require.False(t, o.Attributes()["duration_custom"].(DurationValue).IsUnknown())
	require.False(t, o.Attributes()["duration_custom"].(DurationValue).IsNull())
}

func TestCopyToNested(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "TestString"},
		o.Attributes()["nested"].(types.Object).Attributes()["str"].(types.String),
	)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "TestString"},
		o.Attributes()["nested_nullable"].(types.Object).Attributes()["str"].(types.String),
	)

	require.True(t, o.Attributes()["nested_nullable_with_nil_value"].(types.Object).IsNull())
}

func TestCopyToList(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.String{Null: false, Unknown: false, Value: "el1"},
		types.String{Null: false, Unknown: false, Value: "el2"},
	}, o.Attributes()["string_list"].(types.List).Elements())

	require.Equal(t, types.List{
		Null:     false,
		Unknown:  false,
		Elems:    make([]attr.Value, 0),
		ElemType: types.StringType,
	}, o.Attributes()["string_list_empty"].(types.List))

	require.Equal(t, []attr.Value{
		types.String{Null: false, Unknown: false, Value: "bytes1"},
		types.String{Null: false, Unknown: false, Value: "bytes2"},
	}, o.Attributes()["bytes_list"].(types.List).Elements())

	require.Equal(t, []attr.Value{
		TimeValue{Null: false, Unknown: false, Value: timestamp, Format: time.RFC3339},
		TimeValue{Null: false, Unknown: false, Value: timestamp, Format: time.RFC3339},
	}, o.Attributes()["timestamp_list"].(types.List).Elements())

	require.Equal(t, []attr.Value{
		DurationValue{Null: false, Unknown: false, Value: duration},
		DurationValue{Null: false, Unknown: false, Value: duration},
	}, o.Attributes()["duration_custom_list"].(types.List).Elements())
}

func TestCopyTo_ChangeListSize(t *testing.T) {
	o := copyToTerraformObject(t)

	testObject := createTestObj()

	// Start with two elements.
	diags := CopyTestToTerraform(context.Background(), testObject, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.String{Null: false, Unknown: false, Value: "el1"},
		types.String{Null: false, Unknown: false, Value: "el2"},
	}, o.Attributes()["string_list"].(types.List).Elements())

	// Increase to 3, array access must not panic.
	testObject.StringList = []string{"el1", "el2", "el3"}
	diags = CopyTestToTerraform(context.Background(), testObject, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.String{Null: false, Unknown: false, Value: "el1"},
		types.String{Null: false, Unknown: false, Value: "el2"},
		types.String{Null: false, Unknown: false, Value: "el3"},
	}, o.Attributes()["string_list"].(types.List).Elements())

	// Decrease to a single element, others should be removed.
	testObject.StringList = []string{"elX"}
	diags = CopyTestToTerraform(context.Background(), testObject, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.String{Null: false, Unknown: false, Value: "elX"},
	}, o.Attributes()["string_list"].(types.List).Elements())
}

func TestCopyToNestedList(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	nestedList := o.Attributes()["nested_list"].(types.List)
	firstEl := nestedList.Elements()[0].(types.Object)

	require.Len(t, nestedList.Elements(), 1)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test"},
		firstEl.Attributes()["str"],
	)

	nestedNestedList := o.Attributes()["nested_list"].(types.List).Elements()[0].(types.Object).Attributes()["nested_list"].(types.List)

	require.Len(t, nestedNestedList.Elements(), 2)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test1"},
		nestedNestedList.Elements()[0].(types.Object).Attributes()["str"],
	)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test2"},
		nestedNestedList.Elements()[1].(types.Object).Attributes()["str"],
	)

	nestedMap := firstEl.Attributes()["map"].(types.Map)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "value1"},
		nestedMap.Elements()["key1"].(types.String),
	)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "value2"},
		nestedMap.Elements()["key2"].(types.String),
	)

	nestedMapObject := firstEl.Attributes()["map_object_nested"].(types.Map)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test1"},
		nestedMapObject.Elements()["key1"].(types.Object).Attributes()["str"].(types.String),
	)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test2"},
		nestedMapObject.Elements()["key2"].(types.Object).Attributes()["str"].(types.String),
	)

	nestedListNullable := o.Attributes()["nested_list_nullable"].(types.List)

	require.Len(t, nestedListNullable.Elements(), 1)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test"},
		nestedListNullable.Elements()[0].(types.Object).Attributes()["str"],
	)
}

func TestCopyToMap(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	m := o.Attributes()["map"].(types.Map).Elements()

	require.Equal(t, types.String{Null: false, Unknown: false, Value: "value1"}, m["key1"].(types.String))
	require.Equal(t, types.String{Null: false, Unknown: false, Value: "value2"}, m["key2"].(types.String))
}

func TestCopyToCustom(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(
		t,
		[]attr.Value{
			types.Bool{Null: false, Unknown: false, Value: false},
			types.Bool{Null: false, Unknown: false, Value: false},
			types.Bool{Null: false, Unknown: false, Value: true},
		},
		o.Attributes()["bool_custom_list"].(types.List).Elements(),
	)
}

func TestCopyToOneOfBranch3(t *testing.T) {
	o := copyToTerraformObject(t)
	testObj := createTestObj()
	testObj.OneOf = &Test_Branch3{Branch3: "Test"}

	diags := CopyTestToTerraform(context.Background(), testObj, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test"},
		o.Attributes()["branch3"].(types.String),
	)
}

func TestCopyToOneOfBranch2(t *testing.T) {
	o := copyToTerraformObject(t)
	testObj := createTestObj()
	testObj.OneOf = &Test_Branch2{Branch2: &Branch2{Int32: 5}}

	diags := CopyTestToTerraform(context.Background(), testObj, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(
		t,
		types.Int64{Null: false, Unknown: false, Value: 5},
		o.Attributes()["branch2"].(types.Object).Attributes()["int32"],
	)
}

func TestCopyToOneOfNoBranch(t *testing.T) {
	o := copyToTerraformObject(t)
	testObj := createTestObj()

	diags := CopyTestToTerraform(context.Background(), testObj, &o)
	requireNoDiagErrors(t, diags)

	require.True(t, o.Attributes()["branch1"].(types.Object).IsNull())
	require.True(t, o.Attributes()["branch2"].(types.Object).IsNull())
	require.True(t, o.Attributes()["branch3"].(types.String).IsNull())
}

func TestCopyToEmbeddedField(t *testing.T) {
	o := copyToTerraformObject(t)
	testObj := createTestObj()

	diags := CopyTestToTerraform(context.Background(), testObj, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, "embdtest1", o.Attributes()["embedded_string"].(types.String).ValueString())
	require.False(t, o.Attributes()["embedded_string"].(types.String).IsUnknown())
	require.False(t, o.Attributes()["embedded_string"].(types.String).IsNull())

	require.Equal(t, "embdtest2", o.Attributes()["embedded_nested_field"].(types.Object).Attributes()["embedded_nested_string"].(types.String).ValueString())
}

func TestCopyToOneOfLowercase(t *testing.T) {
	o := copyToTerraformObject(t)
	testObj := createTestObj()
	testObj.LowerSnakeOneof = &Test_Foo{Foo: "1234"}

	diags := CopyTestToTerraform(context.Background(), testObj, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "1234"},
		o.Attributes()["foo"].(types.String),
	)
}

func TestCopyToNestedNullableWithNullTerraformObject(t *testing.T) {
	o := copyToTerraformObject(t)

	nestedNullableType, ok := o.AttributeTypes(t.Context())["nested_nullable"].(types.ObjectType)
	require.True(t, ok)

	o.Attributes()["nested_nullable"] = types.Object{
		Null:      true,
		AttrTypes: nestedNullableType.AttributeTypes(),
	}

	testObj := createTestObj()
	testObj.NestedNullable = &Nested{
		Str: "TestString",
	}

	diags := CopyTestToTerraform(context.Background(), testObj, &o)
	requireNoDiagErrors(t, diags)

	nestedNullable := o.Attributes()["nested_nullable"].(types.Object)
	require.False(t, nestedNullable.IsNull())
	require.False(t, nestedNullable.IsUnknown())
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "TestString"},
		nestedNullable.Attributes()["str"].(types.String),
	)
}

func TestCopyToTerraformPreserveUnknown(t *testing.T) {
	o := copyToTerraformObject(t)
	o.Attributes()["str"] = types.String{Unknown: true, Value: "stale"}

	diags := CopyTestToTerraformPreserveUnknown(context.Background(), createTestObj(), &o, true)
	requireNoDiagErrors(t, diags)

	v := o.Attributes()["str"].(types.String)
	require.True(t, v.IsUnknown())
	require.False(t, v.IsNull())
	require.Equal(t, "TestString", v.ValueString())
}

func TestCopyToTerraformPreserveUnknownNested(t *testing.T) {
	o := copyToTerraformObject(t)

	nestedType, ok := o.AttributeTypes(t.Context())["nested"].(types.ObjectType)
	require.True(t, ok)

	nestedListType, ok := nestedType.AttributeTypes()["nested_list"].(types.ListType)
	require.True(t, ok)

	o.Attributes()["nested"] = types.Object{
		Unknown:   true,
		AttrTypes: nestedType.AttributeTypes(),
		Attrs: map[string]attr.Value{
			"str": types.String{Unknown: true, Value: "stale"},
			"nested_list": types.List{
				Unknown:  true,
				ElemType: nestedListType.ElementType(),
				Elems: []attr.Value{
					types.Object{
						Unknown:   true,
						AttrTypes: nestedListType.ElementType().(types.ObjectType).AttributeTypes(),
						Attrs: map[string]attr.Value{
							"str": types.String{Unknown: true, Value: "stale"},
						},
					},
				},
			},
		},
	}

	diags := CopyTestToTerraformPreserveUnknown(context.Background(), createTestObj(), &o, true)
	requireNoDiagErrors(t, diags)

	nested := o.Attributes()["nested"].(types.Object)
	require.True(t, nested.IsUnknown())

	str := nested.Attributes()["str"].(types.String)
	require.True(t, str.IsUnknown())
	require.Equal(t, "TestString", str.ValueString())

	nestedList := nested.Attributes()["nested_list"].(types.List)
	require.True(t, nestedList.IsUnknown())
	require.Len(t, nestedList.Elements(), 2)

	firstElem := nestedList.Elements()[0].(types.Object)
	require.True(t, firstElem.IsUnknown())
	require.Equal(t, "Test1", firstElem.Attributes()["str"].(types.String).ValueString())
	require.True(t, firstElem.Attributes()["str"].(types.String).IsUnknown())

	secondElem := nestedList.Elements()[1].(types.Object)
	require.False(t, secondElem.IsUnknown())
	require.Equal(t, "Test2", secondElem.Attributes()["str"].(types.String).ValueString())
	require.False(t, secondElem.Attributes()["str"].(types.String).IsUnknown())
}

func TestCopyToCustomPreserveUnknown(t *testing.T) {
	o := copyToTerraformObject(t)
	o.Attributes()["bool_custom_list"] = types.List{
		Unknown: true,
		Elems: []attr.Value{
			types.Bool{Unknown: false, Value: false},
			types.Bool{Unknown: true, Value: false},
			types.Bool{Unknown: true, Value: true},
		},
		ElemType: types.BoolType,
	}

	diags := CopyTestToTerraformPreserveUnknown(context.Background(), createTestObj(), &o, true)
	requireNoDiagErrors(t, diags)

	v := o.Attributes()["bool_custom_list"].(types.List)

	require.True(t, v.IsUnknown())
	require.Equal(t,
		[]attr.Value{
			types.Bool{Null: false, Unknown: false, Value: false},
			types.Bool{Null: false, Unknown: true, Value: false},
			types.Bool{Null: false, Unknown: true, Value: true},
		},
		v.Elements(),
	)
}
