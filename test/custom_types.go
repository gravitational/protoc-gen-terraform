package test

import schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

// Custom duration type
type Duration int64

// Custom bool array
type BoolCustomArray []bool

type SomethingCustom int

// SchemaBoolCustomArray returns schema for custom bool array
func SchemaBoolCustomArray() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Schema{
			Type: schema.TypeBool,
		},
	}
}
