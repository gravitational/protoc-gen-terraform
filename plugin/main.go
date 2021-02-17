package plugin

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/trace"
	"github.com/gzigzigzeo/protoc-gen-terraform/config"
	"github.com/sirupsen/logrus"
)

const (
	name          = "terraform"                                                      // Plugin name
	schemaPkg     = "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"     // Terraform schema package
	validationPkg = "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation" // Terraform validation package
)

// Plugin is terraform generator plugin
type Plugin struct {
	*generator.Generator
	generator.PluginImports

	// Map of reflected messages, public just in case some post analysis is required
	Messages map[string]*Message
}

// NewPlugin creates the new plugin
func NewPlugin() *Plugin {
	return &Plugin{
		Messages: make(map[string]*Message),
	}
}

// Init initializes plugin and sets the generator instance
func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g

	config.ParseTypes(g.Param["types"])
	config.ParseExcludeFields(g.Param["excludeFields"])
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return name
}

// Generate is the plugin body
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	logrus.Printf("Processing: %s", *file.Name)

	p.setImports()

	for _, message := range file.Messages() {
		m := BuildMessage(p.Generator, message, true)
		if m != nil {
			p.Messages[m.GoTypeName] = m
		}
	}

	for _, message := range p.Messages {
		buf, err := message.GoSchemaString()
		if err != nil {
			p.Generator.Fail(trace.Wrap(err).Error())
		}
		p.P(buf.String())
	}

	for _, message := range p.Messages {
		buf, err := message.GoUnmarshalString()
		if err != nil {
			p.Generator.Fail(trace.Wrap(err).Error())
		}
		p.P(buf.String())
	}
}

// setImports sets import definitions for current file
func (p *Plugin) setImports() {
	p.PluginImports = generator.NewPluginImports(p.Generator)

	// So those could be referenced via schema. and validation.
	p.AddImport(schemaPkg)
	p.AddImport(validationPkg)

	// TODO: Temporary
	p.AddImport("github.com/gravitational/teleport/api/types")
	p.AddImport("github.com/gravitational/teleport/api/types/wrappers")
}

// // isMessageRequired returns true if message was marked for export via command-line args
// func (p *Plugin) isMessageRequired(d *generator.Descriptor) bool {
// 	typeName := d.File().GetPackage() + "." + d.GetName()
// 	required := slice.Contains(config.Types, typeName)

// 	return required
// }

// // reflectMessage reflects message type
// func (p *Plugin) reflectMessage(d *generator.Descriptor, nested bool) *Message {
// 	if !nested && !p.isMessageRequired(d) {
// 		return nil
// 	}

// 	logrus.Println(d.GetName())

// 	name := d.GetName()

// 	if p.Messages[name] != nil {
// 		return p.Messages[name]
// 	}

// 	message := p.buildMessage(d)

// 	if !nested {
// 		p.Messages[name] = message
// 	}

// 	return message
// }

// // reflectFields builds array of message.Fields
// func (p *Plugin) reflectFields(m *Message, d *generator.Descriptor) {
// 	for _, f := range d.GetField() {
// 		if !p.isFieldIgnored(d, f) {
// 			f, ok := p.reflectField(d, f)
// 			if ok {
// 				m.Fields = append(m.Fields, f)
// 			}
// 		}
// 	}
// }

// // isMessageRequired returns true if message was marked for export via command-line args
// func (p *Plugin) isFieldIgnored(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) bool {
// 	fieldName := d.File().GetPackage() + "." + d.GetName() + "." + f.GetName()
// 	ignored := slice.Contains(config.ExcludeFields, fieldName)

// 	return ignored
// }

// // reflectField builds field reflection structure, or returns nil in case field must be skipped
// func (p *Plugin) reflectField(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) (*Field, bool) {
// 	b := p.newFieldBuilder(d, f)
// 	ok := b.build()
// 	return b.field, ok
// }
