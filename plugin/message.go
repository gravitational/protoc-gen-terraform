package plugin

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stoewer/go-strcase"
)

// Message holds reflection information about message
type Message struct {
	Name       string // Type name
	NameSnake  string // Type name in snake case
	GoTypeName string // Go type name for this message with package name

	Fields []*Field
}

func (p *Plugin) reflectMessage(d *generator.Descriptor) *Message {
	message := &Message{}

	name := d.GetName()

	message.Name = name
	message.NameSnake = strcase.SnakeCase(name)
	message.GoTypeName = d.File().GoPackageName() + "." + name

	p.reflectFields(message, d)

	return message
}
