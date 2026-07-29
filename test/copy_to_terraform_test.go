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
		types.StringValue("TestString"),
		o.Attributes()["nested"].(types.Object).Attributes()["str"].(types.String),
	)

	require.Equal(
		t,
		types.StringValue("TestString"),
		o.Attributes()["nested_nullable"].(types.Object).Attributes()["str"].(types.String),
	)

	require.True(t, o.Attributes()["nested_nullable_with_nil_value"].(types.Object).IsNull())
}

func TestCopyToList(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.StringValue("el1"),
		types.StringValue("el2"),
	}, o.Attributes()["string_list"].(types.List).Elements())

	require.Equal(t, types.ListValueMust(
		types.StringType,
		make([]attr.Value, 0),
	), o.Attributes()["string_list_empty"].(types.List))

	require.Equal(t, []attr.Value{
		types.StringValue("bytes1"),
		types.StringValue("bytes2"),
	}, o.Attributes()["bytes_list"].(types.List).Elements())

	require.Equal(t, []attr.Value{
		NewTime(timestamp),
		NewTime(timestamp),
	}, o.Attributes()["timestamp_list"].(types.List).Elements())

	require.Equal(t, []attr.Value{
		NewDuration(duration),
		NewDuration(duration),
	}, o.Attributes()["duration_custom_list"].(types.List).Elements())
}

func TestCopyTo_ChangeListSize(t *testing.T) {
	o := copyToTerraformObject(t)

	testObject := createTestObj()

	// Start with two elements.
	diags := CopyTestToTerraform(context.Background(), testObject, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.StringValue("el1"),
		types.StringValue("el2"),
	}, o.Attributes()["string_list"].(types.List).Elements())

	// Increase to 3, array access must not panic.
	testObject.StringList = []string{"el1", "el2", "el3"}
	diags = CopyTestToTerraform(context.Background(), testObject, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.StringValue("el1"),
		types.StringValue("el2"),
		types.StringValue("el3"),
	}, o.Attributes()["string_list"].(types.List).Elements())

	// Decrease to a single element, others should be removed.
	testObject.StringList = []string{"elX"}
	diags = CopyTestToTerraform(context.Background(), testObject, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.StringValue("elX"),
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
		types.StringValue("Test"),
		firstEl.Attributes()["str"],
	)

	nestedNestedList := o.Attributes()["nested_list"].(types.List).Elements()[0].(types.Object).Attributes()["nested_list"].(types.List)

	require.Len(t, nestedNestedList.Elements(), 2)
	require.Equal(
		t,
		types.StringValue("Test1"),
		nestedNestedList.Elements()[0].(types.Object).Attributes()["str"],
	)
	require.Equal(
		t,
		types.StringValue("Test2"),
		nestedNestedList.Elements()[1].(types.Object).Attributes()["str"],
	)

	nestedMap := firstEl.Attributes()["map"].(types.Map)

	require.Equal(
		t,
		types.StringValue("value1"),
		nestedMap.Elements()["key1"].(types.String),
	)
	require.Equal(
		t,
		types.StringValue("value2"),
		nestedMap.Elements()["key2"].(types.String),
	)

	nestedMapObject := firstEl.Attributes()["map_object_nested"].(types.Map)

	require.Equal(
		t,
		types.StringValue("Test1"),
		nestedMapObject.Elements()["key1"].(types.Object).Attributes()["str"].(types.String),
	)
	require.Equal(
		t,
		types.StringValue("Test2"),
		nestedMapObject.Elements()["key2"].(types.Object).Attributes()["str"].(types.String),
	)

	nestedListNullable := o.Attributes()["nested_list_nullable"].(types.List)

	require.Len(t, nestedListNullable.Elements(), 1)
	require.Equal(
		t,
		types.StringValue("Test"),
		nestedListNullable.Elements()[0].(types.Object).Attributes()["str"],
	)
}

func TestCopyToMap(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	m := o.Attributes()["map"].(types.Map).Elements()

	require.Equal(t, types.StringValue("value1"), m["key1"].(types.String))
	require.Equal(t, types.StringValue("value2"), m["key2"].(types.String))
}

func TestCopyToMapPreserveUnknownExistingKeys(t *testing.T) {
	o := copyToTerraformObject(t)

	mapType, ok := o.AttributeTypes(t.Context())["map"].(types.MapType)
	require.True(t, ok)

	o.Attributes()["map"] = types.MapValueMust(
		mapType.ElementType(),
		map[string]attr.Value{
			"key1":  types.StringUnknown(),
			"stale": types.StringUnknown(),
		},
	)

	diags := CopyTestToTerraformPreserveUnknown(context.Background(), createTestObj(), &o, true)
	requireNoDiagErrors(t, diags)

	m := o.Attributes()["map"].(types.Map).Elements()
	require.Len(t, m, 2)

	require.True(t, m["key1"].(types.String).IsUnknown())
	require.Equal(t, types.StringValue("value2"), m["key2"].(types.String))
	require.NotContains(t, m, "stale")
}

func TestCopyToCustom(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(
		t,
		[]attr.Value{
			types.BoolValue(false),
			types.BoolValue(false),
			types.BoolValue(true),
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
		types.StringValue("Test"),
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
		types.Int64Value(5),
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
		types.StringValue("1234"),
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
		types.StringValue("TestString"),
		nestedNullable.Attributes()["str"].(types.String),
	)
}

func TestCopyToTerraformPreserveUnknown(t *testing.T) {
	o := copyToTerraformObject(t)
	o.Attributes()["str"] = types.StringUnknown()

	diags := CopyTestToTerraformPreserveUnknown(context.Background(), createTestObj(), &o, true)
	requireNoDiagErrors(t, diags)

	v := o.Attributes()["str"].(types.String)
	require.True(t, v.IsUnknown())
	require.False(t, v.IsNull())
}

func TestCopyToTerraformPreserveUnknownNested(t *testing.T) {
	o := copyToTerraformObject(t)

	nestedType, ok := o.AttributeTypes(t.Context())["nested"].(types.ObjectType)
	require.True(t, ok)

	nestedListType, ok := nestedType.AttributeTypes()["nested_list"].(types.ListType)
	require.True(t, ok)

	mapType, ok := nestedType.AttributeTypes()["map"].(types.MapType)
	require.True(t, ok)

	mapObjectNestedType, ok := nestedType.AttributeTypes()["map_object_nested"].(types.MapType)
	require.True(t, ok)

	nestedValue, diags := types.ObjectValue(
		nestedType.AttributeTypes(),
		map[string]attr.Value{
			"str":               types.StringUnknown(),
			"nested_list":       types.ListUnknown(nestedListType.ElementType()),
			"map":               types.MapUnknown(mapType.ElementType()),
			"map_object_nested": types.MapUnknown(mapObjectNestedType.ElementType()),
		},
	)
	requireNoDiagErrors(t, diags)
	o.Attributes()["nested"] = nestedValue

	diags = CopyTestToTerraformPreserveUnknown(context.Background(), createTestObj(), &o, true)
	requireNoDiagErrors(t, diags)

	nested := o.Attributes()["nested"].(types.Object)
	require.False(t, nested.IsUnknown())

	str := nested.Attributes()["str"].(types.String)
	require.True(t, str.IsUnknown())

	nestedList := nested.Attributes()["nested_list"].(types.List)
	require.True(t, nestedList.IsUnknown())

	nestedMap := nested.Attributes()["map"].(types.Map)
	require.True(t, nestedMap.IsUnknown())

	nestedMapObject := nested.Attributes()["map_object_nested"].(types.Map)
	require.True(t, nestedMapObject.IsUnknown())
}

func TestCopyToCustomPreserveUnknown(t *testing.T) {
	o := copyToTerraformObject(t)
	o.Attributes()["bool_custom_list"] = types.List{
		Unknown: false,
		Elems: []attr.Value{
			types.BoolValue(false),
			types.BoolUnknown(),
			types.BoolUnknown(),
		},
		ElemType: types.BoolType,
	}

	diags := CopyTestToTerraformPreserveUnknown(context.Background(), createTestObj(), &o, true)
	requireNoDiagErrors(t, diags)

	v := o.Attributes()["bool_custom_list"].(types.List)

	require.False(t, v.IsUnknown())
	require.Equal(t,
		[]attr.Value{
			types.BoolValue(false),
			types.BoolUnknown(),
			types.BoolUnknown(),
		},
		v.Elements(),
	)

	o.Attributes()["bool_custom_list"] = types.List{
		Unknown:  true,
		ElemType: types.BoolType,
	}

	diags = CopyTestToTerraformPreserveUnknown(context.Background(), createTestObj(), &o, true)
	requireNoDiagErrors(t, diags)

	v = o.Attributes()["bool_custom_list"].(types.List)

	require.True(t, v.IsUnknown())
}
