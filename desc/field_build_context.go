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
	// Types represents the path to Terraform Framework types package
	Types = "github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	float64Type = TerraformType{
		Type:             Types + ".Float64Type",
		ValueType:        Types + ".Float64",
		ElemType:         Types + ".Float64Type",
		ElemValueType:    Types + ".Float64",
		IsTypeScalar:     true,
		IsElemTypeScalar: true,
	}

	int64Type = TerraformType{
		Type:             Types + ".Int64Type",
		ValueType:        Types + ".Int64",
		ElemType:         Types + ".Int64Type",
		ElemValueType:    Types + ".Int64",
		IsTypeScalar:     true,
		IsElemTypeScalar: true,
	}

	stringType = TerraformType{
		Type:             Types + ".StringType",
		ValueType:        Types + ".String",
		ElemType:         Types + ".StringType",
		ElemValueType:    Types + ".String",
		IsTypeScalar:     true,
		IsElemTypeScalar: true,
	}

	boolType = TerraformType{
		Type:             Types + ".BoolType",
		ValueType:        Types + ".Bool",
		ElemType:         Types + ".BoolType",
		ElemValueType:    Types + ".Bool",
		IsTypeScalar:     true,
		IsElemTypeScalar: true,
	}

	objectType = TerraformType{
		Type:          Types + ".ObjectType",
		ValueType:     Types + ".Object",
		ElemType:      Types + ".ObjectType",
		ElemValueType: Types + ".Object",
	}
)

// FieldBuildContext is a facade helper struct which facilitates getting field properties
type FieldBuildContext struct {
	MessageBuildContext
	field    *FieldDescriptorProtoExt
	index    int
	typeName string
	imports  *Imports
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
		imports:             m.imports,
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
		imports:             c.imports,
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

// GetTerraformType returns terraform type meta information
func (c *FieldBuildContext) GetTerraformType() (TerraformType, error) {
	t := TerraformType{}

	p := c.field.FieldDescriptorProto

	// gogo protobuf does not support nullable elementary fields
	elemType := strings.ReplaceAll(c.GetGoType(), "[]", "")

	switch {
	case c.field.IsTime():
		if c.config.TimeType == nil {
			return t, trace.Errorf("%v field has time type, but config.time_type is not defined", c.path)
		}
		t = TerraformType{
			Type:          c.config.TimeType.Type,
			ValueType:     c.config.TimeType.ValueType,
			ElemType:      c.config.TimeType.Type,
			ElemValueType: c.config.TimeType.ValueType,
			CastType:      c.config.TimeType.CastType,
		}
	case c.field.IsDuration(c.config.DurationCustomType): // In Terraform Framework special type needs to be defined
		if c.config.DurationType == nil {
			return t, trace.Errorf("%v field has duration type, but config.duration_type is not defined", c.path)
		}
		t = TerraformType{
			Type:          c.config.DurationType.Type,
			ValueType:     c.config.DurationType.ValueType,
			ElemType:      c.config.DurationType.Type,
			ElemValueType: c.config.DurationType.ValueType,
			CastType:      c.config.DurationType.CastType,
		}
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_DOUBLE) || gogoproto.IsStdDouble(p):
		t = float64Type
		t.CastType = "float64"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_FLOAT) || gogoproto.IsStdFloat(p):
		t = float64Type
		t.CastType = "float32"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_INT64) || gogoproto.IsStdInt64(p):
		t = int64Type
		t.CastType = "int64"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT64) || gogoproto.IsStdUInt64(p):
		t = int64Type
		t.CastType = "uint64"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_INT32) || gogoproto.IsStdInt32(p):
		t = int64Type
		t.CastType = "int32"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT32) || gogoproto.IsStdUInt32(p):
		t = int64Type
		t.CastType = "uint32"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED64):
		t = int64Type
		t.CastType = "uint64"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED32):
		t = int64Type
		t.CastType = "uint32"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED32):
		t = int64Type
		t.CastType = "int32"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED64):
		t = int64Type
		t.CastType = "int64"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT32):
		t = int64Type
		t.CastType = "int32"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT64):
		t = int64Type
		t.CastType = "int64"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_BOOL) || gogoproto.IsStdBool(p):
		t = boolType
		t.CastType = "bool"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_STRING) || gogoproto.IsStdString(p):
		t = stringType
		t.CastType = "string"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_BYTES) || gogoproto.IsStdBytes(p):
		t = stringType
		t.CastType = "[]byte"
	case c.field.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_ENUM):
		t = stringType
		t.CastType = "string"
	case c.field.IsMessage():
		t = objectType
		t.IsMessage = true
		t.CastType = elemType
	default:
		return t, trace.Errorf("unknown field type %v", c.GetPath())
	}

	if c.IsRepeated() {
		t.Type = Types + ".ListType"
		t.ValueType = Types + ".List"
	}

	if c.IsMap() {
		t.Type = Types + ".MapType"
		t.ValueType = Types + ".Map"
	}

	t.Type = c.imports.GoString(t.Type, false)
	t.ValueType = c.imports.GoString(t.ValueType, false)
	t.ElemType = c.imports.GoString(t.ElemType, false)
	t.ElemValueType = c.imports.GoString(t.ElemValueType, false)

	if c.IsCastType() {
		t.CastType = elemType
	}

	return t, nil
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

// GetMapValueFieldDescriptorAndType returns field descriptor for a map field
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

// GetPlanModifiers returns field validators
func (c *FieldBuildContext) GetPlanModifiers() []string {
	v, ok := c.config.PlanModifiers[c.GetPath()]
	if !ok {
		v, ok = c.config.PlanModifiers[c.GetNameWithTypeName()]
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
