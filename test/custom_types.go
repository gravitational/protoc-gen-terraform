package test

import (
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Duration custom duration type
type Duration int64

// BoolCustom custom bool array
type BoolCustom bool

// SomethingCustom some custom value
type SomethingCustom int

// SchemaBoolCustom returns schema for custom bool array
func SchemaBoolCustom() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
}

// UnmarshalBoolCustom custom unmarshaller
func UnmarshalBoolCustom(path string, data *schema.ResourceData, target *[]BoolCustom) error {
	rawi, ok := data.GetOk(path)
	if ok {
		arr := rawi.([]interface{})
		*target = make([]BoolCustom, len(arr))

		for i := 0; i < len(arr); i++ {
			(*target)[i] = BoolCustom(arr[i].(bool))
		}
	}

	return nil
}
