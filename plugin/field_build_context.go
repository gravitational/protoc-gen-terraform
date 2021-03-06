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

package plugin

import (
	"strconv"
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/protoc-gen-terraform/config"
	"github.com/gravitational/trace"
)

// FieldBuildContext is a facade helper struct which facilitates getting field information
type FieldBuildContext struct {
	m         *Message
	g         *generator.Generator
	d         *generator.Descriptor
	f         *FieldDescriptorProtoExt
	index     int
	typeName  string
	path      string
	rawGoType string
}

// NewFieldBuildContext creates FieldBuildContext
func NewFieldBuildContext(
	m *Message,
	g *generator.Generator,
	d *generator.Descriptor,
	f *FieldDescriptorProtoExt,
	i int,
) (*FieldBuildContext, error) {
	n := f.GetName()
	typeName := getMessageTypeName(d) + "." + n
	path := m.Path + "." + n

	t, _ := g.GoType(d, f.FieldDescriptorProto)
	if t == "" {
		return nil, trace.Errorf("invalid field go type %v", path)
	}

	c := &FieldBuildContext{
		m:         m,
		g:         g,
		d:         d,
		f:         f,
		index:     i,
		typeName:  typeName,
		path:      path,
		rawGoType: t,
	}

	return c, nil
}

// NewFieldBuildContextWithField creates FieldBuildContext from another context, overriding field
func NewFieldBuildContextWithField(c *FieldBuildContext, f *FieldDescriptorProtoExt, i int) (*FieldBuildContext, error) {
	return NewFieldBuildContext(c.m, c.g, c.d, f, i)
}

// IsExcluded returns true if field is added to config.ExcludeFields
func (c *FieldBuildContext) IsExcluded() bool {
	_, ok1 := config.ExcludeFields[c.GetNameWithTypeName()]
	_, ok2 := config.ExcludeFields[c.GetPath()]

	return ok1 || ok2
}

// GetNameWithTypeName returns field type name with package
func (c *FieldBuildContext) GetNameWithTypeName() string {
	return c.typeName
}

// GetName returns field name
func (c *FieldBuildContext) GetName() string {
	return c.f.GetName()
}

// GetPath returns field path
func (c *FieldBuildContext) GetPath() string {
	return c.path
}

// GetTypeAndIsMessage returns raw schema type, field value type and IsMessage flag for current field
func (c *FieldBuildContext) GetTypeAndIsMessage() (string, string, bool, error) {
	p := c.f.FieldDescriptorProto

	switch {
	case c.f.IsTime():
		return "string", "time.Time", false, nil
	case c.f.IsDuration():
		return "string", "time.Duration", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_DOUBLE) || gogoproto.IsStdDouble(p):
		return "float64", "float64", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_FLOAT) || gogoproto.IsStdFloat(p):
		return "float64", "float32", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_INT64) || gogoproto.IsStdInt64(p):
		return "int", "int64", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT64) || gogoproto.IsStdUInt64(p):
		return "int", "uint64", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_INT32) || gogoproto.IsStdInt32(p):
		return "int", "int32", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT32) || gogoproto.IsStdUInt32(p):
		return "int", "uint32", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED64):
		return "int", "uint64", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED32):
		return "int", "uint32", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_BOOL) || gogoproto.IsStdBool(p):
		return "bool", "bool", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_STRING) || gogoproto.IsStdString(p):
		return "string", "string", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_BYTES) || gogoproto.IsStdBytes(p):
		return "string", "[]byte", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED32):
		return "int", "int32", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED64):
		return "int", "int64", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT32):
		return "int", "int32", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT64):
		return "string", "int64", false, nil
	case c.f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_ENUM):
		return "string", "string", false, nil
	case c.f.IsMessage():
		return "", "", true, nil
	default:
		return "", "", false, trace.Errorf("unknown field type %v", c.GetPath())
	}
}

// IsTime returns true if field is time
func (c *FieldBuildContext) IsTime() bool {
	return c.f.IsTime()
}

// IsDuration returns true if field is duration
func (c *FieldBuildContext) IsDuration() bool {
	return c.f.IsDuration()
}

// IsMessage returns true if field is message
func (c *FieldBuildContext) IsMessage() bool {
	return c.f.IsMessage()
}

// IsCustomType returns true if fields has gogo.custom_type flag
func (c *FieldBuildContext) IsCustomType() bool {
	return c.f.IsCustomType()
}

// GetCustomType returns true if fields has gogo.custom_type flag
func (c *FieldBuildContext) GetCustomType() string {
	return c.f.GetCustomType()
}

// GetComment returns field comment as a single line and as a block comment
func (c *FieldBuildContext) GetComment() (string, string) {
	// ",2," marks that we are extracting comment for a message field. See descriptor.SourceCodeInfo source for details.
	p := c.d.Path() + ",2," + strconv.Itoa(c.index)

	for _, l := range c.d.File().GetSourceCodeInfo().GetLocation() {
		if getLocationPath(l) == p {
			c := strings.Trim(l.GetLeadingComments(), "\n")

			return commentToSingleLine(strings.TrimSpace(c)), appendSlashSlash(c, false)
		}
	}

	return "", ""
}

// GetMessageDescriptor returns underlying field message descriptor
func (c *FieldBuildContext) GetMessageDescriptor() (*generator.Descriptor, error) {
	// Resolve underlying message via protobuf
	x := c.g.ObjectNamed(c.f.GetTypeName())
	desc, ok := x.(*generator.Descriptor)
	if desc == nil || !ok {
		return nil, trace.Errorf("failed to convert %T to *generator.Descriptor", x)
	}

	return desc, nil
}

// IsRepeated returns true if field is repeated
func (c *FieldBuildContext) IsRepeated() bool {
	return !c.g.IsMap(c.f.FieldDescriptorProto) && c.f.IsRepeated()
}

// IsMap returns true if field is map
func (c *FieldBuildContext) IsMap() bool {
	return c.g.IsMap(c.f.FieldDescriptorProto)
}

// GetMapValueFieldDescriptor returns field descriptor for a map field
func (c *FieldBuildContext) GetMapValueFieldDescriptor() (*FieldDescriptorProtoExt, error) {
	m := c.g.GoMapType(nil, c.f.FieldDescriptorProto)

	k, _ := c.g.GoType(c.d, m.KeyField)
	if k != "string" {
		return nil, trace.Errorf("non-string map keys are not supported %v", c.GetPath())
	}

	if m.ValueField == nil {
		return nil, trace.Errorf("map value descriptor is nil %v", c.GetPath())
	}

	return &FieldDescriptorProtoExt{m.ValueField}, nil
}
