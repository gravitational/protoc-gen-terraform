package plugin

import (
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stoewer/go-strcase"
)

// Field represents field reflection struct
// This struct and the following methods know about both schema details and types, and target structs
type Field struct {
	Name      string // Type name
	NameSnake string // Type name in snake case

	// Field properties
	IsRepeated  bool // Is list
	IsAggregate bool // Is aggregate (either list or map)
	IsMessage   bool // Is message (might be list or map in the same time)
	IsRequired  bool // Is required TODO: implement
	IsTime      bool // Contains time, value needs to be parsed from string
	IsDuration  bool // Contains duration, value needs to be parsed from string

	// Type conversion
	TFSchemaType    string // Type which is reflected in Terraform schema (a-la types.TypeString)
	TFSchemaRawType string // Terraform schema raw value type (float64 for types.Float)
	TFSchemaGoType  string // Go type to convert schema raw type to (uint32, []bytes, time.Time, time.Duration)
	GoType          string // Final field type (casttype, customtype, *, [])
	RawGoType       string // Go type without prefixes, but with package name
	GoTypeIsSlice   bool   // Go type is a slice
	GoTypeIsPtr     bool   // Go type is a pointer

	// Auxilary
	TFSchemaValidate      string // Validation applied to tfschema field
	TFSchemaAggregateType string // If current field is aggregate value, it will be rendered via this type
	TFSchemaMaxItems      int    // If current field has nested message, it is list with max items 1

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
	return nil
}

// TODO: move to main package
// reflectField builds field reflection structure, or returns nil in case field must be skipped
func (p *Plugin) reflectField(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) *Field {
	b := fieldBuilder{
		plugin:          p,
		descriptor:      d,
		fieldDescriptor: f,
		field:           &Field{},
	}
	b.build()
	if b.isValid() {
		return b.field
	}
	return nil
}

// build fills in a Field structure
func (b *fieldBuilder) build() {
	b.setName()
	b.setGoType()
	b.resolveType()
}

// isValid returns true if built type is valid
// TODO make build() bool
func (b *fieldBuilder) isValid() bool {
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

// setTypes sets TFSchemaType and TFSchemaGoType, plus TFSchemaRawType according to TFSchemaType
func (b *fieldBuilder) setTypes(schemaType string, goTypeCast string) {
	b.field.TFSchemaType = schemaType

	t := &b.field.TFSchemaRawType

	switch schemaType {
	case "TypeFloat":
		*t = "float64"
	case "TypeInt":
		*t = "int"
	case "TypeBool":
		*t = "bool"
	case "TypeString":
		*t = "string"
	}

	b.field.TFSchemaGoType = goTypeCast
}

// isTime returns true if field contains time at the end
func (b *fieldBuilder) isTime() bool {
	t := b.fieldDescriptor.TypeName

	isStdTime := gogoproto.IsStdTime(b.fieldDescriptor)
	isGoogleTime := (t != nil && strings.HasSuffix(*t, "google.protobuf.Timestamp"))
	isCastToTime := gogoproto.GetCastType(b.fieldDescriptor) == "time.Time"

	return isStdTime || isGoogleTime || isCastToTime
}

// isDuration return true if field contains duration at the end
func (b *fieldBuilder) isDuration() bool {
	isStdDuration := gogoproto.IsStdDuration(b.fieldDescriptor)
	isCastToDuration := gogoproto.GetCastType(b.fieldDescriptor) == "Duration"

	return isStdDuration || isCastToDuration
}

func (b *fieldBuilder) setGoType() {
	f := b.field // shortrut

	// This call is necessary to fill in generator internal structures, regardless of following resolveType result
	goType, _ := b.plugin.GoType(b.descriptor, b.fieldDescriptor)
	f.GoType = goType
	f.RawGoType = goType

	if goType[0] == '[' {
		f.RawGoType = goType[2:]
		f.GoTypeIsSlice = true
	} else if goType[0] == '*' {
		f.RawGoType = goType[1:]
		f.GoTypeIsPtr = true
	}

	f.RawGoType = b.descriptor.File().GetPackage() + "." + f.RawGoType
}

// resolveType analyses field type and sets required fields in Field structure
// This method is pretty much copy & paste from gogo/protobuf generator.GoType
func (b *fieldBuilder) resolveType() {
	d := b.fieldDescriptor // shortcut
	f := b.field           // shortcut

	switch {
	case b.isTime():
		b.setTypes("TypeString", "time.Time")
		f.TFSchemaValidate = "validation.IsRFC3339Time"
		f.IsTime = true
	case b.isDuration():
		b.setTypes("TypeString", "time.Duration")
		f.IsDuration = true
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_DOUBLE) || gogoproto.IsStdDouble(d):
		b.setTypes("TypeFloat", "double64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_FLOAT) || gogoproto.IsStdFloat(d):
		b.setTypes("TypeFloat", "float32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_INT64) || gogoproto.IsStdInt64(d):
		b.setTypes("TypeInt", "int64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT64) || gogoproto.IsStdUInt64(d):
		b.setTypes("TypeInt", "uint64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_INT32) || gogoproto.IsStdInt32(d):
		b.setTypes("TypeInt", "int32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_UINT32) || gogoproto.IsStdUInt32(d):
		b.setTypes("TypeInt", "uint32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED64):
		b.setTypes("TypeInt", "uint64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_FIXED32):
		b.setTypes("TypeInt", "uint32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_BOOL) || gogoproto.IsStdBool(d):
		b.setTypes("TypeBool", "bool")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_STRING) || gogoproto.IsStdString(d):
		b.setTypes("TypeString", "string")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_BYTES) || gogoproto.IsStdBytes(d):
		b.setTypes("TypeString", "[]byte")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED32):
		b.setTypes("TypeString", "int32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED64):
		b.setTypes("TypeString", "int64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT32):
		b.setTypes("TypeString", "int32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT64):
		b.setTypes("TypeString", "int64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_MESSAGE):
		b.setMessage()
		f.TFSchemaAggregateType = "TypeList"
		f.IsMessage = true
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_ENUM):
		b.setTypes("TypeString", "string")
	default:
		b.plugin.Generator.Fail("unknown type for", b.descriptor.GetName(), b.fieldDescriptor.GetName())
		return
	}

	// TODO: switch
	if b.fieldDescriptor.IsRepeated() {
		f.IsRepeated = true
		f.IsAggregate = true
		f.TFSchemaAggregateType = "TypeList"
	} else {
		// Is not repeated message, still a list
		f.TFSchemaMaxItems = 1
	}

	// MAP
	// if g.IsMap(f) {
	// 	gf, _ := g.GoType(d, g.GoMapType(nil, f).ValueField)
	// 	logrus.Println("      ", gf)
	// }

	// Sets types for underlying maps

	// log.Println(d.TypeName)

	// case descriptor.FieldDescriptorProto_TYPE_MESSAGE:
	// 	desc := g.ObjectNamed(field.GetTypeName())
	// 	typ, wire = g.TypeName(desc), "bytes"

	// switch {
	// if needsStar(field, g.file.proto3 && field.Extendee == nil, message != nil && message.allowOneof()) {
	// 	typ = "*" + typ
	// }
	// if isRepeated(field) {
	// 	typ = "[]" + typ
	// }
	// return
}

// setMessage sets reference to nested message
func (b *fieldBuilder) setMessage() {
	if !b.fieldDescriptor.IsMessage() {
		return
	}

	// logrus.Println(b.descriptor.File().GoPackageName())
	// logrus.Println(b.fieldDescriptor.GetTypeName())

	// Resolve underlying message via protobuf
	x := b.plugin.ObjectNamed(b.fieldDescriptor.GetTypeName())
	desc, ok := x.(*generator.Descriptor)
	if desc == nil || !ok {
		return
	}

	// Nested message schema, or nil if message is not whitelisted
	b.field.Message = b.plugin.reflectMessage(desc)
}

// HasNestedMessage returns true if field has complex type
func (f *Field) HasNestedMessage() bool {
	if f.Message != nil {
		return true
	}

	return false
}

// // build builds in fieldReflect structure
// func (b *fieldReflectBuilder) build() {
// 	b.setName()
// 	b.setGoType()
// 	b.setTFSchemaType()
// 	b.setTFSchemaValidate()
// 	b.setTFSchemaCollectionType()
// 	b.setNestedType()
// 	b.setMessage()
// }

// // getTFSchemaType returns terraform schema type and target go type for a field
// func (b *fieldReflectBuilder) getTFSchemaType() (string, string) {
// 	t := b.field.goType

// // setTFSchemaValidate sets validation function for current schema element
// func (b *fieldReflectBuilder) setTFSchemaValidate() {
// 	if strings.Contains(b.field.goType, "time.Time") {
// 		b.field.tfSchemaValidate = "IsRFC3339Time"
// 	}
// }

// // setTFSchemaCollectionType set tf schema type if it represents a collection
// func (b *fieldReflectBuilder) setTFSchemaCollectionType() {
// 	if b.plugin.IsMap(b.fieldDescriptor) {
// 		b.field.tfSchemaCollectionType = "TypeMap"
// 	}

// 	if b.fieldDescriptor.IsRepeated() {
// 		b.field.tfSchemaCollectionType = "TypeList"
// 	}
// }

// // setNestedType sets flag true if field has nested type
// func (b *fieldReflectBuilder) setNestedType() {
// 	b.field.hasNestedType = b.fieldDescriptor.IsMessage()
// }

// IsAggregate returns true if field is either list or map
