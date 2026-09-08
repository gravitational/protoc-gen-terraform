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
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	diag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

var (
	timestamp = time.Now()
	duration  = 5 * time.Minute
)

func createTestObj() *Test {
	return &Test{
		Str:    "TestString",
		Int32:  888,
		Int64:  999,
		Float:  88.5,
		Double: 99.5,
		Bool:   true,
		Bytes:  []byte("TestBytes"),

		Timestamp:                     timestamp,
		TimestampMissing:              time.Time{},
		TimestampNullable:             &timestamp,
		TimestampNullableWithNilValue: nil,
		DurationStandard:              duration,
		DurationStandardMissing:       0,
		DurationCustom:                Duration(duration),
		DurationCustomMissing:         Duration(duration),

		StringList:      []string{"el1", "el2"},
		StringListEmpty: nil,
		BytesList:       [][]byte{[]byte("bytes1"), []byte("bytes2")},

		TimestampList:      []*time.Time{&timestamp, &timestamp},
		DurationCustomList: []Duration{Duration(duration), Duration(duration)},
		BoolCustomList:     []BoolCustom{false, false, true},

		Nested: Nested{
			Str: "TestString",
			NestedList: []*OtherNested{
				{Str: "Test1"},
				{Str: "Test2"},
			},
		},
		NestedNullable: &Nested{
			Str: "TestString",
		},
		NestedNullableWithNilValue: nil,

		NestedList: []Nested{
			{
				Str: "Test",
				NestedList: []*OtherNested{
					{
						Str: "Test1",
					},
					{
						Str: "Test2",
					},
				},
				Map: map[string]string{"key1": "value1", "key2": "value2"},
				MapObjectNested: map[string]OtherNested{
					"key1": {Str: "Test1"},
					"key2": {Str: "Test2"},
				},
			},
		},

		NestedListNullable: []*Nested{
			{
				Str: "Test",
			},
		},

		Map: map[string]string{"key1": "value1", "key2": "value2"},

		EmbeddedField: EmbeddedField{
			EmbeddedString: "embdtest1",
			EmbeddedNestedField: &EmbeddedNestedField{
				EmbeddedNestedString: "embdtest2",
			},
		},

		SchemaOverride: "hello",
	}
}

// copyFromTerraformObject returns a base object used in CopyFrom* tests
func copyFromTerraformObject(t *testing.T) types.Object {
	t.Helper()

	s, d := GenSchemaTestResource(context.Background())

	require.False(t, d.HasError())
	typ := s.Type()

	obj, ok := typ.(types.ObjectType)
	require.True(t, ok)

	must := func(v any, d diag.Diagnostics) attr.Value {
		t.Helper()
		requireNoDiagErrors(t, d)
		return v.(attr.Value)
	}

	attrs := map[string]attr.Value{
		"id":     types.StringNull(),
		"str":    types.StringValue("Test"),
		"int32":  types.Int64Value(98),
		"int64":  types.Int64Value(99),
		"float":  types.Float64Value(0.75),
		"double": types.Float64Value(0.76),
		"bool":   types.BoolValue(true),
		"bytes":  types.StringValue("Test"),

		"timestamp":                         NewTime(timestamp),
		"timestamp_missing":                 UnknownTime(),
		"timestamp_nullable":                NewTime(timestamp),
		"timestamp_nullable_with_nil_value": NullTime(),
		"duration_standard":                 NewDuration(duration),
		"duration_standard_missing":         UnknownDuration(),
		"duration_custom":                   NewDuration(duration),
		"duration_custom_missing":           UnknownDuration(),

		"string_list": must(types.ListValue(
			types.StringType,
			[]attr.Value{types.StringValue("el1"), types.StringValue("el2")},
		)),

		"string_list_empty": must(types.ListValue(types.StringType, []attr.Value{})),
		"bytes_list": must(types.ListValue(
			types.StringType,
			[]attr.Value{types.StringValue("bytes1"), types.StringValue("bytes2")},
		)),

		"timestamp_list": must(types.ListValue(
			UseRFC3339Time(),
			[]attr.Value{NewTime(timestamp), NewTime(timestamp)},
		)),
		"duration_custom_list": must(types.ListValue(
			DurationType{},
			[]attr.Value{NewDuration(duration), NewDuration(duration)},
		)),

		"bool_custom_list": must(types.ListValue(
			types.BoolType,
			[]attr.Value{types.BoolValue(true), types.BoolValue(false), types.BoolValue(true)},
		)),

		"nested": must(types.ObjectValue(
			obj.AttrTypes["nested"].(types.ObjectType).AttrTypes,
			map[string]attr.Value{
				"str": types.StringValue("Test"),
				"map": must(types.MapValue(types.StringType, map[string]attr.Value{
					"key1": types.StringValue("Value1"),
					"key2": types.StringValue("Value2"),
				})),
				"nested_list": types.ListNull(
					obj.AttrTypes["nested"].(types.ObjectType).
						AttrTypes["nested_list"].(types.ListType).
						ElemType,
				),
				"map_object_nested": must(types.MapValue(
					obj.AttrTypes["nested"].(types.ObjectType).
						AttrTypes["map_object_nested"].(types.MapType).
						ElemType,
					map[string]attr.Value{
						"key1": must(types.ObjectValue(
							obj.AttrTypes["nested"].(types.ObjectType).
								AttrTypes["map_object_nested"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes,
							map[string]attr.Value{
								"str": types.StringValue("Test1"),
							},
						)),
						"key2": must(types.ObjectValue(
							obj.AttrTypes["nested"].(types.ObjectType).
								AttrTypes["map_object_nested"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes,
							map[string]attr.Value{
								"str": types.StringValue("Test2"),
							},
						)),
					},
				)),
			},
		)),

		"nested_nullable": must(types.ObjectValue(
			obj.AttrTypes["nested_nullable"].(types.ObjectType).AttrTypes,
			map[string]attr.Value{
				"str": types.StringValue("Test"),
				"map": must(types.MapValue(types.StringType, map[string]attr.Value{
					"key1": types.StringValue("Value1"),
					"key2": types.StringValue("Value2"),
				})),
				"nested_list": types.ListNull(
					obj.AttrTypes["nested_nullable"].(types.ObjectType).
						AttrTypes["nested_list"].(types.ListType).
						ElemType,
				),
				"map_object_nested": must(types.MapValue(
					obj.AttrTypes["nested_nullable"].(types.ObjectType).
						AttrTypes["map_object_nested"].(types.MapType).
						ElemType,
					map[string]attr.Value{
						"key1": must(types.ObjectValue(
							obj.AttrTypes["nested_nullable"].(types.ObjectType).
								AttrTypes["map_object_nested"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes,
							map[string]attr.Value{
								"str": types.StringValue("Test1"),
							},
						)),
					},
				)),
			},
		)),

		"nested_list": must(types.ListValue(
			obj.AttrTypes["nested_list"].(types.ListType).ElemType,
			[]attr.Value{
				must(types.ObjectValue(
					obj.AttrTypes["nested_list"].(types.ListType).
						ElemType.(types.ObjectType).
						AttrTypes,
					map[string]attr.Value{
						"str": types.StringValue("Test"),
						"nested_list": must(types.ListValue(
							obj.AttrTypes["nested_list"].(types.ListType).
								ElemType.(types.ObjectType).
								AttrTypes["nested_list"].(types.ListType).
								ElemType,
							[]attr.Value{
								must(types.ObjectValue(
									obj.AttrTypes["nested_list"].(types.ListType).
										ElemType.(types.ObjectType).
										AttrTypes["nested_list"].(types.ListType).
										ElemType.(types.ObjectType).
										AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test1"),
									},
								)),
								must(types.ObjectValue(
									obj.AttrTypes["nested_list"].(types.ListType).
										ElemType.(types.ObjectType).
										AttrTypes["nested_list"].(types.ListType).
										ElemType.(types.ObjectType).
										AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test2"),
									},
								)),
							},
						)),
						"map": must(types.MapValue(types.StringType, map[string]attr.Value{
							"key1": types.StringValue("Value1"),
							"key2": types.StringValue("Value2"),
						})),
						"map_object_nested": must(types.MapValue(
							obj.AttrTypes["nested_list"].(types.ListType).
								ElemType.(types.ObjectType).
								AttrTypes["map_object_nested"].(types.MapType).
								ElemType,
							map[string]attr.Value{
								"key1": must(types.ObjectValue(
									obj.AttrTypes["nested_list"].(types.ListType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test1"),
									},
								)),
								"key2": must(types.ObjectValue(
									obj.AttrTypes["nested_list"].(types.ListType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test2"),
									},
								)),
							},
						)),
					},
				)),
			},
		)),

		"nested_list_nullable": must(types.ListValue(
			obj.AttrTypes["nested_list_nullable"].(types.ListType).ElemType,
			[]attr.Value{
				must(types.ObjectValue(
					obj.AttrTypes["nested_list_nullable"].(types.ListType).
						ElemType.(types.ObjectType).
						AttrTypes,
					map[string]attr.Value{
						"str": types.StringValue("Test"),
						"nested_list": types.ListNull(
							obj.AttrTypes["nested_list_nullable"].(types.ListType).
								ElemType.(types.ObjectType).
								AttrTypes["nested_list"].(types.ListType).
								ElemType,
						),
						"map": types.MapNull(
							obj.AttrTypes["nested_list_nullable"].(types.ListType).
								ElemType.(types.ObjectType).
								AttrTypes["map"].(types.MapType).
								ElemType,
						),
						"map_object_nested": must(types.MapValue(
							obj.AttrTypes["nested_list_nullable"].(types.ListType).
								ElemType.(types.ObjectType).
								AttrTypes["map_object_nested"].(types.MapType).
								ElemType,
							map[string]attr.Value{
								"key1": must(types.ObjectValue(
									obj.AttrTypes["nested_list_nullable"].(types.ListType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test1"),
									},
								)),
								"key2": must(types.ObjectValue(
									obj.AttrTypes["nested_list_nullable"].(types.ListType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test2"),
									},
								)),
							},
						)),
					},
				)),
			},
		)),

		"map_object": must(types.MapValue(
			obj.AttrTypes["map_object"].(types.MapType).ElemType,
			map[string]attr.Value{
				"key1": must(types.ObjectValue(
					obj.AttrTypes["map_object"].(types.MapType).
						ElemType.(types.ObjectType).
						AttrTypes,
					map[string]attr.Value{
						"str": types.StringValue("Test1"),
						"nested_list": types.ListNull(
							obj.AttrTypes["map_object"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["nested_list"].(types.ListType).
								ElemType,
						),
						"map": types.MapNull(
							obj.AttrTypes["map_object"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["map"].(types.MapType).
								ElemType,
						),
						"map_object_nested": must(types.MapValue(
							obj.AttrTypes["map_object"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["map_object_nested"].(types.MapType).
								ElemType,
							map[string]attr.Value{
								"key1": must(types.ObjectValue(
									obj.AttrTypes["map_object"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test1"),
									},
								)),
								"key2": must(types.ObjectValue(
									obj.AttrTypes["map_object"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test2"),
									},
								)),
							},
						)),
					},
				)),
				"key2": must(types.ObjectValue(
					obj.AttrTypes["map_object"].(types.MapType).
						ElemType.(types.ObjectType).
						AttrTypes,
					map[string]attr.Value{
						"str": types.StringValue("Test2"),
						"nested_list": types.ListNull(
							obj.AttrTypes["map_object"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["nested_list"].(types.ListType).
								ElemType,
						),
						"map": types.MapNull(
							obj.AttrTypes["map_object"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["map"].(types.MapType).
								ElemType,
						),
						"map_object_nested": must(types.MapValue(
							obj.AttrTypes["map_object"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["map_object_nested"].(types.MapType).
								ElemType,
							map[string]attr.Value{
								"key1": must(types.ObjectValue(
									obj.AttrTypes["map_object"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test1"),
									},
								)),
								"key2": must(types.ObjectValue(
									obj.AttrTypes["map_object"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test2"),
									},
								)),
							},
						)),
					},
				)),
			},
		)),

		"map_object_nullable": must(types.MapValue(
			obj.AttrTypes["map_object_nullable"].(types.MapType).ElemType,
			map[string]attr.Value{
				"key1": must(types.ObjectValue(
					obj.AttrTypes["map_object_nullable"].(types.MapType).
						ElemType.(types.ObjectType).
						AttrTypes,
					map[string]attr.Value{
						"str": types.StringValue("Test1"),
						"nested_list": types.ListNull(
							obj.AttrTypes["map_object_nullable"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["nested_list"].(types.ListType).
								ElemType,
						),
						"map": types.MapNull(
							obj.AttrTypes["map_object_nullable"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["map"].(types.MapType).
								ElemType,
						),
						"map_object_nested": must(types.MapValue(
							obj.AttrTypes["map_object_nullable"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["map_object_nested"].(types.MapType).
								ElemType,
							map[string]attr.Value{
								"key1": must(types.ObjectValue(
									obj.AttrTypes["map_object_nullable"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test1"),
									},
								)),
								"key2": must(types.ObjectValue(
									obj.AttrTypes["map_object_nullable"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test2"),
									},
								)),
							},
						)),
					},
				)),
				"key2": must(types.ObjectValue(
					obj.AttrTypes["map_object_nullable"].(types.MapType).
						ElemType.(types.ObjectType).
						AttrTypes,
					map[string]attr.Value{
						"str": types.StringValue("Test2"),
						"nested_list": types.ListNull(
							obj.AttrTypes["map_object_nullable"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["nested_list"].(types.ListType).
								ElemType,
						),
						"map": types.MapNull(
							obj.AttrTypes["map_object_nullable"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["map"].(types.MapType).
								ElemType,
						),
						"map_object_nested": must(types.MapValue(
							obj.AttrTypes["map_object_nullable"].(types.MapType).
								ElemType.(types.ObjectType).
								AttrTypes["map_object_nested"].(types.MapType).
								ElemType,
							map[string]attr.Value{
								"key1": must(types.ObjectValue(
									obj.AttrTypes["map_object_nullable"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test1"),
									},
								)),
								"key2": must(types.ObjectValue(
									obj.AttrTypes["map_object_nullable"].(types.MapType).
										ElemType.(types.ObjectType).
										AttrTypes["map_object_nested"].(types.MapType).
										ElemType.(types.ObjectType).AttrTypes,
									map[string]attr.Value{
										"str": types.StringValue("Test2"),
									},
								)),
							},
						)),
					},
				)),
			},
		)),

		"nested_nullable_with_nil_value": types.ObjectNull(
			obj.AttrTypes["nested_nullable_with_nil_value"].(types.ObjectType).AttrTypes,
		),

		"map": must(types.MapValue(types.StringType, map[string]attr.Value{
			"key1": types.StringValue("Value1"),
			"key2": types.StringValue("Value2"),
		})),

		"mode":                 types.Int64Value(1),
		"branch1":              types.ObjectNull(obj.AttrTypes["branch1"].(types.ObjectType).AttrTypes),
		"branch2":              types.ObjectNull(obj.AttrTypes["branch2"].(types.ObjectType).AttrTypes),
		"branch3":              types.StringNull(),
		"empty_message_branch": types.ObjectNull(obj.AttrTypes["empty_message_branch"].(types.ObjectType).AttrTypes),
		"string_branch":        types.StringNull(),
		"embedded_string":      types.StringValue("embdtest1"),
		"embedded_nested_field": must(types.ObjectValue(
			obj.AttrTypes["embedded_nested_field"].(types.ObjectType).AttrTypes,
			map[string]attr.Value{
				"embedded_nested_string": types.StringValue("embdtest2"),
			},
		)),
		"max_age": NewDuration(duration),
		"string_override": must(types.ListValue(types.StringType, []attr.Value{
			types.StringValue("a"),
			types.StringValue("b"),
			types.StringValue("c"),
		})),
		"foo":             types.StringNull(),
		"bar":             types.StringValue("ham"),
		"schema_override": types.StringValue("hello"),
		"required_str":    types.StringValue("Test"),
	}

	result, diags := types.ObjectValue(obj.AttrTypes, attrs)
	requireNoDiagErrors(t, diags)
	return result
}

func copyToTerraformObject(t *testing.T) types.Object {
	s, d := GenSchemaTestResource(context.Background())
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
	case TimeType:
		return NullTime()
	case DurationType:
		return NullDuration()
	default:
		panic(fmt.Sprintf("unsupported fixture attr type %T", attrType))
	}
}
