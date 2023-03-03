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
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

const (
	// SDK represents the path to Terraform SDK package
	SDK = "github.com/hashicorp/terraform-plugin-framework/tfsdk"
	// Types represents the path to Terraform types package
	Types = "github.com/hashicorp/terraform-plugin-framework/types"
	// Diag represents the path to Terraform diag package
	Diag = "github.com/hashicorp/terraform-plugin-framework/diag"
	// Attr represents the name of Terraform attr package
	Attr = "github.com/hashicorp/terraform-plugin-framework/attr"
	// TFTypes represents the name of Terraform SDK TFTypes package
	TFTypes = "github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Imports represents the collection of imports and acts as the facade for generator.PluginImports
// All normalization methods require full (or empty) package name to be passed.
// They must be called in code generation methods only.
type Imports struct {
	// pluginImports represents the managed plugin imports
	pluginImports generator.PluginImports

	// qualifiers represents the map of package names to qualifiers
	qualifiers map[string]generator.Single

	// importPathOverrides holds import paths that should be used for the
	// included package names.
	importPathOverrides map[string]string
}

// NewImports returns an empty Imports struct with default import set
func NewImports(pluginImports generator.PluginImports, importPathOverrides map[string]string) Imports {
	return Imports{pluginImports, make(map[string]generator.Single), importPathOverrides}
}

// WithPackage concatenates package and typ, returns normalized type name
func (i Imports) WithPackage(pkg, typ string) string {
	return i.WithType(pkg + "." + typ)
}

// WithType takes type name with full package name prepended. Package name is converted to qualifier.
// It handles with array and map elem types, but leaves keys intact.
func (i Imports) WithType(t string) string {
	typ, mod := i.typAndMod(t)

	if !strings.Contains(i.typBeforeBracket(typ), ".") {
		return t
	}

	return i.appendQual(typ, mod)
}

// PrependPackageNameIfMissing prepends package name if it's missing in the type
func (i Imports) PrependPackageNameIfMissing(t, pkg string) string {
	typ, mod := i.typAndMod(t)

	if strings.Contains(i.typBeforeBracket(typ), ".") || pkg == "" || i.isBuiltinType(typ) {
		return t
	}

	return i.appendQual(pkg+"."+typ, mod)
}

// isBuiltinType returns true if t represents built-in type
func (i Imports) isBuiltinType(t string) bool {
	switch t {
	case "bool", "string",
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"byte", "rune", "float32", "float64", "complex64", "complex128":
		return true
	default:
		return false
	}
}

func (i Imports) appendQual(typ, mod string) string {
	pos := strings.LastIndex(i.typBeforeBracket(typ), ".")
	path := typ[0:pos]
	name := typ[pos+1:]

	var qualifier string

	// Qualifier for this import already exists
	s, ok := i.qualifiers[path]
	if ok {
		qualifier = s.Name()
	} else {
		if override, ok := i.importPathOverrides[path]; ok {
			path = override
		}
		// Register new qualifier
		single := i.pluginImports.NewImport(path)
		single.Use()
		qualifier = single.Name()
		i.qualifiers[qualifier] = single
	}

	return mod + qualifier + "." + name
}

// typBeforeBracket returns type of a function
func (i Imports) typBeforeBracket(typ string) string {
	pos := strings.Index(typ, "(")
	if pos > -1 {
		return typ[0:pos]
	}

	return typ
}

// typAndMod return go type name and it's modifiers (*[])
func (i Imports) typAndMod(t string) (string, string) {
	var mod string
	typ := t

	// Split modifiers and type information
	index := strings.LastIndexAny(typ, "[]*")
	if index > -1 {
		typ = t[index+1:]
		mod = t[0 : index+1]
	}

	return typ, mod
}
