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
	time "time"

	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Duration custom duration type
type Duration int64

// String returns duration string representation, must be implemented for custom duration type
func (d Duration) String() string {
	return time.Duration(d).String()
}

// BoolCustom custom bool array
type BoolCustom bool

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

// GetBoolCustomFromResourceData interprets data at path as an array of boolean values.
// The values are returned in target
func GetBoolCustomFromResourceData(path string, data *schema.ResourceData, target *[]BoolCustom) error {
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

// SetBoolCustomToResourceData sets
func SetBoolCustomToResourceData(value *[]BoolCustom) interface{} {
	r := make([]interface{}, len(*value))

	for i, v := range *value {
		r[i] = v
	}

	return r
}
