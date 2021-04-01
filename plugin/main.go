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
	"fmt"

	"github.com/gravitational/protoc-gen-terraform/config"
	"github.com/gravitational/protoc-gen-terraform/render"
	"github.com/gravitational/trace"

	"github.com/disiqueira/gotree"
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

	p.writeSchemaStructureComment()

	err := p.writeSchema()
	if err != nil {
		p.Generator.Fail(err.Error())
	}

	err = p.writeGettersSetters()
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

// writeSchema writes schema definition to target file
func (p *Plugin) writeSchema() error {
	for _, message := range p.Messages {
		var buf bytes.Buffer

		err := render.Template(render.SchemaTpl, message, &buf)
		if err != nil {
			return trace.Wrap(err)
		}
		p.P(buf.String())
	}

	return nil
}

// writeGetters writes unmarshallers definition to target file
func (p *Plugin) writeGettersSetters() error {
	for _, message := range p.Messages {
		var buf bytes.Buffer

		err := render.Template(render.GetTpl, message, &buf)
		if err != nil {
			return trace.Wrap(err)
		}

		err = render.Template(render.SetTpl, message, &buf)
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

	for _, i := range config.CustomImports {
		p.AddImport(generator.GoImportPath(i))
	}
}

// writeSchemaComment writes comment with message structure
func (p *Plugin) writeSchemaStructureComment() {
	for _, m := range p.Messages {
		t := gotree.New(fmt.Sprintf("%s (%s)", m.Name, commentToSingleLine(m.RawComment)))
		p.writeMessageComment(m, t)
		p.P(appendSlashSlash(t.Print(), true))
		p.P()
	}
}

// writeMessageComment generates message structure comment part for a message
func (p *Plugin) writeMessageComment(m *Message, t gotree.Tree) {
	for _, f := range m.Fields {
		s := f.NameSnake

		if f.Kind == "SINGULAR_ELEMENTARY" || f.Kind == "REPEATED_ELEMENTARY" {
			s = s + ":" + f.SchemaRawType
		}

		if f.Kind == "REPEATED_ELEMENTARY" || f.Kind == "OBJECT_MAP" || f.Kind == "REPEATED_MESSAGE" {
			s = "[" + s + "]"
		}

		if f.Kind == "MAP" {
			s = s + ":map"
		}

		if f.Kind == "CUSTOM_TYPE" {
			s = s + " !custom schema, see target code!"
		}

		if f.RawComment != "" {
			s = s + " (" + commentToSingleLine(f.RawComment) + ")"
		}

		i := t.Add(s)

		if f.Kind == "OBJECT_MAP" {
			i.Add("key:string")
			if f.MapValueField.Kind == "CUSTOM_TYPE" {
				i.Add("value !custom schema, see target code!")
			} else {
				p.writeMessageComment(f.MapValueField.Message, i.Add("value"))
			}
		} else if f.Kind == "REPEATED_MESSAGE" || f.Kind == "SINGULAR_MESSAGE" {
			p.writeMessageComment(f.Message, i)
		}
	}
}
