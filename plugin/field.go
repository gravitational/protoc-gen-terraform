package plugin

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gzigzigzeo/protoc-gen-terraform/config"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/stew/slice"
)

// Field represents field reflection struct
// This struct and the following methods know about both schema details and types, and target structs
type Field struct {
	Name      string // Type name
	NameSnake string // Type name in snake case

	// Schema properties
	SchemaRawType string // Terraform schema raw value type (float64 for types.Float)
	SchemaGoType  string // Go type to convert schema raw type to (uint32, []byte, time.Time, time.Duration)

	// Field properties
	RawGoType     string // Field type returned by gogoprotobuf
	GoType        string // Go type without prefixes, but with package name
	GoTypeIsSlice bool   // Go type is a slice
	GoTypeIsPtr   bool   // Go type is a pointer

	// Metadata
	Kind                  string // Field kind (resulting of combination of meta flags)
	IsRepeated            bool   // Is list
	IsMap                 bool   // Is map
	IsMessage             bool   // Is message (might be repeated in the same time)
	IsRequired            bool   // Is required TODO: implement
	IsTime                bool   // Contains time, value needs to be parsed from string
	IsDuration            bool   // Contains duration, value needs to be parsed from string
	IsCustomType          bool   // Custom types require manual marshallers and schemas
	CustomTypeMethodInfix string // Custom type method name

	Message       *Message // Reference to nested message
	MapValueField *Field   // Reference to map value field reflection
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
	defer func() {
		if r := recover(); r != nil {
			e, ok := r.(*buildError)
			if !ok {
				panic(r)
			}
			logrus.Printf("%+v", e)
		}
	}()

	b := newFieldBuilder(g, d, f)
	b.build()
	return b.field
}

// getFieldTypeName returns field name with package
func getFieldTypeName(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) string {
	return getMessageTypeName(d) + "." + f.GetName()
}

func newFieldBuilder(g *generator.Generator, d *generator.Descriptor, f *descriptor.FieldDescriptorProto) *fieldBuilder {
	return &fieldBuilder{
		generator:       g,
		descriptor:      d,
		fieldDescriptor: f,
		field:           &Field{},
	}
}

// build fills in a Field structure
func (b *fieldBuilder) build() {
	b.setName()
	b.setGoType()
	b.resolveType()
	b.setAggregate()
	b.setCustomType()
	b.setKind()
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

// setTypes sets SchemaType and SchemaGoType
func (b *fieldBuilder) setSchemaTypes(schemaRawType string, goTypeCast string) {
	b.field.SchemaRawType = schemaRawType
	b.field.SchemaGoType = goTypeCast
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

// setGoType Sets go type with gogoprotobuf standard method, sets type information flags
func (b *fieldBuilder) setGoType() {
	f := b.field // shortrut

	// This call is necessary to fill in generator internal structures, regardless of following resolveType result
	goType, _ := b.generator.GoType(b.descriptor, b.fieldDescriptor)
	f.RawGoType = goType
	f.GoType = goType

	if goType[0] == '*' {
		f.GoType = goType[1:]
		f.GoTypeIsPtr = true
	}

	if goType == "[]byte" {
		return
	}

	if goType[0] == '[' {
		f.GoType = goType[2:]
		f.GoTypeIsSlice = true
	}
}

// resolveType analyses field type and sets required fields in Field structure
// This method is pretty much copy & paste from gogo/protobuf generator.GoType
func (b *fieldBuilder) resolveType() {
	d := b.fieldDescriptor // shortcut
	f := b.field           // shortcut

	// Basics
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
		b.setMessage()
		f.IsMessage = true
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_ENUM):
		b.setSchemaTypes("string", "string")
	default:
		b.generator.Fail("unknown type for", b.descriptor.GetName(), b.fieldDescriptor.GetName())
		return
	}

	b.prependPackageName()
}

// Append type suffix to cast type, custom type and message
func (b *fieldBuilder) prependPackageName() {
	d := b.fieldDescriptor
	f := b.field

	if b.isMessage() && b.field.Message != nil {
		f.GoType = b.field.Message.GoTypeName
		return
	}

	if gogoproto.IsCastType(d) || gogoproto.IsCustomType(d) {
		// Is cast type is within current package, append default package name to it
		if !strings.Contains(f.GoType, ".") && config.DefaultPkgName != "" {
			f.GoType = config.DefaultPkgName + "." + f.GoType
		}
	}
}

// setMessage sets reference to nested message
func (b *fieldBuilder) setMessage() {
	// Resolve underlying message via protobuf
	x := b.generator.ObjectNamed(b.fieldDescriptor.GetTypeName())
	desc, ok := x.(*generator.Descriptor)
	if desc == nil || !ok {
		return
	}

	// Break dependency
	m := BuildMessage(b.generator, desc, false)
	if m == nil {
		panic(newBuildError("Nested message is invalid"))
	}

	// Nested message schema, or nil if message is not whitelisted
	b.field.Message = m
}

func (b *fieldBuilder) setAggregate() {
	f := b.field

	if b.generator.IsMap(b.fieldDescriptor) {
		f.IsMap = true
		b.setMap()
	} else if b.fieldDescriptor.IsRepeated() {
		f.IsRepeated = true
	}
}

func (b *fieldBuilder) setCustomType() {
	if !gogoproto.IsCustomType(b.fieldDescriptor) {
		return
	}

	b.field.IsCustomType = true
	b.field.CustomTypeMethodInfix = strings.ReplaceAll(strings.ReplaceAll(b.field.GoType, "/", ""), ".", "")
}

// reflectMap sets map value properties
func (b *fieldBuilder) setMap() {
	m := b.generator.GoMapType(nil, b.fieldDescriptor)

	keyGoType, _ := b.generator.GoType(b.descriptor, m.KeyField)
	if keyGoType != "string" {
		b.generator.Fail("Maps with non-string keys are not supported")
	}

	valueField := BuildField(b.generator, b.descriptor, m.ValueField)
	if valueField == nil {
		panic(newBuildError(fmt.Sprintf("Failed to reflect map field %s %s", b.field.GoType, b.field.Name)))
	}
	b.field.MapValueField = valueField
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
