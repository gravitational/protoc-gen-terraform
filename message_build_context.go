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

package main

import (
	"strconv"
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

const (
	rootPackage = "." // indicates that a message belongs to the topmost package
)

// MessageBuildContext represents an utilty struct to pass metadata as a build functions parameter
type MessageBuildContext struct {
	plugin   *Plugin
	imports  *Imports
	config   *Config // Syntax shortcut
	gen      *generator.Generator
	desc     *generator.Descriptor
	path     string
	namePath string
}

// NewMessageBuildContext creates new message build context
func NewMessageBuildContext(plugin *Plugin, desc *generator.Descriptor, path string) MessageBuildContext {
	namePath := path
	if path == "" {
		namePath = desc.GetName()
	}

	// Split package name and all dots
	i := strings.Index(path, ".")
	if i > -1 {
		namePath = namePath[i:]
	}
	namePath = strings.Replace(namePath, ".", "", -1)

	return MessageBuildContext{plugin, plugin.GetImports(), plugin.GetConfig(), plugin.GetGenerator(), desc, path, namePath}
}

// GetGoType returns go type for a message
func (c MessageBuildContext) GetGoType() string {
	name := c.desc.GetName()

	if c.desc.GoImportPath() == rootPackage {
		if c.config.DefaultPackageName == "" {
			return name
		}
		return c.config.DefaultPackageName + "." + name
	}

	return name
}

// GetComment returns the message source code comment and it's raw text
func (c MessageBuildContext) GetComment() string {
	p := c.desc.Path()

	for _, l := range c.desc.File().GetSourceCodeInfo().GetLocation() {
		if c.GetLocationPath(l) == p {
			c := Comment(strings.Trim(l.GetLeadingComments(), "\n"))
			return c.ToSingleLine()
		}
	}

	return ""
}

// IsExcluded returns true if the message is not specified in the type list
func (c MessageBuildContext) IsExcluded() bool {
	_, ok := c.config.Types[c.GetPath()]
	return !ok
}

// GetName returns the message name
func (c MessageBuildContext) GetName() string {
	return c.desc.GetName()
}

// GetPath returns the path to the message
func (c MessageBuildContext) GetPath() string {
	if c.path == "" {
		name := c.desc.GetName()

		if c.desc.GoImportPath() == rootPackage {
			return name
		}

		return c.desc.File().GetPackage() + "." + name
	}

	return c.path
}

// GetNamePath returns the path to the message
func (c MessageBuildContext) GetNamePath() string {
	return c.namePath
}

// GetLocationPath returns the source line path for a given source code location.
func (c MessageBuildContext) GetLocationPath(l *descriptor.SourceCodeInfo_Location) string {
	s := make([]string, len(l.GetPath()))

	for i, v := range l.GetPath() {
		s[i] = strconv.Itoa(int(v))
	}

	return strings.Join(s, ",")
}

// GetInjectedFields returns array of injected fields
func (c *MessageBuildContext) GetInjectedFields() []InjectedField {
	v, ok := c.config.InjectedFields[c.GetPath()]
	if ok {
		return v
	}

	return []InjectedField{}
}

// GetOneOfNames returns the names of OneOf groups in this message
func (c *MessageBuildContext) GetOneOfNames() []string {
	s := make([]string, len(c.desc.OneofDecl))
	for i, d := range c.desc.OneofDecl {
		s[i] = d.GetName()
	}
	return s
}
