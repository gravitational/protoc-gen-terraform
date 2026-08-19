package resource

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

// Generate generates Copy<Name>ToTerraform method.
// Unknown values are overridden.
func (m *MessageCopyToGenerator) Generate(writer io.Writer) (int, error) {
	methodName := "Copy" + m.Name + "ToTerraform"
	helperName := "Copy" + m.Name + "ToTerraformPreserveUnknown"

	// func Copy<name>ToTerraform(ctx context.Context, tf types.Object, obj *apitypes.<name>) (types.Object, diag.Diagnostics)
	// ... statements for a fields
	method :=
		j.Commentf("// %v copies contents of source struct into a Terraform object.\n", methodName).
			Func().Id(methodName).
			Params(
				j.Id("ctx").Id(m.i.WithPackage("context", "Context")),
				j.Id("obj").Op("*").Id(m.i.WithType(m.GoType)),
				j.Id("tf").Op("*").Id(m.i.WithPackage(Types, "Object")),
			).
			Parens(
				j.List(
					j.Id(m.i.WithPackage(Types, "Object")),
					j.Id(m.i.WithPackage(Diag, "Diagnostics")),
				),
			).
			BlockFunc(func(g *j.Group) {
				g.Return(j.Id(helperName).Call(j.Id("ctx"), j.Id("obj"), j.Id("tf"), j.False()))
			})

	return writer.Write([]byte(method.GoString() + "\n"))
}

// GeneratePreserveUnknown generates Copy<Name>ToTerraformPreserveUnknown method.
// If the `preserveUnknown` flag is enabled, Unknown values are preserved.
func (m *MessageCopyToGenerator) GeneratePreserveUnknown(writer io.Writer) (int, error) {
	methodName := "Copy" + m.Name + "ToTerraformPreserveUnknown"
	comment := "// %v copies contents of source struct into a Terraform object.\n" +
		"// Set preserveUnknown to true to preserve unknown values.\n"

	// func Copy<name>ToTerraformPreserveUnknown(ctx context.Context, tf types.Object, obj *apitypes.<name>, preserveUnknown bool) (types.Object, diag.Diagnostics)
	// ... statements for a fields

	method :=
		j.Commentf(comment, methodName).
			Func().Id(methodName).
			Params(
				j.Id("ctx").Id(m.i.WithPackage("context", "Context")),
				j.Id("obj").Op("*").Id(m.i.WithType(m.GoType)),
				j.Id("tf").Op("*").Id(m.i.WithPackage(Types, "Object")),
				j.Id("preserveUnknown").Id("bool"),
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

				// attrType := schema.Type().(types.ObjectType)
				g.Id("attrType").Op(":=").
					Id("schema.Type").Call().
					Assert(j.Id(m.i.WithPackage(Types, "ObjectType")))

				// var attrs map[string]attr.Value
				// if curr == nil || curr.Attributes() == nil {
				// 	attrs = make(map[string]attr.Value)
				// }
				g.Var().Id("attrs").Map(j.String()).Id(m.i.WithPackage(Attr, "Value"))
				g.If(j.Id("tf").Op("==").Nil().Op("||").Id("tf.Attributes").Call().Op("==").Nil()).
					Block(
						j.Id("attrs").Op("=").Make(j.Map(j.String()).Id(m.i.WithPackage(Attr, "Value"))),
					).
					Else().
					Block(
						j.Id("attrs").Op("=").Id("tf.Attributes").Call(),
					)

				m.GenerateFields(g)

				// result, resultDiags := types.ObjectValue(attrType.AttributeTypes(), attrs)
				g.List(j.Id("result"), j.Id("resultDiags")).Op(":=").
					Id(m.i.WithPackage(Types, "ObjectValue")).Call(j.Id("attrType").Dot("AttributeTypes").Call(), j.Id("attrs"))

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
		g.List(j.Id("_"), j.Id("ok")).Op(":=").Id("attrs").Index(j.Lit(f.Name))
		g.If(j.Op("!").Id("ok")).Block(
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

// nextField reads current field value from Terraform object and asserts it's type against expected
func (f *FieldCopyToGenerator) nextField(v string, g func(g *j.Group)) *j.Statement {
	return j.Block(
		// _, ok := ft.AttributeTypes(ctx)["key"]
		j.List(j.Id(v), j.Id("ok")).Op(":=").Id("attrType").Dot("AttributeTypes").Call().Index(j.Lit(f.NameSnake)),
		j.If(j.Id("!ok")).BlockFunc(f.errAttrMissingDiag).Else().BlockFunc(g),
	)
}

// genPrimitiveBody generates a block statement that reads an object field into
// variable "v".
func (f *FieldCopyToGenerator) genPrimitiveBody(g *j.Group, fieldName string) {
	g.If(
		j.Id("preserveUnknown").Op("&&").Id("existing").Op("!=").Nil().Op("&&").Id("existing").Dot("IsUnknown").Call(),
	).Block(
		j.Id("v").Op("=").Id(f.i.WithType(f.UnknownValueMethod)).Call(),
	).Else().Block(
		f.genAssignValue(fieldName),
	)
}

func (f *FieldCopyToGenerator) genAssignValue(fieldName string) *j.Statement {
	switch {
	case f.IsPlaceholder:
		return f.genAssignPlaceholderValue()
	case f.ParentIsOptionalEmbed:
		return f.genAssignOptionalEmbeddedValue(fieldName)
	case f.OneOfName != "":
		return f.genAssignOneOf(fieldName)
	case f.IsNullable, f.IsProto3Optional:
		return f.genAssignNullableValue(fieldName)
	default:
		return f.genAssignNonNullableValue(fieldName)
	}
}

// genAssignPlaceholderValue generates a statement to assign a placeholder value.
//
// Expected format:
//
//	v.Null = true
func (f *FieldCopyToGenerator) genAssignPlaceholderValue() *j.Statement {
	return j.Id("v").Op("=").Id(f.i.WithType(f.NullValueMethod)).Call()
}

// genAssignOptionalEmbeddedValue generates a statement to assign an optional embedded value.
//
// Expected format:
//
//	if obj.<f.ParentIsOptionalEmbedFieldName> == nil {
//		v.Null = true
//	} else {
//		v.Null = false
//		v.Value = <f.ValueCastToType>(<fieldName>)
//	}
func (f *FieldCopyToGenerator) genAssignOptionalEmbeddedValue(fieldName string) *j.Statement {
	return j.If(j.Id("obj." + f.ParentIsOptionalEmbedFieldName).Op("==").Nil()).Block(
		j.Id("v").Op("=").Id(f.i.WithType(f.NullValueMethod)).Call(),
	).Else().Block(
		j.Id("v").Op("=").Id(f.i.WithType(f.ValueToMethod)).Call(j.Id(f.i.WithType(f.ValueCastToType)).Parens(j.Id(fieldName))),
	)
}

// genAssignNullableValue generates a statement to assign a nullable value.
//
// Expected format:
//
//	if <fieldName> == nil {
//		v.Null = true
//	} else {
//		v.Null = false
//		v.Value = <f.GoElemTypeIndirect>(*<fieldName>)
//	}
func (f *FieldCopyToGenerator) genAssignNullableValue(fieldName string) *j.Statement {
	return j.If(j.Id(fieldName).Op("==").Nil()).Block(
		j.Id("v").Op("=").Id(f.i.WithType(f.NullValueMethod)).Call(),
	).Else().Block(
		j.Id("v").Op("=").Id(f.i.WithType(f.ValueToMethod)).Call(j.Id(f.i.WithType(f.GoElemTypeIndirect)).Parens(j.Op("*").Add(j.Id(fieldName)))),
	)
}

// genAssignOneOf generates a statement to assign a oneof value.
//
// Expected format:
//
//	obj, ok := obj.<f.OneOfName>.(*f.OneOfType)
//	if !ok {
//		v.Null = true
//	} else {
//		v.Null = false
//		v.Value = <f.ValueCastToType>(<fieldName>)
//	}
func (f *FieldCopyToGenerator) genAssignOneOf(fieldName string) *j.Statement {
	return j.Block(
		j.List(j.Id("obj"), j.Id("ok")).Op(":=").Id("obj."+f.OneOfName).Assert(j.Id("*"+f.i.WithType(f.OneOfType))),
		j.If(j.Id("!ok")).Block(
			j.Id("v").Op("=").Id(f.i.WithType(f.NullValueMethod)).Call(),
		).Else().Block(
			j.Id("v").Op("=").Id(f.i.WithType(f.ValueToMethod)).Call(j.Id(f.i.WithType(f.ValueCastToType)).Parens(j.Id(fieldName))),
		),
	)
}

// genAssignNonNullableValue generates a statement to assign a non-nullable value.
//
// Expected format:
//
//	v.Null = false
//	v.Value = <f.ValueCastToType>(<fieldName>)
func (f *FieldCopyToGenerator) genAssignNonNullableValue(fieldName string) *j.Statement {
	return j.Id("v").Op("=").Id(f.i.WithType(f.ValueToMethod)).Call(j.Id(f.i.WithType(f.ValueCastToType)).Parens(j.Id(fieldName)))
}

// genObjectBody generates block statement that reads a message into
// variable "v".
func (f *FieldCopyToGenerator) genObjectBody(g *j.Group, m *MessageCopyToGenerator, fieldName string, typ string) {
	copyObj := func(g *j.Group) {
		if len(m.Fields) > 0 {
			if !m.IsEmpty {
				g.Id("obj").Op(":=").Id(fieldName)
			}

			g.Id("attrs").Op(":=").Make(j.Map(j.String()).Id(f.i.WithPackage(Attr, "Value")), j.Len(j.Id("attrType.AttributeTypes").Call()))

			g.List(j.Id("v"), j.Id("ok")).
				Op(":=").
				Id("existing").
				Assert(j.Id(f.i.WithType(typ)))

			g.If(j.Id("ok").Op("&&").Id("v").Dot("Attributes").Call().Op("!=").Nil()).Block(
				j.Id("attrs").Op("=").Id("v").Dot("Attributes").Call(),
			)

			m.GenerateFields(g)

			// result, resultDiags := types.ObjectValue(attrType.AttributeTypes(), attrs)
			g.List(j.Id("result"), j.Id("resultDiags")).Op(":=").
				Id(m.i.WithPackage(Types, "ObjectValue")).Call(j.Id("attrType").Dot("AttributeTypes").Call(), j.Id("attrs"))

			// diags.Append(resultDiags)
			g.Id("diags.Append").Call(j.Id("resultDiags").Op("..."))

			g.Return(j.Id("result"))
		}
	}

	g.Id("v").Op(":=").Func().Params().
		Id(f.i.WithPackage(Attr, "Value")).
		BlockFunc(func(g *j.Group) {
			g.If(
				j.Id("preserveUnknown").Op("&&").Id("existing").Op("!=").Nil().Op("&&").Id("existing").Dot("IsUnknown").Call(),
			).Block(
				j.Return(j.Id(f.i.WithPackage(Types, "ObjectUnknown")).Call(j.Id("attrType").Dot("AttributeTypes").Call())),
			).Else().BlockFunc(func(g *j.Group) {
				if f.IsNullable {
					// if obj.Nested == nil
					g.If(j.Id(fieldName).Op("==").Nil()).Block(
						j.Return(j.Id(f.i.WithPackage(Types, "ObjectNull")).Call(j.Id("attrType").Dot("AttributeTypes").Call())),
					).Else().BlockFunc(
						copyObj,
					)
				} else {
					g.BlockFunc(copyObj)
				}
			})
		}).Call()
}

// assertTo asserts a to typ
func (f *FieldCopyToGenerator) assertTo(typ string, g *j.Group, els func(g *j.Group)) {
	// v, ok := a.(types.ListType)
	g.List(j.Id("attrType"), j.Id("ok")).Op(":=").Id("attrType").Assert(j.Id(f.i.WithType(typ)))
	g.If(j.Id("!ok")).BlockFunc(
		f.errAttrConversionFailure(f.Path, f.Field.Type),
	).Else().BlockFunc(els)
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
		g.Var().Id("v").Id(f.i.WithPackage(Attr, "Value"))
		g.Id("existing").Op(":=").Id("attrs").Index(j.Lit(f.Field.NameSnake))
		f.genPrimitiveBody(g, fieldName)
		g.Id("attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
	})
}

// genObject generates CopyTo statement for a nested message
func (f *FieldCopyToGenerator) genObject() *j.Statement {
	m := NewMessageCopyToGenerator(f.Message, f.i)
	fieldName := "obj." + f.Name

	return f.nextField("attrType", func(g *j.Group) {
		if f.OneOfName != "" {
			f.genOneOfStub(g)
		}

		f.assertTo(f.Field.ElemType, g, func(g *j.Group) {
			g.Id("existing").Op(":=").Id("attrs").Index(j.Lit(f.Field.NameSnake))
			f.genObjectBody(g, m, fieldName, f.Field.ValueType)
			g.Id("attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
		})
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

	var makeValue, makeUnknown, makeElems j.Code

	if f.IsMap {
		makeValue = j.Id(f.i.WithPackage(Types, "MapValue"))
		makeUnknown = j.Id(f.i.WithPackage(Types, "MapUnknown"))
		makeElems = j.Make(j.Map(j.String()).Id(f.i.WithPackage(Attr, "Value")), j.Len(j.Id(fieldName)))
	}

	if f.IsRepeated {
		makeValue = j.Id(f.i.WithPackage(Types, "ListValue"))
		makeUnknown = j.Id(f.i.WithPackage(Types, "ListUnknown"))
		makeElems = j.Make(j.Index().Id(f.i.WithPackage(Attr, "Value")), j.Len(j.Id(fieldName)))
	}

	return f.nextField("attrType", func(g *j.Group) {
		f.assertTo(f.Field.Type, g, func(g *j.Group) {
			g.Var().Id("v").Id(f.i.WithPackage(Attr, "Value"))
			g.Id("existing").Op(":=").Id("attrs").Index(j.Lit(f.Field.NameSnake))

			g.If(
				j.Id("preserveUnknown").Op("&&").Id("existing").Op("!=").Nil().Op("&&").Id("existing").Dot("IsUnknown").Call(),
			).Block(
				j.Id("v").Op("=").Add(makeUnknown).Call(j.Id("attrType").Dot("ElementType").Call()),
			).Else().BlockFunc(func(g *j.Group) {
				g.Id("oldElems").Op(":=").Add(makeElems)

				g.List(j.Id("c"), j.Id("ok")).
					Op(":=").
					Id("existing").
					Assert(j.Id(f.i.WithType(f.Field.ValueType)))

				g.If(j.Id("ok").Op("&&").Id("c").Dot("Elements").Call().Op("!=").Nil()).Block(
					j.Id("oldElems").Op("=").Id("c").Dot("Elements").Call(),
				)

				g.Id("elems").Op(":=").Add(makeElems)

				if f.IsRepeated {
					g.Id("copy").Call(j.Id("elems"), j.Id("oldElems"))
				}

				// for k, a := range obj.List
				g.For(j.List(j.Id("k"), j.Id("a"))).Op(":=").Range().Id(fieldName).BlockFunc(func(g *j.Group) {
					readSelector := "elems"
					if f.IsMap {
						readSelector = "oldElems"
					}
					index := j.Id("k")

					switch f.Kind {
					case PrimitiveListKind, PrimitiveMapKind:
						g.Var().Id("v").Id(f.i.WithPackage(Attr, "Value"))
						g.Id("existing").Op(":=").Id(readSelector).Index(index)
						f.genPrimitiveBody(g, "a")
					default:
						m := NewMessageCopyToGenerator(f.getValueField().Message, f.i)
						g.Id("attrType").Op(":=").Id("attrType.ElementType").Call().Assert(j.Id(f.i.WithType(f.ElemType)))
						g.Id("existing").Op(":=").Id(readSelector).Index(index)
						f.genObjectBody(g, m, "a", f.i.WithType(f.Field.ElemValueType))
					}
					g.Id("elems").Index(j.Id("k")).Op("=").Id("v")
				})

				g.List(j.Id("result"), j.Id("resultDiags")).Op(":=").
					Add(makeValue).
					Call(j.Id("attrType").Dot("ElementType").Call(), j.Id("elems"))

				g.Id("diags.Append").Call(j.Id("resultDiags").Op("..."))

				g.Id("v").Op("=").Id("result")
			})

			g.Id("attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
		})
	})
}

// genCustom generates statement representing custom type
func (f *FieldCopyToGenerator) genCustom() *j.Statement {
	return f.nextField("t", func(g *j.Group) {
		g.Id("v").Op(":=").Id("CopyTo"+f.Suffix).Params(
			j.Id("diags"), j.Id("obj."+f.Name), j.Id("t"), j.Id("attrs").Index(j.Lit(f.NameSnake)), j.Id("preserveUnknown"),
		)
		g.Id("attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
	})
}
