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
	"strconv"
	"strings"

	"github.com/gravitational/protoc-gen-terraform/config"
	"github.com/gravitational/trace"
	"github.com/sirupsen/logrus"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
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

	// RawGoType field type returned by gogoprotobuf with * and []
	RawGoType string

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

	// CustomTypeMethodInfix custom type schema and unmarshal method name infix
	CustomTypeMethodInfix string

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

	// typeName represents full field type name
	typeName string

	// path represents path to current field in schema
	path string
}

// fieldBuilder creates valid Field values
type fieldBuilder struct {
	// generator is the generator instance
	generator *generator.Generator

	// descriptor is the message descriptor
	descriptor *generator.Descriptor

	// fieldDescriptor is the field descriptor
	fieldDescriptor *descriptor.FieldDescriptorProto

	// field is the target Field value
	field *Field
}

// BuildFields builds []*Field from descriptors of specified message
func BuildFields(m *Message, g *generator.Generator, d *generator.Descriptor, path string) error {
	for _, f := range d.GetField() {
		typeName := getFieldTypeName(d, f)
		path := path + "." + d.GetName()

		// Ignore field if it is listed in cli arg (top level)
		_, ok := config.ExcludeFields[typeName]
		if ok {
			continue
		}

		// Ignore field if it's path is listed in cli arg
		_, ok = config.ExcludeFields[path]
		if ok {
			continue
		}

		f, err := BuildField(g, d, f, path)
		if err != nil {
			invErr, ok := trace.Unwrap(err).(*invalidFieldError)

			// invalidFieldError is not considered fatal, we just need to report it to the user and skip it
			if !ok {
				return err
			}

			logrus.Warning(invErr.Error())
			continue
		}

		if f != nil {
			m.Fields = append(m.Fields, f)
		}
	}

	return nil
}

// BuildField builds field reflection structure, or returns nil in case field build failed
func BuildField(g *generator.Generator, d *generator.Descriptor, f *descriptor.FieldDescriptorProto, path string) (*Field, error) {
	b := newFieldBuilder(g, d, f, path)
	err := b.build()
	return b.field, err
}

// getFieldTypeName returns field type name with package
func getFieldTypeName(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) string {
	return getMessageTypeName(d) + "." + f.GetName()
}

// newFieldBuilder constructs an empty fieldBuilder value
func newFieldBuilder(g *generator.Generator, d *generator.Descriptor, f *descriptor.FieldDescriptorProto, path string) *fieldBuilder {
	return &fieldBuilder{
		generator:       g,
		descriptor:      d,
		fieldDescriptor: f,
		field:           &Field{path: path},
	}
}

// build fills in a Field structure
func (b *fieldBuilder) build() error {
	b.resolveName()

	err := b.resolveType()
	if err != nil {
		return trace.Wrap(err)
	}

	err = b.setGoType()
	if err != nil {
		return trace.Wrap(err)
	}

	err = b.setAggregate()
	if err != nil {
		return trace.Wrap(err)
	}

	b.setCustomType()
	b.resolveKind()
	b.setComment()
	b.setRequired()
	b.setComputed()
	b.setDefault()

	return nil
}

// setName sets the field name
func (b *fieldBuilder) resolveName() {
	name := b.fieldDescriptor.GetName()

	b.field.Name = name
	b.field.NameSnake = strcase.SnakeCase(name)
	b.field.typeName = getFieldTypeName(b.descriptor, b.fieldDescriptor)
}

// setGoType sets go type with gogo/protobuf standard method, sets type information flags
// It deconstructs type returned by the gogo package, and then builds a new type string with the package name prepended.
func (b *fieldBuilder) setGoType() error {
	// This call is necessary to fill in generator internal structures, regardless of following resolveType result
	t, _ := b.generator.GoType(b.descriptor, b.fieldDescriptor)

	if t == "" {
		return newInvalidFieldError(b, "invalid field go type")
	}

	b.field.RawGoType = t
	prefix := ""

	// This is an exception: we get all []byte arrays from strings, it is an elementary type on the protobuf side
	// TODO: Param containing list of fields that need to be transformed into byte arrays
	if t == "[]byte" || t == "[]*byte" {
		b.field.GoType = t
		b.field.GoTypeFull = t
		return nil
	}

	// If type is a slice, mark as slice
	if t[0] == '[' {
		t = t[2:]
		prefix = prefix + "[]"
		b.field.GoTypeIsSlice = true
	}

	// If type is a pointer, mark as pointer
	if t[0] == '*' {
		t = t[1:]
		prefix = prefix + "*"
		b.field.GoTypeIsPtr = true
	}

	t = b.prependPackageName(t)

	b.field.GoType = t
	b.field.GoTypeFull = prefix + t

	return nil
}

// isTypeEq returns true if type equals current field descriptor type
func (b *fieldBuilder) isTypeEq(t descriptor.FieldDescriptorProto_Type) bool {
	return *b.fieldDescriptor.Type == t
}

// isTime returns true if field stores a time value (protobuf or standard library)
func (b *fieldBuilder) isTime() bool {
	t := b.fieldDescriptor.TypeName

	isStdTime := gogoproto.IsStdTime(b.fieldDescriptor)
	isGoogleTime := (t != nil && strings.HasSuffix(*t, "google.protobuf.Timestamp"))
	isCastToTime := b.getCastType() == "time.Time"

	return isStdTime || isGoogleTime || isCastToTime
}

// isDuration returns true if field stores a duration value (protobuf or cast to a standard library type)
func (b *fieldBuilder) isDuration() bool {
	ct := b.getCastType()
	t := b.fieldDescriptor.TypeName

	isStdDuration := gogoproto.IsStdDuration(b.fieldDescriptor)
	isGoogleDuration := (t != nil && strings.HasSuffix(*t, "google.protobuf.Duration"))
	isCastToCustomDuration := ct == config.DurationCustomType
	isCastToDuration := ct == "time.Duration"

	return isStdDuration || isGoogleDuration || isCastToDuration || isCastToCustomDuration
}

// isMessage returns true if field is a message
func (b *fieldBuilder) isMessage() bool {
	return b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_MESSAGE)
}

// isCastType returns true if field has gogoprotobuf.casttype flag
func (b *fieldBuilder) isCastType() bool {
	return gogoproto.IsCastType(b.fieldDescriptor)
}

// isCastType returns true if field has gogoprotobuf.customtype flag
func (b *fieldBuilder) isCustomType() bool {
	return gogoproto.IsCustomType(b.fieldDescriptor)
}

// getCastType returns field cast type name
func (b *fieldBuilder) getCastType() string {
	return gogoproto.GetCastType(b.fieldDescriptor)
}

// getCustomType returns field custom type name
func (b *fieldBuilder) getCustomType() string {
	return gogoproto.GetCustomType(b.fieldDescriptor)
}

// setTypes sets SchemaRawType and SchemaGoType
func (b *fieldBuilder) setSchemaTypes(schemaRawType string, goTypeCast string) {
	b.field.SchemaRawType = schemaRawType
	b.field.SchemaGoType = goTypeCast
}

// resolveType analyses field type and sets required fields in Field structure.
func (b *fieldBuilder) resolveType() error {
	d := b.fieldDescriptor // syntax shortcut for gogoproto.IsStd* methods

	switch {
	case b.isTime():
		b.setSchemaTypes("string", "time.Time")
		b.field.IsTime = true
	case b.isDuration():
		b.setSchemaTypes("string", "time.Duration")
		b.field.IsDuration = true
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
			return trace.Wrap(err)
		}
		b.field.IsMessage = true
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_ENUM):
		b.setSchemaTypes("string", "string")
	default:
		return trace.Wrap(newInvalidFieldError(b, "unknown field type"))
	}

	return nil
}

// prependPackageName prepends package name to cast type, custom type or message type names
func (b *fieldBuilder) prependPackageName(t string) (result string) {
	result = t

	if b.isCastType() {
		result = b.getCastType()
	}

	if b.isCustomType() {
		result = b.getCustomType()
	}

	// Prepend package name to overridden field type
	if b.isCastType() || b.isCustomType() {
		// If cast type is within current package, append default package name to it
		if !strings.Contains(result, ".") && config.DefaultPackageName != "" {
			result = config.DefaultPackageName + "." + result
		}
	} else {
		// Get go type from a message
		if b.isMessage() && b.field.Message != nil {
			result = b.field.Message.GoTypeName
		}
	}

	return result
}

// setMessage sets reference to a nested message
func (b *fieldBuilder) setMessage() error {
	// Resolve underlying message via protobuf
	x := b.generator.ObjectNamed(b.fieldDescriptor.GetTypeName())
	desc, ok := x.(*generator.Descriptor)
	if desc == nil || !ok {
		return nil
	}

	// Try to analyse it
	m, err := BuildMessage(b.generator, desc, false, b.field.path)
	if err != nil {
		// If underlying message is invalid, we must consider current field as invalid and not stop
		_, ok := trace.Unwrap(err).(*invalidMessageError)
		if ok {
			return trace.Wrap(
				newInvalidFieldError(b, fmt.Sprintf("failed to reflect message type information: %v", err.Error())),
			)
		}

		return trace.Wrap(err)
	} else if m == nil {
		return trace.Wrap(newInvalidFieldError(b, "field marked as skipped"))
	}

	// Nested message schema, or nil if message is not whitelisted
	b.field.Message = m

	return nil
}

// setAggregate detects and sets IsList and IsMap flags.
func (b *fieldBuilder) setAggregate() error {
	f := b.field

	if b.generator.IsMap(b.fieldDescriptor) {
		f.IsMap = true
		err := b.setMap()
		if err != nil {
			return trace.Wrap(err)
		}
	} else if b.fieldDescriptor.IsRepeated() {
		f.IsRepeated = true
	}

	return nil
}

// setCustomType sets detects and set IsCustomType flag and custom type method infix
func (b *fieldBuilder) setCustomType() {
	if !gogoproto.IsCustomType(b.fieldDescriptor) {
		return
	}

	b.field.IsCustomType = true
	b.field.CustomTypeMethodInfix = strings.ReplaceAll(strings.ReplaceAll(b.field.GoType, "/", ""), ".", "")
}

// setMap sets map value properties
func (b *fieldBuilder) setMap() error {
	m := b.generator.GoMapType(nil, b.fieldDescriptor)

	keyGoType, _ := b.generator.GoType(b.descriptor, m.KeyField)
	if keyGoType != "string" {
		return trace.Wrap(newInvalidFieldError(b, "non-string map keys are not supported"))
	}

	valueField, err := BuildField(b.generator, b.descriptor, m.ValueField, b.field.path)
	if err != nil {
		return err
	}
	b.field.MapValueField = valueField

	return nil
}

// setKind sets field kind
func (b *fieldBuilder) setKind(kind string) {
	b.field.Kind = kind
}

// setKind sets field kind which represents field meta type for generation
func (b *fieldBuilder) resolveKind() {
	f := b.field // shortcut to field flags used in conditions

	switch {
	case f.IsCustomType:
		b.setKind("CUSTOM_TYPE")
	case f.IsMap && f.MapValueField.IsMessage:
		// Terraform does not support map of objects. We replace such field with list of objects having key and value fields.
		b.setKind("OBJECT_MAP") // ex: map[string]struct, requires additional steps to unmarshal
	case f.IsMap:
		b.setKind("MAP") // ex: map[string]string
	case f.IsRepeated && f.IsMessage:
		b.setKind("REPEATED_MESSAGE") // ex: []struct
	case f.IsRepeated:
		b.setKind("REPEATED_ELEMENTARY") // ex: []string
	case f.IsMessage:
		b.setKind("SINGULAR_MESSAGE") // ex: struct
	default:
		b.setKind("SINGULAR_ELEMENTARY") // ex: string
	}
}

// setComment resolves leading comment for this field
func (b *fieldBuilder) setComment() {
	p := b.descriptor.Path() + ",2," + strconv.Itoa(int(b.fieldDescriptor.GetNumber()-1))

	for _, l := range b.descriptor.File().GetSourceCodeInfo().GetLocation() {
		if getLocationPath(l) == p {
			c := strings.Trim(l.GetLeadingComments(), "\n")
			b.field.RawComment = commentToSingleLine(strings.TrimSpace(c))
			b.field.Comment = appendSlashSlash(c, false)
		}
	}
}

// setRequired sets IsRequired flag
func (b *fieldBuilder) setRequired() {
	_, ok := config.RequiredFields[b.field.typeName]
	if ok {
		b.field.IsRequired = true
	}
}

// setComputed sets IsComputed flag
func (b *fieldBuilder) setComputed() {
	_, ok := config.ComputedFields[b.field.typeName]
	if ok {
		b.field.IsComputed = true
	}
}

// setDefault sets default value
func (b *fieldBuilder) setDefault() {
	v, ok := config.Defaults[b.field.typeName]
	if ok {
		b.field.Default = fmt.Sprintf("%s", v)
	}
}
