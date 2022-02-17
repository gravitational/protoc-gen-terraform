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

package main

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
	schema := m.i.WithPackage(SDK, "Schema")
	diags := m.i.WithPackage(Diag, "Diagnostics")
	attr := m.i.WithPackage(SDK, "Attribute")

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
		j.Id("Type"):     j.Id(m.i.WithType(f.Type)),
		j.Id("Required"): j.Lit(f.Required),
		j.Id("Computed"): j.Lit(f.Computed),
		j.Id("Optional"): j.Lit(f.Optional),
	}

	if len(f.Validators) > 0 {
		d[j.Id("Validators")] = generateValidators(m.i, f.Validators)
	}

	if len(f.PlanModifiers) > 0 {
		d[j.Id("PlanModifiers")] = generatePlanModifiers(m.i, f.PlanModifiers)
	}

	return j.Values(d)
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
	if f.Kind == CustomKind {
		return j.Id("GenSchema" + f.Suffix).Call(j.Id("ctx"))
	}

	d := j.Dict{
		j.Id("Description"): j.Lit(f.RawComment),
		j.Id("Type"):        f.schemaType(), // nils are automatically omitted
		j.Id("Attributes"):  f.attributes(),
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
		d[j.Id("Validators")] = generateValidators(f.i, f.Validators)
	}

	// Plan modifiers
	if len(f.PlanModifiers) > 0 {
		d[j.Id("PlanModifiers")] = generatePlanModifiers(f.i, f.PlanModifiers)
	}

	return j.Values(d)
}

// SchemaType returns the schema Type field value
func (f *FieldSchemaGenerator) schemaType() *j.Statement {
	switch f.Kind {
	case PrimitiveKind:
		return f.primitiveSchemaTypeDef()
	case PrimitiveListKind:
		return j.Id(f.i.WithType(f.Type)).Values(j.Dict{
			j.Id("ElemType"): f.primitiveSchemaTypeDef(),
		})
	case PrimitiveMapKind:
		g := NewFieldSchemaGenerator(f.MapValueField, f.i)

		return j.Id(f.i.WithType(f.Type)).Values(j.Dict{
			j.Id("ElemType"): g.primitiveSchemaTypeDef(),
		})
	}

	return nil
}

// Attributes returns a nested attribute definitions
func (f *FieldSchemaGenerator) attributes() *j.Statement {
	switch f.Kind {
	case ObjectKind:
		m := NewMessageSchemaGenerator(f.Message, f.i)

		return f.xNestedAttributes("Single", m)
	case ObjectMapKind:
		m := NewMessageSchemaGenerator(f.MapValueField.Message, f.i)

		return f.xNestedAttributes("Map", m)
	case ObjectListKind:
		m := NewMessageSchemaGenerator(f.Message, f.i)

		return f.xNestedAttributes("List", m)
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

// xNestedAttributes generates *NestedAttributes call
func (f *FieldSchemaGenerator) xNestedAttributes(typ string, m *MessageSchemaGenerator) *j.Statement {
	var options *j.Statement

	if typ != "Single" {
		options = j.Id(f.i.WithPackage(SDK, typ+"NestedAttributesOptions")).Values()
	}

	return j.Id(f.i.WithPackage(SDK, typ+"NestedAttributes")).Params(
		j.Map(j.String()).Id(f.i.WithPackage(SDK, "Attribute")).Values(m.fieldsDictSchema()),
		options,
	)
}

func generatePlanModifiers(imports *Imports, pm []string) j.Code {
	v := make([]jen.Code, len(pm))
	for i, n := range pm {
		v[i] = j.Id(imports.WithType(n))
	}

	return j.Index().Id(imports.WithPackage(SDK, "AttributePlanModifier")).Values(v...)
}

func generateValidators(imports *Imports, vals []string) j.Code {
	v := make([]jen.Code, len(vals))
	for i, n := range vals {
		v[i] = j.Id(imports.WithType(n))
	}

	return j.Index().Id(imports.WithPackage(SDK, "AttributeValidator")).Values(v...)
}
