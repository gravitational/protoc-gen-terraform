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

package desc

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/trace"
)

// Message represents metadata about protobuf message
type Message struct {
	// Name contains type name
	Name string
	// GoType contains Go type name for this message with package name
	GoType string
	// Fields contains the collection of fields
	Fields []*Field
	// RawComment is field comment in proto file without // prepended
	RawComment string
	// Comment is field comment in proto file with // prepended
	Comment string
	// Path represents a path in the schema to the current message (types.UserV2.Metadata)
	Path string
	// NamePath contains type name with nesting (UserV2Metadata)
	NamePath string
	// IsRoot indicates that the message is root in a file
	IsRoot bool
}

// BuildMessage builds Message from its protobuf descriptor.
//
// check must be false for nested messages. Otherwise, we'll have to specify a full list of allowed
// types which would be overexplicit. Use ExcludeFields if you need to skip a nested fields.
//
// nil returned means that operation was successful, but message needs to be skipped.
func BuildMessage(plugin Plugin, desc *generator.Descriptor, isRoot bool, path string) (*Message, error) {
	c := NewMessageBuildContext(plugin, desc, path)

	// Check if message is specified in export type list
	if c.IsExcluded() && isRoot {
		// This is not an error, we must just skip this message
		return nil, nil
	}

	if c.HasOneOf() {
		return nil, trace.Errorf("oneOf messages are not supported yet %v", c.GetGoType())
	}

	fields, err := BuildFields(c)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	// Skip objects with no fields
	if len(fields) == 0 {
		return nil, nil
	}

	message := &Message{
		NamePath: c.GetNamePath(),
		Name:     c.GetName(),
		GoType:   c.GetGoType(),
		Path:     c.GetPath(),
		Fields:   fields,
		IsRoot:   isRoot,
	}

	message.RawComment, message.Comment = c.GetComment()

	return message, nil
}
