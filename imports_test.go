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

package main

import (
	"testing"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	gen := generator.New()
	imports := generator.NewPluginImports(gen)
	importPathOverrides := map[string]string{
		"types": "example.com/package/types",
	}
	s := NewImports(imports, importPathOverrides)

	require.Equal(t, s.WithType("string"), "string")
	require.Equal(t, s.WithType("[]string"), "[]string")
	require.Equal(t, s.WithType("[]Type"), "[]Type")
	require.Equal(t, s.WithType("[]*Type"), "[]*Type")

	require.Equal(t, s.WithType("github.com/test/ext.SomeType"), "github_com_test_ext.SomeType")
	require.Equal(t, s.WithType("github.com/test/ext.SomeOtherType"), "github_com_test_ext.SomeOtherType")
	require.Equal(t, s.WithType("github.com/other_test/ext.SomeType"), "github_com_other_test_ext.SomeType")

	require.Equal(t, s.PrependPackageNameIfMissing("string", "test"), "string")
	require.Equal(t, s.PrependPackageNameIfMissing("[]Test", "test"), "[]test.Test")

	require.Equal(t, s.WithType("UseValidator(1)"), "UseValidator(1)")
	require.Equal(t, s.WithType("github.com/hashicorp/terraform-plugin-framework/tfsdk.UseStateForUnknown()"), "github_com_hashicorp_terraform_plugin_framework_tfsdk.UseStateForUnknown()")
	require.Equal(t, s.WithType(`UseValidator("teleport.dev/origin")`), `UseValidator("teleport.dev/origin")`)

	require.Equal(t, s.WithType("types.Values"), "example_com_package_types.Values")
	require.Equal(t, s.WithType("*types.Values"), "*example_com_package_types.Values")
	require.Equal(t, s.PrependPackageNameIfMissing("[]Test", "types"), "[]example_com_package_types.Test")
}
