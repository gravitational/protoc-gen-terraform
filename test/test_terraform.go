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

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: test.proto

package test

import (
	fmt "fmt"
	math "math"
	time "time"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	"github.com/gravitational/protoc-gen-terraform/accessors"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

var (
	// SchemaTest is schema for Test message definition.
	SchemaTest = GenSchemaTest()

	// SchemaMetaTest is schema metadata for Test message definition.
	SchemaMetaTest = GenSchemaMetaTest()
)

// SupressDurationChange supresses change for equal durations written differently, ex.: "1h" and "1h0m"
func SupressDurationChange(k string, old string, new string, d *schema.ResourceData) bool {
	o, err := time.ParseDuration(old)
	if err != nil {
		return false
	}

	n, err := time.ParseDuration(new)
	if err != nil {
		return false
	}

	return o == n
}

// SchemaTest returns schema for Test
//
// Test message definition.
func GenSchemaTest() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Str string field
		"str": {
			Type:        schema.TypeString,
			Description: "Str string field",
			Optional:    true,
		},
		// Int32 int32 field
		"int32": {
			Type:        schema.TypeInt,
			Description: "Int32 int32 field",
			Optional:    true,
		},
		// Int64 int64 field
		"int64": {
			Type:        schema.TypeInt,
			Description: "Int64 int64 field",
			Optional:    true,
		},
		// Float float field
		"float": {
			Type:        schema.TypeFloat,
			Description: "Float float field",
			Optional:    true,
		},
		// Double double field
		"double": {
			Type:        schema.TypeFloat,
			Description: "Double double field",
			Optional:    true,
		},
		// Bool bool field
		"bool": {
			Type:        schema.TypeBool,
			Description: "Bool bool field",
			Optional:    true,
		},
		// Bytest byte[] field
		"bytes": {
			Type:        schema.TypeString,
			Description: "Bytest byte[] field",
			Optional:    true,
		},
		// Timestamp time.Time field
		"timestamp": {
			Type:         schema.TypeString,
			Description:  "Timestamp time.Time field",
			ValidateFunc: validation.IsRFC3339Time,
			Optional:     true,
		},
		// Timestamp time.Time field
		"timestamp_missing": {
			Type:         schema.TypeString,
			Description:  "Timestamp time.Time field",
			ValidateFunc: validation.IsRFC3339Time,
			Optional:     true,
		},
		// TimestampNullable *time.Time field
		"timestamp_nullable": {
			Type:         schema.TypeString,
			Description:  "TimestampNullable *time.Time field",
			ValidateFunc: validation.IsRFC3339Time,
			Optional:     true,
		},
		// TimestampNullableWithNilValue *time.Time field
		"timestamp_nullable_with_nil_value": {
			Type:         schema.TypeString,
			Description:  "TimestampNullableWithNilValue *time.Time field",
			ValidateFunc: validation.IsRFC3339Time,
			Optional:     true,
		},
		// DurationStandard time.Duration field (standard)
		"duration_standard": {
			Type:             schema.TypeString,
			Description:      "DurationStandard time.Duration field (standard)",
			DiffSuppressFunc: SupressDurationChange,
			Optional:         true,
		},
		// DurationStandardMissing time.Duration field (standard) missing in input data
		"duration_standard_missing": {
			Type:             schema.TypeString,
			Description:      "DurationStandardMissing time.Duration field (standard) missing in input data",
			DiffSuppressFunc: SupressDurationChange,
			Optional:         true,
		},
		// DurationCustom time.Duration field (with casttype)
		"duration_custom": {
			Type:             schema.TypeString,
			Description:      "DurationCustom time.Duration field (with casttype)",
			DiffSuppressFunc: SupressDurationChange,
			Optional:         true,
		},
		// DurationCustomMissing time.Duration field (with casttype) missing in input data
		"duration_custom_missing": {
			Type:             schema.TypeString,
			Description:      "DurationCustomMissing time.Duration field (with casttype) missing in input data",
			DiffSuppressFunc: SupressDurationChange,
			Optional:         true,
		},
		// StringList []string field
		"string_list": {

			Optional:    true,
			Type:        schema.TypeList,
			Description: "StringList []string field",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// StringListEmpty []string field
		"string_list_empty": {

			Optional:    true,
			Type:        schema.TypeList,
			Description: "StringListEmpty []string field",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// BoolCustomList []bool field
		"bool_custom_list": SchemaBoolCustom(),
		// BytesList [][]byte field
		"bytes_list": {

			Optional:    true,
			Type:        schema.TypeList,
			Description: "BytesList [][]byte field",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// TimestampList []time.Time field
		"timestamp_list": {

			Optional:    true,
			Type:        schema.TypeList,
			Description: "TimestampList []time.Time field",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// DurationCustomList []time.Duration field
		"duration_custom_list": {

			Optional:    true,
			Type:        schema.TypeList,
			Description: "DurationCustomList []time.Duration field",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// Nested nested message field, non-nullable
		"nested": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Description: "Nested message definition",

			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Str string field
					"str": {
						Type:        schema.TypeString,
						Description: "Str string field",
						Optional:    true,
					},
					// Nested repeated nested messages
					"nested_list": {

						Optional:    true,
						Type:        schema.TypeList,
						Description: "Nested repeated nested messages",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str string field
								"str": {
									Type:        schema.TypeString,
									Description: "Str string field",
									Optional:    true,
								},
							},
						},
					},
					// Nested map repeated nested messages
					"map": {

						Optional:    true,
						Type:        schema.TypeMap,
						Description: "Nested map repeated nested messages",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// MapObjectNested nested object map
					"map_object_nested": {

						Optional:    true,
						Type:        schema.TypeSet,
						Description: "MapObjectNested nested object map",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Type:        schema.TypeList,
									MaxItems:    1,
									Description: "OtherNested message nested into nested message",

									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str string field
											"str": {
												Type:        schema.TypeString,
												Description: "Str string field",
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// NestedNullable nested message field, nullabel
		"nested_nullable": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Description: "Nested message definition",

			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Str string field
					"str": {
						Type:        schema.TypeString,
						Description: "Str string field",
						Optional:    true,
					},
					// Nested repeated nested messages
					"nested_list": {

						Optional:    true,
						Type:        schema.TypeList,
						Description: "Nested repeated nested messages",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str string field
								"str": {
									Type:        schema.TypeString,
									Description: "Str string field",
									Optional:    true,
								},
							},
						},
					},
					// Nested map repeated nested messages
					"map": {

						Optional:    true,
						Type:        schema.TypeMap,
						Description: "Nested map repeated nested messages",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// MapObjectNested nested object map
					"map_object_nested": {

						Optional:    true,
						Type:        schema.TypeSet,
						Description: "MapObjectNested nested object map",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Type:        schema.TypeList,
									MaxItems:    1,
									Description: "OtherNested message nested into nested message",

									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str string field
											"str": {
												Type:        schema.TypeString,
												Description: "Str string field",
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// NestedNullableWithNilValue nested message field, with no value set
		"nested_nullable_with_nil_value": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Description: "Nested message definition",

			Optional: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Str string field
					"str": {
						Type:        schema.TypeString,
						Description: "Str string field",
						Optional:    true,
					},
					// Nested repeated nested messages
					"nested_list": {

						Optional:    true,
						Type:        schema.TypeList,
						Description: "Nested repeated nested messages",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str string field
								"str": {
									Type:        schema.TypeString,
									Description: "Str string field",
									Optional:    true,
								},
							},
						},
					},
					// Nested map repeated nested messages
					"map": {

						Optional:    true,
						Type:        schema.TypeMap,
						Description: "Nested map repeated nested messages",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// MapObjectNested nested object map
					"map_object_nested": {

						Optional:    true,
						Type:        schema.TypeSet,
						Description: "MapObjectNested nested object map",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Type:        schema.TypeList,
									MaxItems:    1,
									Description: "OtherNested message nested into nested message",

									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str string field
											"str": {
												Type:        schema.TypeString,
												Description: "Str string field",
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// NestedList nested message array
		"nested_list": {

			Optional:    true,
			Type:        schema.TypeList,
			Description: "NestedList nested message array",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Str string field
					"str": {
						Type:        schema.TypeString,
						Description: "Str string field",
						Optional:    true,
					},
					// Nested repeated nested messages
					"nested_list": {

						Optional:    true,
						Type:        schema.TypeList,
						Description: "Nested repeated nested messages",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str string field
								"str": {
									Type:        schema.TypeString,
									Description: "Str string field",
									Optional:    true,
								},
							},
						},
					},
					// Nested map repeated nested messages
					"map": {

						Optional:    true,
						Type:        schema.TypeMap,
						Description: "Nested map repeated nested messages",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// MapObjectNested nested object map
					"map_object_nested": {

						Optional:    true,
						Type:        schema.TypeSet,
						Description: "MapObjectNested nested object map",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Type:        schema.TypeList,
									MaxItems:    1,
									Description: "OtherNested message nested into nested message",

									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str string field
											"str": {
												Type:        schema.TypeString,
												Description: "Str string field",
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// NestedListNullable nested message array
		"nested_list_nullable": {

			Optional:    true,
			Type:        schema.TypeList,
			Description: "NestedListNullable nested message array",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Str string field
					"str": {
						Type:        schema.TypeString,
						Description: "Str string field",
						Optional:    true,
					},
					// Nested repeated nested messages
					"nested_list": {

						Optional:    true,
						Type:        schema.TypeList,
						Description: "Nested repeated nested messages",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str string field
								"str": {
									Type:        schema.TypeString,
									Description: "Str string field",
									Optional:    true,
								},
							},
						},
					},
					// Nested map repeated nested messages
					"map": {

						Optional:    true,
						Type:        schema.TypeMap,
						Description: "Nested map repeated nested messages",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// MapObjectNested nested object map
					"map_object_nested": {

						Optional:    true,
						Type:        schema.TypeSet,
						Description: "MapObjectNested nested object map",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Type:        schema.TypeList,
									MaxItems:    1,
									Description: "OtherNested message nested into nested message",

									Optional: true,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str string field
											"str": {
												Type:        schema.TypeString,
												Description: "Str string field",
												Optional:    true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// Map normal map
		"map": {

			Optional:    true,
			Type:        schema.TypeMap,
			Description: "Map normal map",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// MapObject is the object map
		"map_object": {

			Optional:    true,
			Type:        schema.TypeSet,
			Description: "MapObject is the object map",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:        schema.TypeList,
						MaxItems:    1,
						Description: "Nested message definition",

						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str string field
								"str": {
									Type:        schema.TypeString,
									Description: "Str string field",
									Optional:    true,
								},
								// Nested repeated nested messages
								"nested_list": {

									Optional:    true,
									Type:        schema.TypeList,
									Description: "Nested repeated nested messages",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str string field
											"str": {
												Type:        schema.TypeString,
												Description: "Str string field",
												Optional:    true,
											},
										},
									},
								},
								// Nested map repeated nested messages
								"map": {

									Optional:    true,
									Type:        schema.TypeMap,
									Description: "Nested map repeated nested messages",
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								// MapObjectNested nested object map
								"map_object_nested": {

									Optional:    true,
									Type:        schema.TypeSet,
									Description: "MapObjectNested nested object map",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"key": {
												Type:     schema.TypeString,
												Required: true,
											},
											"value": {
												Type:        schema.TypeList,
												MaxItems:    1,
												Description: "OtherNested message nested into nested message",

												Optional: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														// Str string field
														"str": {
															Type:        schema.TypeString,
															Description: "Str string field",
															Optional:    true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		// MapObjectNullable is the object map with nullable values
		"map_object_nullable": {

			Optional:    true,
			Type:        schema.TypeSet,
			Description: "MapObjectNullable is the object map with nullable values",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Type:        schema.TypeList,
						MaxItems:    1,
						Description: "Nested message definition",

						Optional: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str string field
								"str": {
									Type:        schema.TypeString,
									Description: "Str string field",
									Optional:    true,
								},
								// Nested repeated nested messages
								"nested_list": {

									Optional:    true,
									Type:        schema.TypeList,
									Description: "Nested repeated nested messages",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str string field
											"str": {
												Type:        schema.TypeString,
												Description: "Str string field",
												Optional:    true,
											},
										},
									},
								},
								// Nested map repeated nested messages
								"map": {

									Optional:    true,
									Type:        schema.TypeMap,
									Description: "Nested map repeated nested messages",
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								// MapObjectNested nested object map
								"map_object_nested": {

									Optional:    true,
									Type:        schema.TypeSet,
									Description: "MapObjectNested nested object map",
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"key": {
												Type:     schema.TypeString,
												Required: true,
											},
											"value": {
												Type:        schema.TypeList,
												MaxItems:    1,
												Description: "OtherNested message nested into nested message",

												Optional: true,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														// Str string field
														"str": {
															Type:        schema.TypeString,
															Description: "Str string field",
															Optional:    true,
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// GenSchemaMetaTest returns schema for Test
//
// Test message definition.
func GenSchemaMetaTest() map[string]*accessors.SchemaMeta {
	return map[string]*accessors.SchemaMeta{
		// Str string field
		"str": {
			Name:       "Str",
			IsTime:     false,
			IsDuration: false,
		},
		// Int32 int32 field
		"int32": {
			Name:       "Int32",
			IsTime:     false,
			IsDuration: false,
		},
		// Int64 int64 field
		"int64": {
			Name:       "Int64",
			IsTime:     false,
			IsDuration: false,
		},
		// Float float field
		"float": {
			Name:       "Float",
			IsTime:     false,
			IsDuration: false,
		},
		// Double double field
		"double": {
			Name:       "Double",
			IsTime:     false,
			IsDuration: false,
		},
		// Bool bool field
		"bool": {
			Name:       "Bool",
			IsTime:     false,
			IsDuration: false,
		},
		// Bytest byte[] field
		"bytes": {
			Name:       "Bytes",
			IsTime:     false,
			IsDuration: false,
		},
		// Timestamp time.Time field
		"timestamp": {
			Name:       "Timestamp",
			IsTime:     true,
			IsDuration: false,
		},
		// Timestamp time.Time field
		"timestamp_missing": {
			Name:       "TimestampMissing",
			IsTime:     true,
			IsDuration: false,
		},
		// TimestampNullable *time.Time field
		"timestamp_nullable": {
			Name:       "TimestampNullable",
			IsTime:     true,
			IsDuration: false,
		},
		// TimestampNullableWithNilValue *time.Time field
		"timestamp_nullable_with_nil_value": {
			Name:       "TimestampNullableWithNilValue",
			IsTime:     true,
			IsDuration: false,
		},
		// DurationStandard time.Duration field (standard)
		"duration_standard": {
			Name:       "DurationStandard",
			IsTime:     false,
			IsDuration: true,
		},
		// DurationStandardMissing time.Duration field (standard) missing in input data
		"duration_standard_missing": {
			Name:       "DurationStandardMissing",
			IsTime:     false,
			IsDuration: true,
		},
		// DurationCustom time.Duration field (with casttype)
		"duration_custom": {
			Name:       "DurationCustom",
			IsTime:     false,
			IsDuration: true,
		},
		// DurationCustomMissing time.Duration field (with casttype) missing in input data
		"duration_custom_missing": {
			Name:       "DurationCustomMissing",
			IsTime:     false,
			IsDuration: true,
		},
		// StringList []string field
		"string_list": {
			Name:       "StringList",
			IsTime:     false,
			IsDuration: false,
		},
		// StringListEmpty []string field
		"string_list_empty": {
			Name:       "StringListEmpty",
			IsTime:     false,
			IsDuration: false,
		},
		// BoolCustomList []bool field
		"bool_custom_list": {
			Name:       "BoolCustomList",
			IsTime:     false,
			IsDuration: false,
			Getter:     GetBoolCustom,
		},
		// BytesList [][]byte field
		"bytes_list": {
			Name:       "BytesList",
			IsTime:     false,
			IsDuration: false,
		},
		// TimestampList []time.Time field
		"timestamp_list": {
			Name:       "TimestampList",
			IsTime:     true,
			IsDuration: false,
		},
		// DurationCustomList []time.Duration field
		"duration_custom_list": {
			Name:       "DurationCustomList",
			IsTime:     false,
			IsDuration: true,
		},
		// Nested nested message field, non-nullable
		"nested": {
			Name:       "Nested",
			IsTime:     false,
			IsDuration: false,
			Nested: map[string]*accessors.SchemaMeta{
				// Str string field
				"str": {
					Name:       "Str",
					IsTime:     false,
					IsDuration: false,
				},
				// Nested repeated nested messages
				"nested_list": {
					Name:       "NestedList",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
				// Nested map repeated nested messages
				"map": {
					Name:       "Map",
					IsTime:     false,
					IsDuration: false,
				},
				// MapObjectNested nested object map
				"map_object_nested": {
					Name:       "MapObjectNested",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
			},
		},
		// NestedNullable nested message field, nullabel
		"nested_nullable": {
			Name:       "NestedNullable",
			IsTime:     false,
			IsDuration: false,
			Nested: map[string]*accessors.SchemaMeta{
				// Str string field
				"str": {
					Name:       "Str",
					IsTime:     false,
					IsDuration: false,
				},
				// Nested repeated nested messages
				"nested_list": {
					Name:       "NestedList",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
				// Nested map repeated nested messages
				"map": {
					Name:       "Map",
					IsTime:     false,
					IsDuration: false,
				},
				// MapObjectNested nested object map
				"map_object_nested": {
					Name:       "MapObjectNested",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
			},
		},
		// NestedNullableWithNilValue nested message field, with no value set
		"nested_nullable_with_nil_value": {
			Name:       "NestedNullableWithNilValue",
			IsTime:     false,
			IsDuration: false,
			Nested: map[string]*accessors.SchemaMeta{
				// Str string field
				"str": {
					Name:       "Str",
					IsTime:     false,
					IsDuration: false,
				},
				// Nested repeated nested messages
				"nested_list": {
					Name:       "NestedList",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
				// Nested map repeated nested messages
				"map": {
					Name:       "Map",
					IsTime:     false,
					IsDuration: false,
				},
				// MapObjectNested nested object map
				"map_object_nested": {
					Name:       "MapObjectNested",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
			},
		},
		// NestedList nested message array
		"nested_list": {
			Name:       "NestedList",
			IsTime:     false,
			IsDuration: false,
			Nested: map[string]*accessors.SchemaMeta{
				// Str string field
				"str": {
					Name:       "Str",
					IsTime:     false,
					IsDuration: false,
				},
				// Nested repeated nested messages
				"nested_list": {
					Name:       "NestedList",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
				// Nested map repeated nested messages
				"map": {
					Name:       "Map",
					IsTime:     false,
					IsDuration: false,
				},
				// MapObjectNested nested object map
				"map_object_nested": {
					Name:       "MapObjectNested",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
			},
		},
		// NestedListNullable nested message array
		"nested_list_nullable": {
			Name:       "NestedListNullable",
			IsTime:     false,
			IsDuration: false,
			Nested: map[string]*accessors.SchemaMeta{
				// Str string field
				"str": {
					Name:       "Str",
					IsTime:     false,
					IsDuration: false,
				},
				// Nested repeated nested messages
				"nested_list": {
					Name:       "NestedList",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
				// Nested map repeated nested messages
				"map": {
					Name:       "Map",
					IsTime:     false,
					IsDuration: false,
				},
				// MapObjectNested nested object map
				"map_object_nested": {
					Name:       "MapObjectNested",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
			},
		},
		// Map normal map
		"map": {
			Name:       "Map",
			IsTime:     false,
			IsDuration: false,
		},
		// MapObject is the object map
		"map_object": {
			Name:       "MapObject",
			IsTime:     false,
			IsDuration: false,
			Nested: map[string]*accessors.SchemaMeta{
				// Str string field
				"str": {
					Name:       "Str",
					IsTime:     false,
					IsDuration: false,
				},
				// Nested repeated nested messages
				"nested_list": {
					Name:       "NestedList",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
				// Nested map repeated nested messages
				"map": {
					Name:       "Map",
					IsTime:     false,
					IsDuration: false,
				},
				// MapObjectNested nested object map
				"map_object_nested": {
					Name:       "MapObjectNested",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
			},
		},
		// MapObjectNullable is the object map with nullable values
		"map_object_nullable": {
			Name:       "MapObjectNullable",
			IsTime:     false,
			IsDuration: false,
			Nested: map[string]*accessors.SchemaMeta{
				// Str string field
				"str": {
					Name:       "Str",
					IsTime:     false,
					IsDuration: false,
				},
				// Nested repeated nested messages
				"nested_list": {
					Name:       "NestedList",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
				// Nested map repeated nested messages
				"map": {
					Name:       "Map",
					IsTime:     false,
					IsDuration: false,
				},
				// MapObjectNested nested object map
				"map_object_nested": {
					Name:       "MapObjectNested",
					IsTime:     false,
					IsDuration: false,
					Nested: map[string]*accessors.SchemaMeta{
						// Str string field
						"str": {
							Name:       "Str",
							IsTime:     false,
							IsDuration: false,
						},
					},
				},
			},
		},
	}
}
