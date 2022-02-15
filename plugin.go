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
package main

import (
	"bytes"
	"io"
	"sort"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
)

const (
	// pluginName contains plugin name
	pluginName = "terraform"
)

// Plugin is terraform generator plugin
type Plugin struct {
	*generator.Generator
	generator.PluginImports
	// Config represents the plugin configuration
	Config *Config
	// Messages represents list of the messages in a protoc file
	Messages []*Message
	// Imports represents import package->qualifier dictionary
	Imports Imports
}

// NewPlugin creates the new plugin
func NewPlugin() *Plugin {
	return &Plugin{Messages: make([]*Message, 0)}
}

// Init initializes plugin and sets the generator instance
func (p *Plugin) Init(g *generator.Generator) {
	var err error

	p.Generator = g
	p.PluginImports = generator.NewPluginImports(p.Generator)

	p.Config, err = ReadConfig(g.Param)
	if err != nil {
		p.Generator.Fail(err.Error())
	}
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return pluginName
}

// RegisterMessage adds a new entry to AllMessages
func (p *Plugin) RegisterMessage(m *Message) {
	p.Messages = append(p.Messages, m)
}

// Generate goes over messages in the file passed from gogo, builds reflection structs and writes a target file
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	var buf bytes.Buffer

	log.Printf("Processing: %s", *file.Name)

	p.Imports = NewImports(p.PluginImports)
	p.build(file)

	err := p.write(p.Messages, &buf)
	if err != nil {
		p.Generator.Fail(err.Error())
	}

	p.P(buf.String())
}

// GetConfig returns plugin config
func (p *Plugin) GetConfig() *Config {
	return p.Config
}

// GetGenerator returns plugin generator
func (p *Plugin) GetGenerator() *generator.Generator {
	return p.Generator
}

// GetImports returns plugin imports
func (p *Plugin) GetImports() *Imports {
	return &p.Imports
}

// build builds the message dictionary from a messages in protoc file
func (p *Plugin) build(file *generator.FileDescriptor) {
	for _, message := range file.Messages() {
		m, err := BuildMessage(p, message, true, "")
		if err != nil {
			log.WithError(err).Warningf("failed to build the message %v", message.GetName())
			continue
		}

		// A message is nil if it is not required in the output
		if m != nil {
			p.RegisterMessage(m)
		}
	}
	// Sort messages if required
	if p.Config.Sort {
		sort.Slice(p.Messages, func(i, j int) bool {
			return p.Messages[i].Name < p.Messages[j].Name
		})
	}
}

// write writes schema and meta to output file
func (p *Plugin) write(m []*Message, out io.Writer) error {
	for _, message := range m {
		if !message.IsRoot {
			continue
		}

		g := NewMessageSchemaGenerator(message, &p.Imports)
		_, err := g.Generate(out)
		if err != nil {
			return trace.Wrap(err)
		}
	}

	for _, message := range m {
		if !message.IsRoot {
			continue
		}

		f := NewMessageCopyFromGenerator(message, &p.Imports)
		_, err := f.Generate(out)
		if err != nil {
			return trace.Wrap(err)
		}

		t := NewMessageCopyToGenerator(message, &p.Imports)
		_, err = t.Generate(out)
		if err != nil {
			return trace.Wrap(err)
		}
	}

	_, err := NewSharedCodeGenerator(&p.Imports).Write(out)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}
