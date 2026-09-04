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

type schemaTarget struct {
	schemaPackage         string
	schemaDescription     string
	functionSuffix        string
	supportsPlanModifiers bool
}

var (
	resourceSchemaTarget = schemaTarget{
		schemaPackage:         ResourceSchema,
		schemaDescription:     "resource schema",
		functionSuffix:        "Resource",
		supportsPlanModifiers: true,
	}
	dataSourceSchemaTarget = schemaTarget{
		schemaPackage:         DataSourceSchema,
		schemaDescription:     "datasource schema",
		functionSuffix:        "DataSource",
		supportsPlanModifiers: false,
	}
)

func (t schemaTarget) functionName(name string) string {
	return "GenSchema" + name + t.functionSuffix
}

func (t schemaTarget) attributeType(attributeType string) string {
	return t.schemaPackage + "." + attributeType
}

// MessageSchemaGenerator is the decorator struct to generate tfsdk.Schema of a message
type MessageSchemaGenerator struct {
	*Message
	i      *Imports
	target schemaTarget
}

// NewMessageSchemaGenerator returns new MessageSchemaGenerator struct
func NewMessageSchemaGenerator(m *Message, i *Imports, target schemaTarget) *MessageSchemaGenerator {
	return &MessageSchemaGenerator{Message: m, i: i, target: target}
}

// Generate returns Go code for message schema
func (m *MessageSchemaGenerator) Generate(writer io.Writer) (int, error) {
	id := m.target.functionName(m.Name)
	schema := m.i.WithPackage(m.target.schemaPackage, "Schema")
	diags := m.i.WithPackage(Diag, "Diagnostics")
	attr := m.i.WithPackage(m.target.schemaPackage, "Attribute")

	j := j.Commentf("// %v returns Terraform Framework %v definition for %v\n", id, m.target.schemaDescription, m.Name).
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
		f := NewFieldSchemaGenerator(f, m.i, m.target)
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
		d[j.Id("Validators")] = generateValidatorsWithType(
			m.i,
			f.ValidatorType,
			f.Validators)
	}

	if m.target.supportsPlanModifiers && len(f.PlanModifiers) > 0 {
		d[j.Id("PlanModifiers")] = generatePlanModifiersWithType(
			m.i,
			f.PlanModifierType,
			f.PlanModifiers)
	}

	return j.Id(m.i.WithType(m.target.attributeType(f.AttributeType))).Values(d)
}

// FieldSchemaGenerator represents the decorator for Field code generation
type FieldSchemaGenerator struct {
	*Field
	i      *Imports
	target schemaTarget
}

// NewFieldSchemaGenerator returns new FieldSchemaGenerator struct
func NewFieldSchemaGenerator(f *Field, i *Imports, target schemaTarget) *FieldSchemaGenerator {
	return &FieldSchemaGenerator{Field: f, i: i, target: target}
}

// Generate returns field schema
func (f *FieldSchemaGenerator) Generate() *j.Statement {
	if f.AttributeType == "" {
		return f.genLegacy()
	}

	dict := f.baseAttributeDict()

	switch f.Kind {
	case ObjectKind:
		return f.genSingleNestedAttribute(dict)
	case ObjectListKind:
		return f.genListNestedAttribute(dict)
	case ObjectMapKind:
		return f.genMapNestedAttribute(dict)
	case PrimitiveKind, PrimitiveListKind, PrimitiveMapKind:
		return f.genPrimitiveAttribute(dict)
	default:
		return f.genCustomAttribute(dict)
	}
}

func (f *FieldSchemaGenerator) baseAttributeDict() j.Dict {
	d := j.Dict{
		j.Id("Description"): j.Lit(f.Comment),
		j.Id("CustomType"):  f.genCustomType(),
		j.Id("ElementType"): f.genElemType(),
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
		d[j.Id("Validators")] = generateValidatorsWithType(
			f.i,
			f.ValidatorType,
			f.Validators)
	}

	// Plan modifiers
	if f.target.supportsPlanModifiers && len(f.PlanModifiers) > 0 {
		d[j.Id("PlanModifiers")] = generatePlanModifiersWithType(
			f.i,
			f.PlanModifierType,
			f.PlanModifiers)
	}

	return d
}

func (f *FieldSchemaGenerator) genPrimitiveAttribute(d j.Dict) *j.Statement {
	return j.Id(f.i.WithType(f.target.attributeType(f.AttributeType))).Values(d)
}

func (f *FieldSchemaGenerator) nestedAttributes(m *MessageSchemaGenerator) *j.Statement {
	fieldsDict := m.fieldsDictSchema()
	return j.Map(j.String()).Id(f.i.WithPackage(f.target.schemaPackage, "Attribute")).Values(fieldsDict)
}

func (f *FieldSchemaGenerator) genSingleNestedAttribute(d j.Dict) *j.Statement {
	m := NewMessageSchemaGenerator(f.Message, f.i, f.target)
	nestedAttributes := f.nestedAttributes(m)
	d[j.Id("Attributes")] = nestedAttributes
	return j.Id(f.i.WithPackage(f.target.schemaPackage, "SingleNestedAttribute")).Values(d)
}

func (f *FieldSchemaGenerator) genListNestedAttribute(d j.Dict) *j.Statement {
	m := NewMessageSchemaGenerator(f.Message, f.i, f.target)
	nestedAttributes := f.nestedAttributes(m)
	d[j.Id("NestedObject")] = j.Id(f.i.WithPackage(f.target.schemaPackage, "NestedAttributeObject")).
		Values(j.Dict{j.Id("Attributes"): nestedAttributes})
	return j.Id(f.i.WithPackage(f.target.schemaPackage, "ListNestedAttribute")).Values(d)
}

func (f *FieldSchemaGenerator) genMapNestedAttribute(d j.Dict) *j.Statement {
	m := NewMessageSchemaGenerator(f.MapValueField.Message, f.i, f.target)
	nestedAttributes := f.nestedAttributes(m)
	d[j.Id("NestedObject")] = j.Id(f.i.WithPackage(f.target.schemaPackage, "NestedAttributeObject")).
		Values(j.Dict{j.Id("Attributes"): nestedAttributes})
	return j.Id(f.i.WithPackage(f.target.schemaPackage, "MapNestedAttribute")).Values(d)
}

func (f *FieldSchemaGenerator) genCustomAttribute(d j.Dict) *j.Statement {
	return j.Id(f.target.functionName(f.Suffix)).
		Call(
			j.Id("ctx"),
			j.Id(f.i.WithType(f.target.attributeType(f.AttributeType))).Values(d),
		)
}

func (f *FieldSchemaGenerator) genCustomType() *j.Statement {
	if f.Kind == PrimitiveKind && f.TerraformType.TypeConstructor != "" {
		return j.Id(f.i.WithType(f.TerraformType.TypeConstructor))
	}
	return nil
}

func (f *FieldSchemaGenerator) genLegacy() *j.Statement {
	d := j.Dict{
		j.Id("Description"): j.Lit(f.Comment),
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
	if f.target.supportsPlanModifiers && len(f.PlanModifiers) > 0 {
		d[j.Id("PlanModifiers")] = generatePlanModifiers(f.i, f.PlanModifiers)
	}

	if f.Kind == CustomKind {
		return j.Id(f.target.functionName(f.Suffix)).Call(j.Id("ctx"), j.Id(f.i.WithPackage(SDK, "Attribute")).Values(d))
	}

	return j.Id(f.i.WithPackage(SDK, "Attribute")).Values(d)
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
		g := NewFieldSchemaGenerator(f.MapValueField, f.i, f.target)
		return j.Id(f.i.WithType(f.Type)).Values(j.Dict{
			j.Id("ElemType"): g.primitiveSchemaTypeDef(),
		})
	}

	return nil
}

func (f *FieldSchemaGenerator) genElemType() *j.Statement {
	switch f.Kind {
	case PrimitiveListKind:
		return f.primitiveSchemaTypeDef()
	case PrimitiveMapKind:
		g := NewFieldSchemaGenerator(f.MapValueField, f.i, f.target)
		return g.primitiveSchemaTypeDef()
	default:
		return nil
	}
}

// Attributes returns a nested attribute definitions
func (f *FieldSchemaGenerator) attributes() *j.Statement {
	switch f.Kind {
	case ObjectKind:
		m := NewMessageSchemaGenerator(f.Message, f.i, f.target)

		return f.xNestedAttributes("Single", m)
	case ObjectMapKind:
		m := NewMessageSchemaGenerator(f.MapValueField.Message, f.i, f.target)

		return f.xNestedAttributes("Map", m)
	case ObjectListKind:
		m := NewMessageSchemaGenerator(f.Message, f.i, f.target)

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

func generatePlanModifiersWithType(imports *Imports, planModifierType string, pm []string) j.Code {
	v := make([]jen.Code, len(pm))
	for i, n := range pm {
		v[i] = j.Id(imports.WithType(n))
	}

	return j.Index().Id(imports.WithType(planModifierType)).Values(v...)
}

func generateValidatorsWithType(imports *Imports, validatorType string, vals []string) j.Code {
	v := make([]jen.Code, len(vals))
	for i, n := range vals {
		v[i] = j.Id(imports.WithType(n))
	}

	return j.Index().Id(imports.WithType(validatorType)).Values(v...)
}
