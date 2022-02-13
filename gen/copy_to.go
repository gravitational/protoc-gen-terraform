package gen

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"github.com/gravitational/protoc-gen-terraform/desc"
)

const (
	errWriting = "Error writing value to Terraform"
)

// MessageCopyToGenerator is the visitor struct to generate tfsdk.Schema of a message
type MessageCopyToGenerator struct {
	*desc.Message
	i *desc.Imports
}

// NewMessageCopyToGenerator returns new MessageCopyToGenerator struct
func NewMessageCopyToGenerator(m *desc.Message, i *desc.Imports) *MessageCopyToGenerator {
	return &MessageCopyToGenerator{Message: m, i: i}
}

// Generate generates CopyToTF<Name> method
func (m *MessageCopyToGenerator) Generate() []byte {
	methodName := "Copy" + m.Name + "ToTerraform"
	tf := Id("tf").Op("*").Id(m.i.WithPackage(Types, "Object"))
	obj := Id("obj").Id(m.GoType)
	diags := Var().Id("diags").Id(m.i.WithPackage(Diag, "Diagnostics"))

	// func Copy<name>ToTerraform(tf types.Object, obj *apitypes.<name>)
	// ... statements for a fields
	method :=
		Commentf("// %v copies contents of the source Terraform object into a target struct\n", methodName).
			Func().Id(methodName).
			Params(obj, tf).
			Id(m.i.WithPackage(Diag, "Diagnostics")).
			BlockFunc(func(g *Group) {
				g.Add(diags)
				g.Id("tf.Null").Op("=").False()
				g.Id("tf.Unknown").Op("=").False()
				g.If(Id("tf.Attrs").Op("==").Nil()).Block(
					Id("tf.Attrs").Op("=").Make(Map(String()).Id(m.i.WithPackage(Attr, "Value"))),
				)
				m.GenerateFields(g)
				g.Return(Id("diags"))
			})

	return []byte(method.GoString() + "\n")
}

// GenerateFields generates specific statements for CopyToTF<name> methods
func (m *MessageCopyToGenerator) GenerateFields(g *Group) {
	for _, f := range m.Fields {
		g.Add(NewFieldCopyToGenerator(f, m.i).Generate())
	}
}

// FieldCopyToGenerator is a visitor for a field
type FieldCopyToGenerator struct {
	*desc.Field
	i *desc.Imports
}

// NewFieldCopyToGenerator returns new FieldCopyToGenerator struct
func NewFieldCopyToGenerator(f *desc.Field, i *desc.Imports) *FieldCopyToGenerator {
	return &FieldCopyToGenerator{Field: f, i: i}
}

// Generate generates CopyTo fragment for a field of different kind
func (f *FieldCopyToGenerator) Generate() *Statement {
	switch f.Kind {
	case desc.Primitive:
		return f.genPrimitive()
	case desc.Nested:
		return f.genNested()
	case desc.PrimitiveList, desc.PrimitiveMap, desc.NestedList, desc.NestedMap:
		return f.genListOrMap()
	case desc.Custom:
		return f.genCustom()
	}
	return nil
}

// nextField reads current field value from Terraform object and asserts it's type against expected
func (f *FieldCopyToGenerator) nextField(v string, g func(g *Group)) *Statement {
	return Block(
		// _, ok := ft.AttrsTypes["key"]
		List(Id(v), Id("ok")).Op(":=").Id("tf.AttrTypes").Index(Lit(f.NameSnake)),
		If(Id("!ok")).Block(
			Id("diags.AddError").Call(
				Lit(errWriting),
				Lit(fmt.Sprintf("A value type for %v is missing in the target Terraform object AttrTypes", f.Path)),
			),
		).Else().BlockFunc(g),
	)
}

// getAttr v, ok := tf.Attrs["name"]
func (f *FieldCopyToGenerator) getAttr(v string, typ string, g *Group) {
	g.List(Id(v), Id("ok")).Op(":=").Id("tf.Attrs").Index(Lit(f.Field.NameSnake)).Assert(Id(typ))
}

// genPrimitiveBody generates block which reads object field into v
func (f *FieldCopyToGenerator) genPrimitiveBody(fieldName string, g *Group) {
	f.getAttr("v", f.Field.ElemValueType, g)
	g.If(Id("!ok")).Block(
		Id("v").Op("=").Id(f.Field.ElemValueType).Values(),
	)

	if f.IsNullable {
		g.If(Id(fieldName).Op("==").Nil()).Block(
			Id("v.Null").Op("=").True(),
		).Else().Block(
			Id("v.Value").Op("=").Id(f.GoElemTypeIndirect).Parens(Op("*").Add(Id(fieldName))),
		)
	} else {
		// Non-nullable fields always have value
		g.Id("v.Value").Op("=").Id(f.ValueCastToType).Parens(Id(fieldName))
	}

	g.Id("v.Unknown").Op("=").False()
}

// genNestedBody generates block which reads message into v
func (f *FieldCopyToGenerator) genNestedBody(m *MessageCopyToGenerator, fieldName string, typ string, g *Group) {
	copyObj := func(g *Group) {
		g.Id("obj").Op(":=").Id(fieldName)
		g.Id("tf").Op(":=").Id("&v")
		m.GenerateFields(g)
	}

	f.getAttr("v", f.Field.ElemValueType, g)
	g.If(Id("!ok")).Block(
		// v := types.Object{Attrs: make(map[string]attr.Value, len(o.AttrTypes)), AttrTypes: o.AttrTypes}
		Id("v").Op("=").Id(typ).Block(Dict{
			Id("Attrs"):     Make(Map(String()).Id(f.i.WithPackage(Attr, "Value")), Len(Id("o.AttrTypes"))),
			Id("AttrTypes"): Id("o.AttrTypes"),
		}),
	).Else().Block(
		If(Id("v.Attrs").Op("==").Nil()).Block(
			Id("v.Attrs").Op("=").Make(Map(String()).Id(f.i.WithPackage(Attr, "Value")), Len(Id("tf.AttrTypes"))),
		),
	)
	if f.IsNullable {
		// if obj.Nested == nil
		g.If(Id(fieldName).Op("==").Nil()).Block(
			Id("v.Null").Op("=").True(),
		).Else().BlockFunc(
			copyObj,
		)
	} else {
		g.BlockFunc(copyObj)
	}
	g.Id("v.Unknown").Op("=").False()
}

// assertTo asserts a to typ
func (f *FieldCopyToGenerator) assertTo(typ string, g *Group, els func(g *Group)) {
	// v, ok := a.(types.ListType)
	g.List(Id("o"), Id("ok")).Op(":=").Id("a").Assert(Id(typ))
	g.If(Id("!ok")).Block(
		Id("diags.AddError").Call(
			Lit(errWriting),
			Lit(fmt.Sprintf("A type for %v can not be converted to %v", f.Path, f.Field.Type)),
		),
	).Else().BlockFunc(els)
}

// getValueField returns list/map value field
func (f *FieldCopyToGenerator) getValueField() *desc.Field {
	if f.IsMap {
		return f.MapValueField
	}

	return f.Field
}

// genPrimitive generates CopyTo statement for a primitive type
func (f *FieldCopyToGenerator) genPrimitive() *Statement {
	fieldName := "obj." + f.Name

	return f.nextField("_", func(g *Group) {
		f.genPrimitiveBody(fieldName, g)
		g.Id("tf.Attrs").Index(Lit(f.NameSnake)).Op("=").Id("v")
	})
}

// genNested generates CopyTo statement for a nested message
func (f *FieldCopyToGenerator) genNested() *Statement {
	m := NewMessageCopyToGenerator(f.Message, f.i)
	fieldName := "obj." + f.Name

	return f.nextField("a", func(g *Group) {
		f.assertTo(f.Field.ElemType, g, func(g *Group) {
			f.genNestedBody(m, fieldName, f.Field.ValueType, g)
			g.Id("tf.Attrs").Index(Lit(f.NameSnake)).Op("=").Id("v")
		})
	})
}

func (f *FieldCopyToGenerator) genListOrMap() *Statement {
	fieldName := "obj." + f.Name

	var mk Code

	if f.IsMap {
		// make(map[string]attr.Value, len(obj.Map))
		mk = Make(Map(String()).Id(f.i.WithPackage(Attr, "Value")), Len(Id(fieldName)))
	}
	if f.IsRepeated {
		// make(map[string]attr.Value, len(obj.List))
		mk = Make(Index().Id(f.i.WithPackage(Attr, "Value")), Len(Id(fieldName)))
	}

	return f.nextField("a", func(g *Group) {
		f.assertTo(f.Field.Type, g, func(g *Group) {
			f.getAttr("c", f.Field.ValueType, g)

			g.If(Id("!ok")).Block(
				// c := types.Object{Elems: make([]attr.Value, ElemType: o.ElemType}
				Id("c").Op("=").Id(f.Field.ValueType).Block(Dict{
					Id("Elems"):    mk,
					Id("ElemType"): Id("o.ElemType"),
				}),
			).Else().Block(
				If(Id("c.Elems").Op("==").Nil()).Block(
					Id("c.Elems").Op("=").Add(mk),
				),
			)

			g.If(Id(fieldName)).Op("!=").Nil().BlockFunc(func(g *Group) {
				if (f.Kind != desc.PrimitiveList) && (f.Kind != desc.PrimitiveMap) {
					g.Id("o").Op(":=").Id("o.ElemType").Assert(Id(f.ElemType))
				}
				// for k, a := range obj.List
				g.For(List(Id("k"), Id("a"))).Op(":=").Range().Id(fieldName).BlockFunc(func(g *Group) {
					if (f.Kind == desc.PrimitiveList) || (f.Kind == desc.PrimitiveMap) {
						f.genPrimitiveBody("a", g)
					} else {
						m := NewMessageCopyToGenerator(f.getValueField().Message, f.i)
						f.genNestedBody(m, "a", f.Field.ElemValueType, g)
					}
					g.Id("c.Elems").Index(Id("k")).Op("=").Id("v")
				})
			})

			g.Id("c.Unknown").Op("=").False()
			g.Id("tf.Attrs").Index(Lit(f.NameSnake)).Op("=").Id("c")
		})
	})
}

// genCustom generates statement representing custom type
func (f *FieldCopyToGenerator) genCustom() *Statement {
	return f.nextField("_", func(g *Group) {
		g.Id("v").Op(":=").Id("CopyTo"+f.Suffix).Params(Id("diags"), Id("obj."+f.Name))
		g.Id("tf.Attrs").Index(Lit(f.NameSnake)).Op("=").Id("v")
	})
}
