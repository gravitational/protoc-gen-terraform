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
	"testing"
	time "time"

	"github.com/stretchr/testify/require"
)

func TestCopyFromTerraformPrimitives(t *testing.T) {
	target := Test{}
	require.NoError(t, CopyTestFromTerraform(copyFromTerraformObject, &target))

	require.Equal(t, "Test", target.Str)
	require.Equal(t, int32(98), target.Int32)
	require.Equal(t, int64(99), target.Int64)
	require.Equal(t, float32(0.75), target.Float)
	require.Equal(t, float64(0.76), target.Double)
	require.Equal(t, true, target.Bool)
	require.Equal(t, []byte("Test"), target.Bytes)
}

func TestTestCopyFromTerraformTimestamps(t *testing.T) {
	target := Test{TimestampNullableWithNilValue: &timestamp}
	require.NoError(t, CopyTestFromTerraform(copyFromTerraformObject, &target))

	require.Equal(t, timestamp, target.Timestamp)
	require.Equal(t, timestamp, *target.TimestampNullable)
	require.Nil(t, target.TimestampNullableWithNilValue)

	require.Equal(t, duration, target.DurationStandard)
	require.Equal(t, Duration(duration), target.DurationCustom)
}

func TestCopyFromNested(t *testing.T) {
	target := Test{NestedNullableWithNilValue: &Nested{}}
	require.NoError(t, CopyTestFromTerraform(copyFromTerraformObject, &target))

	require.Equal(t, "Test", target.Nested.Str)
	require.Equal(t, "Test", target.NestedNullable.Str)

	require.Nil(t, target.NestedNullableWithNilValue)
}

func TestCopyFromList(t *testing.T) {
	target := Test{}
	require.NoError(t, CopyTestFromTerraform(copyFromTerraformObject, &target))

	require.Equal(t, []string{"el1", "el2"}, target.StringList)
	require.Empty(t, target.StringListEmpty)
	require.Equal(t, [][]byte{[]byte("bytes1"), []byte("bytes2")}, target.BytesList)
	require.Equal(t, []*time.Time{&timestamp, &timestamp}, target.TimestampList)
	require.Equal(t, []Duration{Duration(duration), Duration(duration)}, target.DurationCustomList)
}

func TestCopyFromNestedList(t *testing.T) {
	target := Test{}
	require.NoError(t, CopyTestFromTerraform(copyFromTerraformObject, &target))

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
	target := Test{}
	require.NoError(t, CopyTestFromTerraform(copyFromTerraformObject, &target))

	require.Equal(t, map[string]string{"key1": "Value1", "key2": "Value2"}, target.Map)
	require.Equal(t, map[string]string{"key1": "Value1", "key2": "Value2"}, target.Nested.Map)
	require.Equal(t, map[string]string{"key1": "Value1", "key2": "Value2"}, target.NestedNullable.Map)
}

func TestCopyFromNestedMap(t *testing.T) {
	target := Test{}
	require.NoError(t, CopyTestFromTerraform(copyFromTerraformObject, &target))

	require.Equal(t, "Test1", target.MapObject["key1"].Str)
	require.Equal(t, "Test2", target.MapObject["key2"].Str)
	require.Equal(t, "Test1", target.MapObjectNullable["key1"].Str)
	require.Equal(t, "Test2", target.MapObjectNullable["key2"].Str)
}

func TestCustom(t *testing.T) {
	target := Test{}
	require.NoError(t, CopyTestFromTerraform(copyFromTerraformObject, &target))

	require.Equal(t, []BoolCustom{true, false, true}, target.BoolCustomList)
}
