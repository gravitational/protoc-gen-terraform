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

	require.Equal(t, "TestString", o.Attrs["str"].(types.String).Value)
	require.False(t, o.Attrs["str"].(types.String).Unknown)
	require.False(t, o.Attrs["str"].(types.String).Null)

	require.Equal(t, int64(888), o.Attrs["int32"].(types.Int64).Value)
	require.False(t, o.Attrs["int32"].(types.Int64).Unknown)
	require.False(t, o.Attrs["int32"].(types.Int64).Null)

	require.Equal(t, int64(999), o.Attrs["int64"].(types.Int64).Value)
	require.False(t, o.Attrs["int64"].(types.Int64).Unknown)
	require.False(t, o.Attrs["int64"].(types.Int64).Null)

	require.Equal(t, float64(88.5), o.Attrs["float"].(types.Float64).Value)
	require.False(t, o.Attrs["float"].(types.Float64).Unknown)
	require.False(t, o.Attrs["float"].(types.Float64).Null)

	require.Equal(t, float64(99.5), o.Attrs["double"].(types.Float64).Value)
	require.False(t, o.Attrs["double"].(types.Float64).Unknown)
	require.False(t, o.Attrs["double"].(types.Float64).Null)

	require.True(t, o.Attrs["bool"].(types.Bool).Value)
	require.False(t, o.Attrs["bool"].(types.Bool).Unknown)
	require.False(t, o.Attrs["bool"].(types.Bool).Null)

	require.Equal(t, "TestBytes", o.Attrs["bytes"].(types.String).Value)
	require.False(t, o.Attrs["bytes"].(types.String).Unknown)
	require.False(t, o.Attrs["bytes"].(types.String).Null)
}

func TestCopyToTime(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, timestamp, o.Attrs["timestamp"].(TimeValue).Value)
	require.False(t, o.Attrs["timestamp"].(TimeValue).Unknown)
	require.False(t, o.Attrs["timestamp"].(TimeValue).Null)

	require.Equal(t, time.Time{}, o.Attrs["timestamp_missing"].(TimeValue).Value)
	require.False(t, o.Attrs["timestamp_missing"].(TimeValue).Unknown)
	// Handle empty time value
	// require.True(t, o.Attrs["timestamp_missing"].(TimeValue).Null)
}

func TestCopyToDuration(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, duration, o.Attrs["duration_standard"].(DurationValue).Value)
	require.False(t, o.Attrs["duration_standard"].(DurationValue).Unknown)
	require.False(t, o.Attrs["duration_standard"].(DurationValue).Null)

	require.Equal(t, duration, o.Attrs["duration_custom"].(DurationValue).Value)
	require.False(t, o.Attrs["duration_custom"].(DurationValue).Unknown)
	require.False(t, o.Attrs["duration_custom"].(DurationValue).Null)
}

func TestCopyToNested(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "TestString"},
		o.Attrs["nested"].(types.Object).Attrs["str"].(types.String),
	)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "TestString"},
		o.Attrs["nested_nullable"].(types.Object).Attrs["str"].(types.String),
	)

	require.True(t, o.Attrs["nested_nullable_with_nil_value"].(types.Object).Null)
}

func TestCopyToList(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.String{Null: false, Unknown: false, Value: "el1"},
		types.String{Null: false, Unknown: false, Value: "el2"},
	}, o.Attrs["string_list"].(types.List).Elems)

	require.Equal(t, types.List{
		Null:     true,
		Unknown:  false,
		Elems:    make([]attr.Value, 0),
		ElemType: types.StringType,
	}, o.Attrs["string_list_empty"].(types.List))

	require.Equal(t, []attr.Value{
		types.String{Null: false, Unknown: false, Value: "bytes1"},
		types.String{Null: false, Unknown: false, Value: "bytes2"},
	}, o.Attrs["bytes_list"].(types.List).Elems)

	require.Equal(t, []attr.Value{
		TimeValue{Null: false, Unknown: false, Value: timestamp, Format: time.RFC3339},
		TimeValue{Null: false, Unknown: false, Value: timestamp, Format: time.RFC3339},
	}, o.Attrs["timestamp_list"].(types.List).Elems)

	require.Equal(t, []attr.Value{
		DurationValue{Null: false, Unknown: false, Value: duration},
		DurationValue{Null: false, Unknown: false, Value: duration},
	}, o.Attrs["duration_custom_list"].(types.List).Elems)
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
	}, o.Attrs["string_list"].(types.List).Elems)

	// Increase to 3, array access must not panic.
	testObject.StringList = []string{"el1", "el2", "el3"}
	diags = CopyTestToTerraform(context.Background(), testObject, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.String{Null: false, Unknown: false, Value: "el1"},
		types.String{Null: false, Unknown: false, Value: "el2"},
		types.String{Null: false, Unknown: false, Value: "el3"},
	}, o.Attrs["string_list"].(types.List).Elems)

	// Decrease to a single element, others should be removed.
	testObject.StringList = []string{"elX"}
	diags = CopyTestToTerraform(context.Background(), testObject, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, []attr.Value{
		types.String{Null: false, Unknown: false, Value: "elX"},
	}, o.Attrs["string_list"].(types.List).Elems)
}

func TestCopyToNestedList(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	nestedList := o.Attrs["nested_list"].(types.List)
	firstEl := nestedList.Elems[0].(types.Object)

	require.Len(t, nestedList.Elems, 1)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test"},
		firstEl.Attrs["str"],
	)

	nestedNestedList := o.Attrs["nested_list"].(types.List).Elems[0].(types.Object).Attrs["nested_list"].(types.List)

	require.Len(t, nestedNestedList.Elems, 2)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test1"},
		nestedNestedList.Elems[0].(types.Object).Attrs["str"],
	)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test2"},
		nestedNestedList.Elems[1].(types.Object).Attrs["str"],
	)

	nestedMap := firstEl.Attrs["map"].(types.Map)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "value1"},
		nestedMap.Elems["key1"].(types.String),
	)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "value2"},
		nestedMap.Elems["key2"].(types.String),
	)

	nestedMapObject := firstEl.Attrs["map_object_nested"].(types.Map)

	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test1"},
		nestedMapObject.Elems["key1"].(types.Object).Attrs["str"].(types.String),
	)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test2"},
		nestedMapObject.Elems["key2"].(types.Object).Attrs["str"].(types.String),
	)

	nestedListNullable := o.Attrs["nested_list_nullable"].(types.List)

	require.Len(t, nestedListNullable.Elems, 1)
	require.Equal(
		t,
		types.String{Null: false, Unknown: false, Value: "Test"},
		nestedListNullable.Elems[0].(types.Object).Attrs["str"],
	)
}

func TestCopyToMap(t *testing.T) {
	o := copyToTerraformObject(t)

	diags := CopyTestToTerraform(context.Background(), createTestObj(), &o)
	requireNoDiagErrors(t, diags)

	m := o.Attrs["map"].(types.Map).Elems

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
		o.Attrs["bool_custom_list"].(types.List).Elems,
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
		o.Attrs["branch3"].(types.String),
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
		o.Attrs["branch2"].(types.Object).Attrs["int32"],
	)
}

func TestCopyToOneOfNoBranch(t *testing.T) {
	o := copyToTerraformObject(t)
	testObj := createTestObj()

	diags := CopyTestToTerraform(context.Background(), testObj, &o)
	requireNoDiagErrors(t, diags)

	require.True(t, o.Attrs["branch1"].(types.Object).Null)
	require.True(t, o.Attrs["branch2"].(types.Object).Null)
	require.True(t, o.Attrs["branch3"].(types.String).Null)
}

func TestCopyToEmbeddedField(t *testing.T) {
	o := copyToTerraformObject(t)
	testObj := createTestObj()

	diags := CopyTestToTerraform(context.Background(), testObj, &o)
	requireNoDiagErrors(t, diags)

	require.Equal(t, "embdtest1", o.Attrs["embedded_string"].(types.String).Value)
	require.False(t, o.Attrs["embedded_string"].(types.String).Unknown)
	require.False(t, o.Attrs["embedded_string"].(types.String).Null)

	require.Equal(t, "embdtest2", o.Attrs["embedded_nested_field"].(types.Object).Attrs["embedded_nested_string"].(types.String).Value)
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
		o.Attrs["foo"].(types.String),
	)
}
