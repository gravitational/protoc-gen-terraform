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
	fmt "fmt"
	"reflect"
	time "time"

	"github.com/gravitational/protoc-gen-terraform/accessors"
	"github.com/gravitational/trace"
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

// GetBoolCustom interprets data at path as an array of boolean values.
// The values are returned in target
func GetBoolCustom(
	path string,
	target reflect.Value,
	meta *accessors.SchemaMeta,
	sch *schema.Schema,
	data *schema.ResourceData,
) error {
	len, err := accessors.GetLen(path, data)
	if err != nil {
		return err
	}
	if len == 0 {
		// TODO: Share with accessors
		target.Set(reflect.Zero(target.Type()))
		return nil
	}

	t := make([]BoolCustom, len)

	for i := 0; i < len; i++ {
		p := fmt.Sprintf("%v.%v", path, i)

		raw := data.Get(p)
		v, ok := raw.(bool)
		if !ok {
			return trace.Errorf("can not convert %T to bool", raw)
		}

		t[i] = BoolCustom(v)
	}

	target.Set(reflect.ValueOf(t))

	return nil
}

// SetBoolCustomToResourceData sets
func SetBoolCustom(value *[]BoolCustom) (interface{}, error) {
	r := make([]interface{}, len(*value))

	for i, v := range *value {
		r[i] = v
	}

	return r, nil
}
