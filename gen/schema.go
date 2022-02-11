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

package gen

import (
	// . significantly improves readability of a generator statements.
	// That's also the reason to extract it into the separate package.

	"github.com/dave/jennifer/jen"
	. "github.com/dave/jennifer/jen"
	"github.com/gravitational/protoc-gen-terraform/desc"
)

// MessageSchemaGenerator is the decorator struct to generate tfsdk.Schema of a message
type MessageSchemaGenerator struct {
	*desc.Message
	i *desc.Imports
}

// NewMessageSchemaGenerator returns new MessageSchemaGenerator struct
func NewMessageSchemaGenerator(m *desc.Message, i *desc.Imports) *MessageSchemaGenerator {
	return &MessageSchemaGenerator{Message: m, i: i}
}

// Generate returns Go code for message schema
func (m *MessageSchemaGenerator) Generate() []byte {
	id := "GenSchema" + m.Name
	schema := m.i.WithPackage(SDK, "Schema")
	diags := m.i.WithPackage(Diag, "Diagnostics")
	attr := m.i.WithPackage(SDK, "Attribute")

	j := Commentf("// %v returns tfsdk.Schema definition for %v\n", id, m.Name).
		Func().
		Id(id).
		Params(Id("ctx").Id(m.i.GoString("context.Context", false))).
		Params( // return params
			Id(schema),
			Id(diags),
		).
		Block(
			Return(
				Id(schema).Values(
					Dict{
						Id("Attributes"): Map(String()).Id(attr).Values(
							m.fieldsDictSchema(),
						),
					},
				),
				Nil(),
			),
		)

	return []byte(j.GoString() + "\n")
}

// FieldsDictSchema reutrns jen.Dict of the generated message fields
func (m *MessageSchemaGenerator) fieldsDictSchema() Dict {
	d := Dict{}

	for _, f := range m.Fields {
		f := NewFieldSchemaGenerator(f, m.i)
		d[Lit(f.NameSnake)] = f.Generate()
	}

	return d
}

// FieldSchemaGenerator represents the decorator for Field code generation
type FieldSchemaGenerator struct {
	*desc.Field
	i *desc.Imports
}

// NewFieldSchemaGenerator returns new FieldSchemaGenerator struct
func NewFieldSchemaGenerator(f *desc.Field, i *desc.Imports) *FieldSchemaGenerator {
	return &FieldSchemaGenerator{Field: f, i: i}
}

// Generate returns field schema
func (f *FieldSchemaGenerator) Generate() *Statement {
	if f.Kind == desc.Custom {
		return Id("GenSchema" + f.Suffix).Call(Id("ctx"))
	}

	d := Dict{
		Id("Description"): Lit(f.RawComment),
		Id("Type"):        f.schemaType(), // nils are automatically omitted
		Id("Attributes"):  f.attributes(),
	}

	// Required/Optional
	if f.IsRequired {
		d[Id("Required")] = True()
	} else {
		d[Id("Optional")] = True()
	}

	// Sensitive
	if f.IsSensitive {
		d[Id("Sensitive")] = True()
	}

	// Computed
	if f.IsComputed {
		d[Id("Computed")] = True()
	}

	// Validators
	if len(f.Validators) > 0 {
		v := make([]jen.Code, len(f.Validators))
		for i, n := range f.Validators {
			v[i] = Id(f.i.GoString(n, false))
		}

		d[Id("Validators")] = Index().String().Values(v...)
	}

	// Validators
	if len(f.PlanModifiers) > 0 {
		v := make([]jen.Code, len(f.PlanModifiers))
		for i, n := range f.PlanModifiers {
			v[i] = Id(f.i.GoString(n, false))
		}

		d[Id("PlanModifiers")] = Index().Id(f.i.WithPackage(SDK, "AttributePlanModifier")).Values(v...)
	}

	return Values(d)
}

// SchemaType returns the schema Type field value
func (f *FieldSchemaGenerator) schemaType() *Statement {
	switch f.Kind {
	case desc.Primitive:
		return f.primitiveSchemaTypeDef()
	case desc.PrimitiveList:
		return Id(f.Type).Values(Dict{
			Id("ElemType"): f.primitiveSchemaTypeDef(),
		})
	case desc.PrimitiveMap:
		g := NewFieldSchemaGenerator(f.MapValueField, f.i)

		return Id(f.Type).Values(Dict{
			Id("ElemType"): g.primitiveSchemaTypeDef(),
		})
	}

	return nil
}

// Attributes returns a nested attribute definitions
func (f *FieldSchemaGenerator) attributes() *Statement {
	switch f.Kind {
	case desc.Nested:
		m := NewMessageSchemaGenerator(f.Message, f.i)

		return f.xNestedAttributes("Single", m)
	case desc.NestedMap:
		m := NewMessageSchemaGenerator(f.MapValueField.Message, f.i)

		return f.xNestedAttributes("Map", m)
	case desc.NestedList:
		m := NewMessageSchemaGenerator(f.Message, f.i)

		return f.xNestedAttributes("List", m)
	}
	return nil
}

// primitiveSchemaTypeDef returns the primitive type
func (f *FieldSchemaGenerator) primitiveSchemaTypeDef() *Statement {
	if f.IsTypeScalar {
		return Id(f.ElemType)
	}

	return Id(f.ElemType).Values()
}

// xNestedAttributes generates *NestedAttributes call
func (f *FieldSchemaGenerator) xNestedAttributes(typ string, m *MessageSchemaGenerator) *Statement {
	var options *Statement

	if typ != "Single" {
		options = Id(f.i.WithPackage(SDK, typ+"NestedAttributesOptions")).Values()
	}

	return Id(f.i.WithPackage(SDK, typ+"NestedAttributes")).Params(
		Map(String()).Id(f.i.WithPackage(SDK, "Attribute")).Values(m.fieldsDictSchema()),
		options,
	)
}
