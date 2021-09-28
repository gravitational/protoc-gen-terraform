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
	. "github.com/dave/jennifer/jen"
	"github.com/gravitational/protoc-gen-terraform/desc"
)

const (
	// TFSDK represents the path to Terraform SDK package
	TFSDK = "github.com/hashicorp/terraform-plugin-framework/tfsdk"
	// Types represents the path to Terraform types package
	Types = "github.com/hashicorp/terraform-plugin-framework/types"
	// Diag represents the path to Terraform diag package
	Diag = "github.com/hashicorp/terraform-plugin-framework/diag"
	// TFAttrURL represents the name of Terraform attr package
	Attr = "github.com/hashicorp/terraform-plugin-framework/attr"
)

// GeneratorContext represents the struct which shares common methods along the generator structs
//
// Conventions:
// * obj - API object
// * tf - Terraform object
// * a - a variable which stores attr.Value
// * v - a.(specific type)
// * el - a map or list element
type GeneratorContext struct {
	i desc.Imports
}

// NewGeneratorContext returns the new initialized context for Generators
func NewGeneratorContext(i desc.Imports) GeneratorContext {
	return GeneratorContext{i}
}

// Qual returns qualified name
func (c GeneratorContext) Qual(path string, name string) *Statement {
	return Qual(c.i.GetQual(path), name)
}

// ParamCtx returns ctx context.Context param
func (c GeneratorContext) ParamCtx() *Statement {
	return Id("ctx").Add(c.Qual("context", "Context"))
}

// DiagDiagnostics returns diag.Diagnostice qualifier
func (c GeneratorContext) DiagDiagnostics() *Statement {
	return c.Qual(Diag, "Diagnostics")
}

// SDKSchema returns tfsdk.Schema qualifier
func (c GeneratorContext) SDKSchema() *Statement {
	return c.Qual(TFSDK, "Schema")
}

// SDKAttribute returns tfsdk.Attribute qualifier
func (c GeneratorContext) SDKAttribute() *Statement {
	return c.Qual(TFSDK, "Attribute")
}

// TFType returns qualified Terraform SDK type name
func (c GeneratorContext) TFType(name string) *Statement {
	return c.Qual(Types, name)
}

// AttrValue returns qualified Terraform SDK attr.Value
func (c GeneratorContext) AttrValue() *Statement {
	return c.Qual(Attr, "Value")
}

// GetPrimitiveType returns Terraform Go type for a primitive field
func (c GeneratorContext) GetPrimitiveType(f *desc.Field) *Statement {
	switch f.TFSchemaType {
	case desc.Int64Type:
		return c.Qual(Types, "Int64")
	case desc.Float64Type:
		return c.Qual(Types, "Float64")
	case desc.StringType:
		return c.Qual(Types, "String")
	case desc.BoolType:
		return c.Qual(Types, "Bool")
	}

	if f.TFSchemaValueType != "" {
		return Id(f.TFSchemaValueType)
	}

	return nil
}

// CastValueToGoElemType generates type(v.Value) statement
func (c GeneratorContext) CastValueToGoElemType(f *desc.Field) *Statement {
	return Id(f.GoElemType).Parens(Id("v.Value"))
}

// GetNestedType returns Terraform nested value Go type
func (c GeneratorContext) GetNestedType(f *desc.Field) *Statement {
	if f.IsMap {
		return c.TFType("Map")
	}

	if f.IsRepeated {
		return c.TFType("List")
	}

	return c.TFType("Object")
}

// getTFAttr reads an attribute from the Terraform object and checks that attr exists
//
// a, ok := o.Attrs["name"]; if !ok { return err }
func (c *GeneratorContext) getTFAttr(name, path string) Statement {
	return Statement{
		List(Id("a"), Id("ok")).Op(":=").Id("tf.Attrs").Index(Lit(name)),
		If(Id("!ok")).
			Block(
				Return(Qual("fmt", "Errorf").Call(Lit("Attr " + name + " is missing in the Terraform object (" + path + ")"))),
			),
	}
}

// AssertTFAttrToPrimitiveType generates v.(to) statement to primitive attr type
func (c *GeneratorContext) AssertTFAttrToPrimitiveType(f *desc.Field) Statement {
	return c.AssertTFAttrTo(c.GetPrimitiveType(f), f.Path)
}

// AssertTFAttrToNestedType generate v.(to) statement to attr nested type (List or Map)
func (c *GeneratorContext) AssertTFAttrToNestedType(f *desc.Field) Statement {
	return c.AssertTFAttrTo(c.GetNestedType(f), f.Path)
}

// AssertTFAttrToObject generates v.(to) statement to attr.Object
func (c *GeneratorContext) AssertTFAttrToObject(f *desc.Field) Statement {
	return c.AssertTFAttrTo(c.TFType("Object"), f.Path)
}

// AssertTFAttrTo returns the v.(to) statement
//
// v, ok := a.(to); if !ok { return err }
func (c *GeneratorContext) AssertTFAttrTo(to *Statement, path string) Statement {
	return Statement{
		List(Id("v"), Id("ok")).Op(":=").Id("a").Assert(to),
		If(Id("!ok")).
			Block(
				Return(Qual("fmt", "Errorf").Call(Lit("Failed to convert " + path + " to " + to.GoString()))),
			),
	}
}
