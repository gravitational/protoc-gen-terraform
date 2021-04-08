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

// Package accessors contains Get and Set methods for ResourceData
package accessors

import (
	"reflect"

	"github.com/gravitational/trace"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Set assigns object data from object to schema.ResourceData
//
// Example:
//   user := UserV2{Name: "example"}
//   Set(&user, data, SchemaUserV2, MetaUserV2)
func Set(
	obj interface{},
	data *schema.ResourceData,
	sch map[string]*schema.Schema,
	meta map[string]*SchemaMeta,
) error {
	if obj == nil {
		return trace.Errorf("obj must not be nil")
	}

	// Obj must be reference
	t := reflect.Indirect(reflect.ValueOf(obj))
	return getFragment("", &t, meta, sch, data)
}
