// Package render is the utility package with embeds templates and render method
package render

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"

	// go:embed won't work otherwise
	_ "embed"
)

var (
	// SchemaTpl is schema template
	//go:embed _tpl/message_schema.tpl
	SchemaTpl string

	// UnmarshalTpl is unmarshaller template
	//go:embed _tpl/message_unmarshal.tpl
	UnmarshalTpl string
)

// Template renders template from string to bytes.buffer
func Template(content string, pipeline interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	tpl, err := template.New("template").Funcs(sprig.TxtFuncMap()).Parse(content)
	if err != nil {
		return nil, err
	}

	err = tpl.ExecuteTemplate(&buf, "template", pipeline)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
