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
	"sort"
	"strings"

	"github.com/gravitational/protoc-gen-terraform/config"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/trace"
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

	// IsSensitive field has Sensitive flag
	IsSensitive bool

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

	// StateFunc is field state func name
	StateFunc string
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

	// Sort fields if required
	if config.Sort {
		sort.Slice(m.Fields, func(i, j int) bool {
			return m.Fields[i].NameSnake < m.Fields[j].NameSnake
		})
	}

	return nil
}

// BuildField builds Field structure
func BuildField(c *FieldBuildContext) (*Field, error) {
	var err error

	if c.IsExcluded() {
		return nil, nil
	}

	f := &Field{Name: c.GetName()}

	f.setNameSnake(c)

	f.RawComment, f.Comment = c.GetComment()
	f.SchemaRawType, f.SchemaGoType, f.IsMessage, err = c.GetTypeAndIsMessage()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	f.IsTime = c.IsTime()
	f.IsDuration = c.IsDuration()

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
	f.setSensitive(c)

	f.setDefault(c)

	err = f.setConfigMode(c)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	f.setStateFunc(c)
	f.setKind()

	return f, nil
}

// setNameSnake sets snake name for the field taking exceptions into place
func (f *Field) setNameSnake(c *FieldBuildContext) {
	v1, ok1 := config.FieldNameReplacements[c.GetPath()]
	v2, ok2 := config.FieldNameReplacements[c.GetNameWithTypeName()]

	if ok1 {
		f.NameSnake = v1
	} else if ok2 {
		f.NameSnake = v2
	} else {
		f.setNameSnakeWithJSONTag(c)
	}
}

// setNameSnakeWithJSONTag sets snake name for the field
func (f *Field) setNameSnakeWithJSONTag(c *FieldBuildContext) {
	if config.UseJSONTag {
		n := c.f.GetJSONName()
		if n != "" {
			f.NameSnake = n
			return
		}
	}

	f.NameSnake = c.GetSnakeName()
}

// setKind resolves and sets kind the field
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

	v, ok := config.Suffixes[c.GetCustomType()]

	if ok {
		f.Suffix = v
		return nil
	}

	f.Suffix = strings.ReplaceAll(strings.ReplaceAll(c.GetCustomType(), "/", ""), ".", "")

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

// setSensitive sets IsForceNew flag
func (f *Field) setSensitive(c *FieldBuildContext) {
	_, ok1 := config.SensitiveFields[c.GetNameWithTypeName()]
	_, ok2 := config.SensitiveFields[c.GetPath()]

	if ok1 || ok2 {
		f.IsSensitive = true
	}
}

// setDefault sets field default value
func (f *Field) setDefault(c *FieldBuildContext) {
	v1, ok1 := config.Defaults[c.GetPath()]
	v2, ok2 := config.Defaults[c.GetNameWithTypeName()]

	if ok1 {
		f.Default = fmt.Sprintf("%#v", v1)
	} else if ok2 {
		f.Default = fmt.Sprintf("%#v", v2)
	}
}

// setConfigMode sets field config mode
func (f *Field) setConfigMode(c *FieldBuildContext) error {
	_, a1 := config.ConfigModeAttrFields[c.GetNameWithTypeName()]
	_, a2 := config.ConfigModeAttrFields[c.GetPath()]

	_, b1 := config.ConfigModeBlockFields[c.GetNameWithTypeName()]
	_, b2 := config.ConfigModeBlockFields[c.GetPath()]

	if (a1 || a2) && (b1 || b2) {
		return trace.Errorf("field %v can not have SchemaConfigModeAttrs and SchemaConfigModeBlock "+
			"in the same time, check config_mode_attrs/config_mode_block configuration variables", c.GetPath())
	}

	if a1 || a2 {
		f.ConfigMode = "SchemaConfigModeAttr"
	}

	if b1 || b2 {
		f.ConfigMode = "SchemaConfigModeBlock"
	}

	return nil
}

// setStateFunc sets field state func
func (f *Field) setStateFunc(c *FieldBuildContext) {
	v1, ok1 := config.StateFunc[c.GetPath()]
	v2, ok2 := config.StateFunc[c.GetNameWithTypeName()]

	if ok1 {
		f.StateFunc = v1
	} else if ok2 {
		f.StateFunc = v2
	}
}
