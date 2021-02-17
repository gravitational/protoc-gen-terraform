package plugin

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stoewer/go-strcase"

	// go:embed won't work otherwise
	_ "embed"
)

var (
	//go:embed _tpl/message_schema.tpl
	schemaTpl string
	//go:embed _tpl/message_unmarshal.tpl
	unmarshalTpl string
)

// Message holds reflection information about message
type Message struct {
	Name       string   // Type name
	NameSnake  string   // Type name in snake case, schema field name
	GoTypeName string   // Go type name for this message with package name
	Fields     []*Field // Collection of fields
	// TODO: Comments
}

// build Builds message
func (p *Plugin) buildMessage(d *generator.Descriptor) *Message {
	name := d.GetName()

	message := &Message{}

	message.Name = name
	message.NameSnake = strcase.SnakeCase(name)
	message.GoTypeName = d.File().GetPackage() + "." + name

	p.reflectFields(message, d)

	return message
}

// GoUnmarshalString returns go code for this message as unmarshaller
func (m *Message) GoUnmarshalString() (*bytes.Buffer, error) {
	return m.renderTemplate(unmarshalTpl, "unmarshal")
}

// GoSchemaString returns go code for this message as terraform schema
func (m *Message) GoSchemaString() (*bytes.Buffer, error) {
	return m.renderTemplate(schemaTpl, "schema")
}

// renderTemplate renders template from embedded template
func (m *Message) renderTemplate(content string, name string) (*bytes.Buffer, error) {
	var buf bytes.Buffer

	tpl, err := template.New(name).Funcs(sprig.TxtFuncMap()).Parse(content)
	if err != nil {
		return nil, err
	}

	err = tpl.ExecuteTemplate(&buf, name, m)
	if err != nil {
		return nil, err
	}

	return &buf, nil
}
