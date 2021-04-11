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
	"fmt"
	"strings"

	"github.com/gravitational/protoc-gen-terraform/config"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/trace"
	"github.com/stoewer/go-strcase"
)

// Field represents metadata about protobuf message field descriptor.
// This struct and the following methods know about both schema details and target field details.
type Field struct {
	// Name type name
	Name string

	// NameSnake type name in snake case (Terraform schema field name)
	NameSnake string

	// SchemaRawType Terraform schema raw value go type (float64 for types.Float), used in schema generation
	SchemaRawType string

	// SchemaGoType Go type to convert schema raw type to (uint32, []byte, time.Time, time.Duration)
	SchemaGoType string

	// GoType Go type without [] and *, but with package name prepended
	GoType string

	// GoTypeIsSlice specifies whether Go type is a slice
	GoTypeIsSlice bool

	// GoTypeIsPtr specifies whether Go type is a pointer
	GoTypeIsPtr bool

	// GoTypeFull Go type with [] and * and package name prepended
	GoTypeFull string

	// Kind Field kind (resulting of combination of flags below, see setKind method)
	Kind string

	// IsRepeated field is list?
	IsRepeated bool

	// IsMap field is map?
	IsMap bool

	// IsMessage field is message? (might be map or list at the same time)
	IsMessage bool

	// IsRequired field is required
	IsRequired bool

	// IsComputed field is computed
	IsComputed bool

	// IsTime field contains time? (needs to be parsed from string)
	IsTime bool

	// IsDuration field contains duration? (needs to be parsed from string)
	IsDuration bool

	// IsCustomType field has gogo.customtype flag?
	IsCustomType bool

	// IsForceNew field has ForceNew flag
	IsForceNew bool

	// Suffix custom type schema and unmarshal method name suffix
	Suffix string

	// Message reference to nested message
	Message *Message

	// MapValueField reference to map value field reflection
	MapValueField *Field

	// RawComment is field comment in proto file without // prepended
	RawComment string

	// Comment is field comment in proto file with // prepended
	Comment string

	// Default is field default value
	Default string

	// ConfigMode is field config mode
	ConfigMode string
}

// BuildFields builds []*Field from descriptors of specified message
func BuildFields(m *Message, g *generator.Generator, d *generator.Descriptor) error {
	for i, fd := range d.GetField() {
		fd := &FieldDescriptorProtoExt{fd}

		c, err := NewFieldBuildContext(m, g, d, fd, i)
		if err != nil {
			return trace.Wrap(err)
		}

		f, err := BuildField(c)
		if err != nil {
			return trace.Wrap(err)
		}

		if f != nil {
			m.Fields = append(m.Fields, f)
		}
	}

	return nil
}

// BuildField builds Field structure
func BuildField(c *FieldBuildContext) (*Field, error) {
	var err error

	if c.IsExcluded() {
		return nil, nil
	}

	n := c.GetName()

	f := &Field{
		Name:      n,
		NameSnake: strcase.SnakeCase(n),
	}

	f.RawComment, f.Comment = c.GetComment()
	f.SchemaRawType, f.SchemaGoType, f.IsMessage, err = c.GetTypeAndIsMessage()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	f.IsTime = c.IsTime()
	f.IsDuration = c.IsDuration()

	// Byte slice is an exception, is should be treated as string
	if c.IsByteSlice() {
		f.GoType = c.GetBytesExceptionGoType()
		f.GoTypeFull = c.GetRawGoType()
	} else {
		f.GoTypeIsPtr = c.GetGoTypeIsPtr()
		f.GoTypeIsSlice = c.GetGoTypeIsSlice()

		if f.IsMessage {
			d, err := c.GetMessageDescriptor()
			if err != nil {
				return nil, trace.Wrap(err)
			}

			m, err := BuildMessage(c.g, d, false, c.path)
			if err != nil {
				return nil, trace.Wrap(err)
			}
			if m == nil {
				return nil, nil
			}

			f.Message = m
			f.GoType = c.GetGoType(m.GoTypeName)
		} else {
			f.GoType = c.GetGoType("")
		}

		f.setGoTypeFull()
	}

	f.IsRepeated = c.IsRepeated()

	if c.IsMap() {
		f.IsMap = true

		d, err := c.GetMapValueFieldDescriptor()
		if err != nil {
			return nil, trace.Wrap(err)
		}

		ctx, err := NewFieldBuildContextWithField(c, d, -1)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		vf, err := BuildField(ctx)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		f.MapValueField = vf
	}

	err = f.setCustomType(c)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	f.setRequired(c)
	f.setComputed(c)
	f.setForceNew(c)

	err = f.setDefault(c)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	err = f.setConfigMode(c)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	f.setKind()

	return f, nil
}

// setGoTypeFull constructs and sets full go type name for the field
func (f *Field) setGoTypeFull() {
	r := f.GoType

	if f.GoTypeIsPtr {
		r = "*" + r
	}
	if f.GoTypeIsSlice {
		r = "[]" + r
	}

	f.GoTypeFull = r
}

// setKind resolves and sets kind for current field
func (f *Field) setKind() {
	switch {
	case f.IsCustomType:
		f.Kind = "CUSTOM_TYPE"
	case f.IsMap && f.MapValueField.IsMessage:
		// Terraform does not support map of objects. We replace such field with list of objects having key and value fields.
		f.Kind = "MESSSAGE_MAP" // ex: map[string]struct, requires additional steps to unmarshal
	case f.IsMap:
		f.Kind = "MAP" // ex: map[string]string
	case f.IsRepeated && f.IsMessage:
		f.Kind = "REPEATED_MESSAGE" // ex: []struct
	case f.IsRepeated:
		f.Kind = "REPEATED_ELEMENTARY" // ex: []string
	case f.IsMessage:
		f.Kind = "SINGULAR_MESSAGE" // ex: struct
	default:
		f.Kind = "SINGULAR_ELEMENTARY" // ex: string
	}
}

// setCustomType sets custom type information
func (f *Field) setCustomType(c *FieldBuildContext) error {
	if !c.IsCustomType() {
		return nil
	}

	f.IsCustomType = true

	v, ok := config.Suffixes[c.f.GetCustomType()]

	if ok {
		f.Suffix = v
		return nil
	}

	f.Suffix = strings.ReplaceAll(strings.ReplaceAll(f.GoType, "/", ""), ".", "")

	return nil
}

// setRequired sets IsRequired flag
func (f *Field) setRequired(c *FieldBuildContext) {
	_, ok1 := config.RequiredFields[c.GetNameWithTypeName()]
	_, ok2 := config.RequiredFields[c.GetPath()]

	if ok1 || ok2 {
		f.IsRequired = true
	}
}

// setComputed sets IsComputed flag
func (f *Field) setComputed(c *FieldBuildContext) {
	_, ok1 := config.ComputedFields[c.GetNameWithTypeName()]
	_, ok2 := config.ComputedFields[c.GetPath()]

	if ok1 || ok2 {
		f.IsComputed = true
	}
}

// setForceNew sets IsForceNew flag
func (f *Field) setForceNew(c *FieldBuildContext) {
	_, ok1 := config.ForceNewFields[c.GetNameWithTypeName()]
	_, ok2 := config.ForceNewFields[c.GetPath()]

	if ok1 || ok2 {
		f.IsForceNew = true
	}
}

// setDefault returns field default value
func (f *Field) setDefault(c *FieldBuildContext) error {
	v1, ok1 := config.Defaults[c.GetPath()]
	v2, ok2 := config.Defaults[c.GetNameWithTypeName()]

	if ok1 && ok2 && c.GetPath() != c.GetNameWithTypeName() {
		return trace.Errorf("field has default value set by path " + c.GetPath() + " and by name " + c.GetNameWithTypeName())
	}

	if ok1 {
		f.Default = fmt.Sprintf("%#v", v1)
	}

	if ok2 {
		f.Default = fmt.Sprintf("%#v", v2)
	}

	return nil
}

// setConfigMode sets field config mode
func (f *Field) setConfigMode(c *FieldBuildContext) error {
	_, a1 := config.ConfigModeAttrFields[c.GetNameWithTypeName()]
	_, a2 := config.ConfigModeAttrFields[c.GetPath()]

	_, b1 := config.ConfigModeBlockFields[c.GetNameWithTypeName()]
	_, b2 := config.ConfigModeBlockFields[c.GetPath()]

	if (a1 || a2) && (b1 || b2) {
		return trace.Errorf("field " + c.GetPath() + " can not have SchemaConfigModeAttrs and SchemaConfigModeBlock " +
			"in the same time, check config_mode_attrs/config_mode_block configuration variables")
	}

	if a1 || a2 {
		f.ConfigMode = "SchemaConfigModeAttr"
	}

	if b1 || b2 {
		f.ConfigMode = "SchemaConfigModeBlock"
	}

	return nil
}
