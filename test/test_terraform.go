// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: test.proto

package test

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
	math "math"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// Type full name: Test
func SchemaTest() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Str SINGULAR_ELEMENTARY
		"str": {
			Type:     schema.TypeString,
			Optional: true,
		},
		// Int32 SINGULAR_ELEMENTARY
		"int32": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		// Int64 SINGULAR_ELEMENTARY
		"int64": {
			Type:     schema.TypeInt,
			Optional: true,
		},
		// Float SINGULAR_ELEMENTARY
		"float": {
			Type:     schema.TypeFloat,
			Optional: true,
		},
		// Double SINGULAR_ELEMENTARY
		"double": {
			Type:     schema.TypeFloat,
			Optional: true,
		},
		// Bool SINGULAR_ELEMENTARY
		"bool": {
			Type:     schema.TypeBool,
			Optional: true,
		},
		// Bytes SINGULAR_ELEMENTARY
		"bytes": {
			Type:     schema.TypeString,
			Optional: true,
		},
		// Timestamp SINGULAR_ELEMENTARY
		"timestamp": {
			Type:         schema.TypeString,
			ValidateFunc: validation.IsRFC3339Time,
			Optional:     true,
		},
		// DurationStd SINGULAR_ELEMENTARY
		"duration_std": {
			Type:     schema.TypeString,
			Optional: true,
		},
		// DurationCustom SINGULAR_ELEMENTARY
		"duration_custom": {
			Type:     schema.TypeString,
			Optional: true,
		},
		// TimestampN SINGULAR_ELEMENTARY
		"timestamp_n": {
			Type:         schema.TypeString,
			ValidateFunc: validation.IsRFC3339Time,
			Optional:     true,
		},
		// StringA REPEATED_ELEMENTARY
		"string_a": {
			Optional: true,
			Type:     schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// BoolA CUSTOM_TYPE
		"bool_a": SchemaBoolCustom(),
		// BytesA REPEATED_ELEMENTARY
		"bytes_a": {
			Optional: true,
			Type:     schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// TimestampA REPEATED_ELEMENTARY
		"timestamp_a": {
			Optional: true,
			Type:     schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// GracePeriodA REPEATED_ELEMENTARY
		"grace_period_a": {
			Optional: true,
			Type:     schema.TypeList,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// Nested SINGULAR_MESSAGE
		"nested": {
			Optional: true,
			Type:     schema.TypeList,
			MaxItems: 1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Str SINGULAR_ELEMENTARY
					"str": {
						Type:     schema.TypeString,
						Optional: true,
					},
					// Nested REPEATED_MESSAGE
					"nested": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str SINGULAR_ELEMENTARY
								"str": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					// NestedS MAP
					"nested_s": {
						Optional: true,
						Type:     schema.TypeMap,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// NestedM OBJECT_MAP
					"nested_m": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Optional: true,
									Type:     schema.TypeList,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str SINGULAR_ELEMENTARY
											"str": {
												Type:     schema.TypeString,
												Optional: true,
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
		// NestedA REPEATED_MESSAGE
		"nested_a": {
			Optional: true,
			Type:     schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Str SINGULAR_ELEMENTARY
					"str": {
						Type:     schema.TypeString,
						Optional: true,
					},
					// Nested REPEATED_MESSAGE
					"nested": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str SINGULAR_ELEMENTARY
								"str": {
									Type:     schema.TypeString,
									Optional: true,
								},
							},
						},
					},
					// NestedS MAP
					"nested_s": {
						Optional: true,
						Type:     schema.TypeMap,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// NestedM OBJECT_MAP
					"nested_m": {
						Optional: true,
						Type:     schema.TypeList,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Optional: true,
									Type:     schema.TypeList,
									MaxItems: 1,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str SINGULAR_ELEMENTARY
											"str": {
												Type:     schema.TypeString,
												Optional: true,
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
		// NestedS MAP
		"nested_s": {
			Optional: true,
			Type:     schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// NestedM OBJECT_MAP
		"nested_m": {
			Optional: true,
			Type:     schema.TypeList,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Optional: true,
						Type:     schema.TypeList,
						MaxItems: 1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str SINGULAR_ELEMENTARY
								"str": {
									Type:     schema.TypeString,
									Optional: true,
								},
								// Nested REPEATED_MESSAGE
								"nested": {
									Optional: true,
									Type:     schema.TypeList,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											// Str SINGULAR_ELEMENTARY
											"str": {
												Type:     schema.TypeString,
												Optional: true,
											},
										},
									},
								},
								// NestedS MAP
								"nested_s": {
									Optional: true,
									Type:     schema.TypeMap,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								// NestedM OBJECT_MAP
								"nested_m": {
									Optional: true,
									Type:     schema.TypeList,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"key": {
												Type:     schema.TypeString,
												Required: true,
											},
											"value": {
												Optional: true,
												Type:     schema.TypeList,
												MaxItems: 1,
												Elem: &schema.Resource{
													Schema: map[string]*schema.Schema{
														// Str SINGULAR_ELEMENTARY
														"str": {
															Type:     schema.TypeString,
															Optional: true,
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
func UnmarshalTest(d *schema.ResourceData, t *Test) error {
	p := ""

	{

		_raw, ok := d.GetOk(p + "str")
		if ok {
			_value := _raw.(string)
			t.Str = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "int32")
		if ok {
			_value := int32(int32(_raw.(int)))
			t.Int32 = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "int64")
		if ok {
			_value := int64(int64(_raw.(int)))
			t.Int64 = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "float")
		if ok {
			_value := float32(float32(_raw.(float64)))
			t.Float = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "double")
		if ok {
			_value := _raw.(float64)
			t.Double = _value
		}
	}
	{

		_raw, ok := d.GetOkExists(p + "bool")
		if ok {
			_value := _raw.(bool)
			t.Bool = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "bytes")
		if ok {
			_value := []byte([]byte(_raw.(string)))
			t.Bytes = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "timestamp")
		if ok {
			_value, err := time.Parse(time.RFC3339, _raw.(string))
			if err != nil {
				return fmt.Errorf("Malformed time value for field Timestamp : %w", err)
			}
			t.Timestamp = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "duration_std")
		if ok {
			_valued, err := time.ParseDuration(_raw.(string))
			if err != nil {
				return fmt.Errorf("Malformed duration value for field DurationStd : %w", err)
			}
			_value := time.Duration(_valued)
			t.DurationStd = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "duration_custom")
		if ok {
			_valued, err := time.ParseDuration(_raw.(string))
			if err != nil {
				return fmt.Errorf("Malformed duration value for field DurationCustom : %w", err)
			}
			_value := Duration(_valued)
			t.DurationCustom = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "timestamp_n")
		if ok {
			_value, err := time.Parse(time.RFC3339, _raw.(string))
			if err != nil {
				return fmt.Errorf("Malformed time value for field TimestampN : %w", err)
			}
			t.TimestampN = &_value
		}
	}
	{
		_rawi, ok := d.GetOk(p + "string_a")
		if ok {
			_rawi := _rawi.([]interface{})
			t.StringA = make([]string, len(_rawi))
			for i := 0; i < len(_rawi); i++ {
				_raw := _rawi[i]
				_value := _raw.(string)
				t.StringA[i] = _value
			}
		}
	}
	{
		err := UnmarshalBoolCustom(p+"bool_a", d, &t.BoolA)
		if err != nil {
			return err
		}
	}
	{
		_rawi, ok := d.GetOk(p + "bytes_a")
		if ok {
			_rawi := _rawi.([]interface{})
			t.BytesA = make([][]byte, len(_rawi))
			for i := 0; i < len(_rawi); i++ {
				_raw := _rawi[i]
				_value := []byte([]byte(_raw.(string)))
				t.BytesA[i] = _value
			}
		}
	}
	{
		_rawi, ok := d.GetOk(p + "timestamp_a")
		if ok {
			_rawi := _rawi.([]interface{})
			t.TimestampA = make([]*time.Time, len(_rawi))
			for i := 0; i < len(_rawi); i++ {
				_raw := _rawi[i]
				_value, err := time.Parse(time.RFC3339, _raw.(string))
				if err != nil {
					return fmt.Errorf("Malformed time value for field TimestampA : %w", err)
				}
				t.TimestampA[i] = &_value
			}
		}
	}
	{
		_rawi, ok := d.GetOk(p + "grace_period_a")
		if ok {
			_rawi := _rawi.([]interface{})
			t.GracePeriodA = make([]Duration, len(_rawi))
			for i := 0; i < len(_rawi); i++ {
				_raw := _rawi[i]
				_valued, err := time.ParseDuration(_raw.(string))
				if err != nil {
					return fmt.Errorf("Malformed duration value for field GracePeriodA : %w", err)
				}
				_value := Duration(_valued)
				t.GracePeriodA[i] = _value
			}
		}
	}

	return nil
}
