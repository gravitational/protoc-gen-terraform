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
	"sort"
	"strings"

	"github.com/gravitational/trace"
)

// Kind type of a field
type Kind uint8

const (
	// PrimitiveKind represents kind of a field containing an elementary value (string, int, etc)
	PrimitiveKind Kind = 1
	// PrimitiveListKind represents kind of a field containing a list of elementary values ([]string, []int, etc)
	PrimitiveListKind Kind = 2
	// ObjectKind represents kind of a field containing a nested object
	ObjectKind Kind = 3
	// ObjectListKind represents kind of a field containing a list of a nested object
	ObjectListKind Kind = 4
	// PrimitiveMapKind represents kind of a field containing a map of elementary values (map[string]string)
	PrimitiveMapKind Kind = 5
	// ObjectMapKind represents kind of a field containing a map of a nested objects (map[string]T)
	ObjectMapKind Kind = 6
	// CustomKind represents kind of a field containing a custom type
	CustomKind Kind = 7
)

// TerraformType represents Terraform type information
type TerraformType struct {
	// Type represents Terraform attr.Type name
	Type string
	// ValueType represents Terraform attr.Value name
	ValueType string
	// ElemType represents Terraform attr.Type name for list/map Elem, equals Type by default
	ElemType string
	// ElemValueType represents Terraform attr.Value name for list/map Elem, equals ValueType by default
	ElemValueType string
	// IsTypeScalar is true when Type is not a real struct and is represented with numeric constant on Terraform side
	IsTypeScalar bool
	// IsTypeScalar is true when ElemType is not a real struct and is represented with numeric constant on Terraform side
	IsElemTypeScalar bool
	// ValueCastToType represents Go type of either ValueType.Value or ElemValueType.Value
	ValueCastToType string
	// ValueCastFromType represents Go type of a counterpart object field or field elem to cast from .Value
	ValueCastFromType string
	// ZeroValue represents zero value for a type
	ZeroValue string
	// IsMessage field is a nested message? (might be map or list at the same time)
	IsMessage bool
	// TypeConstructor represents expression which is used to initialize type in schema
	TypeConstructor string
}

// ProtobufType represents protobuf object field type information
type ProtobufType struct {
	// GoType represents raw go type of a source protobuf object field (builtin, struct, slice, map, pointer)
	GoType string
	// GoElemType represents raw go type of a slice/map element (with possible *), otherwise equals GoType
	GoElemType string
	// GoElemTypeIndirect string represents raw go type slice/map element without *, otherwise equals GoElemType
	GoElemTypeIndirect string
	// OneOfType represents go type for OneOf type wrapper
	OneOfType string
	// OneOfName represents OneOf field name within the parent struct
	OneOfName string
	// IsPlaceholder represents flag, which indicates that this field is used as a placeholder for a message with no fields
	IsPlaceholder bool
}

// Field represents metadata of protobuf message field descriptor
type Field struct {
	// Name field name
	Name string
	// NameSnake represents Terraform schema field name. It is taken from json_tag or generated or explicitly specified using NameOverrides
	NameSnake string
	// Kind represents field kind: resulting combination of the flags below. Refer to setKind method
	Kind Kind

	TerraformType
	ProtobufType

	// Suffix represents a custom type suffix used to refer to custom methods (GenSchema<Suffix>)
	Suffix string

	// IsRepeated field is a list?
	IsRepeated bool
	// IsMap field is a map?
	IsMap bool

	// IsRequired field is required?
	IsRequired bool
	// IsComputed field is computed?
	IsComputed bool
	// IsCustomType field has gogo.customtype flag?
	IsCustomType bool
	// IsNullable represents field nullable state
	IsNullable bool
	// IsSensitive is field sensitive? (password, token)
	IsSensitive bool
	// IsEmbed is field embedded?
	IsEmbed bool
	// Validators represents the array of field validators for a field
	Validators []string
	// PlanModifiers represents the array of plan modifiers for a field
	PlanModifiers []string

	// Message represents a nested message
	Message *Message
	// MapValueField represents a Field of map value
	MapValueField *Field

	// Comment is field comment in proto file with // prepended
	Comment string

	// Path represents the path to the current field in proto message (Metadata.Name)
	Path string
}

// BuildFields builds []*Field from a descriptors of the specified message
func BuildFields(m MessageBuildContext) ([]*Field, error) {
	messageFields := m.desc.GetField()

	fields := make([]*Field, 0, len(messageFields))

	// Inject artificial field when message has no fields
	if len(messageFields) == 0 {
		fields = append(fields, BuildPlaceholderField(m.GetPath()))
		return fields, nil
	}

	for i, field := range messageFields {
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
		Name:          c.GetName(),
		NameSnake:     c.GetNameSnake(),
		IsRequired:    c.GetFlagValue(c.config.RequiredFields),
		IsComputed:    c.IsComputed(),
		IsSensitive:   c.GetFlagValue(c.config.SensitiveFields),
		IsRepeated:    c.IsRepeated(),
		IsMap:         c.IsMap(),
		IsNullable:    c.GetNullable(),
		IsEmbed:       c.IsEmbed(),
		Validators:    c.GetValidators(),
		PlanModifiers: c.GetPlanModifiers(),
		Path:          c.GetPath(),
		Comment:       c.GetComment(),
	}

	f.GoType = c.GetGoType()
	f.GoElemType = f.GoType

	f.TerraformType, err = c.GetTerraformType()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	if f.IsMessage && !c.IsMap() {
		err = f.setMessage(c)
		if err != nil {
			return nil, trace.Wrap(err)
		}
	}

	if c.IsRepeated() {
		f.setRepeatedGoElemType()
	}

	if c.IsMap() {
		err = f.setMapValues(c)
		if err != nil {
			return nil, trace.Wrap(err)
		}
	}

	f.setTerraformTypeOverride(c)
	f.setCustomType(c)

	f.GoElemTypeIndirect = strings.Replace(f.GoElemType, "*", "", -1)

	f.Kind = f.getKind()

	if c.IsOneOf() {
		f.OneOfName = c.GetOneOfFieldName()
		f.OneOfType = c.GetOneOfTypeName()
	}

	return f, nil
}

// BuildPlaceholderField represents no-field-message single field placeholder
func BuildPlaceholderField(basePath string) *Field {

	return &Field{
		Kind:       PrimitiveKind,
		Name:       "active",
		NameSnake:  "active",
		IsComputed: true,
		Comment:    "Automatically generated field preventing empty message errors",

		ProtobufType: ProtobufType{
			GoType:             "bool",
			GoElemType:         "bool",
			GoElemTypeIndirect: "bool",
			IsPlaceholder:      true,
		},

		TerraformType: boolType,
		Path:          basePath + ".active",
	}

}

// setRepeatedGoElemType fixes GoElemType for current field if it's repeated
func (f *Field) setRepeatedGoElemType() {
	f.GoElemType = f.GoType[strings.Index(f.GoType, "]")+1:]
}

// setMapValues sets required attributes for map in the current field
func (f *Field) setMapValues(c *FieldBuildContext) error {
	var typ string
	var err error

	// gogoprotobuf returns incorrect elem type for maps. It always contains "*", we have to override.
	typ, f.MapValueField, err = f.getMapValueField(c)
	if err != nil {
		return trace.Wrap(err)
	}

	// Otherwise, that would contain artificial protobuf Map_Entry type information
	f.GoType = typ
	f.IsNullable = strings.Contains(typ, "*")

	f.ElemType = f.MapValueField.ElemType
	f.ElemValueType = f.MapValueField.ElemValueType
	f.ValueCastToType = f.MapValueField.ValueCastToType
	f.ValueCastFromType = f.MapValueField.ValueCastFromType

	f.GoElemType = f.MapValueField.GoElemType

	return nil
}

// setMessage sets nested message for current field
func (f *Field) setMessage(c *FieldBuildContext) error {
	var err error

	f.Message, err = f.getMessage(c)
	if err != nil {
		return trace.Wrap(err)
	}
	if f.Message == nil {
		return nil
	}
	c.plugin.RegisterMessage(f.Message)

	return nil
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
	// For some reason, gogoprotobuf incorrectly treats nullable status when map value is a message.
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
		return CustomKind
	case f.IsMap && f.MapValueField.IsMessage:
		return ObjectMapKind // ex: map[string]struct
	case f.IsMap:
		return PrimitiveMapKind // ex: map[string]string
	case f.IsRepeated && f.IsMessage:
		return ObjectListKind // ex: []struct
	case f.IsRepeated:
		return PrimitiveListKind // ex: []string
	case f.IsMessage:
		return ObjectKind // ex: struct
	}
	return PrimitiveKind // ex: string
}

// setSchemaCustomType sets schema type override
func (f *Field) setTerraformTypeOverride(c *FieldBuildContext) {
	o := c.GetTerraformTypeOverride()
	if o != nil {
		f.Type = o.Type
		f.ValueType = o.ValueType
		f.ValueCastToType = o.CastToType
		f.ValueCastFromType = o.CastFromType
		f.TypeConstructor = o.TypeConstructor

		f.ElemType = f.Type
		f.ElemValueType = f.ValueType
		f.ValueCastFromType = f.ValueCastToType
	}
}

// setCustomType sets IsCustomType, GoCustomType and Suffix.
// Please note that CustomType overrides the whole field type.
// Repeated customtype and map customtype would be the same type.
func (f *Field) setCustomType(c *FieldBuildContext) {
	if !c.IsCustomType() {
		return
	}

	f.IsCustomType = true

	v, ok := c.config.Suffixes[c.GetCustomType()]
	if ok {
		f.Suffix = v
		return
	}

	// Default suffix: package and type name without / and .
	f.Suffix = strings.ReplaceAll(strings.ReplaceAll(c.GetCustomType(), "/", ""), ".", "")
}
