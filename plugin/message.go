/*
Copyright 2015-2020 Gravitational, Inc.

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

package plugin

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/protoc-gen-terraform/config"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/stew/slice"
)

var (
	cache map[string]*Message = make(map[string]*Message)
)

// Message holds reflection information about message
type Message struct {
	// Name type name
	Name string

	// NameSnake type name in snake case (Terraform schema field name)
	NameSnake string

	// GoTypeName Go type name for this message with package name
	GoTypeName string

	// Fields Collection of fields
	Fields []*Field
}

// BuildMessage builds Message from its protobuf descriptor
// checkValiditiy should be false for nested messages. We do not check them over allowed types,
// otherwise it will be overexplicit. Use excludeFields to skip fields.
func BuildMessage(g *generator.Generator, d *generator.Descriptor, checkValidity bool) *Message {
	typeName := getMessageTypeName(d)

	// Check if message is specified in export type list
	if !slice.Contains(config.Types, typeName) && checkValidity {
		return nil
	}

	if cache[typeName] != nil {
		return cache[typeName]
	}

	for _, field := range d.GetField() {
		if field.OneofIndex != nil {
			logrus.Println("Oneof messages are not supported yet")
			return nil
		}
	}

	name := d.GetName()

	message := &Message{
		Name:       name,
		NameSnake:  strcase.SnakeCase(name),
		GoTypeName: typeName,
	}

	BuildFields(message, g, d)

	return message
}

// getMessageTypeName returns full message name, with prepended DefaultPkgName if needed
func getMessageTypeName(d *generator.Descriptor) string {
	if d.GoImportPath() != "." {
		return d.File().GetPackage() + "." + d.GetName()
	}
	if config.DefaultPkgName != "" {
		return config.DefaultPkgName + "." + d.GetName()
	}
	return d.GetName()
}
