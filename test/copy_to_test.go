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

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
)

func TestCopyToElementaries(t *testing.T) {
	o := types.Object{}
	require.NoError(t, copier.CopyWithOption(&o, &copyToTerraformObject, copier.Option{DeepCopy: true}))

	err := CopyTestToTerraform(&o, testObj)
	require.NoError(t, err)

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
	o := types.Object{}
	require.NoError(t, copier.CopyWithOption(&o, &copyToTerraformObject, copier.Option{DeepCopy: true}))

	err := CopyTestToTerraform(&o, testObj)
	require.NoError(t, err)

	require.Equal(t, timestamp, o.Attrs["timestamp"].(TimeValue).Value)
	require.False(t, o.Attrs["timestamp"].(TimeValue).Unknown)
	require.False(t, o.Attrs["timestamp"].(TimeValue).Null)

	require.Equal(t, time.Time{}, o.Attrs["timestamp_missing"].(TimeValue).Value)
	require.False(t, o.Attrs["timestamp_missing"].(TimeValue).Unknown)
	require.False(t, o.Attrs["timestamp_missing"].(TimeValue).Null) // We use Terraform as the source of truth
}

// func TestTimeCopyFrom(t *testing.T) {
// 	source := createObj()
// 	target := &Test{}

// 	err := target.CopyFrom(&source)
// 	require.NoError(t, err)

// 	require.Equal(t, timestamp, target.Timestamp.Value)
// 	require.Equal(t, false, target.Timestamp.Null)
// 	require.Equal(t, timestamp, target.TimestampNullable.Value)
// 	require.Equal(t, false, target.TimestampNullable.Null)
// 	require.Equal(t, true, target.TimestampNullableWithNilValue.Null)
// }

// func TestDurationCopyFrom(t *testing.T) {
// 	source := createObj()
// 	target := &Test{}

// 	err := target.CopyFrom(&source)
// 	require.NoError(t, err)

// 	require.Equal(t, duration, target.DurationStandard.Value)
// 	require.Equal(t, time.Duration(0), target.DurationStandardMissing.Value)
// 	require.Equal(t, duration, target.DurationCustom.Value)
// 	require.Equal(t, time.Duration(0), target.DurationCustomMissing.Value)
// }

// func TestSliceCopyFrom(t *testing.T) {
// 	source := createObj()
// 	target := &Test{}

// 	err := target.CopyFrom(&source)
// 	require.NoError(t, err)

// 	require.Equal(t, []string{"Element1", "Element2"}, []string{target.StringList[0].Value, target.StringList[1].Value})
// 	require.Equal(t, []bool{false, false}, []bool{target.StringList[0].Null, target.StringList[1].Null})

// 	require.Len(t, target.StringListEmpty, 0)

// 	require.Equal(t, []bool{true, false, true}, []bool{bool(target.BoolCustomList[0]), bool(target.BoolCustomList[1]), bool(target.BoolCustomList[2])})

// 	require.Equal(t, []string{"Element1", "Element2"}, []string{target.BytesList[0].Value, target.BytesList[1].Value})
// 	require.Equal(t, []bool{false, false}, []bool{target.BytesList[0].Null, target.BytesList[1].Null})

// 	require.Equal(
// 		t,
// 		[]time.Time{timestamp, timestamp},
// 		[]time.Time{target.TimestampList[0].Value, target.TimestampList[1].Value},
// 	)

// 	require.Equal(
// 		t,
// 		[]time.Duration{duration, duration},
// 		[]time.Duration{target.DurationCustomList[0].Value, target.DurationCustomList[1].Value},
// 	)

// 	require.Equal(t, "NestedElement1", target.NestedListNullable[0].Str.Value)
// 	require.Equal(t, "NestedElement2", target.NestedListNullable[1].Str.Value)
// }

// func TestNestedCopyFrom(t *testing.T) {
// 	source := createObj()
// 	target := &Test{}

// 	err := target.CopyFrom(&source)
// 	require.NoError(t, err)

// 	require.Equal(t, "TestStr", target.Nested.Str.Value)
// 	require.Equal(t, "OtherStr1", target.Nested.NestedList[0].Str.Value)
// 	require.Equal(t, "OtherStr2", target.Nested.NestedList[1].Str.Value)
// 	require.Equal(t, "value1", target.Nested.Map["key1"].Value)
// 	require.Equal(t, "value2", target.Nested.Map["key2"].Value)

// 	require.Equal(t, "TestStr", target.NestedNullable.Str.Value)

// 	require.Equal(t, "OtherStr1", target.Nested.NestedList[0].Str.Value)
// 	require.Equal(t, "OtherStr2", target.Nested.NestedList[1].Str.Value)
// }

// func TestMapObjectCopyFrom(t *testing.T) {
// 	source := createObj()
// 	target := &Test{}

// 	err := target.CopyFrom(&source)
// 	require.NoError(t, err)

// 	require.Equal(t, "Value1", target.MapObject["key1"].Str.Value)
// 	require.Equal(t, "Value2", target.MapObject["key2"].Str.Value)
// 	require.Equal(t, "OtherStr1", target.Nested.MapObjectNested["key1"].Str.Value)
// 	require.Equal(t, "OtherStr2", target.Nested.MapObjectNested["key2"].Str.Value)
// }

// // FromTerraformObject() (terraform -> apitypes)
// //   - copies known values to an object
// //   - nullifies unknown
// // ToTerraformObject() (apitypes -> terraform)
// //   - updates values of a terraform object
