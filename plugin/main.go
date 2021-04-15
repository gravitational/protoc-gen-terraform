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

// Package plugin is gogoprotobuf package for Terraform code generation
package plugin

import (
	"bytes"

	"github.com/gravitational/protoc-gen-terraform/config"
	"github.com/gravitational/protoc-gen-terraform/render"
	"github.com/gravitational/trace"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/sirupsen/logrus"
)

const (
	// name contains plugin name
	name = "terraform"

	// schemaPkg contains name of Terraform schema package
	schemaPkg = "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	// validationPkg contains name of Terraform validation package
	validationPkg = "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	// accessorsPkg contains SchemaMeta definition
	accessorsPkg = "github.com/gravitational/protoc-gen-terraform/accessors"
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

	err := config.Read(g.Param)
	if err != nil {
		p.Generator.Fail(err.Error())
	}
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

	p.build(file)

	err := p.writeVars()
	if err != nil {
		p.Generator.Fail(err.Error())
	}

	err = p.writeSchema()
	if err != nil {
		p.Generator.Fail(err.Error())
	}
}

// reflect builds message dictionary from a messages in protoc file
func (p *Plugin) build(file *generator.FileDescriptor) {
	for _, message := range file.Messages() {
		m, err := BuildMessage(p.Generator, message, true, "")

		if err != nil {
			logrus.Warning(err)
			continue
		}

		if m != nil {
			p.Messages[m.GoTypeName] = m
		}
	}
}

// writeVars writes top level variables and methods
func (p *Plugin) writeVars() error {
	var buf bytes.Buffer

	err := render.Template(render.VarsTpl, p.Messages, &buf)
	if err != nil {
		return trace.Wrap(err)
	}
	p.P(buf.String())

	return nil
}

// writeSchema writes schema definition to target file
func (p *Plugin) writeSchema() error {
	for _, message := range p.Messages {
		var buf bytes.Buffer

		err := render.Template(render.SchemaTpl, message, &buf)
		if err != nil {
			return trace.Wrap(err)
		}
		p.P(buf.String())

		buf.Reset()

		err = render.Template(render.MetaTpl, message, &buf)
		if err != nil {
			return trace.Wrap(err)
		}
		p.P(buf.String())
	}

	return nil
}

// setImports sets import definitions for current file
func (p *Plugin) setImports() {
	p.PluginImports = generator.NewPluginImports(p.Generator)

	// So those could be referenced via schema. and validation.
	p.AddImport(schemaPkg)
	p.AddImport(validationPkg)
	p.AddImport(accessorsPkg)

	for _, i := range config.CustomImports {
		p.AddImport(generator.GoImportPath(i))
	}
}
