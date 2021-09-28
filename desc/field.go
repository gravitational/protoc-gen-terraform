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
	"sort"
	"strings"

	"github.com/gravitational/trace"
)

// Kind type of a field
type Kind uint8

const (
	// Primitive represents kind of a field containing an elementary value (string, int, etc)
	Primitive Kind = 1
	// PrimitiveList represents kind of a field containing a list of elementary values ([]string, []int, etc)
	PrimitiveList Kind = 2
	// Nested represents kind of a field containing a nested object
	Nested Kind = 3
	// NestedList represents kind of a field containing a list of a nested object
	NestedList Kind = 4
	// PrimitiveMap represents kind of a field containing a map of elementary values (map[string]string)
	PrimitiveMap Kind = 5
	// NestedMap represents kind of a field containing a map of a nested objects (map[string]T)
	NestedMap Kind = 6
	// Custom represents kind of a field containing a custom type
	Custom Kind = 7
)

// Field represents metadata of protobuf message field descriptor
type Field struct {
	// Name field name
	Name string
	// NameSnake represents Terraform schema field name. It is taken from json_tag or generated or explicitly specified using NameOverrides
	NameSnake string
	// Kind represents field kind: resulting combination of the flags below. Refer to setKind method
	Kind Kind
	// TFSchemaType represents Terraform schema type name (Int64Type)
	TFSchemaType string
	// SchemaValueType represents Terraform schema value type name (Int64)
	TFSchemaValueType string
	// TFSchemaValueCastType represents the Go type of the .Value member
	TFSchemaValueCastType string

	// GoType represents raw go type of a target protobuf field (with [], *)
	GoType string
	// GoElemType represents raw go type of an element of a target protobuf field (no [] for maps/arrays, no * for pointers)
	GoElemType string
	// GoElemTypeZeroValue returns Go zero value for a type
	GoElemTypeZeroValue string
	// GoCustomType represents a raw Go type for custom type field
	GoCustomType string
	// Suffix represents a custom type suffix used to refer to custom methods (GenSchema<Suffix>)
	Suffix string

	// IsRepeated field is a list?
	IsRepeated bool
	// IsMap field is a map?
	IsMap bool
	// IsMessage field is a nested message? (might be map or list at the same time)
	IsMessage bool

	// IsRequired field is required?
	IsRequired bool
	// IsComputed field is computed?
	IsComputed bool
	// IsCustomType field has gogo.customtype flag?
	IsCustomType bool
	// IsNullable represents field nullable state
	IsNullable bool
	// IsSensitive is field sensitive? (password, token) (TFSDK only)
	IsSensitive bool
	// Validators represents an array of the field validators (TFSDK only)
	Validators []string

	// Message represents a nested message
	Message *Message
	// MapValueField represents a map value
	MapValueField *Field

	// RawComment is field comment in proto file without // prepended
	RawComment string
	// Comment is field comment in proto file with // prepended
	Comment string

	// Path represents the path to the current field
	Path string
}

// BuildFields builds []*Field from a descriptors of the specified message
func BuildFields(m MessageBuildContext) ([]*Field, error) {
	fields := make([]*Field, 0)

	for i, field := range m.desc.GetField() {
		fieldExt := &FieldDescriptorProtoExt{field}

		c, err := NewFieldBuildContext(m, fieldExt, i)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		f, err := BuildField(c)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		if f != nil {
			fields = append(fields, f)
		}
	}

	// Sort fields if required
	if m.config.Sort {
		sort.Slice(fields, func(i, j int) bool {
			return fields[i].Name < fields[j].Name
		})
	}

	return fields, nil
}

// BuildField builds Field structure
func BuildField(c *FieldBuildContext) (*Field, error) {
	var err error

	if c.IsExcluded() {
		return nil, nil
	}

	f := &Field{
		Name:        c.GetName(),
		NameSnake:   c.GetNameSnake(),
		IsRequired:  c.GetFlagValue(c.config.RequiredFields),
		IsComputed:  c.GetFlagValue(c.config.ComputedFields),
		IsSensitive: c.GetFlagValue(c.config.SensitiveFields),
		IsRepeated:  c.IsRepeated(),
		IsMap:       c.IsMap(),
		IsNullable:  c.GetNullable(),
		GoType:      c.GetGoType(),
		Validators:  c.GetValidators(),
		Path:        c.GetPath(),
	}

	f.RawComment, f.Comment = c.GetComment()
	f.TFSchemaType, f.GoElemType, f.GoElemTypeZeroValue, f.IsMessage, err = c.GetTypeAndIsMessage()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	f.TFSchemaValueCastType = f.getSchemaCastType(f.TFSchemaType)

	if f.IsMessage && !c.IsMap() {
		f.Message, err = f.getMessage(c)
		if err != nil {
			return nil, trace.Wrap(err)
		}
		if f.Message == nil {
			return nil, nil
		}
		c.plugin.RegisterMessage(f.Message)
	}

	if c.IsMap() {
		var typ string

		// gogoprotobuf returns incorrect elem type for maps. It always contains "*", we have to override.
		typ, f.MapValueField, err = f.getMapValueField(c)
		if err != nil {
			return nil, trace.Wrap(err)
		}

		f.GoType = typ
		f.IsNullable = strings.Contains(typ, "*")
		f.GoElemType = f.MapValueField.GoElemType
		f.GoElemTypeZeroValue = f.MapValueField.GoElemTypeZeroValue
	}

	t, v, a := c.GetTFSchemaTypes()
	if t != "" && v != "" && a != "" {
		f.TFSchemaType = c.imports.GoString(t, false)
		f.TFSchemaValueType = c.imports.GoString(v, false)
		f.TFSchemaValueCastType = c.imports.GoString(a, false)
	}

	f.setCustomType(c)

	f.Kind = f.getKind()

	return f, nil
}

// getMessage returns a nested message definition
func (f *Field) getMessage(c *FieldBuildContext) (*Message, error) {
	d, err := c.GetMessageDescriptor()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	m, err := BuildMessage(c.plugin, d, false, c.path)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	if m == nil {
		return nil, nil
	}

	return m, nil
}

// getMapValueField returns map value field for a field
func (f *Field) getMapValueField(c *FieldBuildContext) (string, *Field, error) {
	// For some reason, gogoprotobuf incorrectly treats nullable status of when map value is a message.
	// We have to override it.
	typ, d, err := c.GetMapValueFieldDescriptorAndType()
	if err != nil {
		return "", nil, trace.Wrap(err)
	}

	ctx, err := NewMapValueFieldBuildContext(c, d, -1, typ)
	if err != nil {
		return "", nil, trace.Wrap(err)
	}

	vf, err := BuildField(ctx)
	if err != nil {
		return "", nil, trace.Wrap(err)
	}

	return typ, vf, nil
}

// getKind resolves and sets kind the field
func (f *Field) getKind() Kind {
	switch {
	case f.IsCustomType:
		return Custom
	case f.IsMap && f.MapValueField.IsMessage:
		// Terraform SDKV2 does not support map of objects.
		// We replace such field with list of objects having key and value fields.
		return NestedMap // ex: map[string]struct
	case f.IsMap:
		return PrimitiveMap // ex: map[string]string
	case f.IsRepeated && f.IsMessage:
		return NestedList // ex: []struct
	case f.IsRepeated:
		return PrimitiveList // ex: []string
	case f.IsMessage:
		return Nested // ex: struct
	}
	return Primitive // ex: string
}

// setCustomType sets IsCustomType, GoCustomType and Suffix.
// Please note that CustomType overrides the whole field type.
// Repeated customtype and map customtype would be the same type.
func (f *Field) setCustomType(c *FieldBuildContext) {
	if !c.IsCustomType() {
		return
	}

	f.IsCustomType = true
	f.GoCustomType = c.GetCustomType()

	v, ok := c.config.Suffixes[c.GetCustomType()]
	if ok {
		f.Suffix = v
		return
	}

	f.Suffix = strings.ReplaceAll(strings.ReplaceAll(c.GetCustomType(), "/", ""), ".", "")
}

// Returns .Value go type from a type name
func (f *Field) getSchemaCastType(t string) string {
	switch t {
	case Int64Type:
		return "int64"
	case Float64Type:
		return "float64"
	case StringType:
		return "string"
	case BoolType:
		return "bool"
	case TimeType:
		return "time.Time"
	case DurationType:
		return "time.Duration"
	}

	return ""
}
