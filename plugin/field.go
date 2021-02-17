package plugin

import (
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stoewer/go-strcase"
)

const (
	durationCustomType = "Duration" // This type name will be treated as extendee of time.Duration, TODO: Parameterize?
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
	Kind                       string // Field kind (resulting of combination of meta flags)
	IsRepeated                 bool   // Is list
	IsMap                      bool
	IsAggregate                bool // Is aggregate (either list or map)
	IsMessage                  bool // Is message (might be repeated in the same time)
	IsRequired                 bool // Is required TODO: implement
	IsTime                     bool // Contains time, value needs to be parsed from string
	IsDuration                 bool // Contains duration, value needs to be parsed from string
	IsElementaryValueContainer bool // Field contains single field

	Message *Message // Reference to nested message
}

// fieldBuilder is axilarry struct responsible for building Field
type fieldBuilder struct {
	plugin          *Plugin
	descriptor      *generator.Descriptor
	fieldDescriptor *descriptor.FieldDescriptorProto
	field           *Field
}

func (p *Plugin) newFieldBuilder(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) *fieldBuilder {
	return &fieldBuilder{
		plugin:          p,
		descriptor:      d,
		fieldDescriptor: f,
		field:           &Field{},
	}
}

// build fills in a Field structure
func (b *fieldBuilder) build() bool {
	b.setName()
	b.setGoType()
	b.resolveType()
	b.setKind()

	return b.isValid()
}

// isValid returns true if built type is valid
func (b *fieldBuilder) isValid() bool {
	// Maps are temporary disabled
	if b.field.IsMap {
		return false
	}

	// No kind == invalid
	if b.field.Kind == "" {
		return false
	}

	// If field is message, but underlying message failed to reflect (is invalid)
	if b.field.IsMessage && b.field.Message == nil {
		return false
	}

	return true
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
	isCastToCustomDuration := ct == durationCustomType
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
	goType, _ := b.plugin.GoType(b.descriptor, b.fieldDescriptor)
	f.RawGoType = goType
	f.GoType = goType

	if goType[0] == '[' {
		f.GoType = goType[2:]
		f.GoTypeIsSlice = true
	} else if goType[0] == '*' {
		f.GoType = goType[1:]
		f.GoTypeIsPtr = true
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
		b.setSchemaTypes("float64", "double64")
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
		b.plugin.Generator.Fail("unknown type for", b.descriptor.GetName(), b.fieldDescriptor.GetName())
		return
	}

	// Field value is repeated (slice)
	if b.fieldDescriptor.IsRepeated() {
		f.IsRepeated = true
		f.IsAggregate = true
	}

	if b.plugin.IsMap(b.fieldDescriptor) {
		f.IsMap = true

		x := b.plugin.GoMapType(nil, b.fieldDescriptor)
		x = x
		//logrus.Println(x)
	}

	// Append type suffix to cast type, custom type and message
	// Note on custom type: it is guaranteed that custom type has the same subset of fields as protobuf schema type
	// However, it does not guarantee that custom type can be directly casted to schema type
	// This is clearly gogoprotobuf antipattern or lack of proper implementation
	if gogoproto.IsCastType(d) || gogoproto.IsCustomType(d) || b.isMessage() {
		l := f.GoType[0:1]

		// In other words, if the first letter of a type name is uppercase, this means that type is not prefixed with
		// package name.
		if strings.ToLower(l) != l {
			f.GoType = b.descriptor.File().GoPackageName() + "." + f.GoType
		}
	}
}

// setMessage sets reference to nested message
func (b *fieldBuilder) setMessage() {
	// Resolve underlying message via protobuf
	x := b.plugin.ObjectNamed(b.fieldDescriptor.GetTypeName())
	desc, ok := x.(*generator.Descriptor)
	if desc == nil || !ok {
		return
	}

	m := b.plugin.reflectMessage(desc, true)

	if m != nil {
		// Nested message schema, or nil if message is not whitelisted
		b.field.Message = m
		b.setIsContainer()
	}
}

// setIsFold sets folding flag. This flag means that field is a container for a single field
// For instance, that could be custom BoolValue with only bool Value field, which could be set directly.
func (b *fieldBuilder) setIsContainer() {
	f := b.field

	if f.IsAggregate {
		return
	}

	m := f.Message

	if len(m.Fields) == 1 && !m.Fields[0].IsAggregate {
		f.IsElementaryValueContainer = true
	}
}

// setKind sets field kind which represents field kind for generation
func (b *fieldBuilder) setKind() {
	f := b.field

	switch {
	case f.IsAggregate && f.IsRepeated && f.IsMessage:
		f.Kind = "REPEATED_MESSAGE"
	case f.IsAggregate && f.IsRepeated:
		f.Kind = "REPEATED_ELEMENTARY"
	case f.IsElementaryValueContainer:
		f.Kind = "ELEMENTARY_CONTAINER"
	case f.IsMessage:
		f.Kind = "SINGULAR_MESSAGE"
	default:
		f.Kind = "SINGULAR_ELEMENTARY"
	}
}
