package render

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
)

// Template renders template from string template
func Template(content string, name string, pipeline interface{}) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	tpl, err := template.New(name).Funcs(sprig.TxtFuncMap()).Parse(content)
	if err != nil {
		return nil, err
	}

	err = tpl.ExecuteTemplate(&buf, name, pipeline)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
