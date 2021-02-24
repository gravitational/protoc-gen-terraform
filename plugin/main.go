// Package plugin is gogoprotobuf package for Terraform code generation
package plugin

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/protoc-gen-terraform/config"
	"github.com/gravitational/protoc-gen-terraform/render"
	"github.com/gravitational/trace"
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
	config.ParseExcludeFields(g.Param["exclude_fields"])
	config.ParseDefaultPkgName(g.Param["pkg"])
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return name
}

// Generate goes over messages in the file passed from gogo, builds reflection structs and writes a target file
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	logrus.Printf("Processing: %s", *file.Name)

	// Adds Terraform package imports to target file
	p.setImports()

	p.reflect(file)
	p.writeSchema()
	p.writeUnmarshallers()
}

// reflect builds message dictionary from a messages in protoc file
func (p *Plugin) reflect(file *generator.FileDescriptor) {
	for _, message := range file.Messages() {
		m := BuildMessage(p.Generator, message, true)
		if m != nil {
			p.Messages[m.GoTypeName] = m
		}
	}
}

// writeSchema writes schema definition to target file
func (p *Plugin) writeSchema() {
	for _, message := range p.Messages {
		buf, err := render.Template(render.SchemaTpl, message)
		if err != nil {
			p.Generator.Fail(trace.Wrap(err).Error())
		}
		p.P(buf.String())
	}
}

// writeUnmarshallers writes unmarshallers definition to target file
func (p *Plugin) writeUnmarshallers() {
	for _, message := range p.Messages {
		buf, err := render.Template(render.UnmarshalTpl, message)
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
}
