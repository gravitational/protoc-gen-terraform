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
