package main

import (
	"io"

	j "github.com/dave/jennifer/jen"
)

// MessageCopyToGenerator is the visitor struct to generate tfsdk.Schema of a message
type MessageCopyToGenerator struct {
	*Message
	i *Imports
}

// NewMessageCopyToGenerator returns new MessageCopyToGenerator struct
func NewMessageCopyToGenerator(m *Message, i *Imports) *MessageCopyToGenerator {
	return &MessageCopyToGenerator{Message: m, i: i}
}

// Generate generates CopyToTF<Name> method
func (m *MessageCopyToGenerator) Generate(writer io.Writer) (int, error) {
	methodName := "Copy" + m.Name + "ToTerraform"

	// func Copy<name>ToTerraform(ctx context.Context, obj *apitypes.<name>) (*types.Object, diag.Diagnostics)
	// ... statements for a fields
	method :=
		j.Commentf("// %v copies contents of source struct into a Terraform object.\n", methodName).
			Func().Id(methodName).
			Params(
				j.Id("ctx").Id(m.i.WithPackage("context", "Context")),
				j.Id("obj").Op("*").Id(m.i.WithType(m.GoType)),
				j.Id("curr").Op("*").Id(m.i.WithPackage(Types, "Object")),
			).
			Parens(
				j.List(
					j.Id(m.i.WithPackage(Types, "Object")),
					j.Id(m.i.WithPackage(Diag, "Diagnostics")),
				),
			).
			BlockFunc(func(g *j.Group) {
				// schema, diags := GenSchemaFoo(ctx)
				g.List(j.Id("schema"), j.Id("diags")).Op(":=").Id("GenSchema" + m.Name).Call(j.Id("ctx"))

				// if diags.HasError() {
				// 	return types.Object{}, diags
				// }
				g.If(j.Id("diags.HasError").Call()).Block(
					j.Return(
						j.Id(m.i.WithPackage(Types, "Object")).Values(),
						j.Id("diags"),
					),
				)

				// schemaObj := schema.Type().(types.ObjectType)
				g.Id("schemaObj").Op(":=").
					Id("schema.Type").Call().
					Assert(j.Id(m.i.WithPackage(Types, "ObjectType")))

				// if obj == nil {
				// 	return types.ObjectNull(schemaObj.AttrTypes), diags
				// }
				g.If(j.Id("obj").Op("==").Nil()).Block(
					j.Return(
						j.Id(m.i.WithPackage(Types, "ObjectNull")).Call(j.Id("schemaObj.AttrTypes")),
						j.Id("diags"),
					),
				)

				// var attrs map[string]attr.Value
				// if curr == nil || curr.Attributes() == nil {
				// 	attrs = make(map[string]attr.Value)
				// }
				g.Var().Id("attrs").Map(j.String()).Id(m.i.WithPackage(Attr, "Value"))
				g.If(j.Id("curr").Op("==").Nil().Op("||").Id("curr.Attributes").Call().Op("==").Nil()).
					Block(
						j.Id("attrs").Op("=").Make(j.Map(j.String()).Id(m.i.WithPackage(Attr, "Value"))),
					).
					Else().
					Block(
						j.Id("attrs").Op("=").Id("curr.Attributes").Call(),
					)

				m.GenerateFields(g)

				// result, resultDiags := types.ObjectValue(schemaObj.AttrTypes, attrs)
				g.List(j.Id("result"), j.Id("resultDiags")).Op(":=").
					Id(m.i.WithPackage(Types, "ObjectValue")).Call(j.Id("schemaObj.AttrTypes"), j.Id("attrs"))

				// diags.Append(resultDiags)
				g.Id("diags.Append").Call(j.Id("resultDiags").Op("..."))

				// return result, diags
				g.Return(j.Id("result"), j.Id("diags"))
			})

	return writer.Write([]byte(method.GoString() + "\n"))
}

// GenerateFields generates specific statements for CopyToTF<name> methods
func (m *MessageCopyToGenerator) GenerateFields(g *j.Group) {
	for _, f := range m.Fields {
		g.Add(NewFieldCopyToGenerator(f, m.i).Generate())
	}
	for _, f := range m.InjectedFields {
		// if _, ok := attrs["foo"]; !ok {
		// 	attrs["foo"] = types.NullString()
		// }
		g.If(
			j.List(j.Id("_"), j.Id("ok")).Op(":=").
				Id("attrs").Index(j.Lit(f.Name))).Op(";").
			Op("!").Id("ok").
			Block(
				j.Id("attrs").Index(j.Lit(f.Name)).
					Op("=").
					Id(m.i.WithType(f.ValueMethod)).Call(),
			)
	}
}

// FieldCopyToGenerator is a visitor for a field
type FieldCopyToGenerator struct {
	*Field
	i *Imports
}

// NewFieldCopyToGenerator returns new FieldCopyToGenerator struct
func NewFieldCopyToGenerator(f *Field, i *Imports) *FieldCopyToGenerator {
	return &FieldCopyToGenerator{Field: f, i: i}
}

// errMissingDiag diags.Append(attrMissingDiag{path})
func (f *FieldCopyToGenerator) errAttrMissingDiag(g *j.Group) {
	g.Id("diags.Append").Call(
		j.Id("attrWriteMissingDiag").Values(j.Lit(f.Path)),
	)
}

// errAttrConversionFailure diags.Append(attrConversionFailureDiag{path, typ})
func (f *FieldCopyToGenerator) errAttrConversionFailure(path string, typ string) func(g *j.Group) {
	return func(g *j.Group) {
		g.Id("diags.Append").Call(
			j.Id("attrWriteConversionFailureDiag").Values(j.Lit(path), j.Lit(typ)),
		)
	}
}

// Generate generates CopyTo fragment for a field of different kind
func (f *FieldCopyToGenerator) Generate() *j.Statement {
	switch f.Kind {
	case PrimitiveKind:
		return f.genPrimitive()
	case ObjectKind:
		return f.genObject()
	case PrimitiveListKind, PrimitiveMapKind, ObjectListKind, ObjectMapKind:
		return f.genListOrMap()
	case CustomKind:
		return f.genCustom()
	}
	return nil
}

// genPrimitiveBody generates block which reads object field into v
func (f *FieldCopyToGenerator) genPrimitiveBody(fieldName string, g *j.Group) {
	// var v attr.Value
	g.Var().Id("v").Id(f.i.WithPackage(Attr, "Value"))

	if f.IsPlaceholder {
		// Set placeholder fields (for empty structs) to their null value.
		g.Id("v").Op("=").Id(f.i.WithType(f.NullValueMethod)).Call()
		return
	}

	if f.ParentIsOptionalEmbed {
		g.If(j.Id("obj." + f.ParentIsOptionalEmbedFieldName).Op("==").Nil()).Block(
			j.Id("v").Op("=").Id(f.i.WithType(f.NullValueMethod)).Call(),
		).Else().BlockFunc(func(g *j.Group) {
			f.genAssignValue(g, fieldName)
		})
	} else {
		f.genAssignValue(g, fieldName)
	}
}

func (f *FieldCopyToGenerator) genAssignValue(g *j.Group, fieldName string) {
	if f.IsNullable {
		// if obj.Foo == nil {
		// 	v = types.NullString()
		// } else {
		// 	v = types.StringValue(string(*obj.Foo))
		// }
		g.If(j.Id(fieldName).Op("==").Nil()).Block(
			j.Id("v").Op("=").Id(f.i.WithType(f.NullValueMethod)).Call(),
		).Else().Block(
			j.Id("v").Op("=").
				Id(f.i.WithType(f.ValueToMethod)).
				Call(
					j.Id(f.i.WithType(f.GoElemTypeIndirect)).Parens(j.Op("*").Add(j.Id(fieldName))),
				),
		)
		return
	}

	// val := string(obj.Foo)
	g.Id("val").Op(":=").
		Id(f.i.WithType(f.ValueCastToType)).Parens(j.Id(fieldName))

	if f.ZeroValue == "" {
		// v = types.StringValue(val)
		g.Id("v").Op("=").Id(f.i.WithType(f.ValueToMethod)).Call(j.Id("val"))
	} else {
		// For non-nullable fields, treat the zero value as null.
		//
		// if val == "" {
		g.If(j.Id("val").Op("==").Id(f.ZeroValue)).
			Block(
				// v = types.StringNull()
				j.Id("v").Op("=").Id(f.i.WithType(f.NullValueMethod)).Call(),
			).
			// } else {
			Else().
			Block(
				// v = types.StringValue(val)
				j.Id("v").Op("=").Id(f.i.WithType(f.ValueToMethod)).Call(j.Id("val")),
			)
	}
}

// genObjectBody generates block which reads message into v
func (f *FieldCopyToGenerator) genObjectBody(m *MessageCopyToGenerator, fieldName string, typ string, g *j.Group) {
	// Wrap object conversion in an anonymous function so we don't shadow the
	// parent object's attrs map etc.
	//
	// v := func() attr.Value { ... }()
	g.Id("v").Op(":=").Func().
		Params().
		Id(f.i.WithPackage(Attr, "Value")).
		BlockFunc(func(g *j.Group) {
			if f.OneOfName != "" {
				f.genOneOfStub(g)
			}

			if f.IsNullable {
				// if obj.Nested == nil
				g.If(j.Id(fieldName).Op("==").Nil()).Block(
					// return types.ObjectNull(schemaObj.AttrTypes)
					j.Return(j.Id(f.i.WithPackage(Types, "ObjectNull")).Call(j.Id("schemaObj.AttrTypes"))),
				)
			}

			// attrs := make(map[string]attr.Value)
			g.Id("attrs").Op(":=").Make(j.Map(j.String()).Id(f.i.WithPackage(Attr, "Value")))
			if len(m.Fields) > 0 {
				g.BlockFunc(func(g *j.Group) {
					if !m.IsEmpty {
						g.Id("obj").Op(":=").Id(fieldName)
					}
					m.GenerateFields(g)
				})
			}

			// result, objDiags := types.ObjectValue(schema.AttrTypes, attrs)
			g.List(
				j.Id("result"),
				j.Id("objDiags"),
			).Op(":=").Id(
				f.i.WithPackage(Types, "ObjectValue"),
			).Call(j.Id("schemaObj.AttrTypes"), j.Id("attrs"))

			// diags.Append(objDiags...)
			g.Id("diags.Append").Call(j.Id("objDiags").Op("..."))

			// return result
			g.Return(j.Id("result"))
		}).
		Call()
}

// getValueField returns list/map value field
func (f *FieldCopyToGenerator) getValueField() *Field {
	if f.IsMap {
		return f.MapValueField
	}

	return f.Field
}

// genPrimitive generates CopyTo statement for a primitive type
func (f *FieldCopyToGenerator) genPrimitive() *j.Statement {
	fieldName := "obj." + f.Name

	return j.BlockFunc(func(g *j.Group) {
		if f.OneOfName != "" {
			f.genOneOfStub(g)
		}

		f.genPrimitiveBody(fieldName, g)
		g.Id("attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
	})
}

// genObject generates CopyTo statement for a nested message
func (f *FieldCopyToGenerator) genObject() *j.Statement {
	m := NewMessageCopyToGenerator(f.Message, f.i)
	fieldName := "obj." + f.Name

	return j.BlockFunc(func(g *j.Group) {
		// schemaObj := schemaObj.Attributes["foo"].Type.(types.ObjectType)
		g.Id("schemaObj").Op(":=").
			Id("schemaObj.AttrTypes").
			Index(j.Lit(f.NameSnake)).
			Assert(j.Id(f.i.WithPackage(Types, "ObjectType")))

		f.genObjectBody(m, fieldName, f.Field.ValueType, g)
		g.Id("attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
	})
}

func (f *FieldCopyToGenerator) genOneOfStub(g *j.Group) {
	// {
	//     obj, ok := obj.OneOf.(*Test_Branch3)
	//     if !ok { obj = &Test_Branch3{} }
	// }
	g.List(j.Id("obj"), j.Id("ok")).Op(":=").Id("obj." + f.OneOfName).Assert(j.Id("*" + f.i.WithType(f.OneOfType)))
	g.If(j.Id("!ok")).Block(
		j.Id("obj").Op("=").Id("&" + f.i.WithType(f.OneOfType)).Values(),
	)
}

func (f *FieldCopyToGenerator) genListOrMap() *j.Statement {
	fieldName := "obj." + f.Name

	var makeElems, constructor, nullValue, elemType j.Code

	if f.IsMap {
		// elemType := schemaObj.AttrTypes["foo"].(types.MapType).ElemType
		elemType = j.Id("elemType").Op(":=").
			Id("schemaObj.AttrTypes").
			Index(j.Lit(f.NameSnake)).
			Assert(j.Id(f.i.WithPackage(Types, "MapType"))).
			Dot("ElemType")

		// make(map[string]attr.Value, len(obj.Map))
		makeElems = j.Make(j.Map(j.String()).Id(f.i.WithPackage(Attr, "Value")), j.Len(j.Id(fieldName)))

		// types.MapValue
		constructor = j.Id(f.i.WithPackage(Types, "MapValue"))

		// types.MapNull
		nullValue = j.Id(f.i.WithPackage(Types, "MapNull"))
	}

	if f.IsRepeated {
		// elemType := schemaObj.AttrTypes["foo"].(types.ListType).ElemType
		elemType = j.Id("elemType").Op(":=").
			Id("schemaObj.AttrTypes").
			Index(j.Lit(f.NameSnake)).
			Assert(j.Id(f.i.WithPackage(Types, "ListType"))).
			Dot("ElemType")

		// make([]attr.Value, len(obj.List))
		makeElems = j.Make(j.Index().Id(f.i.WithPackage(Attr, "Value")), j.Len(j.Id(fieldName)))

		// types.ListValue
		constructor = j.Id(f.i.WithPackage(Types, "ListValue"))

		// types.ListNull
		nullValue = j.Id(f.i.WithPackage(Types, "ListNull"))
	}

	return j.BlockFunc(func(g *j.Group) {
		g.Add(elemType)

		// var v attr.Value
		g.Var().Id("v").Id(f.i.WithPackage(Attr, "Value"))

		// if len(obj.Foo) == 0 {
		g.If(j.Id("len").Call(j.Id(fieldName)).Op("==").Lit(0)).Block(
			j.Id("v").Op("=").Add(nullValue).Call(j.Id("elemType")),
		).
			Else().
			BlockFunc(func(g *j.Group) {
				g.Id("elems").Op(":=").Add(makeElems)

				if f.Kind == ObjectListKind || f.Kind == ObjectMapKind {
					// Schema of the objects inside the list or map, referred to
					// by genObjectBody.
					//
					// schemaObj := elemType.(types.ObjectType)
					g.Id("schemaObj").Op(":=").
						Id("elemType").Assert(j.Id(f.i.WithPackage(Types, "ObjectType")))
				}

				// for k, a := range obj.List
				g.For(j.List(j.Id("k"), j.Id("a"))).Op(":=").Range().Id(fieldName).BlockFunc(func(g *j.Group) {
					if (f.Kind == PrimitiveListKind) || (f.Kind == PrimitiveMapKind) {
						f.genPrimitiveBody("a", g)
					} else {
						m := NewMessageCopyToGenerator(f.getValueField().Message, f.i)
						f.genObjectBody(m, "a", f.i.WithType(f.Field.ElemValueType), g)
					}
					g.Id("elems").Index(j.Id("k")).Op("=").Id("v")
				})

				// result, resultDiags := types.ListValue(elemType, elems)
				g.List(j.Id("result"), j.Id("resultDiags")).Op(":=").
					Add(constructor).
					Call(j.Id("elemType"), j.Id("elems"))

				// diags.Append(resultDiags...)
				g.Id("diags.Append").Call(j.Id("resultDiags").Op("..."))

				// v = result
				g.Id("v").Op("=").Id("result")
			})

		g.Id("attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
	})
}

// genCustom generates statement representing custom type
func (f *FieldCopyToGenerator) genCustom() *j.Statement {
	return j.BlockFunc(func(g *j.Group) {
		g.Id("v").Op(":=").Id("CopyTo"+f.Suffix).Params(
			j.Id("diags"), j.Id("obj."+f.Name),
		)
		g.Id("attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
	})
}
