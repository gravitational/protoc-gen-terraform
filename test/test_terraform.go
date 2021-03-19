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
		// DurationCustomA REPEATED_ELEMENTARY
		"duration_custom_a": {
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
					// NestedM MAP
					"nested_m": {
						Optional: true,
						Type:     schema.TypeMap,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// NestedMObj OBJECT_MAP
					"nested_m_obj": {
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
					// NestedM MAP
					"nested_m": {
						Optional: true,
						Type:     schema.TypeMap,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
					// NestedMObj OBJECT_MAP
					"nested_m_obj": {
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
		// NestedM MAP
		"nested_m": {
			Optional: true,
			Type:     schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		// NestedMObj OBJECT_MAP
		"nested_m_obj": {
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
								// NestedM MAP
								"nested_m": {
									Optional: true,
									Type:     schema.TypeMap,
									Elem: &schema.Schema{
										Type: schema.TypeString,
									},
								},
								// NestedMObj OBJECT_MAP
								"nested_m_obj": {
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
			_value, err := time.Parse(time.RFC3339, _raws)
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
			_value, err := time.Parse(time.RFC3339, _raws)
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
				_value, err := time.Parse(time.RFC3339, _raws)
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

func SetTestToResourceData(d *schema.ResourceData, t *Test) error {
	p := ""

	{
		err := d.Set(p+"str", t.Str)
		if err != nil {
			return err
		}
	}
	// {
	// 	err := d.Set(p+"int32", t.Int32)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// {
	// 	err := d.Set(p+"int64", t.Int64)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// {
	// 	err := d.Set(p+"float", t.Float)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// {
	// 	err := d.Set(p+"double", t.Double)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// {
	// 	err := d.Set(p+"bool", t.Bool)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// {
	// 	err := d.Set(p+"bytes", t.Bytes)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// {
	// 	err := d.Set(p+"timestamp", t.Timestamp)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// {
	// 	err := d.Set(p+"duration_std", t.DurationStd)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// {
	// 	err := d.Set(p+"duration_custom", t.DurationCustom)
	// 	if err != nil {
	// 		return err
	// 	}
	// }
	// {
	// 	err := d.Set(p+"timestamp_n", t.TimestampN)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	return nil
}
