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

// Imports represents a collection of imports for a current file
type Imports struct {
	imp map[string]string
	def string
}

// NewImports returns an empty Imports struct with default import set
func NewImports() Imports {
	return Imports{make(map[string]string), ""}
}

// SetDefault sets the default package name
func (s *Imports) SetDefault(qual, path string) {
	s.AddImport(qual, path)
	s.def = qual
}

// AddImport adds a new import to the set
func (s Imports) AddImport(qual, path string) {
	for k := range s.imp {
		if k == path {
			return
		}
	}

	s.AddQual(qual, path)
}

// AddQual adds qualifier for a paht
func (s Imports) AddQual(qual string, path string) {
	s.imp[path] = qual
}

// GetQual returns qualifier by id or by package path
func (s Imports) GetQual(path string) string {
	return s.imp[path]
}

// GoString normalizes Go type passed as an argument. Adds the default package qualifier if required.
// Replaces url to a package with it's qualifier if known.
func (s Imports) GoString(t string, prependDefault bool) string {
	typ, mod := TypAndMod(t)

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
			p, ok := s.imp[path]
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

// TypAndMod return go type name and modifiers (*[])
func TypAndMod(t string) (string, string) {
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
