package main

import (
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stoewer/go-strcase"
)

type fieldReflect struct {
	name      string // Field name
	snakeName string // Snake cased field name

	goType string // Target entity go type

	tfSchemaType           string // Terraform schema type (e.g schema.TypeBool without "schema" part)
	tfSchemaGoType         string // Terraform schema go type to cast value read from schema into (e.g. .(string) for TypeString)
	tfSchemaValidate       string // Terraform validator if needed (for time, for now)
	tfSchemaCollectionType string // Terraform collection type (list or map)

	hasNestedType bool            // Field represents nested structure
	message       *messageReflect // Nested message definition
}

type fieldReflectBuilder struct {
	plugin          *Plugin
	descriptor      *generator.Descriptor
	fieldDescriptor *descriptor.FieldDescriptorProto
	field           fieldReflect
}

// reflectFields generates slice of reflect structures for message fields
func (p *Plugin) reflectFields(m *messageReflect, d *generator.Descriptor) {
	m.fields = make([]*fieldReflect, len(d.Field))

	for index, f := range d.Field {
		m.fields[index] = p.reflectField(d, f)
	}
}

// reflectField builds fieldReflect for specific field
func (p *Plugin) reflectField(d *generator.Descriptor, f *descriptor.FieldDescriptorProto) *fieldReflect {
	b := fieldReflectBuilder{
		plugin:          p,
		descriptor:      d,
		fieldDescriptor: f,
	}
	b.build()
	return &b.field
}

// build builds in fieldReflect structure
func (b *fieldReflectBuilder) build() {
	b.setName()
	b.setGoType()
	b.setTFSchemaType()
	b.setTFSchemaValidate()
	b.setTFSchemaCollectionType()
	b.setNestedType()
	b.setMessage()
}

// setName sets field name and + snake cased
func (b *fieldReflectBuilder) setName() {
	b.field.name = b.fieldDescriptor.GetName()
	b.field.snakeName = strcase.SnakeCase(b.field.name)
}

// setGoType sets target structure field go type
func (b *fieldReflectBuilder) setGoType() {
	goType, _ := b.plugin.GoType(b.descriptor, b.fieldDescriptor)
	b.field.goType = goType
}

// setTFSchemaType sets terraform schema type and target go type for a field
func (b *fieldReflectBuilder) setTFSchemaType() {
	t, g := b.getTFSchemaType()
	b.field.tfSchemaType = t
	b.field.tfSchemaGoType = g
}

// getTFSchemaType returns terraform schema type and target go type for a field
func (b *fieldReflectBuilder) getTFSchemaType() (string, string) {
	t := b.field.goType

	if strings.Contains(t, "float") || strings.Contains(t, "fixed") {
		return "TypeFloat", "float64"
	}

	if strings.Contains(t, "string") {
		return "TypeString", "string"
	}

	if strings.Contains(t, "int") {
		return "TypeInt", "int"
	}

	if strings.Contains(t, "bool") {
		return "TypeBool", "bool"
	}

	if strings.Contains(t, "byte") {
		return "TypeString", "string"
	}

	if strings.Contains(t, "time.Time") {
		return "TypeString", "string"
	}

	if strings.Contains(t, "time.Duration") {
		return "TypeInt", "int"
	}

	return "", ""
}

// setTFSchemaValidate sets validation function for current schema element
func (b *fieldReflectBuilder) setTFSchemaValidate() {
	if strings.Contains(b.field.goType, "time.Time") {
		b.field.tfSchemaValidate = "IsRFC3339Time"
	}
}

// setTFSchemaCollectionType set tf schema type if it represents a collection
func (b *fieldReflectBuilder) setTFSchemaCollectionType() {
	if b.plugin.IsMap(b.fieldDescriptor) {
		b.field.tfSchemaCollectionType = "TypeMap"
	}

	if b.fieldDescriptor.IsRepeated() {
		b.field.tfSchemaCollectionType = "TypeList"
	}
}

// setNestedType sets flag true if field has nested type
func (b *fieldReflectBuilder) setNestedType() {
	b.field.hasNestedType = b.fieldDescriptor.IsMessage()
}

// setMessage sets reference to complex
func (b *fieldReflectBuilder) setMessage() {
	if !b.fieldDescriptor.IsMessage() {
		return
	}

	x := b.plugin.ObjectNamed(b.fieldDescriptor.GetTypeName())
	desc, ok := x.(*generator.Descriptor)
	if desc == nil || !ok {
		return
	}
	b.field.message = b.plugin.reflectMessage(desc)
}
