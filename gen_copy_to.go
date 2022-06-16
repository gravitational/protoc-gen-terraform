package main

import (
	"io"

	j "github.com/dave/jennifer/jen"
)

const (
	errWriting = "Error writing value to Terraform"
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
	tf := j.Id("tf").Op("*").Id(m.i.WithPackage(Types, "Object"))
	obj := j.Id("obj").Id(m.i.WithType(m.GoType))
	diags := j.Var().Id("diags").Id(m.i.WithPackage(Diag, "Diagnostics"))
	ctx := j.Id("ctx").Id(m.i.WithPackage("context", "Context"))

	// func Copy<name>ToTerraform(ctx context.Context, tf types.Object, obj *apitypes.<name>)
	// ... statements for a fields
	method :=
		j.Commentf("// %v copies contents of the source Terraform object into a target struct\n", methodName).
			Func().Id(methodName).
			Params(ctx, obj, tf).
			Id(m.i.WithPackage(Diag, "Diagnostics")).
			BlockFunc(func(g *j.Group) {
				g.Add(diags)
				g.Id("tf.Null").Op("=").False()
				g.Id("tf.Unknown").Op("=").False()
				g.If(j.Id("tf.Attrs").Op("==").Nil()).Block(
					j.Id("tf.Attrs").Op("=").Make(j.Map(j.String()).Id(m.i.WithPackage(Attr, "Value"))),
				)
				m.GenerateFields(g)
				g.Return(j.Id("diags"))
			})

	return writer.Write([]byte(method.GoString() + "\n"))
}

// GenerateFields generates specific statements for CopyToTF<name> methods
func (m *MessageCopyToGenerator) GenerateFields(g *j.Group) {
	for _, f := range m.Fields {
		g.Add(NewFieldCopyToGenerator(f, m.i).Generate())
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
		// _, ok := ft.AttrsTypes["key"]
		j.List(j.Id(v), j.Id("ok")).Op(":=").Id("tf.AttrTypes").Index(j.Lit(f.NameSnake)),
		j.If(j.Id("!ok")).BlockFunc(f.errAttrMissingDiag).Else().BlockFunc(g),
	)
}

// getAttr v, ok := tf.Attrs["name"]
func (f *FieldCopyToGenerator) getAttr(v string, typ string, g *j.Group) {
	g.List(
		j.Id(v), j.Id("ok"),
	).Op(":=").Id("tf.Attrs").Index(
		j.Lit(f.Field.NameSnake),
	).Assert(j.Id(f.i.WithType(typ)))
}

// genZeroValue generates zero value from an empty AttrType
func (f *FieldCopyToGenerator) genZeroValue(fieldName string) func(*j.Group) {
	return func(g *j.Group) {
		// This generates an empty attr.Value from a Terraform type
		// v, err = t.ValueFromTerraform(ctx, tftypes.NewValue(t.TerraformType(ctx, nil)))
		g.List(j.Id("i"), j.Id("err")).Op(":=").Id("t.ValueFromTerraform").Call(
			j.Id("ctx"),
			j.Id(f.i.WithPackage(TFTypes, "NewValue")).Call(
				j.Id("t.TerraformType").Call(j.Id("ctx")), j.Nil(),
			),
		)

		// if err != nil { diags.AddError }
		g.If(j.Id("err != nil")).Block(
			j.Id("diags.Append").Call(j.Id("attrWriteGeneralError").Values(j.Lit(f.Path), j.Id("err"))),
		)

		// v, ok = i.(types.Time)
		g.List(j.Id("v"), j.Id("ok")).Op("=").Id("i").Assert(j.Id(f.i.WithType(f.ElemValueType)))

		// if !ok { diags.AddError }
		g.If(j.Id("!ok")).BlockFunc(f.errAttrConversionFailure(f.Path, f.ElemValueType))

		// v.Null = v.Value == ""
		if f.ZeroValue != "" {
			g.Id("v.Null").Op("=").Id(f.i.WithType(f.ValueCastToType)).Parens(j.Id(fieldName)).Op("==").Id(f.ZeroValue)
		} else {
			g.Id("v.Null").Op("=").False()
		}
	}
}

// genPrimitiveBody generates block which reads object field into v
func (f *FieldCopyToGenerator) genPrimitiveBody(fieldName string, g *j.Group) {
	f.getAttr("v", f.i.WithType(f.Field.ElemValueType), g)
	g.If(j.Id("!ok")).BlockFunc(f.genZeroValue(fieldName))

	if f.IsNullable {
		g.If(j.Id(fieldName).Op("==").Nil()).Block(
			j.Id("v.Null").Op("=").True(),
		).Else().Block(
			j.Id("v.Null").Op("=").False(),
			j.Id("v.Value").Op("=").Id(f.i.WithType(f.GoElemTypeIndirect)).Parens(j.Op("*").Add(j.Id(fieldName))),
		)
	} else {
		// Non-nullable fields always have value
		g.Id("v.Value").Op("=").Id(f.i.WithType(f.ValueCastToType)).Parens(j.Id(fieldName))
	}

	g.Id("v.Unknown").Op("=").False()
}

// genObjectBody generates block which reads message into v
func (f *FieldCopyToGenerator) genObjectBody(m *MessageCopyToGenerator, fieldName string, typ string, g *j.Group) {
	copyObj := func(g *j.Group) {
		g.Id("obj").Op(":=").Id(fieldName)
		g.Id("tf").Op(":=").Id("&v")
		m.GenerateFields(g)
	}

	f.getAttr("v", f.Field.ElemValueType, g)
	g.If(j.Id("!ok")).Block(
		// v := types.Object{Attrs: make(map[string]attr.Value, len(o.AttrTypes)), AttrTypes: o.AttrTypes}
		j.Id("v").Op("=").Id(f.i.WithType(typ)).Block(j.Dict{
			j.Id("Attrs"):     j.Make(j.Map(j.String()).Id(f.i.WithPackage(Attr, "Value")), j.Len(j.Id("o.AttrTypes"))),
			j.Id("AttrTypes"): j.Id("o.AttrTypes"),
		}),
	).Else().Block(
		j.If(j.Id("v.Attrs").Op("==").Nil()).Block(
			j.Id("v.Attrs").Op("=").Make(j.Map(j.String()).Id(f.i.WithPackage(Attr, "Value")), j.Len(j.Id("tf.AttrTypes"))),
		),
	)
	if f.IsNullable {
		// if obj.Nested == nil
		g.If(j.Id(fieldName).Op("==").Nil()).Block(
			j.Id("v.Null").Op("=").True(),
		).Else().BlockFunc(
			copyObj,
		)
	} else {
		g.BlockFunc(copyObj)
	}
	g.Id("v.Unknown").Op("=").False()
}

// assertTo asserts a to typ
func (f *FieldCopyToGenerator) assertTo(typ string, g *j.Group, els func(g *j.Group)) {
	// v, ok := a.(types.ListType)
	g.List(j.Id("o"), j.Id("ok")).Op(":=").Id("a").Assert(j.Id(f.i.WithType(typ)))
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

	return f.nextField("t", func(g *j.Group) {
		if f.OneOfName != "" {
			f.genOneOfStub(g)
		}

		f.genPrimitiveBody(fieldName, g)
		g.Id("tf.Attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
	})
}

// genObject generates CopyTo statement for a nested message
func (f *FieldCopyToGenerator) genObject() *j.Statement {
	m := NewMessageCopyToGenerator(f.Message, f.i)
	fieldName := "obj." + f.Name

	return f.nextField("a", func(g *j.Group) {
		if f.OneOfName != "" {
			f.genOneOfStub(g)
		}

		f.assertTo(f.Field.ElemType, g, func(g *j.Group) {
			f.genObjectBody(m, fieldName, f.Field.ValueType, g)
			g.Id("tf.Attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
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

	var mk j.Code

	if f.IsMap {
		// make(map[string]attr.Value, len(obj.Map))
		mk = j.Make(j.Map(j.String()).Id(f.i.WithPackage(Attr, "Value")), j.Len(j.Id(fieldName)))
	}
	if f.IsRepeated {
		// make(map[string]attr.Value, len(obj.List))
		mk = j.Make(j.Index().Id(f.i.WithPackage(Attr, "Value")), j.Len(j.Id(fieldName)))
	}

	return f.nextField("a", func(g *j.Group) {
		f.assertTo(f.Field.Type, g, func(g *j.Group) {
			f.getAttr("c", f.Field.ValueType, g)

			g.If(j.Id("!ok")).Block(
				// c := types.Object{Elems: make([]attr.Value, ElemType: o.ElemType}
				j.Id("c").Op("=").Id(f.i.WithType(f.Field.ValueType)).Block(j.Dict{
					j.Id("Elems"):    mk,
					j.Id("ElemType"): j.Id("o.ElemType"),
					j.Id("Null"):     j.True(),
				}),
			).Else().Block(
				j.If(j.Id("c.Elems").Op("==").Nil()).Block(
					j.Id("c.Elems").Op("=").Add(mk),
				),
			)

			g.If(j.Id(fieldName)).Op("!=").Nil().BlockFunc(func(g *j.Group) {
				if (f.Kind == PrimitiveListKind) || (f.Kind == PrimitiveMapKind) {
					g.Id("t").Op(":=").Id("o.ElemType")
				} else {
					g.Id("o").Op(":=").Id("o.ElemType").Assert(j.Id(f.i.WithType(f.ElemType)))
				}
				// for k, a := range obj.List
				g.For(j.List(j.Id("k"), j.Id("a"))).Op(":=").Range().Id(fieldName).BlockFunc(func(g *j.Group) {
					if (f.Kind == PrimitiveListKind) || (f.Kind == PrimitiveMapKind) {
						f.genPrimitiveBody("a", g)
					} else {
						m := NewMessageCopyToGenerator(f.getValueField().Message, f.i)
						f.genObjectBody(m, "a", f.i.WithType(f.Field.ElemValueType), g)
					}
					g.Id("c.Elems").Index(j.Id("k")).Op("=").Id("v")
				})

				// if len(obj.Test) > 0
				g.If(j.Len(j.Id(fieldName))).Op(">").Lit(0).Block(
					j.Id("c.Null").Op("=").False(),
				)
			})

			g.Id("c.Unknown").Op("=").False()
			g.Id("tf.Attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("c")
		})
	})
}

// genCustom generates statement representing custom type
func (f *FieldCopyToGenerator) genCustom() *j.Statement {
	return f.nextField("t", func(g *j.Group) {
		g.Id("v").Op(":=").Id("CopyTo"+f.Suffix).Params(
			j.Id("diags"), j.Id("obj."+f.Name), j.Id("t"), j.Id("tf.Attrs").Index(j.Lit(f.NameSnake)),
		)
		g.Id("tf.Attrs").Index(j.Lit(f.NameSnake)).Op("=").Id("v")
	})
}
