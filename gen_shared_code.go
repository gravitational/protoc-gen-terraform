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
	"io"
	"strings"

	_ "embed"
)

//go:embed shared_code.go.tpl
var shared string

// Generator interface represents generic generator interface
type Generator interface {
	Write(io.Writer) (int, error)
}

type SharedCodeGenerator struct {
	i *Imports
}

func NewSharedCodeGenerator(i *Imports) SharedCodeGenerator {
	return SharedCodeGenerator{i}
}

// Generate returns generated code for globals
func (s SharedCodeGenerator) Write(writer io.Writer) (int, error) {
	value := strings.NewReplacer(
		"diag.Severity", s.i.WithPackage(Diag, "Severity"),
		"diag.SeverityError", s.i.WithPackage(Diag, "SeverityError"),
		"fmt.Sprintf", s.i.WithPackage("fmt", "Sprintf"),
		"diag.Diagnostic", s.i.WithPackage(Diag, "Diagnostic"),
	).Replace(shared)

	return writer.Write([]byte(value + "\n"))
}
