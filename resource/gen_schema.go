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

package resource

import (
	// . significantly improves readability of a generator statements.
	// That's also the reason to extract it into the separate package.

	"io"

	"github.com/dave/jennifer/jen"
	j "github.com/dave/jennifer/jen"
)

// MessageSchemaGenerator is the decorator struct to generate tfsdk.Schema of a message
type MessageSchemaGenerator struct {
	*Message
	i *Imports
}

// NewMessageSchemaGenerator returns new MessageSchemaGenerator struct
func NewMessageSchemaGenerator(m *Message, i *Imports) *MessageSchemaGenerator {
	return &MessageSchemaGenerator{Message: m, i: i}
}

// Generate returns Go code for message schema
func (m *MessageSchemaGenerator) Generate(writer io.Writer) (int, error) {
	id := "GenSchema" + m.Name
	schema := m.i.WithPackage(Schema, "Schema")
	diags := m.i.WithPackage(Diag, "Diagnostics")
	attr := m.i.WithPackage(Schema, "Attribute")

	j := j.Commentf("// %v returns tfsdk.Schema definition for %v\n", id, m.Name).
		Func().
		Id(id).
		Params(j.Id("ctx").Id(m.i.WithType("context.Context"))).
		Params( // return params
			j.Id(schema),
			j.Id(diags),
		).
		Block(
			j.Return(
				j.Id(schema).Values(
					j.Dict{
						j.Id("Attributes"): j.Map(j.String()).Id(attr).Values(
							m.fieldsDictSchema(),
						),
					},
				),
				j.Nil(),
			),
		)

	return writer.Write([]byte(j.GoString() + "\n"))
}

// FieldsDictSchema reutrns jen.Dict of the generated message fields
func (m *MessageSchemaGenerator) fieldsDictSchema() j.Dict {
	d := j.Dict{}

	for _, f := range m.Fields {
		f := NewFieldSchemaGenerator(f, m.i)
		d[j.Lit(f.NameSnake)] = f.Generate()
	}

	if len(m.Message.InjectedFields) > 0 {
		for _, f := range m.Message.InjectedFields {
			d[j.Lit(f.Name)] = m.generateInjectedField(f)
		}
	}

	return d
}

// generateInjectedField generates code for injected field
func (m *MessageSchemaGenerator) generateInjectedField(f InjectedField) j.Code {
	d := j.Dict{
		j.Id("Required"): j.Lit(f.Required),
		j.Id("Computed"): j.Lit(f.Computed),
		j.Id("Optional"): j.Lit(f.Optional),
	}

	if len(f.Validators) > 0 {
		valType := validatorTypeForTerraformType(f.Type)
		d[j.Id("Validators")] = generateValidators(m.i, valType, f.Validators)
	}

	if len(f.PlanModifiers) > 0 {
		pmType := planModifierTypeForTerraformType(f.Type)
		d[j.Id("PlanModifiers")] = generatePlanModifiers(m.i, pmType, f.PlanModifiers)
	}

	return j.Id(m.i.WithPackage(Schema, attributeTypeForTerraformType(f.Type))).Values(d)
}

// FieldSchemaGenerator represents the decorator for Field code generation
type FieldSchemaGenerator struct {
	*Field
	i *Imports
}

// NewFieldSchemaGenerator returns new FieldSchemaGenerator struct
func NewFieldSchemaGenerator(f *Field, i *Imports) *FieldSchemaGenerator {
	return &FieldSchemaGenerator{Field: f, i: i}
}

// Generate returns field schema
func (f *FieldSchemaGenerator) Generate() *j.Statement {
	d := f.baseAttributeDict()

	switch f.Kind {
	case ObjectKind:
		return f.singleNestedAttribute(d)
	case ObjectListKind:
		return f.listNestedAttribute(d)
	case ObjectMapKind:
		return f.mapNestedAttribute(d)
	default:
		return f.primitiveAttributeType().Values(d)
	}
}

func (f *FieldSchemaGenerator) baseAttributeDict() j.Dict {
	d := j.Dict{
		j.Id("Description"): j.Lit(f.Comment),
		j.Id("ElementType"): f.elemType(), // nils are automatically omitted
	}

	// Required/Optional
	if f.IsRequired {
		d[j.Id("Required")] = j.True()
	} else {
		d[j.Id("Optional")] = j.True()
	}

	// Sensitive
	if f.IsSensitive {
		d[j.Id("Sensitive")] = j.True()
	}

	// Computed
	if f.IsComputed {
		d[j.Id("Computed")] = j.True()
	}

	// Validators
	if len(f.Validators) > 0 {
		valType := validatorTypeForTerraformType(f.Type)
		d[j.Id("Validators")] = generateValidators(f.i, valType, f.Validators)
	}

	// Plan modifiers
	if len(f.PlanModifiers) > 0 {
		pmType := planModifierTypeForTerraformType(f.Type)
		d[j.Id("PlanModifiers")] = generatePlanModifiers(f.i, pmType, f.PlanModifiers)
	}

	return d
}

// elemType returns the element type
func (f *FieldSchemaGenerator) elemType() *j.Statement {
	switch f.Kind {
	case PrimitiveKind:
		return nil
	case PrimitiveListKind:
		return f.primitiveSchemaTypeDef()
	case PrimitiveMapKind:
		g := NewFieldSchemaGenerator(f.MapValueField, f.i)
		return g.primitiveSchemaTypeDef()
	}

	return nil
}

// primitiveSchemaTypeDef returns the primitive type
func (f *FieldSchemaGenerator) primitiveSchemaTypeDef() *j.Statement {
	if f.IsTypeScalar {
		return j.Id(f.i.WithType(f.ElemType))
	}

	if f.TypeConstructor != "" {
		return j.Id(f.i.WithType(f.TypeConstructor))
	}

	return j.Id(f.i.WithType(f.ElemType)).Values()
}

func (f *FieldSchemaGenerator) primitiveAttributeType() *j.Statement {
	return j.Id(f.i.WithPackage(Schema, attributeTypeForTerraformType(f.Type)))
}

func (f *FieldSchemaGenerator) nestedAttributes(m *MessageSchemaGenerator) *j.Statement {
	return j.Map(j.String()).Id(f.i.WithPackage(Schema, "Attribute")).Values(m.fieldsDictSchema())
}

func (f *FieldSchemaGenerator) singleNestedAttribute(d j.Dict) *j.Statement {
	m := NewMessageSchemaGenerator(f.Message, f.i)
	d[j.Id("Attributes")] = f.nestedAttributes(m)
	return j.Id(f.i.WithPackage(Schema, "SingleNestedAttribute")).Values(d)
}

func (f *FieldSchemaGenerator) listNestedAttribute(d j.Dict) *j.Statement {
	m := NewMessageSchemaGenerator(f.Message, f.i)
	d[j.Id("NestedObject")] = j.Id(f.i.WithPackage(Schema, "NestedAttributeObject")).Values(j.Dict{
		j.Id("Attributes"): f.nestedAttributes(m),
	})
	return j.Id(f.i.WithPackage(Schema, "ListNestedAttribute")).Values(d)
}

func (f *FieldSchemaGenerator) mapNestedAttribute(d j.Dict) *j.Statement {
	m := NewMessageSchemaGenerator(f.MapValueField.Message, f.i)
	d[j.Id("NestedObject")] = j.Id(f.i.WithPackage(Schema, "NestedAttributeObject")).Values(j.Dict{
		j.Id("Attributes"): f.nestedAttributes(m),
	})
	return j.Id(f.i.WithPackage(Schema, "MapNestedAttribute")).Values(d)
}

func generatePlanModifiers(imports *Imports, pmType string, pm []string) j.Code {
	v := make([]jen.Code, len(pm))
	for i, n := range pm {
		v[i] = j.Id(imports.WithType(n))
	}

	return j.Index().Id(imports.WithType(pmType)).Values(v...)
}

func generateValidators(imports *Imports, valType string, vals []string) j.Code {
	v := make([]jen.Code, len(vals))
	for i, n := range vals {
		v[i] = j.Id(imports.WithType(n))
	}

	return j.Index().Id(imports.WithType(valType)).Values(v...)
}

func attributeTypeForTerraformType(t string) string {
	switch t {
	case Types + ".StringType":
		return "StringAttribute"
	case Types + ".BoolType":
		return "BoolAttribute"
	case Types + ".Int64Type":
		return "Int64Attribute"
	case Types + ".Float64Type":
		return "Float64Attribute"
	case Types + ".ListType":
		return "ListAttribute"
	case Types + ".MapType":
		return "MapAttribute"
	case Types + ".ObjectType":
		return "ObjectAttribute"
	default:
		return "ObjectAttribute"
	}
}

func planModifierTypeForTerraformType(t string) string {
	switch t {
	case Types + ".StringType":
		return PlanModifier + ".String"
	case Types + ".BoolType":
		return PlanModifier + ".Bool"
	case Types + ".Int64Type":
		return PlanModifier + ".Int64"
	case Types + ".Float64Type":
		return PlanModifier + ".Float64"
	case Types + ".ListType":
		return PlanModifier + ".List"
	case Types + ".MapType":
		return PlanModifier + ".Map"
	case Types + ".ObjectType":
		return PlanModifier + ".Object"
	default:
		return PlanModifier + ".Object"
	}
}

func validatorTypeForTerraformType(t string) string {
	switch t {
	case Types + ".StringType":
		return Validator + ".String"
	case Types + ".BoolType":
		return Validator + ".Bool"
	case Types + ".Int64Type":
		return Validator + ".Int64"
	case Types + ".Float64Type":
		return Validator + ".Float64"
	case Types + ".ListType":
		return Validator + ".List"
	case Types + ".MapType":
		return Validator + ".Map"
	case Types + ".ObjectType":
		return Validator + ".Object"
	default:
		return ""
	}
}
