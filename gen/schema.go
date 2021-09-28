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
	c GeneratorContext
}

// NewMessageSchemaGenerator returns new MessageSchemaGenerator struct
func NewMessageSchemaGenerator(m *desc.Message, c GeneratorContext) *MessageSchemaGenerator {
	return &MessageSchemaGenerator{Message: m, c: c}
}

// Generate returns Go code for message schema
func (m *MessageSchemaGenerator) Generate() []byte {
	id := "GenSchema" + m.Name

	j := Commentf("// %v returns tfsdk.Schema definition for %v\n", id, m.Name).
		Func().
		Id(id).
		Params(m.c.ParamCtx()).
		Params( // return params
			m.c.SDKSchema(),
			m.c.DiagDiagnostics(),
		).
		Block(
			Return(
				Add(m.c.SDKSchema()).Values(
					Dict{
						Id("Attributes"): Map(String()).Add(m.c.SDKAttribute()).Values(
							m.FieldsDictSchema(),
						),
					},
				),
				Nil(),
			),
		)

	return []byte(j.GoString() + "\n")
}

// FieldsDictSchema reutrns jen.Dict of the generated message fields
func (m *MessageSchemaGenerator) FieldsDictSchema() Dict {
	d := Dict{}

	for _, f := range m.Fields {
		f := NewFieldSchemaGenerator(f, m.c)
		d[Lit(f.NameSnake)] = f.Generate()
	}

	return d
}

// FieldSchemaGenerator represents the decorator for Field code generation
type FieldSchemaGenerator struct {
	*desc.Field
	c GeneratorContext
}

// NewFieldSchemaGenerator returns new FieldSchemaGenerator struct
func NewFieldSchemaGenerator(f *desc.Field, c GeneratorContext) *FieldSchemaGenerator {
	return &FieldSchemaGenerator{Field: f, c: c}
}

// Generate returns field schema
func (f *FieldSchemaGenerator) Generate() *Statement {
	if f.Kind == desc.Custom {
		return Id("GenSchema" + f.Suffix).Call(Id("ctx"))
	}

	d := Dict{
		Id("Description"): Lit(f.RawComment),
		Id("Type"):        f.SchemaType(), // nils are automatically omitted
		Id("Attributes"):  f.Attributes(),
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
			v[i] = Id(n)
		}

		d[Id("Validators")] = Index().String().Values(v...)
	}

	return Values(d)
}

// SchemaType returns the schema Type field value
func (f *FieldSchemaGenerator) SchemaType() *Statement {
	switch f.Kind {
	case desc.Primitive:
		return f.primitiveSchemaTypeDef()
	case desc.PrimitiveList:
		return Add(f.c.TFType("ListType")).Values(Dict{
			Id("ElemType"): f.primitiveSchemaTypeDef(),
		})
	case desc.PrimitiveMap:
		f := NewFieldSchemaGenerator(f.MapValueField, f.c)

		return Add(f.c.TFType("MapType")).Values(Dict{
			Id("ElemType"): f.primitiveSchemaTypeDef(),
		})
	}

	return nil
}

// Attributes returns a nested attribute definitions
func (f *FieldSchemaGenerator) Attributes() *Statement {
	switch f.Kind {
	case desc.Nested:
		m := NewMessageSchemaGenerator(f.Message, f.c)

		return f.xNestedAttributes("Single", m)
	case desc.NestedMap:
		m := NewMessageSchemaGenerator(f.MapValueField.Message, f.c)

		return f.xNestedAttributes("Map", m)
	case desc.NestedList:
		m := NewMessageSchemaGenerator(f.Message, f.c)

		return f.xNestedAttributes("List", m)
	}
	return nil
}

// primitiveSchemaTypeDef returns the primitive type
func (f *FieldSchemaGenerator) primitiveSchemaTypeDef() *Statement {
	switch f.TFSchemaType {
	case desc.Int64Type, desc.Float64Type, desc.StringType, desc.BoolType:
		return f.c.Qual(Types, f.TFSchemaType)
	default:
		return Id(f.TFSchemaType + "{}")
	}
}

// xNestedAttributes generates *NestedAttributes call
func (f *FieldSchemaGenerator) xNestedAttributes(typ string, m *MessageSchemaGenerator) *Statement {
	var options *Statement

	if typ != "Single" {
		options = f.c.Qual(TFSDK, typ+"NestedAttributesOptions").Values()
	}

	return f.c.Qual(TFSDK, typ+"NestedAttributes").Params(
		Map(String()).Add(f.c.SDKAttribute()).Values(m.FieldsDictSchema()),
		options,
	)
}
