package main

import (
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stoewer/go-strcase"
)

type fieldReflect struct {
	name             string // Field name
	snakeName        string // Snake cased field name
	tfSchemaType     string // Terraform schema type (e.g schema.TypeBool without "schema" part)
	tfSchemaGoType   string // Terraform schema go type to cast value read from schema into
	tfSchemaValidate string // Terraform validator if needed

	//entityGoType string // Go type to convert schema value from
}

// newFieldsReflect generates slice of reflect structures for message fields
func (p *Plugin) newFieldsReflect(m *messageReflect, d *generator.Descriptor) {
	m.fields = make([]*fieldReflect, len(d.Field))

	for index, f := range d.Field {
		m.fields[index] = p.newFieldReflect(m, d, f)
	}
}

// newFieldReflect generates reflect structure for message field
func (p *Plugin) newFieldReflect(
	m *messageReflect,
	d *generator.Descriptor,
	f *descriptor.FieldDescriptorProto,
) *fieldReflect {
	field := &fieldReflect{}

	// Set base parameters
	field.name = f.GetName()
	field.snakeName = strcase.SnakeCase(field.name)

	// Set terraform schema type parameters. We'll use them to read values from ResourceData object.
	goType, _ := p.GoType(d, f)
	tfSchemaType, tfSchemaGoType := tfSchemaTypeFromGoType(goType)

	field.tfSchemaType = tfSchemaType
	field.tfSchemaGoType = tfSchemaGoType

	// Validation, required for time field
	field.tfSchemaValidate = tfSchemaValidateFromGoType(goType)

	return field
}

// tfSchemaTypeFromGoType returns Terraform type and go type for that Terraform type
func tfSchemaTypeFromGoType(goType string) (string, string) {
	if strings.Contains(goType, "float") || strings.Contains(goType, "fixed") {
		return "TypeFloat", "float64"
	}

	if strings.Contains(goType, "string") {
		return "TypeString", "string"
	}

	if strings.Contains(goType, "int") {
		return "TypeInt", "int"
	}

	if strings.Contains(goType, "bool") {
		return "TypeBool", "bool"
	}

	if strings.Contains(goType, "time.Time") {
		return "TypeString", "string"
	}

	if strings.Contains(goType, "time.Duration") {
		return "TypeInt", "int"
	}

	return "", ""
}

// tfSchemaTypeFromGoType returns validation function for current schema element
func tfSchemaValidateFromGoType(goType string) string {
	if strings.Contains(goType, "time.Time") {
		return "IsRFC3339Time"
	}

	return ""
}
