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

func TestCopyFromPrimitives(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, "Test", target.Str)
	require.Equal(t, int32(98), target.Int32)
	require.Equal(t, int64(99), target.Int64)
	require.Equal(t, float32(0.75), target.Float)
	require.Equal(t, float64(0.76), target.Double)
	require.Equal(t, true, target.Bool)
	require.Equal(t, []byte("Test"), target.Bytes)
	require.Equal(t, Mode_ON, target.Mode)
}

func TestCopyFromTime(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{TimestampNullableWithNilValue: &timestamp}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, timestamp, target.Timestamp)
	require.Equal(t, timestamp, *target.TimestampNullable)
	require.Nil(t, target.TimestampNullableWithNilValue)

	require.Equal(t, duration, target.DurationStandard)
	require.Equal(t, Duration(duration), target.DurationCustom)
}

func TestCopyFromNested(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{NestedNullableWithNilValue: &Nested{}}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, "Test", target.Nested.Str)
	require.Equal(t, "Test", target.NestedNullable.Str)

	require.Nil(t, target.NestedNullableWithNilValue)
}

func TestCopyFromList(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, []string{"el1", "el2"}, target.StringList)
	require.Empty(t, target.StringListEmpty)
	require.Equal(t, [][]byte{[]byte("bytes1"), []byte("bytes2")}, target.BytesList)
	require.Equal(t, []*time.Time{&timestamp, &timestamp}, target.TimestampList)
	require.Equal(t, []Duration{Duration(duration), Duration(duration)}, target.DurationCustomList)
}

func TestCopyFromNestedList(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Len(t, target.NestedList, 1)
	require.Equal(t, "Test", target.NestedList[0].Str)
	require.Len(t, target.NestedList[0].NestedList, 2)
	require.Equal(t, "Test1", target.NestedList[0].NestedList[0].Str)
	require.Equal(t, "Test2", target.NestedList[0].NestedList[1].Str)
	require.Equal(t, "Value1", target.NestedList[0].Map["key1"])
	require.Equal(t, "Value2", target.NestedList[0].Map["key2"])

	require.Equal(t, "Test", target.NestedListNullable[0].Str)
}

func TestCopyFromMap(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, map[string]string{"key1": "Value1", "key2": "Value2"}, target.Map)
	require.Equal(t, map[string]string{"key1": "Value1", "key2": "Value2"}, target.Nested.Map)
	require.Equal(t, map[string]string{"key1": "Value1", "key2": "Value2"}, target.NestedNullable.Map)
}

func TestCopyFromNestedMap(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, "Test1", target.MapObject["key1"].Str)
	require.Equal(t, "Test2", target.MapObject["key2"].Str)
	require.Equal(t, "Test1", target.MapObjectNullable["key1"].Str)
	require.Equal(t, "Test2", target.MapObjectNullable["key2"].Str)
}

func TestCopyFromCustom(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, []BoolCustom{true, false, true}, target.BoolCustomList)
}

func TestCopyFromOneOfScalarBranch(t *testing.T) {
	obj := copyFromTerraformObject(t)
	obj.Attrs["branch3"] = types.String{Value: "Test"}

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, "Test", target.OneOf.(*Test_Branch3).Branch3)
}

func TestCopyFromOneOfObjectBranch(t *testing.T) {
	obj := copyFromTerraformObject(t)
	obj.Attrs["branch2"] = types.Object{
		Attrs: map[string]attr.Value{
			"int32": types.Int64{Value: 5},
		},
	}

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, int32(5), target.OneOf.(*Test_Branch2).Branch2.Int32)
}

func TestCopyFromOneOfObjectNoBranch(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, nil, target.OneOf)
}

func TestCopyFromEmbeddedField(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, "embdtest1", target.EmbeddedString)
	require.Equal(t, "embdtest2", target.EmbeddedNestedField.EmbeddedNestedString)
}

func TestCopyFromNullableEmbeddedField(t *testing.T) {
	obj := copyFromTerraformObject(t)

	target := Test{}
	require.False(t, CopyTestFromTerraform(context.Background(), obj, &target).HasError())

	require.Equal(t, Duration(5*time.Minute), target.Value)
}
