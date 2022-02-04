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
	// SDK represents the path to Terraform SDK package
	SDK = "github.com/hashicorp/terraform-plugin-framework/tfsdk"
	// Types represents the path to Terraform types package
	Types = "github.com/hashicorp/terraform-plugin-framework/types"
	// Diag represents the path to Terraform diag package
	Diag = "github.com/hashicorp/terraform-plugin-framework/diag"
	// Attr represents the name of Terraform attr package
	Attr = "github.com/hashicorp/terraform-plugin-framework/attr"
)

// Generator interface represents generic generator interface
type Generator interface {
	Generate() []byte
}

// GeneratorContext represents the struct which shares common methods along the generator structs
type GeneratorContext struct {
	i desc.Imports
}

// NewGeneratorContext returns the new initialized context for Generators
func NewGeneratorContext(i desc.Imports) GeneratorContext {
	return GeneratorContext{i}
}

// Qual returns qualified name of a type
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

// Schema returns tfsdk.Schema qualifier
func (c GeneratorContext) Schema() *Statement {
	return c.Qual(SDK, "Schema")
}

// AttributePlanModifier returns tfsdk.AttributePlanModifier qualifier
func (c GeneratorContext) AttributePlanModifier() *Statement {
	return c.Qual(SDK, "AttributePlanModifier")
}

// Attribute returns tfsdk.Attribute qualifier
func (c GeneratorContext) Attribute() *Statement {
	return c.Qual(SDK, "Attribute")
}

// Types returns qualified Terraform SDK type name
func (c GeneratorContext) Types(name string) *Statement {
	return c.Qual(Types, name)
}

// AttrValue returns qualified Terraform SDK attr.Value
func (c GeneratorContext) AttrValue() *Statement {
	return c.Qual(Attr, "Value")
}

// func (c GeneratorContext) getType(f *desc.Field, suffix string) *Statement {
// 	if f.IsMap {
// 		return c.Types("Map" + suffix)
// 	}

// 	if f.IsRepeated {
// 		return c.Types("List" + suffix)
// 	}

// 	return c.getElemType(f, suffix)
// }

// // GetElemAttrValueType returns Terraform Elem type for a List/Map field in Attrs map
// func (c GeneratorContext) getElemType(f *desc.Field, suffix string) *Statement {
// 	if f.TerraformValueType != "" {
// 		return Id(f.TerraformValueType)
// 	}

// 	switch f.TerraformElemType {
// 	case desc.Int64Type:
// 		return c.Qual(Types, "Int64"+suffix)
// 	case desc.Float64Type:
// 		return c.Qual(Types, "Float64"+suffix)
// 	case desc.StringType:
// 		return c.Qual(Types, "String"+suffix)
// 	case desc.BoolType:
// 		return c.Qual(Types, "Bool"+suffix)
// 	}

// 	return c.Qual(Types, "Object"+suffix)
// }

// // GetAttrValueType returns Terraform attr.Value type for a field in Attrs map
// func (c GeneratorContext) GetAttrValueType(f *desc.Field) *Statement {
// 	return c.getType(f, "")
// }

// // GetElemAttrValueType returns Terraform Elem type for a List/Map field in Attrs map
// func (c GeneratorContext) GetElemAttrValueType(f *desc.Field) *Statement {
// 	return c.getElemType(f, "")
// }

// // GetAttrTypeType returns Terraform attr.Value type for a field in Attrs map
// func (c GeneratorContext) GetAttrTypeType(f *desc.Field) *Statement {
// 	return c.getType(f, "Type")
// }

// // GetElemAttrTypeType returns Terraform Elem type for a List/Map field in Attrs map
// func (c GeneratorContext) GetElemAttrTypeType(f *desc.Field) *Statement {
// 	return c.getType(f, "Type")
// }

// GetNestedValueType returns Terraform nested value Go type
func (c GeneratorContext) getNestedTypeName(f *desc.Field, suffix string) *Statement {
	if f.IsMap {
		return c.Types("Map" + suffix)
	}

	if f.IsRepeated {
		return c.Types("List" + suffix)
	}

	return c.Types("Object" + suffix)
}

// GetNestedValueType returns Terraform nested value Go type
func (c GeneratorContext) GetNestedValueType(f *desc.Field) *Statement {
	return c.getNestedTypeName(f, "")
}

// GetNestedType returns Terraform nested value Go type
func (c GeneratorContext) GetNestedType(f *desc.Field) *Statement {
	return c.getNestedTypeName(f, "Type")
}

// GetMapEl reads map element from the Terraform object into the variable.
//
// a, ok := o.Attrs["name"]; if !ok { return err }
func (c *GeneratorContext) GetMapEl(a, m, key, path string) Statement {
	return Statement{
		List(Id(a), Id("ok")).Op(":=").Id("tf." + m).Index(Lit(key)),
		If(Id("!ok")).
			Block(
				Return(Qual("fmt", "Errorf").Call(Lit(m + "[\"" + key + "\"] is missing in the Terraform object (" + path + ")"))),
			),
	}
}

// AssertTo returns the v.(to) statement
//
// v, ok := a.(to); if !ok { return err }
func (c *GeneratorContext) AssertTo(to *Statement, v, op, path string) Statement {
	return Statement{
		List(Id(v), Id("ok")).Op(op).Id("a").Assert(to),
		If(Id("!ok")).
			Block(
				Return(Qual("fmt", "Errorf").Call(Lit("Failed to convert " + path + " to " + to.GoString()))),
			),
	}
}

// // AssertToPrimitiveValueType generates v.(to) statement to primitive attr type
// func (c *GeneratorContext) AssertToPrimitiveValueType(f *desc.Field, op string) Statement {
// 	return c.AssertTo(c.GetAttrValueType(f), "v", op, f.Path)
// }

// AssertToNestedValueType generate v.(to) statement to attr nested type (List or Map)
func (c *GeneratorContext) AssertToNestedValueType(f *desc.Field, op string) Statement {
	return c.AssertTo(c.GetNestedValueType(f), "v", op, f.Path)
}

// AssertToObjectValueType generates v.(to) statement to attr.Object
func (c *GeneratorContext) AssertToObjectValueType(f *desc.Field, op string) Statement {
	return c.AssertTo(c.Types("Object"), "v", op, f.Path)
}
