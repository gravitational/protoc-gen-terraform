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
	"strconv"
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/trace"
	"github.com/stoewer/go-strcase"
)

const (
	// Int64Type represents SchemaType for an integer field
	Int64Type = "Int64Type"
	// Float64Type represents SchemaType for a float field
	Float64Type = "Float64Type"
	// StringType represents SchemaType for a string field
	StringType = "StringType"
	// BoolType represents SchemaType for a bool field
	BoolType = "BoolType"
	// TimeType represents SchemaType for a time field
	TimeType = "TimeType"
	// DurationType represents SchemaType for a duration field
	DurationType = "DurationType"
)

// FieldBuildContext is a facade helper struct which facilitates getting a field properties
type FieldBuildContext struct {
	MessageBuildContext
	field    *FieldDescriptorProtoExt
	index    int
	typeName string
	path     string
	goType   string
}

// NewFieldBuildContext creates FieldBuildContext
func NewFieldBuildContext(m MessageBuildContext, field *FieldDescriptorProtoExt, index int) (*FieldBuildContext, error) {
	typeName := m.GetName() + "." + field.GetName()
	path := m.GetPath() + "." + field.GetName()

	// Needs to be called explicitly because of gogo implementation details
	t, _ := m.gen.GoType(m.desc, field.FieldDescriptorProto)
	if t == "" {
		return nil, trace.Errorf("invalid field go type %v", path)
	}

	c := &FieldBuildContext{
		MessageBuildContext: m,
		field:               field,
		index:               index,
		typeName:            typeName,
		path:                path,
		goType:              m.imports.GoString(t, true),
	}

	return c, nil
}

// NewMapValueFieldBuildContext creates FieldBuildContext for MapValueField
func NewMapValueFieldBuildContext(c *FieldBuildContext, field *FieldDescriptorProtoExt, index int, typ string) (*FieldBuildContext, error) {
	// We've gen.GoType always returns *type here, have to override
	i := strings.LastIndex(typ, "]")
	t := typ[i+1:]

	return &FieldBuildContext{
		MessageBuildContext: c.MessageBuildContext,
		field:               field,
		index:               index,
		typeName:            c.typeName,
		path:                c.path,
		goType:              c.imports.GoString(t, true),
	}, nil
}

// GetGoType returns Go type for the field
func (c *FieldBuildContext) GetGoType() string {
	return c.goType
}

// IsExcluded returns true if field is added to config.ExcludeFields
func (c *FieldBuildContext) IsExcluded() bool {
	return c.GetFlagValue(c.config.ExcludeFields)
}

// GetNameWithTypeName returns field type name with package
func (c *FieldBuildContext) GetNameWithTypeName() string {
	return c.typeName
}

// GetName returns field name
func (c *FieldBuildContext) GetName() string {
	return c.field.GetName()
}

// GetPath returns a field path
func (c *FieldBuildContext) GetPath() string {
	return c.path
}

// GetTypeAndIsMessage returns schema type, field value type and IsMessage flag for current field.
//
// It returns:
// - Terraform schema type (with artificial exceptions for time and duration)
// - Go elem value type (including cast type if set)
// - Go signature for zero value (either "", 0.0 or types.Bool(false))
// - flag indicating that this field contains a nested message
func (c *FieldBuildContext) GetTypeAndIsMessage() (string, string, string, bool, error) {
	p := c.field.FieldDescriptorProto

	// gogo protobuf does not support nullable elementary fields
	elemType := strings.ReplaceAll(c.GetGoType(), "[]", "")

	// orCastType returns type passed as an argument or cast type if defined
	orCastType := func(v string) string {
		if c.IsCastType() {
			return elemType
		}

		return v
	}

	// zeroOrCast returns zero value passed as an arument or nil or zero value in cast type
	zeroOrCast := func(v string) string {
		if c.IsCastType() {
			return elemType + "(" + v + ")"
		}

		return v
	}

	switch {
	case c.field.IsTime():
		if c.config.TimeType == nil {
			return "", "", "", false, trace.Errorf("%v field has time type, but config.time_type is not defined", c.path)
		}
		return TimeType, orCastType("time.Time"), zeroOrCast("time.Time{}"), false, nil // In SDKV2 time and duration are represented as strings in RFC3339 format
	case c.field.IsDuration(c.config.DurationCustomType): // In Terraform Framework special type needs to be defined
		if c.config.DurationType == nil {
			return "", "", "", false, trace.Errorf("%v field has duration type, but config.duration_type is not defined", c.path)
		}
		return DurationType, orCastType("time.Duration"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_DOUBLE) || gogoproto.IsStdDouble(p):
		return Float64Type, orCastType("float64"), zeroOrCast("0.0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_FLOAT) || gogoproto.IsStdFloat(p):
		return Float64Type, orCastType("float32"), zeroOrCast("0.0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_INT64) || gogoproto.IsStdInt64(p):
		return Int64Type, orCastType("int64"), zeroOrCast("0"), false, nil // Terraform uses the single types.Int64Type to represent all integers
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT64) || gogoproto.IsStdUInt64(p):
		return Int64Type, orCastType("uint64"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_INT32) || gogoproto.IsStdInt32(p):
		return Int64Type, orCastType("int32"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT32) || gogoproto.IsStdUInt32(p):
		return Int64Type, orCastType("uint32"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED64):
		return Int64Type, orCastType("uint64"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED32):
		return Int64Type, orCastType("uint32"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED32):
		return Int64Type, orCastType("int32"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED64):
		return Int64Type, orCastType("int64"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT32):
		return Int64Type, orCastType("int32"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT64):
		return StringType, orCastType("int64"), zeroOrCast("0"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_BOOL) || gogoproto.IsStdBool(p):
		return BoolType, orCastType("bool"), zeroOrCast("false"), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_STRING) || gogoproto.IsStdString(p):
		return StringType, orCastType("string"), zeroOrCast(`""`), false, nil
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_BYTES) || gogoproto.IsStdBytes(p):
		return StringType, orCastType("[]byte"), zeroOrCast("nil"), false, nil // Byte arrays are represented as strings
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_ENUM):
		return StringType, orCastType("string"), zeroOrCast(`""`), false, nil
	case c.field.IsMessage():
		return "", elemType, "<undefined>", true, nil
	default:
		return "", "", "", false, trace.Errorf("unknown field type %v", c.GetPath())
	}
}

// IsTime returns true if field is time
func (c *FieldBuildContext) IsTime() bool {
	return c.field.IsTime()
}

// IsDuration returns true if field is duration
func (c *FieldBuildContext) IsDuration() bool {
	return c.field.IsDuration(c.config.DurationCustomType)
}

// IsMessage returns true if field is message
func (c *FieldBuildContext) IsMessage() bool {
	return c.field.IsMessage()
}

// IsCustomType returns true if fields has gogo.custom_type flag
func (c *FieldBuildContext) IsCustomType() bool {
	return c.field.IsCustomType()
}

// GetCustomType returns true if fields has gogo.custom_type flag
func (c *FieldBuildContext) GetCustomType() string {
	return c.field.GetCustomType()
}

// IsCastType returns true if fields has gogo.cast_type flag
func (c *FieldBuildContext) IsCastType() bool {
	return c.field.IsCastType()
}

// GetComment returns field comment as a single line and as a block comment
func (c *FieldBuildContext) GetComment() (string, string) {
	// ",2," marks that we are extracting comment for a message field. See descriptor.SourceCodeInfo source for details.
	p := c.desc.Path() + ",2," + strconv.Itoa(c.index)

	for _, l := range c.desc.File().GetSourceCodeInfo().GetLocation() {
		if c.GetLocationPath(l) == p {
			c := Comment(strings.TrimSpace(strings.Trim(l.GetLeadingComments(), "\n")))
			return c.ToSingleLine(), c.SlashSlash(false)
		}
	}

	return "", ""
}

// GetMessageDescriptor returns underlying field message descriptor
func (c *FieldBuildContext) GetMessageDescriptor() (*generator.Descriptor, error) {
	// Resolve underlying message via protobuf
	x := c.gen.ObjectNamed(c.field.GetTypeName())
	desc, ok := x.(*generator.Descriptor)
	if desc == nil || !ok {
		return nil, trace.Errorf("failed to convert %T to *generator.Descriptor", x)
	}

	return desc, nil
}

// IsRepeated returns true if field is repeated
func (c *FieldBuildContext) IsRepeated() bool {
	return !c.gen.IsMap(c.field.FieldDescriptorProto) && c.field.IsRepeated()
}

// IsMap returns true if field is map
func (c *FieldBuildContext) IsMap() bool {
	return c.gen.IsMap(c.field.FieldDescriptorProto)
}

// GetMapValueFieldDescriptor returns field descriptor for a map field
func (c *FieldBuildContext) GetMapValueFieldDescriptorAndType() (string, *FieldDescriptorProtoExt, error) {
	m := c.gen.GoMapType(nil, c.field.FieldDescriptorProto)

	k, _ := c.gen.GoType(c.desc, m.KeyField)
	if k != "string" {
		return "", nil, trace.Errorf("non-string map keys are not supported %v", c.GetPath())
	}

	if m.ValueField == nil {
		return "", nil, trace.Errorf("map value descriptor is nil %v", c.GetPath())
	}

	return m.GoType, &FieldDescriptorProtoExt{m.ValueField}, nil
}

// GetNameSnake returns field snake name
func (c *FieldBuildContext) GetNameSnake() string {
	v, ok := c.config.NameOverrides[c.GetPath()]
	if !ok {
		v, ok = c.config.NameOverrides[c.GetNameWithTypeName()]
	}

	if ok {
		return v
	}

	n := c.field.GetJSONName()
	if n != "" {
		return n
	}

	return strcase.SnakeCase(c.field.GetName())
}

// GetFlagValue returns a flag value for the field
func (c *FieldBuildContext) GetFlagValue(f map[string]struct{}) bool {
	_, ok1 := f[c.GetNameWithTypeName()]
	_, ok2 := f[c.GetPath()]

	return ok1 || ok2
}

// GetValidators returns field validators
func (c *FieldBuildContext) GetValidators() []string {
	v, ok := c.config.Validators[c.GetPath()]
	if !ok {
		v, ok = c.config.Validators[c.GetNameWithTypeName()]
	}

	if ok {
		return v
	}

	return []string{}
}

// GetNullable returns the nullable flag
func (c *FieldBuildContext) GetNullable() bool {
	// E_NULLABLE option is not applicable here because by default all fields must be nullable by protobuf specs.
	// We need to examine a target Go type for this specific field.
	return strings.Contains(c.GetGoType(), "*")
}

// GetTFSchemaTypes returns TFSchemaType overrides
func (c *FieldBuildContext) GetTFSchemaTypes() (string, string, string) {
	v, ok := c.config.SchemaTypes[c.GetPath()]
	if !ok {
		v, ok = c.config.SchemaTypes[c.GetNameWithTypeName()]
	}

	switch {
	case ok:
		return v.Type, v.ValueType, v.CastType
	case c.IsTime() && c.config.TimeType != nil:
		v := c.config.TimeType
		return v.Type, v.ValueType, v.CastType
	case c.IsDuration() && c.config.DurationType != nil:
		v := c.config.DurationType
		return v.Type, v.ValueType, v.CastType
	}

	return "", "", ""
}
