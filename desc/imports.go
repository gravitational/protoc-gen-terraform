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

package desc

import (
	"strings"
)

// TODO: Return unknown type if package/qual is unknown
// TODO: Return unknown type if package/qual is duplicate

// Imports represents the collection of file imports.
// It is responsible for type names normalization and allows types to be referenced by full package names.
type Imports struct {
	quals map[string]string
	def   string
}

// NewImports returns an empty Imports struct with default import set
func NewImports() Imports {
	return Imports{make(map[string]string), ""}
}

// SetDefaultQual sets the default package name and path
func (s *Imports) SetDefaultQual(qual, path string) {
	s.AddQual(qual, path)
	s.def = qual
}

// AddQual adds a new package to the set
func (s Imports) AddQual(qual, path string) {
	for k := range s.quals {
		if k == path {
			return
		}
	}

	s.quals[path] = qual
}

// GetQual returns qualifier by id or by package path
func (s Imports) GetQual(path string) string {
	return s.quals[path]
}

// GoString normalizes Go type passed as an argument.
//   * Package qualifiers must be added explicitly using AddImport().
//   * Accepts type name in a form of "full/package/name.TypeName" and replaces full package name with it's qualifier.
//   * If prependDefault is true and package name is missing it appends the default package qualifier.
//   * Built-in primitive types are returned as is.
//   * Slice and map value types are converted as well.
func (s Imports) GoString(t string, prependDefault bool) string {
	typ, mod := typAndMod(t)

	switch typ {
	case "bool", "string", // No need to store it dynamically
		"int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
		"byte", "rune", "float32", "float64", "complex64", "complex128":
		return t // It's a primitive type, do nothing
	default:
		// There is a package prepended to this type already
		if strings.Contains(typ, ".") {
			i := strings.LastIndex(typ, ".")
			path := typ[0:i]
			name := typ[i+1:]

			// Qualifier for this import already exists
			p, ok := s.quals[path]
			if ok {
				path = p
			}

			return mod + path + "." + name
		}

		// Default package name should be prepended if it exists and is required
		if s.def != "" && prependDefault {
			return mod + s.def + "." + typ
		}

		// Leave intact
		return t
	}
}

// WithPackage concatenates package and typ, returns normalized type name
func (s Imports) WithPackage(path string, typ string) string {
	return s.GoString(path+"."+typ, false)
}

// typAndMod return go type name and it's modifiers (*[])
func typAndMod(t string) (string, string) {
	var mod string
	typ := t

	// Split modifiers and type information
	index := strings.LastIndexAny(typ, "[]*")
	if index > -1 {
		mod = typ[0 : index+1]
		typ = typ[index+1:]
	}

	return typ, mod
}
