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

var (
	timestamp = time.Now()
	duration  = 5 * time.Minute
)

func createTestObj() Test {
	return Test{
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
		},
		NestedNullable: &Nested{
			Str: "TestString",
		},
		NestedNullableWithNilValue: nil,

		NestedList: []Nested{
			{
				Str: "Test",
				NestedList: []*OtherNested{
					&OtherNested{
						Str: "Test1",
					},
					&OtherNested{
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
			EmbStr: "embdtest1",
			EmbeddedNestedField: &EmbeddedNestedField{
				EmbNStr: "embdtest2",
			},
		},
	}
}

// copyFromTerraformObject returns a base object used in CopyFrom* tests
func copyFromTerraformObject(t *testing.T) types.Object {
	s, d := GenSchemaTest(context.Background())

	require.False(t, d.HasError())
	typ := s.AttributeType()

	obj, ok := typ.(types.ObjectType)
	require.True(t, ok)

	return types.Object{
		Null:    false,
		Unknown: false,
		Attrs: map[string]attr.Value{
			"str":    types.String{Value: "Test"},
			"int32":  types.Int64{Value: 98},
			"int64":  types.Int64{Value: 99},
			"float":  types.Float64{Value: 0.75},
			"double": types.Float64{Value: 0.76},
			"bool":   types.Bool{Value: true},
			"bytes":  types.String{Value: "Test"},

			"timestamp":                         TimeValue{Value: timestamp},
			"timestamp_missing":                 TimeValue{Unknown: true},
			"timestamp_nullable":                TimeValue{Value: timestamp},
			"timestamp_nullable_with_nil_value": TimeValue{Null: true},
			"duration_standard":                 DurationValue{Value: duration},
			"duration_standard_missing":         DurationValue{Unknown: true},
			"duration_custom":                   DurationValue{Value: duration},
			"duration_custom_missing":           DurationValue{Unknown: true},

			"string_list": types.List{
				Elems: []attr.Value{types.String{Value: "el1"}, types.String{Value: "el2"}},
			},
			"string_list_empty": types.List{Null: false},
			"bytes_list": types.List{
				Elems: []attr.Value{types.String{Value: "bytes1"}, types.String{Value: "bytes2"}},
			},

			"timestamp_list": types.List{
				Elems: []attr.Value{TimeValue{Value: timestamp}, TimeValue{Value: timestamp}},
			},
			"duration_custom_list": types.List{
				Elems: []attr.Value{DurationValue{Value: duration}, DurationValue{Value: duration}},
			},

			"bool_custom_list": types.List{
				Elems: []attr.Value{types.Bool{Value: true}, types.Bool{Value: false}, types.Bool{Value: true}},
			},

			"nested": types.Object{
				Attrs: map[string]attr.Value{
					"str": types.String{Value: "Test"},
					"map": types.Map{
						Elems: map[string]attr.Value{
							"key1": types.String{Value: "Value1"},
							"key2": types.String{Value: "Value2"},
						},
					},
					"nested_list": types.List{Null: true},
					"map_object_nested": types.Map{
						Elems: map[string]attr.Value{
							"key1": types.Object{
								Attrs: map[string]attr.Value{
									"str":         types.String{Value: "Test1"},
									"nested_list": types.List{Null: true},
									"map":         types.Map{Null: true},
								},
							},
							"key2": types.Object{
								Attrs: map[string]attr.Value{
									"str":         types.String{Value: "Test2"},
									"nested_list": types.List{Null: true},
									"map":         types.Map{Null: true},
								},
							},
						},
					},
				},
			},

			"nested_nullable": types.Object{
				Attrs: map[string]attr.Value{
					"str": types.String{Value: "Test"},
					"map": types.Map{
						Elems: map[string]attr.Value{
							"key1": types.String{Value: "Value1"},
							"key2": types.String{Value: "Value2"},
						},
					},
					"nested_list": types.List{Null: true},
					"map_object_nested": types.Map{
						Elems: map[string]attr.Value{
							"key1": types.Object{
								Attrs: map[string]attr.Value{
									"str":         types.String{Value: "Test1"},
									"nested_list": types.List{Null: true},
									"map":         types.Map{Null: true},
								},
							},
							"key2": types.Object{
								Attrs: map[string]attr.Value{
									"str":         types.String{Value: "Test2"},
									"nested_list": types.List{Null: true},
									"map":         types.Map{Null: true},
								},
							},
						},
					},
				},
			},

			"nested_nullable_with_nil_value": types.Object{Null: true},

			"nested_list": types.List{
				Elems: []attr.Value{
					types.Object{
						Attrs: map[string]attr.Value{
							"str": types.String{Value: "Test"},
							"nested_list": types.List{
								Elems: []attr.Value{
									types.Object{
										Attrs: map[string]attr.Value{
											"str": types.String{Value: "Test1"},
										},
									},
									types.Object{
										Attrs: map[string]attr.Value{
											"str": types.String{Value: "Test2"},
										},
									},
								},
							},
							"map": types.Map{
								Elems: map[string]attr.Value{
									"key1": types.String{Value: "Value1"},
									"key2": types.String{Value: "Value2"},
								},
							},
							"map_object_nested": types.Map{
								Elems: map[string]attr.Value{
									"key1": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test1"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
									"key2": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test2"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
								},
							},
						},
					},
				},
			},

			"nested_list_nullable": types.List{
				Elems: []attr.Value{
					types.Object{
						Attrs: map[string]attr.Value{
							"str":         types.String{Value: "Test"},
							"nested_list": types.List{Null: true},
							"map":         types.Map{Null: true},
							"map_object_nested": types.Map{
								Elems: map[string]attr.Value{
									"key1": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test1"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
									"key2": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test2"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
								},
							},
						},
					},
				},
			},

			"map": types.Map{
				Elems: map[string]attr.Value{
					"key1": types.String{Value: "Value1"},
					"key2": types.String{Value: "Value2"},
				},
			},

			"map_object": types.Map{
				Elems: map[string]attr.Value{
					"key1": types.Object{
						Attrs: map[string]attr.Value{
							"str":         types.String{Value: "Test1"},
							"nested_list": types.List{Null: true},
							"map":         types.Map{Null: true},
							"map_object_nested": types.Map{
								Elems: map[string]attr.Value{
									"key1": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test1"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
									"key2": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test2"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
								},
							},
						},
					},
					"key2": types.Object{
						Attrs: map[string]attr.Value{
							"str":         types.String{Value: "Test2"},
							"nested_list": types.List{Null: true},
							"map":         types.Map{Null: true},
							"map_object_nested": types.Map{
								Elems: map[string]attr.Value{
									"key1": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test1"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
									"key2": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test2"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
								},
							},
						},
					},
				},
			},

			"map_object_nullable": types.Map{
				Elems: map[string]attr.Value{
					"key1": types.Object{
						Attrs: map[string]attr.Value{
							"str":         types.String{Value: "Test1"},
							"nested_list": types.List{Null: true},
							"map":         types.Map{Null: true},
							"map_object_nested": types.Map{
								Elems: map[string]attr.Value{
									"key1": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test1"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
									"key2": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test2"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
								},
							},
						},
					},
					"key2": types.Object{
						Attrs: map[string]attr.Value{
							"str":         types.String{Value: "Test2"},
							"nested_list": types.List{Null: true},
							"map":         types.Map{Null: true},
							"map_object_nested": types.Map{
								Elems: map[string]attr.Value{
									"key1": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test1"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
									"key2": types.Object{
										Attrs: map[string]attr.Value{
											"str":         types.String{Value: "Test2"},
											"nested_list": types.List{Null: true},
											"map":         types.Map{Null: true},
										},
									},
								},
							},
						},
					},
				},
			},
			"mode":                 types.Int64{Value: 1},
			"branch1":              types.Object{Null: true},
			"branch2":              types.Object{Null: true},
			"branch3":              types.String{Null: true},
			"empty_message_branch": types.Object{Null: true},
			"string_branch":        types.String{Null: true},
			"embedded_string":      types.String{Value: "embdtest1"},
			"embedded_nested_field": types.Object{
				Attrs: map[string]attr.Value{
					"embedded_nested_string": types.String{Value: "embdtest2"},
				},
			},
		},
		AttrTypes: obj.AttrTypes,
	}
}

// copyToTerraformObject returns the base object used in CopyTo* tests
func copyToTerraformObject(t *testing.T) types.Object {
	s, d := GenSchemaTest(context.Background())

	require.False(t, d.HasError())
	typ := s.AttributeType()

	obj, ok := typ.(types.ObjectType)
	require.True(t, ok)

	return types.Object{
		Unknown:   false,
		Null:      false,
		Attrs:     make(map[string]attr.Value),
		AttrTypes: obj.AttrTypes,
	}
}
