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

// Package render is the utility package with embeds templates and render method
package render

import (
	"io"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/gravitational/trace"

	// go:embed won't work otherwise
	_ "embed"
)

var (
	// SchemaTpl is schema template
	//go:embed _tpl/message_schema.tpl
	SchemaTpl string

	// MetaTpl is schema metadata template
	//go:embed _tpl/message_meta.tpl
	MetaTpl string

	// LicenseTpl is license message template
	//go:embed _tpl/license.tpl
	LicenseTpl string

	// VarsTpl is the template for global variables, type definitions and shared methods
	//go:embed _tpl/vars.tpl
	VarsTpl string

	// NOTE: soon be obsolete
	// GetTpl is unmarshaller template
	//go:embed _tpl/message_get.tpl
	GetTpl string

	// SetTpl is unmarshaller template
	//go:embed _tpl/message_set.tpl
	SetTpl string
)

// Template renders template from string to the specified writer
func Template(content string, pipeline interface{}, w io.Writer) error {
	tpl, err := template.New("template").Funcs(sprig.TxtFuncMap()).Parse(content)
	if err != nil {
		return trace.Wrap(err)
	}

	err = tpl.ExecuteTemplate(w, "template", pipeline)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
