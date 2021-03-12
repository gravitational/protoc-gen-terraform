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
	"fmt"
	"strings"

	"github.com/gravitational/protoc-gen-terraform/config"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/stew/slice"
)

// Field represents field reflection struct
// This struct and the following methods know about both schema details and target field details
type Field struct {
	// Name type name
	Name string

	// NameSnake type name in snake case (Terraform schema field name)
	NameSnake string

	// SchemaRawType Terraform schema raw value go type (float64 for types.Float), used in schema generation
	SchemaRawType string

	// SchemaGoType Go type to convert schema raw type to (uint32, []byte, time.Time, time.Duration)
	SchemaGoType string

	// RawGoType field type returned by gogoprotobuf with * and []
	RawGoType string

	// GoType Go type without [] and *, but with package name prepended
	GoType string

	// GoTypeIsSlice Go type is a slice?
	GoTypeIsSlice bool

	// GoTypeIsPtr Go type is a pointer?
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

	// IsRequired field is required? TODO: implement via params?
	IsRequired bool

	// IsTime field contains time? (needs to be parsed from string)
	IsTime bool

	// IsDuration field contains duration? (needs to be parsed from string)
	IsDuration bool

	// IsCustomType field has gogo.customtype flag?
	IsCustomType bool

	// CustomTypeMethodInfix custom type schema and unmarshal method name infix
	CustomTypeMethodInfix string

	// Message reference to nested message
	Message *Message

	// MapValueField reference to map value field reflection
	MapValueField *Field
}

// fieldBuilder is axilarry struct responsible for building Field
type fieldBuilder struct {
	generator       *generator.Generator
	descriptor      *generator.Descriptor
	fieldDescriptor *descriptor.FieldDescriptorProto
	field           *Field
}

// BuildFields builds []*Field from descriptors of specified message
func BuildFields(m *Message, g *generator.Generator, d *generator.Descriptor) {
	for _, f := range d.GetField() {
		typeName := getFieldTypeName(d, f)

		// Ignore field if it is listed in cli arg
		if slice.Contains(config.ExcludeFields, typeName) {
			continue
		}

		f := BuildField(g, d, f)
		if f != nil {
			m.Fields = append(m.Fields, f)
		}
	}
}

// BuildField builds field reflection structure, or returns nil in case field build failed
func BuildField(g *generator.Generator, d *generator.Descriptor, f *descriptor.FieldDescriptorProto) *Field {
	b := newFieldBuilder(g, d, f)
	err := b.build()

	if err != nil {
		// Display error in the log, proceed further
		logrus.Printf("%+v", err)
		return nil
	}
	return b.field
}

// getFieldTypeName returns field name with package
func getFieldTypeName(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) string {
	return getMessageTypeName(d) + "." + f.GetName()
}

// newFieldBuilder constructs an empty fieldBuilder struct
func newFieldBuilder(g *generator.Generator, d *generator.Descriptor, f *descriptor.FieldDescriptorProto) *fieldBuilder {
	return &fieldBuilder{
		generator:       g,
		descriptor:      d,
		fieldDescriptor: f,
		field:           &Field{},
	}
}

// build fills in a Field structure
func (b *fieldBuilder) build() error {
	b.setName()

	err := b.resolveType()
	if err != nil {
		return err
	}

	b.setGoType()

	err = b.setAggregate()
	if err != nil {
		return err
	}

	b.setCustomType()
	b.setKind()

	return nil
}

// setName sets the field name
func (b *fieldBuilder) setName() {
	name := b.fieldDescriptor.GetName()

	b.field.Name = name
	b.field.NameSnake = strcase.SnakeCase(name)
}

// isTypeEq returns true if type
func (b *fieldBuilder) isTypeEq(t descriptor.FieldDescriptorProto_Type) bool {
	return *b.fieldDescriptor.Type == t
}

// isTime returns true if field stores time (is standard, google or golang time)
func (b *fieldBuilder) isTime() bool {
	t := b.fieldDescriptor.TypeName

	isStdTime := gogoproto.IsStdTime(b.fieldDescriptor)
	isGoogleTime := (t != nil && strings.HasSuffix(*t, "google.protobuf.Timestamp"))
	isCastToTime := gogoproto.GetCastType(b.fieldDescriptor) == "time.Time"

	return isStdTime || isGoogleTime || isCastToTime
}

// isDuration return true if field stores duration (is standard duration, or casted to duration)
func (b *fieldBuilder) isDuration() bool {
	ct := gogoproto.GetCastType(b.fieldDescriptor)
	t := b.fieldDescriptor.TypeName

	isStdDuration := gogoproto.IsStdDuration(b.fieldDescriptor)
	isGoogleDuration := (t != nil && strings.HasSuffix(*t, "google.protobuf.Duration"))
	isCastToCustomDuration := ct == config.DurationCustomType
	isCastToDuration := ct == "time.Duration"

	return isStdDuration || isGoogleDuration || isCastToDuration || isCastToCustomDuration
}

// isMessage returns true if field is message
func (b *fieldBuilder) isMessage() bool {
	return b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_MESSAGE)
}

// setTypes sets SchemaType and SchemaGoType
func (b *fieldBuilder) setSchemaTypes(schemaRawType string, goTypeCast string) {
	b.field.SchemaRawType = schemaRawType
	b.field.SchemaGoType = goTypeCast
}

// resolveType analyses field type and sets required fields in Field structure
// This method is pretty much copy & paste from gogo/protobuf generator.GoType
func (b *fieldBuilder) resolveType() error {
	d := b.fieldDescriptor // shortcut
	f := b.field           // shortcut

	switch {
	case b.isTime():
		b.setSchemaTypes("string", "time.Time")
		f.IsTime = true
	case b.isDuration():
		b.setSchemaTypes("string", "time.Duration")
		f.IsDuration = true
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_DOUBLE) || gogoproto.IsStdDouble(d):
		b.setSchemaTypes("float64", "float64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_FLOAT) || gogoproto.IsStdFloat(d):
		b.setSchemaTypes("float64", "float32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_INT64) || gogoproto.IsStdInt64(d):
		b.setSchemaTypes("int", "int64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT64) || gogoproto.IsStdUInt64(d):
		b.setSchemaTypes("int", "uint64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_INT32) || gogoproto.IsStdInt32(d):
		b.setSchemaTypes("int", "int32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT32) || gogoproto.IsStdUInt32(d):
		b.setSchemaTypes("int", "uint32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED64):
		b.setSchemaTypes("int", "uint64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED32):
		b.setSchemaTypes("int", "uint32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_BOOL) || gogoproto.IsStdBool(d):
		b.setSchemaTypes("bool", "bool")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_STRING) || gogoproto.IsStdString(d):
		b.setSchemaTypes("string", "string")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_BYTES) || gogoproto.IsStdBytes(d):
		b.setSchemaTypes("string", "[]byte")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED32):
		b.setSchemaTypes("int", "int32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED64):
		b.setSchemaTypes("int", "int64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT32):
		b.setSchemaTypes("int", "int32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT64):
		b.setSchemaTypes("string", "int64")
	case b.isMessage():
		err := b.setMessage()
		if err != nil {
			return err
		}
		f.IsMessage = true
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_ENUM):
		b.setSchemaTypes("string", "string")
	default:
		b.generator.Fail("unknown type for", b.descriptor.GetName(), b.fieldDescriptor.GetName())
		return nil
	}

	return nil
}

// Append type suffix to cast type, custom type or message
func (b *fieldBuilder) prependPackageName() {
	d := b.fieldDescriptor
	f := b.field

	// Override field type
	if gogoproto.IsCastType(d) {
		f.GoType = gogoproto.GetCastType(d)
	} else if gogoproto.IsCustomType(d) {
		f.GoType = gogoproto.GetCustomType(d)
	}

	// Prepend package name to overridden field type
	if gogoproto.IsCastType(d) || gogoproto.IsCustomType(d) {
		// If cast type is within current package, append default package name to it
		if !strings.Contains(f.GoType, ".") && config.DefaultPkgName != "" {
			f.GoType = config.DefaultPkgName + "." + f.GoType
		}
	} else {
		// Get go type from a message
		if b.isMessage() && b.field.Message != nil {
			f.GoType = b.field.Message.GoTypeName
		}
	}
}

// setGoType Sets go type with gogoprotobuf standard method, sets type information flags
// It deconstructs type returned by gogo, and than reconstructs it back using type with prepend package.
func (b *fieldBuilder) setGoType() {
	f := b.field // shortrut

	// This call is necessary to fill in generator internal structures, regardless of following resolveType result
	goType, _ := b.generator.GoType(b.descriptor, b.fieldDescriptor)
	f.RawGoType = goType
	f.GoType = goType

	// If type is a slice, mark as slice
	if f.GoType[0] == '[' {
		f.GoType = f.GoType[2:]
		f.GoTypeIsSlice = true
		f.GoTypeFull = "[]"
	}

	// If type is a pointer, mark as pointer
	if f.GoType[0] == '*' {
		f.GoType = f.GoType[1:]
		f.GoTypeIsPtr = true
		f.GoTypeFull = f.GoTypeFull + "*"
	}

	b.prependPackageName()

	// This is an exception: we get all []byte arrays from strings, it is an elementary type on protobuf side
	// TODO: Param containing list of fields which want to be real byte arrays
	if goType == "[]byte" || goType == "[]*byte" {
		f.GoType = goType
		f.GoTypeFull = goType
		return
	}

	f.GoTypeFull = f.GoTypeFull + f.GoType
}

// setMessage sets reference to nested message
func (b *fieldBuilder) setMessage() error {
	// Resolve underlying message via protobuf
	x := b.generator.ObjectNamed(b.fieldDescriptor.GetTypeName())
	desc, ok := x.(*generator.Descriptor)
	if desc == nil || !ok {
		return nil
	}

	// Try to analyse it
	m := BuildMessage(b.generator, desc, false)
	if m == nil {
		return fmt.Errorf("Nested message is invalid for field %v", b.field.Name)
	}

	// Nested message schema, or nil if message is not whitelisted
	b.field.Message = m

	return nil
}

// Sets IsList and IsMap flags
func (b *fieldBuilder) setAggregate() error {
	f := b.field

	if b.generator.IsMap(b.fieldDescriptor) {
		f.IsMap = true
		err := b.setMap()
		if err != nil {
			return err
		}
	} else if b.fieldDescriptor.IsRepeated() {
		f.IsRepeated = true
	}

	return nil
}

// Sets gogo.customtype flag
func (b *fieldBuilder) setCustomType() {
	if !gogoproto.IsCustomType(b.fieldDescriptor) {
		return
	}

	b.field.IsCustomType = true
	b.field.CustomTypeMethodInfix = strings.ReplaceAll(strings.ReplaceAll(b.field.GoType, "/", ""), ".", "")
}

// reflectMap sets map value properties
func (b *fieldBuilder) setMap() error {
	m := b.generator.GoMapType(nil, b.fieldDescriptor)

	keyGoType, _ := b.generator.GoType(b.descriptor, m.KeyField)
	if keyGoType != "string" {
		b.generator.Fail("Maps with non-string keys are not supported")
	}

	valueField := BuildField(b.generator, b.descriptor, m.ValueField)
	if valueField == nil {
		return fmt.Errorf("Failed to reflect map field %s %s", b.field.GoType, b.field.Name)
	}
	b.field.MapValueField = valueField

	return nil
}

// setKind sets field kind which represents field meta type for generation
func (b *fieldBuilder) setKind() {
	f := b.field

	switch {
	case f.IsCustomType:
		f.Kind = "CUSTOM_TYPE"
	case f.IsMap && f.MapValueField.IsMessage:
		// Terraform does not support map of objects. We replace such field with list of objects having key and value fields.
		f.Kind = "OBJECT_MAP" // ex: map[string]struct, requires additional steps to unmarshal
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
