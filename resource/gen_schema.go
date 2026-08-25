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
	"github.com/gravitational/trace"
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

	fieldsDict, err := m.fieldsDictSchema()
	if err != nil {
		return 0, trace.Wrap(err)
	}

	j := j.Commentf("// %v returns tfsdk.Schema definition for %v\n", id, m.Name).
		Func().
		Id(id).
		Params(j.Id("ctx").Id(m.i.WithType("context.Context"))).
		Params( // return params
			j.Id(schema),
			j.Id(diags),
		).
		Block(
			j.Var().Id("diags").Id(diags),
			j.Return(
				j.Id(schema).Values(
					j.Dict{
						j.Id("Attributes"): j.Map(j.String()).Id(attr).Values(
							fieldsDict,
						),
					},
				),
				j.Id("diags"),
			),
		)

	return writer.Write([]byte(j.GoString() + "\n"))
}

// FieldsDictSchema reutrns jen.Dict of the generated message fields
func (m *MessageSchemaGenerator) fieldsDictSchema() (j.Dict, error) {
	d := j.Dict{}

	for _, f := range m.Fields {
		f := NewFieldSchemaGenerator(f, m.i)
		field, err := f.Generate()
		if err != nil {
			return nil, trace.Wrap(err)
		}
		d[j.Lit(f.NameSnake)] = field
	}

	if len(m.Message.InjectedFields) > 0 {
		for _, f := range m.Message.InjectedFields {
			injected, err := m.generateInjectedField(f)
			if err != nil {
				return nil, trace.Wrap(err)
			}
			d[j.Lit(f.Name)] = injected
		}
	}

	return d, nil
}

// generateInjectedField generates code for injected field
func (m *MessageSchemaGenerator) generateInjectedField(f InjectedField) (*j.Statement, error) {
	d := j.Dict{
		j.Id("Required"): j.Lit(f.Required),
		j.Id("Computed"): j.Lit(f.Computed),
		j.Id("Optional"): j.Lit(f.Optional),
	}

	if len(f.Validators) > 0 {
		baseType, err := baseTypeForAttributeType(f.AttributeType)
		if err != nil {
			return nil, trace.Wrap(err, "failed to get base type")
		}
		d[j.Id("Validators")] = generateValidators(m.i, j.Id(m.i.WithType(Validator+baseType)), f.Validators)
	}

	if len(f.PlanModifiers) > 0 {
		baseType, err := baseTypeForAttributeType(f.AttributeType)
		if err != nil {
			return nil, trace.Wrap(err, "failed to get base type")
		}
		d[j.Id("PlanModifiers")] = generatePlanModifiers(m.i, j.Id(m.i.WithType(PlanModifier+baseType)), f.PlanModifiers)
	}

	return j.Id(m.i.WithType(f.AttributeType)).Values(d), nil
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
func (f *FieldSchemaGenerator) Generate() (j.Code, error) {
	d, err := f.baseAttributeDict()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	switch f.Kind {
	case ObjectKind:
		return f.singleNestedAttribute(d)
	case ObjectListKind:
		return f.listNestedAttribute(d)
	case ObjectMapKind:
		return f.mapNestedAttribute(d)
	case CustomKind:
		return f.customAttribute(d), nil
	default:
		return f.primitiveAttribute().Values(d), nil
	}
}

func (f *FieldSchemaGenerator) baseAttributeDict() (j.Dict, error) {
	d := j.Dict{
		j.Id("Description"): j.Lit(f.Comment),
		j.Id("CustomType"):  f.customType(), // nils are automatically omitted
		j.Id("ElementType"): f.elemType(),   // nils are automatically omitted
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
		validators, err := f.generateValidators()
		if err != nil {
			return nil, trace.Wrap(err)
		}
		d[j.Id("Validators")] = validators
	}

	// Plan modifiers
	if len(f.PlanModifiers) > 0 {
		planModifiers, err := f.generatePlanModifiers()
		if err != nil {
			return nil, trace.Wrap(err)
		}
		d[j.Id("PlanModifiers")] = planModifiers
	}

	return d, nil
}

func (f *FieldSchemaGenerator) customType() *j.Statement {
	if f.Kind != PrimitiveKind {
		return nil
	}

	if f.TerraformType.TypeConstructor != "" {
		return j.Id(f.TerraformType.TypeConstructor)
	}
	return nil
}

// elemType returns the element type
func (f *FieldSchemaGenerator) elemType() *j.Statement {
	switch f.Kind {
	case PrimitiveListKind:
		return f.primitiveSchemaTypeDef()
	case PrimitiveMapKind:
		g := NewFieldSchemaGenerator(f.MapValueField, f.i)
		return g.primitiveSchemaTypeDef()
	default:
		return nil
	}
}

// primitiveSchemaTypeDef returns the primitive type
func (f *FieldSchemaGenerator) primitiveSchemaTypeDef() *j.Statement {
	if f.IsTypeScalar {
		return j.Id(f.i.WithType(f.ElemType))
	}

	if f.TerraformType.TypeConstructor != "" {
		return j.Id(f.TerraformType.TypeConstructor)
	}

	return j.Id(f.i.WithType(f.ElemType)).Values()
}

func (f *FieldSchemaGenerator) primitiveAttribute() *j.Statement {
	return j.Id(f.i.WithType(f.AttributeType))
}

func (f *FieldSchemaGenerator) nestedAttributes(m *MessageSchemaGenerator) (*j.Statement, error) {
	fieldsDict, err := m.fieldsDictSchema()
	if err != nil {
		return nil, trace.Wrap(err)
	}
	return j.Map(j.String()).Id(f.i.WithPackage(Schema, "Attribute")).Values(fieldsDict), nil
}

func (f *FieldSchemaGenerator) singleNestedAttribute(d j.Dict) (*j.Statement, error) {
	m := NewMessageSchemaGenerator(f.Message, f.i)
	nestedAttributes, err := f.nestedAttributes(m)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	d[j.Id("Attributes")] = nestedAttributes
	return j.Id(f.i.WithPackage(Schema, "SingleNestedAttribute")).Values(d), nil
}

func (f *FieldSchemaGenerator) listNestedAttribute(d j.Dict) (*j.Statement, error) {
	m := NewMessageSchemaGenerator(f.Message, f.i)
	nestedAttributes, err := f.nestedAttributes(m)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	d[j.Id("NestedObject")] = j.Id(f.i.WithPackage(Schema, "NestedAttributeObject")).Values(j.Dict{
		j.Id("Attributes"): nestedAttributes,
	})
	return j.Id(f.i.WithPackage(Schema, "ListNestedAttribute")).Values(d), nil
}

func (f *FieldSchemaGenerator) mapNestedAttribute(d j.Dict) (*j.Statement, error) {
	m := NewMessageSchemaGenerator(f.MapValueField.Message, f.i)
	nestedAttributes, err := f.nestedAttributes(m)
	if err != nil {
		return nil, trace.Wrap(err)
	}
	d[j.Id("NestedObject")] = j.Id(f.i.WithPackage(Schema, "NestedAttributeObject")).Values(j.Dict{
		j.Id("Attributes"): nestedAttributes,
	})
	return j.Id(f.i.WithPackage(Schema, "MapNestedAttribute")).Values(d), nil
}

func (f *FieldSchemaGenerator) customAttribute(d j.Dict) *j.Statement {
	return j.Id("GenSchema"+f.Suffix).
		Call(
			j.Id("ctx"),
			j.Op("&").Id("diags"),
			j.Id(f.i.WithType(f.AttributeType)).Values(d),
		)
}

func (f *FieldSchemaGenerator) generatePlanModifiers() (*j.Statement, error) {
	planModifierType, err := f.planModifierTypeForAttributeType(f.AttributeType)
	if err != nil {
		return nil, trace.Wrap(err, "failed to get plan modifier type")
	}

	return generatePlanModifiers(f.i, planModifierType, f.PlanModifiers), nil
}

func (f *FieldSchemaGenerator) generateValidators() (*j.Statement, error) {
	validatorType, err := f.validatorTypeForAttributeType(f.AttributeType)
	if err != nil {
		return nil, trace.Wrap(err, "failed to get validator type")
	}

	return generateValidators(f.i, validatorType, f.Validators), nil
}

func generatePlanModifiers(imports *Imports, planModiferType *j.Statement, pm []string) *j.Statement {
	v := make([]jen.Code, len(pm))
	for i, n := range pm {
		v[i] = j.Id(imports.WithType(n))
	}

	return j.Index().Add(planModiferType).Values(v...)
}

func generateValidators(imports *Imports, validatorType *j.Statement, vals []string) *j.Statement {
	v := make([]jen.Code, len(vals))
	for i, n := range vals {
		v[i] = j.Id(imports.WithType(n))
	}

	return j.Index().Add(validatorType).Values(v...)
}

var baseByAttributeType = map[string]string{
	Schema + ".StringAttribute":  ".String",
	Schema + ".BoolAttribute":    ".Bool",
	Schema + ".Int64Attribute":   ".Int64",
	Schema + ".Float64Attribute": ".Float64",
	Schema + ".ListAttribute":    ".List",
	Schema + ".MapAttribute":     ".Map",
	Schema + ".ObjectAttribute":  ".Object",
}

func baseTypeForAttributeType(t string) (string, error) {
	baseType, ok := baseByAttributeType[t]
	if !ok {
		return "", trace.BadParameter("unexpected attribute type %q", t)
	}
	return baseType, nil
}

func (f *FieldSchemaGenerator) planModifierTypeForAttributeType(t string) (*j.Statement, error) {
	baseType, ok := baseByAttributeType[t]
	if !ok {
		return nil, trace.BadParameter("unexpected attribute type %q", t)
	}

	return j.Id(f.i.WithType(PlanModifier + baseType)), nil
}

func (f *FieldSchemaGenerator) validatorTypeForAttributeType(t string) (*j.Statement, error) {
	baseType, ok := baseByAttributeType[t]
	if !ok {
		return nil, trace.BadParameter("unexpected attribute type %q", t)
	}

	return j.Id(f.i.WithType(Validator + baseType)), nil
}
