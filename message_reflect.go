package main

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stoewer/go-strcase"
)

type messageReflect struct {
	name      string // Type name
	snakeName string // Type name in snake case
	goType    string // Go type name for this message with package name

	fields []*fieldReflect // Fields
}

func (p *Plugin) newMessageReflect(d *generator.Descriptor) *messageReflect {
	message := &messageReflect{}

	message.name = d.GetName()
	message.snakeName = strcase.SnakeCase(message.name)
	message.goType = p.pkg.Use() + "." + message.name

	p.newFieldsReflect(message, d)

	return message
}
