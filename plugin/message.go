package plugin

import (
	"bytes"
	"errors"
	"path"
	"runtime"
	"text/template"

	"github.com/Masterminds/sprig"
)

var (
	templateFilename = "message.tpl"
	templatesDir     = "_tpl"
)

// Message holds reflection information about message
type Message struct {
	Name       string   // Type name
	NameSnake  string   // Type name in snake case, schema field name
	GoTypeName string   // Go type name for this message with package name
	Fields     []*Field // Collection of fields
	// TODO: Comments
}

// GoString returns go code for this message as terraform schema
func (m *Message) GoString() (*bytes.Buffer, error) {
	var buf bytes.Buffer

	_, filename, _, ok := runtime.Caller(1)

	if !ok {
		return nil, errors.New("Can't get path to runtime file")
	}

	filepath := path.Join(path.Dir(filename), templatesDir, templateFilename)

	tpl, err := template.New("message").Funcs(sprig.TxtFuncMap()).ParseFiles(filepath)

	if err != nil {
		return nil, err
	}

	err = tpl.ExecuteTemplate(&buf, templateFilename, m)
	if err != nil {
		return nil, err
	}

	return &buf, nil

}
