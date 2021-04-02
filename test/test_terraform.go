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

// Test (Test message definition.)
// └── str:string (Str string field)
// └── int32:int (Int32 int32 field)
// └── int64:int (Int64 int64 field)
// └── float:float64 (Float float field)
// └── double:float64 (Double double field)
// └── bool:bool (Bool bool field)
// └── bytes:string (Bytest byte[] field)
// └── timestamp:string (Timestamp time.Time field)
// └── duration_std:string (DurationStd time.Duration field (standard))
// └── duration_custom:string (DurationCustom time.Duration field (custom))
// └── timestamp_n:string (TimestampN *time.Time field)
// └── [string_a:string] (StringA []string field)
// └── bool_a !custom schema, see target code! (BoolA []bool field)
// └── [bytes_a:string] (BytesA [][]byte field)
// └── [timestamp_a:string] (TimestampA []time.Time field)
// └── [duration_custom_a:string] (DurationCustomA []time.Duration field)
// └── nested (Nested nested message field)
// │   ├── str:string (Str string field)
// │   ├── [nested] (Nested repeated nested messages)
// │   │   ├── str:string (Str string field)
// │   ├── nested_m:map (Nested map repeated nested messages)
// │   ├── [nested_m_obj] (NestedMObj nested object map)
// │       └── key:string
// │       └── value
// │           └── str:string (Str string field)
// └── [nested_a] (NestedA nested message array)
// │   ├── str:string (Str string field)
// │   ├── [nested] (Nested repeated nested messages)
// │   │   ├── str:string (Str string field)
// │   ├── nested_m:map (Nested map repeated nested messages)
// │   ├── [nested_m_obj] (NestedMObj nested object map)
// │       └── key:string
// │       └── value
// │           └── str:string (Str string field)
// └── nested_m:map (NestedM normal map)
// └── [nested_m_obj] (NestedMObj object map)
//     └── key:string
//     └── value
//         └── str:string (Str string field)
//         └── [nested] (Nested repeated nested messages)
//         │   ├── str:string (Str string field)
//         └── nested_m:map (Nested map repeated nested messages)
//         └── [nested_m_obj] (NestedMObj nested object map)
//             └── key:string
//             └── value
//                 └── str:string (Str string field)

// SchemaTest returns schema for Test
//
// Test message definition.
func SchemaTest() map[string]*schema.Schema {
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
		// DurationStd time.Duration field (standard)
		"duration_std": {
			Type:        schema.TypeString,
			Description: "DurationStd time.Duration field (standard)",
			DiffSuppressFunc: func(k string, old string, new string, d *schema.ResourceData) bool {
				o, err := time.ParseDuration(old)
				if err != nil {
					return false
				}

				n, err := time.ParseDuration(new)
				if err != nil {
					return false
				}

				return o == n
			},
			Optional: true,
		},
		// DurationCustom time.Duration field (custom)
		"duration_custom": {
			Type:        schema.TypeString,
			Description: "DurationCustom time.Duration field (custom)",
			DiffSuppressFunc: func(k string, old string, new string, d *schema.ResourceData) bool {
				o, err := time.ParseDuration(old)
				if err != nil {
					return false
				}

				n, err := time.ParseDuration(new)
				if err != nil {
					return false
				}

				return o == n
			},
			Optional: true,
		},
		// TimestampN *time.Time field
		"timestamp_n": {
			Type:         schema.TypeString,
			Description:  "TimestampN *time.Time field",
			ValidateFunc: validation.IsRFC3339Time,
			Optional:     true,
		},
		// StringA []string field
		"string_a": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "StringA []string field",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// BoolA []bool field
		"bool_a": SchemaBoolCustom(),
		// BytesA [][]byte field
		"bytes_a": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "BytesA [][]byte field",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// TimestampA []time.Time field
		"timestamp_a": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "TimestampA []time.Time field",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// DurationCustomA []time.Duration field
		"duration_custom_a": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "DurationCustomA []time.Duration field",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// Nested nested message field
		"nested": {
			Optional:    true,
			Type:        schema.TypeList,
			MaxItems:    1,
			ConfigMode:  schema.SchemaConfigModeAttr,
			Description: "Nested message definition",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Str string field
					"str": {
						Type:        schema.TypeString,
						Description: "Str string field",
						Optional:    true,
					},
					// Nested repeated nested messages
					"nested": {
						Optional:    true,
						Type:        schema.TypeList,
						Description: "Message nested into nested message",
						ConfigMode:  schema.SchemaConfigModeAttr,
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
					"nested_m": {
						Optional:    true,
						Type:        schema.TypeMap,
						Description: "Nested map repeated nested messages",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// NestedMObj nested object map
					"nested_m_obj": {
						Optional:    true,
						Type:        schema.TypeList,
						Description: "NestedMObj nested object map",
						ConfigMode:  schema.SchemaConfigModeAttr,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Optional:    true,
									Type:        schema.TypeList,
									MaxItems:    1,
									ConfigMode:  schema.SchemaConfigModeAttr,
									Description: "Message nested into nested message",
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
		// NestedA nested message array
		"nested_a": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "Nested message definition",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					// Str string field
					"str": {
						Type:        schema.TypeString,
						Description: "Str string field",
						Optional:    true,
					},
					// Nested repeated nested messages
					"nested": {
						Optional:    true,
						Type:        schema.TypeList,
						Description: "Message nested into nested message",
						ConfigMode:  schema.SchemaConfigModeAttr,
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
					"nested_m": {
						Optional:    true,
						Type:        schema.TypeMap,
						Description: "Nested map repeated nested messages",
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// NestedMObj nested object map
					"nested_m_obj": {
						Optional:    true,
						Type:        schema.TypeList,
						Description: "NestedMObj nested object map",
						ConfigMode:  schema.SchemaConfigModeAttr,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"key": {
									Type:     schema.TypeString,
									Required: true,
								},
								"value": {
									Optional:    true,
									Type:        schema.TypeList,
									MaxItems:    1,
									ConfigMode:  schema.SchemaConfigModeAttr,
									Description: "Message nested into nested message",
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
		// NestedM normal map
		"nested_m": {
			Optional:    true,
			Type:        schema.TypeMap,
			Description: "NestedM normal map",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// NestedMObj object map
		"nested_m_obj": {
			Optional:    true,
			Type:        schema.TypeList,
			Description: "NestedMObj object map",
			ConfigMode:  schema.SchemaConfigModeAttr,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"key": {
						Type:     schema.TypeString,
						Required: true,
					},
					"value": {
						Optional:    true,
						Type:        schema.TypeList,
						MaxItems:    1,
						ConfigMode:  schema.SchemaConfigModeAttr,
						Description: "Nested message definition",
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								// Str string field
								"str": {
									Type:        schema.TypeString,
									Description: "Str string field",
									Optional:    true,
								},
								// Nested repeated nested messages
								"nested": {
									Optional:    true,
									Type:        schema.TypeList,
									Description: "Message nested into nested message",
									ConfigMode:  schema.SchemaConfigModeAttr,
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
								"nested_m": {
									Optional:    true,
									Type:        schema.TypeMap,
									Description: "Nested map repeated nested messages",
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								// NestedMObj nested object map
								"nested_m_obj": {
									Optional:    true,
									Type:        schema.TypeList,
									Description: "NestedMObj nested object map",
									ConfigMode:  schema.SchemaConfigModeAttr,
									Elem: &schema.Resource{
										Schema: map[string]*schema.Schema{
											"key": {
												Type:     schema.TypeString,
												Required: true,
											},
											"value": {
												Optional:    true,
												Type:        schema.TypeList,
												MaxItems:    1,
												ConfigMode:  schema.SchemaConfigModeAttr,
												Description: "Message nested into nested message",
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
func GetTestFromResourceData(d *schema.ResourceData, t *Test) error {
	p := ""

	{

		_raw, ok := d.GetOk(p + "str")

		if ok {
			_raws, ok := _raw.(string)
			if !ok {
				return fmt.Errorf("can not convert %T to string", _raws)
			}
			_value := string(string(_raws))
			t.Str = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "int32")

		if ok {
			_raws, ok := _raw.(int)
			if !ok {
				return fmt.Errorf("can not convert %T to int", _raws)
			}
			_value := int32(int32(_raws))
			t.Int32 = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "int64")

		if ok {
			_raws, ok := _raw.(int)
			if !ok {
				return fmt.Errorf("can not convert %T to int", _raws)
			}
			_value := int64(int64(_raws))
			t.Int64 = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "float")

		if ok {
			_raws, ok := _raw.(float64)
			if !ok {
				return fmt.Errorf("can not convert %T to float64", _raws)
			}
			_value := float32(float32(_raws))
			t.Float = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "double")

		if ok {
			_raws, ok := _raw.(float64)
			if !ok {
				return fmt.Errorf("can not convert %T to float64", _raws)
			}
			_value := float64(float64(_raws))
			t.Double = _value
		}
	}
	{

		_raw, ok := d.GetOkExists(p + "bool")

		if ok {
			_raws, ok := _raw.(bool)
			if !ok {
				return fmt.Errorf("can not convert %T to bool", _raws)
			}
			_value := bool(bool(_raws))
			t.Bool = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "bytes")

		if ok {
			_raws, ok := _raw.(string)
			if !ok {
				return fmt.Errorf("can not convert %T to string", _raws)
			}
			_value := []byte([]byte(_raws))
			t.Bytes = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "timestamp")

		if ok {
			_raws, ok := _raw.(string)
			if !ok {
				return fmt.Errorf("can not convert %T to string", _raws)
			}
			_value, err := time.Parse(time.RFC3339Nano, _raws)
			if err != nil {
				return fmt.Errorf("malformed time value for field Timestamp : %w", err)
			}
			t.Timestamp = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "duration_std")

		if ok {
			_raws, ok := _raw.(string)
			if !ok {
				return fmt.Errorf("can not convert %T to string", _raws)
			}
			_valued, err := time.ParseDuration(_raws)
			if err != nil {
				return fmt.Errorf("malformed duration value for field DurationStd : %w", err)
			}
			_value := time.Duration(_valued)
			t.DurationStd = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "duration_custom")

		if ok {
			_raws, ok := _raw.(string)
			if !ok {
				return fmt.Errorf("can not convert %T to string", _raws)
			}
			_valued, err := time.ParseDuration(_raws)
			if err != nil {
				return fmt.Errorf("malformed duration value for field DurationCustom : %w", err)
			}
			_value := Duration(_valued)
			t.DurationCustom = _value
		}
	}
	{

		_raw, ok := d.GetOk(p + "timestamp_n")

		if ok {
			_raws, ok := _raw.(string)
			if !ok {
				return fmt.Errorf("can not convert %T to string", _raws)
			}
			_value, err := time.Parse(time.RFC3339Nano, _raws)
			if err != nil {
				return fmt.Errorf("malformed time value for field TimestampN : %w", err)
			}
			t.TimestampN = &_value
		}
	}
	{
		_rawi, ok := d.GetOk(p + "string_a")
		if ok {
			_rawi, ok := _rawi.([]interface{})
			if !ok {
				return fmt.Errorf("count not convert %T to []interface{}", _rawi)
			}
			t.StringA = make([]string, len(_rawi))
			for i := 0; i < len(_rawi); i++ {
				_raw := _rawi[i]
				_raws, ok := _raw.(string)
				if !ok {
					return fmt.Errorf("can not convert %T to string", _raws)
				}
				_value := string(string(_raws))
				t.StringA[i] = _value
			}
		}
	}
	{
		err := GetBoolCustomFromResourceData(p+"bool_a", d, &t.BoolA)
		if err != nil {
			return err
		}
	}
	{
		_rawi, ok := d.GetOk(p + "bytes_a")
		if ok {
			_rawi, ok := _rawi.([]interface{})
			if !ok {
				return fmt.Errorf("count not convert %T to []interface{}", _rawi)
			}
			t.BytesA = make([][]byte, len(_rawi))
			for i := 0; i < len(_rawi); i++ {
				_raw := _rawi[i]
				_raws, ok := _raw.(string)
				if !ok {
					return fmt.Errorf("can not convert %T to string", _raws)
				}
				_value := []byte([]byte(_raws))
				t.BytesA[i] = _value
			}
		}
	}
	{
		_rawi, ok := d.GetOk(p + "timestamp_a")
		if ok {
			_rawi, ok := _rawi.([]interface{})
			if !ok {
				return fmt.Errorf("count not convert %T to []interface{}", _rawi)
			}
			t.TimestampA = make([]*time.Time, len(_rawi))
			for i := 0; i < len(_rawi); i++ {
				_raw := _rawi[i]
				_raws, ok := _raw.(string)
				if !ok {
					return fmt.Errorf("can not convert %T to string", _raws)
				}
				_value, err := time.Parse(time.RFC3339Nano, _raws)
				if err != nil {
					return fmt.Errorf("malformed time value for field TimestampA : %w", err)
				}
				t.TimestampA[i] = &_value
			}
		}
	}
	{
		_rawi, ok := d.GetOk(p + "duration_custom_a")
		if ok {
			_rawi, ok := _rawi.([]interface{})
			if !ok {
				return fmt.Errorf("count not convert %T to []interface{}", _rawi)
			}
			t.DurationCustomA = make([]Duration, len(_rawi))
			for i := 0; i < len(_rawi); i++ {
				_raw := _rawi[i]
				_raws, ok := _raw.(string)
				if !ok {
					return fmt.Errorf("can not convert %T to string", _raws)
				}
				_valued, err := time.ParseDuration(_raws)
				if err != nil {
					return fmt.Errorf("malformed duration value for field DurationCustomA : %w", err)
				}
				_value := Duration(_valued)
				t.DurationCustomA[i] = _value
			}
		}
	}
	{
		p := p + "nested" + ".0."

		_obj := Nested{}
		t.Nested = &_obj
		t := &_obj

		{

			_raw, ok := d.GetOk(p + "str")

			if ok {
				_raws, ok := _raw.(string)
				if !ok {
					return fmt.Errorf("can not convert %T to string", _raws)
				}
				_value := string(string(_raws))
				t.Str = _value
			}
		}
		{
			p := p + "nested"

			_raw, ok := d.GetOk(p)
			if ok {
				_rawi, ok := _raw.([]interface{})
				if !ok {
					return fmt.Errorf("can not convert %T to []interface{}", _raw)
				}

				t.Nested = make([]*NestedLevel2, len(_rawi))
				for i := 0; i < len(_rawi); i++ {

					_obj := NestedLevel2{}
					t.Nested[i] = &_obj

					{
						t := t.Nested[i]
						p := p + fmt.Sprintf(".%v.", i)
						{

							_raw, ok := d.GetOk(p + "str")

							if ok {
								_raws, ok := _raw.(string)
								if !ok {
									return fmt.Errorf("can not convert %T to string", _raws)
								}
								_value := string(string(_raws))
								t.Str = _value
							}
						}

					}
				}
			}
		}
		{

			p := p + "nested_m"
			_rawm, ok := d.GetOk(p)
			if ok {
				_rawmi, ok := _rawm.(map[string]interface{})
				if !ok {
					return fmt.Errorf("can not convert %T to map[string]interface{}", _rawm)
				}
				t.NestedM = make(map[string]string, len(_rawmi))
				for _k, _v := range _rawmi {
					_raw := _v
					_raws, ok := _raw.(string)
					if !ok {
						return fmt.Errorf("can not convert %T to string", _raws)
					}
					_value := string(string(_raws))
					t.NestedM[_k] = _value
				}
			}
		}
		{
			p := p + "nested_m_obj"

			_raw, ok := d.GetOk(p)
			if ok {
				_rawi, ok := _raw.([]interface{})
				if !ok {
					return fmt.Errorf("can not convert %T to []interface{}", _raw)
				}

				_value := make(map[string]*NestedLevel2)

				for i := range _rawi {
					_rawkey := d.Get(fmt.Sprintf("%v.%v.", p, i) + "key")
					_key, ok := _rawkey.(string)
					if !ok {
						return fmt.Errorf("can not convert %T to string", _rawkey)
					}
					if _key == "" {
						return fmt.Errorf("missing key field in object map NestedMObj")
					}

					_obj := NestedLevel2{}
					_value[_key] = &_obj
					t := &_obj

					{
						p := fmt.Sprintf("%v.%v.value.0.", p, i)
						{

							_raw, ok := d.GetOk(p + "str")

							if ok {
								_raws, ok := _raw.(string)
								if !ok {
									return fmt.Errorf("can not convert %T to string", _raws)
								}
								_value := string(string(_raws))
								t.Str = _value
							}
						}

					}
				}

				t.NestedMObj = _value
			}
		}

	}
	{
		p := p + "nested_a"

		_raw, ok := d.GetOk(p)
		if ok {
			_rawi, ok := _raw.([]interface{})
			if !ok {
				return fmt.Errorf("can not convert %T to []interface{}", _raw)
			}

			t.NestedA = make([]*Nested, len(_rawi))
			for i := 0; i < len(_rawi); i++ {

				_obj := Nested{}
				t.NestedA[i] = &_obj

				{
					t := t.NestedA[i]
					p := p + fmt.Sprintf(".%v.", i)
					{

						_raw, ok := d.GetOk(p + "str")

						if ok {
							_raws, ok := _raw.(string)
							if !ok {
								return fmt.Errorf("can not convert %T to string", _raws)
							}
							_value := string(string(_raws))
							t.Str = _value
						}
					}
					{
						p := p + "nested"

						_raw, ok := d.GetOk(p)
						if ok {
							_rawi, ok := _raw.([]interface{})
							if !ok {
								return fmt.Errorf("can not convert %T to []interface{}", _raw)
							}

							t.Nested = make([]*NestedLevel2, len(_rawi))
							for i := 0; i < len(_rawi); i++ {

								_obj := NestedLevel2{}
								t.Nested[i] = &_obj

								{
									t := t.Nested[i]
									p := p + fmt.Sprintf(".%v.", i)
									{

										_raw, ok := d.GetOk(p + "str")

										if ok {
											_raws, ok := _raw.(string)
											if !ok {
												return fmt.Errorf("can not convert %T to string", _raws)
											}
											_value := string(string(_raws))
											t.Str = _value
										}
									}

								}
							}
						}
					}
					{

						p := p + "nested_m"
						_rawm, ok := d.GetOk(p)
						if ok {
							_rawmi, ok := _rawm.(map[string]interface{})
							if !ok {
								return fmt.Errorf("can not convert %T to map[string]interface{}", _rawm)
							}
							t.NestedM = make(map[string]string, len(_rawmi))
							for _k, _v := range _rawmi {
								_raw := _v
								_raws, ok := _raw.(string)
								if !ok {
									return fmt.Errorf("can not convert %T to string", _raws)
								}
								_value := string(string(_raws))
								t.NestedM[_k] = _value
							}
						}
					}
					{
						p := p + "nested_m_obj"

						_raw, ok := d.GetOk(p)
						if ok {
							_rawi, ok := _raw.([]interface{})
							if !ok {
								return fmt.Errorf("can not convert %T to []interface{}", _raw)
							}

							_value := make(map[string]*NestedLevel2)

							for i := range _rawi {
								_rawkey := d.Get(fmt.Sprintf("%v.%v.", p, i) + "key")
								_key, ok := _rawkey.(string)
								if !ok {
									return fmt.Errorf("can not convert %T to string", _rawkey)
								}
								if _key == "" {
									return fmt.Errorf("missing key field in object map NestedMObj")
								}

								_obj := NestedLevel2{}
								_value[_key] = &_obj
								t := &_obj

								{
									p := fmt.Sprintf("%v.%v.value.0.", p, i)
									{

										_raw, ok := d.GetOk(p + "str")

										if ok {
											_raws, ok := _raw.(string)
											if !ok {
												return fmt.Errorf("can not convert %T to string", _raws)
											}
											_value := string(string(_raws))
											t.Str = _value
										}
									}

								}
							}

							t.NestedMObj = _value
						}
					}

				}
			}
		}
	}
	{

		p := p + "nested_m"
		_rawm, ok := d.GetOk(p)
		if ok {
			_rawmi, ok := _rawm.(map[string]interface{})
			if !ok {
				return fmt.Errorf("can not convert %T to map[string]interface{}", _rawm)
			}
			t.NestedM = make(map[string]string, len(_rawmi))
			for _k, _v := range _rawmi {
				_raw := _v
				_raws, ok := _raw.(string)
				if !ok {
					return fmt.Errorf("can not convert %T to string", _raws)
				}
				_value := string(string(_raws))
				t.NestedM[_k] = _value
			}
		}
	}
	{
		p := p + "nested_m_obj"

		_raw, ok := d.GetOk(p)
		if ok {
			_rawi, ok := _raw.([]interface{})
			if !ok {
				return fmt.Errorf("can not convert %T to []interface{}", _raw)
			}

			_value := make(map[string]*Nested)

			for i := range _rawi {
				_rawkey := d.Get(fmt.Sprintf("%v.%v.", p, i) + "key")
				_key, ok := _rawkey.(string)
				if !ok {
					return fmt.Errorf("can not convert %T to string", _rawkey)
				}
				if _key == "" {
					return fmt.Errorf("missing key field in object map NestedMObj")
				}

				_obj := Nested{}
				_value[_key] = &_obj
				t := &_obj

				{
					p := fmt.Sprintf("%v.%v.value.0.", p, i)
					{

						_raw, ok := d.GetOk(p + "str")

						if ok {
							_raws, ok := _raw.(string)
							if !ok {
								return fmt.Errorf("can not convert %T to string", _raws)
							}
							_value := string(string(_raws))
							t.Str = _value
						}
					}
					{
						p := p + "nested"

						_raw, ok := d.GetOk(p)
						if ok {
							_rawi, ok := _raw.([]interface{})
							if !ok {
								return fmt.Errorf("can not convert %T to []interface{}", _raw)
							}

							t.Nested = make([]*NestedLevel2, len(_rawi))
							for i := 0; i < len(_rawi); i++ {

								_obj := NestedLevel2{}
								t.Nested[i] = &_obj

								{
									t := t.Nested[i]
									p := p + fmt.Sprintf(".%v.", i)
									{

										_raw, ok := d.GetOk(p + "str")

										if ok {
											_raws, ok := _raw.(string)
											if !ok {
												return fmt.Errorf("can not convert %T to string", _raws)
											}
											_value := string(string(_raws))
											t.Str = _value
										}
									}

								}
							}
						}
					}
					{

						p := p + "nested_m"
						_rawm, ok := d.GetOk(p)
						if ok {
							_rawmi, ok := _rawm.(map[string]interface{})
							if !ok {
								return fmt.Errorf("can not convert %T to map[string]interface{}", _rawm)
							}
							t.NestedM = make(map[string]string, len(_rawmi))
							for _k, _v := range _rawmi {
								_raw := _v
								_raws, ok := _raw.(string)
								if !ok {
									return fmt.Errorf("can not convert %T to string", _raws)
								}
								_value := string(string(_raws))
								t.NestedM[_k] = _value
							}
						}
					}
					{
						p := p + "nested_m_obj"

						_raw, ok := d.GetOk(p)
						if ok {
							_rawi, ok := _raw.([]interface{})
							if !ok {
								return fmt.Errorf("can not convert %T to []interface{}", _raw)
							}

							_value := make(map[string]*NestedLevel2)

							for i := range _rawi {
								_rawkey := d.Get(fmt.Sprintf("%v.%v.", p, i) + "key")
								_key, ok := _rawkey.(string)
								if !ok {
									return fmt.Errorf("can not convert %T to string", _rawkey)
								}
								if _key == "" {
									return fmt.Errorf("missing key field in object map NestedMObj")
								}

								_obj := NestedLevel2{}
								_value[_key] = &_obj
								t := &_obj

								{
									p := fmt.Sprintf("%v.%v.value.0.", p, i)
									{

										_raw, ok := d.GetOk(p + "str")

										if ok {
											_raws, ok := _raw.(string)
											if !ok {
												return fmt.Errorf("can not convert %T to string", _raws)
											}
											_value := string(string(_raws))
											t.Str = _value
										}
									}

								}
							}

							t.NestedMObj = _value
						}
					}

				}
			}

			t.NestedMObj = _value
		}
	}

	return nil
}

func SetTestToResourceData(d *schema.ResourceData, t *Test, skip bool) error {
	obj := make(map[string]interface{})

	{
		_v := t.Str

		_value := string(_v)
		obj["str"] = _value
	}
	{
		_v := t.Int32

		_value := int(_v)
		obj["int32"] = _value
	}
	{
		_v := t.Int64

		_value := int(_v)
		obj["int64"] = _value
	}
	{
		_v := t.Float

		_value := float64(_v)
		obj["float"] = _value
	}
	{
		_v := t.Double

		_value := float64(_v)
		obj["double"] = _value
	}
	{
		_v := t.Bool

		_value := bool(_v)
		obj["bool"] = _value
	}
	{
		_v := t.Bytes

		_value := string(_v)
		obj["bytes"] = _value
	}
	{
		_v := t.Timestamp

		_value := _v.Format(time.RFC3339Nano)
		obj["timestamp"] = _value
	}
	{
		_v := t.DurationStd

		_value := time.Duration(_v).String()
		obj["duration_std"] = _value
	}
	{
		_v := t.DurationCustom

		_value := time.Duration(_v).String()
		obj["duration_custom"] = _value
	}
	{
		_v := t.TimestampN
		if _v != nil {

			_value := _v.Format(time.RFC3339Nano)
			obj["timestamp_n"] = _value
		}
	}
	{
		_arr := t.StringA
		if len(_arr) > 0 {
			_raw := make([]string, len(_arr))

			for i, _v := range _arr {
				_value := string(_v)
				_raw[i] = _value
			}

			obj["string_a"] = _raw
		}
	}
	{
		_v, err := SetBoolCustomToResourceData(&t.BoolA)
		if err != nil {
			return err
		}
		obj["bool_a"] = _v
	}
	{
		_arr := t.BytesA
		if len(_arr) > 0 {
			_raw := make([]string, len(_arr))

			for i, _v := range _arr {
				_value := string(_v)
				_raw[i] = _value
			}

			obj["bytes_a"] = _raw
		}
	}
	{
		_arr := t.TimestampA
		if len(_arr) > 0 {
			_raw := make([]string, len(_arr))

			for i, _v := range _arr {
				_value := _v.Format(time.RFC3339Nano)
				_raw[i] = _value
			}

			obj["timestamp_a"] = _raw
		}
	}
	{
		_arr := t.DurationCustomA
		if len(_arr) > 0 {
			_raw := make([]string, len(_arr))

			for i, _v := range _arr {
				_value := time.Duration(_v).String()
				_raw[i] = _value
			}

			obj["duration_custom_a"] = _raw
		}
	}
	{

		if t.Nested != nil {

			msg := make(map[string]interface{})
			obj["nested"] = []interface{}{msg}
			{
				obj := msg
				t := t.Nested

				{
					_v := t.Str

					_value := string(_v)
					obj["str"] = _value
				}
				{
					arr := make([]interface{}, len(t.Nested))

					if len(arr) > 0 {
						for i, t := range t.Nested {
							obj := make(map[string]interface{})
							{
								_v := t.Str

								_value := string(_v)
								obj["str"] = _value
							}

							arr[i] = obj
						}

						obj["nested"] = arr
					}
				}
				{

					m := make(map[string]interface{})
					v := t.NestedM

					if len(v) > 0 {
						for key, _v := range v {
							_value := string(_v)
							m[key] = _value
						}

						obj["nested_m"] = m
					}
				}
				{

					a := make([]interface{}, len(t.NestedMObj))
					n := 0

					for k, v := range t.NestedMObj {
						i := make(map[string]interface{})
						i["key"] = k

						obj := make(map[string]interface{})
						t := v
						{
							_v := t.Str

							_value := string(_v)
							obj["str"] = _value
						}

						i["value"] = []interface{}{obj}

						a[n] = i
						n++
					}

					if len(a) > 0 {
						obj["nested_m_obj"] = a
					}
				}

			}

		}

	}
	{
		arr := make([]interface{}, len(t.NestedA))

		if len(arr) > 0 {
			for i, t := range t.NestedA {
				obj := make(map[string]interface{})
				{
					_v := t.Str

					_value := string(_v)
					obj["str"] = _value
				}
				{
					arr := make([]interface{}, len(t.Nested))

					if len(arr) > 0 {
						for i, t := range t.Nested {
							obj := make(map[string]interface{})
							{
								_v := t.Str

								_value := string(_v)
								obj["str"] = _value
							}

							arr[i] = obj
						}

						obj["nested"] = arr
					}
				}
				{

					m := make(map[string]interface{})
					v := t.NestedM

					if len(v) > 0 {
						for key, _v := range v {
							_value := string(_v)
							m[key] = _value
						}

						obj["nested_m"] = m
					}
				}
				{

					a := make([]interface{}, len(t.NestedMObj))
					n := 0

					for k, v := range t.NestedMObj {
						i := make(map[string]interface{})
						i["key"] = k

						obj := make(map[string]interface{})
						t := v
						{
							_v := t.Str

							_value := string(_v)
							obj["str"] = _value
						}

						i["value"] = []interface{}{obj}

						a[n] = i
						n++
					}

					if len(a) > 0 {
						obj["nested_m_obj"] = a
					}
				}

				arr[i] = obj
			}

			obj["nested_a"] = arr
		}
	}
	{

		m := make(map[string]interface{})
		v := t.NestedM

		if len(v) > 0 {
			for key, _v := range v {
				_value := string(_v)
				m[key] = _value
			}

			obj["nested_m"] = m
		}
	}
	{

		a := make([]interface{}, len(t.NestedMObj))
		n := 0

		for k, v := range t.NestedMObj {
			i := make(map[string]interface{})
			i["key"] = k

			obj := make(map[string]interface{})
			t := v
			{
				_v := t.Str

				_value := string(_v)
				obj["str"] = _value
			}
			{
				arr := make([]interface{}, len(t.Nested))

				if len(arr) > 0 {
					for i, t := range t.Nested {
						obj := make(map[string]interface{})
						{
							_v := t.Str

							_value := string(_v)
							obj["str"] = _value
						}

						arr[i] = obj
					}

					obj["nested"] = arr
				}
			}
			{

				m := make(map[string]interface{})
				v := t.NestedM

				if len(v) > 0 {
					for key, _v := range v {
						_value := string(_v)
						m[key] = _value
					}

					obj["nested_m"] = m
				}
			}
			{

				a := make([]interface{}, len(t.NestedMObj))
				n := 0

				for k, v := range t.NestedMObj {
					i := make(map[string]interface{})
					i["key"] = k

					obj := make(map[string]interface{})
					t := v
					{
						_v := t.Str

						_value := string(_v)
						obj["str"] = _value
					}

					i["value"] = []interface{}{obj}

					a[n] = i
					n++
				}

				if len(a) > 0 {
					obj["nested_m_obj"] = a
				}
			}

			i["value"] = []interface{}{obj}

			a[n] = i
			n++
		}

		if len(a) > 0 {
			obj["nested_m_obj"] = a
		}
	}

	for key, value := range obj {
		err := d.Set(key, value)
		if err != nil {
			return err
		}
	}

	return nil
}
