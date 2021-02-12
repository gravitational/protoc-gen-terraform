package plugin

import (
	"log"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stoewer/go-strcase"
)

// Field represents field reflection struct
type Field struct {
	Name       string // Type name
	NameSnake  string // Type name in snake case
	GoType     string // Field go type
	IsRepeated bool   // Field is list
	IsNullable bool   // Field is nullable and has *

	TFSchemaType          string // Type which is reflected in Terraform schema
	TFSchemaTypeCast      string // Type which must Terraform schema value cast to
	TFSchemaGoTypeCast    string // Final type to cast to in go (would be time.Duration, while TypeCast is int)
	TFSchemaValidate      string // Validation applied to tfschema field
	TFSchemaAggregateType string // If current field is aggregate value, it will be rendered via this type

	Message *Message // Nested message
}

type fieldBuilder struct {
	plugin          *Plugin
	descriptor      *generator.Descriptor
	fieldDescriptor *descriptor.FieldDescriptorProto
	field           *Field
}

func (p *Plugin) reflectFields(m *Message, d *generator.Descriptor) {
	for _, f := range d.GetField() {
		if p.isFieldRequired(f) {
			m.Fields = append(m.Fields, p.reflectField(d, f))
		}
	}
}

// isFieldRequired returns true if field type is listed in allowed types
func (p *Plugin) isFieldRequired(f *descriptor.FieldDescriptorProto) bool {
	return true
}

func (p *Plugin) reflectField(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) *Field {
	b := fieldBuilder{
		plugin:          p,
		descriptor:      d,
		fieldDescriptor: f,
		field:           &Field{},
	}
	b.build()
	return b.field
}

// build fills in a Field structure
func (b *fieldBuilder) build() {
	b.setName()
	b.resolveType()
	b.setNullable()
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

// setTypes is utility setter method
func (b *fieldBuilder) setTypes(schemaType string, goTypeCast string) {
	b.field.TFSchemaType = schemaType

	t := &b.field.TFSchemaTypeCast

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

	b.field.TFSchemaGoTypeCast = goTypeCast
}

// resolveType analyses field type and sets required fields in Field structure
// This method is pretty much copy & paste from gogo/protobuf generator.GoType
func (b *fieldBuilder) resolveType() {
	d := b.fieldDescriptor

	// This call is necessary to fill in generator internal structures, regardless of following resolveType result
	goType, _ := b.plugin.GoType(b.descriptor, d)
	b.field.GoType = goType

	switch {
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
		b.setTypes("TypeString", "string")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED32):
		b.setTypes("TypeString", "int32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SFIXED64):
		b.setTypes("TypeString", "int64")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT32):
		b.setTypes("TypeString", "int32")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_SINT64):
		b.setTypes("TypeString", "int64")
	case gogoproto.IsStdTime(d) || *d.TypeName == ".google.protobuf.Timestamp":
		b.setTypes("TypeString", "time.Time")
		b.field.TFSchemaValidate = "validation.IsRFC3339Time"
	case gogoproto.IsStdDuration(d):
		b.setTypes("TypeInt", "time.Duration")
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_MESSAGE):
		b.setMessage()
	case b.isTypeEq(descriptor.FieldDescriptorProto_TYPE_ENUM):
		b.setTypes("TypeString", "string")
	default:
		log.Fatal("unknown type for", d.GetName())
	}

	if b.fieldDescriptor.IsRepeated() {
		b.field.IsRepeated = true
		b.field.TFSchemaAggregateType = "TypeList"
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
	// case gogoproto.IsCustomType(field) && gogoproto.IsCastType(field):
	// 	g.Fail(field.GetName() + " cannot be custom type and cast type")
	// case gogoproto.IsCustomType(field):
	// 	var packageName string
	// 	var err error
	// 	packageName, typ, err = getCustomType(field)
	// 	if err != nil {
	// 		g.Fail(err.Error())
	// 	}
	// 	if len(packageName) > 0 {
	// 		g.customImports = append(g.customImports, packageName)
	// 	}
	// case gogoproto.IsCastType(field):
	// 	var packageName string
	// 	var err error
	// 	packageName, typ, err = getCastType(field)
	// 	if err != nil {
	// 		g.Fail(err.Error())
	// 	}
	// 	if len(packageName) > 0 {
	// 		g.customImports = append(g.customImports, packageName)
	// 	}
	// }
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

	x := b.plugin.ObjectNamed(b.fieldDescriptor.GetTypeName())
	desc, ok := x.(*generator.Descriptor)
	if desc == nil || !ok {
		return
	}
	b.field.Message = b.plugin.reflectMessage(desc)
	b.plugin.registerMessage(b.field.Message)
}

// setNullable sets nullable flag
func (b *fieldBuilder) setNullable() {
	b.field.IsNullable = b.field.GoType[0] == '*'
}

// 	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
// 	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
// 	"github.com/stoewer/go-strcase"
// )

// type fieldReflect struct {
// 	name      string // Field name
// 	snakeName string // Snake cased field name

// 	goType string // Target entity go type

// 	tfSchemaType           string // Terraform schema type (e.g schema.TypeBool without "schema" part)
// 	tfSchemaGoType         string // Terraform schema go type to cast value read from schema into (e.g. .(string) for TypeString)
// 	tfSchemaValidate       string // Terraform validator if needed (for time, for now)
// 	tfSchemaCollectionType string // Terraform collection type (list or map)

// 	hasNestedType bool            // Field represents nested structure
// 	message       *messageReflect // Nested message definition
// }

// type fieldReflectBuilder struct {
// 	plugin          *Plugin
// 	descriptor      *generator.Descriptor
// 	fieldDescriptor *descriptor.FieldDescriptorProto
// 	field           fieldReflect
// }

// // reflectFields generates slice of reflect structures for message fields
// func (p *Plugin) reflectFields(m *messageReflect, d *generator.Descriptor) {
// 	m.fields = make([]*fieldReflect, len(d.Field))

// 	for index, f := range d.Field {
// 		m.fields[index] = p.reflectField(d, f)
// 	}
// }

// // reflectField builds fieldReflect for specific field
// func (p *Plugin) reflectField(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) *fieldReflect {
// 	b := fieldReflectBuilder{
// 		plugin:          p,
// 		descriptor:      d,
// 		fieldDescriptor: f,
// 	}
// 	b.build()
// 	return &b.field
// }

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

// // setName sets field name and + snake cased
// func (b *fieldReflectBuilder) setName() {
// 	b.field.name = b.fieldDescriptor.GetName()
// 	b.field.snakeName = strcase.SnakeCase(b.field.name)
// }

// // setGoType sets target structure field go type

// // setTFSchemaType sets terraform schema type and target go type for a field
// func (b *fieldReflectBuilder) setTFSchemaType() {
// 	t, g := b.getTFSchemaType()
// 	b.field.tfSchemaType = t
// 	b.field.tfSchemaGoType = g
// }

// // getTFSchemaType returns terraform schema type and target go type for a field
// func (b *fieldReflectBuilder) getTFSchemaType() (string, string) {
// 	t := b.field.goType

// 	if strings.Contains(t, "float") || strings.Contains(t, "fixed") {
// 		return "TypeFloat", "float64"
// 	}

// 	if strings.Contains(t, "string") {
// 		return "TypeString", "string"
// 	}

// 	if strings.Contains(t, "int") {
// 		return "TypeInt", "int"
// 	}

// 	if strings.Contains(t, "bool") {
// 		return "TypeBool", "bool"
// 	}

// 	if strings.Contains(t, "byte") {
// 		return "TypeString", "string"
// 	}

// 	if strings.Contains(t, "time.Time") {
// 		return "TypeString", "string"
// 	}

// 	if strings.Contains(t, "time.Duration") {
// 		return "TypeInt", "int"
// 	}

// 	return "", ""
// }

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

func (f *Field) IsAggregate() bool {
	if f.IsRepeated {
		return true
	}
	return false
}
