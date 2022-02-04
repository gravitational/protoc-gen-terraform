package gen

// import (
// 	"fmt"
// 	"strings"

// 	. "github.com/dave/jennifer/jen"
// 	"github.com/gravitational/protoc-gen-terraform/desc"
// )

// const (
// 	errWriting = "Error writing value to Terraform"
// )

// // MessageCopyToGenerator is the visitor struct to generate tfsdk.Schema of a message
// type MessageCopyToGenerator struct {
// 	*desc.Message
// 	c GeneratorContext
// }

// // NewMessageCopyToGenerator returns new MessageCopyToGenerator struct
// func NewMessageCopyToGenerator(m *desc.Message, c GeneratorContext) *MessageCopyToGenerator {
// 	return &MessageCopyToGenerator{Message: m, c: c}
// }

// // Generate generates CopyToTF<Name> method
// func (m *MessageCopyToGenerator) Generate() []byte {
// 	methodName := "Copy" + m.Name + "ToTerraform"
// 	tf := Id("tf").Op("*").Add(m.c.Types("Object"))
// 	obj := Id("obj").Id(m.GoType)
// 	diags := Var().Id("diags").Add(m.c.DiagDiagnostics())

// 	// func Copy<name>ToTerraform(tf types.Object, obj *apitypes.<name>)
// 	// ... statements for a fields
// 	method :=
// 		Commentf("// %v copies contents of the source Terraform object into a target struct\n", methodName).
// 			Func().Id(methodName).
// 			Params(obj, tf).
// 			Add(m.c.DiagDiagnostics()).
// 			BlockFunc(func(g *Group) {
// 				g.Add(diags)
// 				m.GenerateFields(g)
// 				g.Return(Id("diags"))
// 			})

// 	return []byte(method.GoString() + "\n")
// }

// // GenerateFields generates specific statements for CopyToTF<name> methods
// func (m *MessageCopyToGenerator) GenerateFields(g *Group) {
// 	for _, f := range m.Fields {
// 		g.Add(NewFieldCopyToGenerator(f, m.c).Generate())
// 	}
// }

// // FieldCopyToGenerator is a visitor for a field
// type FieldCopyToGenerator struct {
// 	*desc.Field
// 	c GeneratorContext
// }

// // NewFieldCopyToGenerator returns new FieldCopyToGenerator struct
// func NewFieldCopyToGenerator(f *desc.Field, c GeneratorContext) *FieldCopyToGenerator {
// 	return &FieldCopyToGenerator{Field: f, c: c}
// }

// // Generate generates CopyTo fragment for a field of different kind
// func (f *FieldCopyToGenerator) Generate() *Statement {
// 	switch f.Kind {
// 	case desc.Primitive:
// 		return f.genPrimitive()
// 	case desc.Nested:
// 		return f.genNested()
// 	case desc.PrimitiveList, desc.PrimitiveMap, desc.NestedList, desc.NestedMap:
// 		return f.genPrimitiveListOrMap()
// 	case desc.Custom:
// 		return f.genCustom()
// 	}
// 	return nil
// }

// // nextField reads current field value from Terraform object and asserts it's type against expected
// func (f *FieldCopyToGenerator) nextField(v string, g func(g *Group)) *Statement {
// 	return Block(
// 		// _, ok := ft.AttrsTypes["key"]
// 		List(Id(v), Id("ok")).Op(":=").Id("tf.AttrTypes").Index(Lit(f.NameSnake)),
// 		If(Id("!ok")).Block(
// 			Id("diags.AddError").Call(
// 				Lit(errWriting),
// 				Lit(fmt.Sprintf("A value type for %v is missing in the target Terraform object AttrTypes", f.Path)),
// 			),
// 		).Else().BlockFunc(g),
// 	)
// }

// // genPrimitiveBody generates block which reads object field into v
// func (f *FieldCopyToGenerator) genPrimitiveBody(fieldName string, g *Group) {
// 	g.Var().Id("v").Add(f.c.GetElemAttrValueType(f.Field))

// 	if f.IsNullable {
// 		g.If(Id(fieldName).Op("==").Nil()).Block(
// 			Id("v.Null").Op("=").True(),
// 		).Else().Block(
// 			Id("v.Value").Op("=").Id(strings.Replace(f.GoElemType, "*", "", -1)).Parens(Op("*").Add(Id(fieldName))),
// 			Id("v.Null").Op("=").False(),
// 		)
// 	} else {
// 		// Non-nullable fields always have value
// 		g.Id("v.Value").Op("=").Id(f.SchemaValueCastType).Parens(Id(fieldName))
// 		g.Id("v.Null").Op("=").False()
// 	}
// 	g.Id("v.Unknown").Op("=").False()
// }

// // genNestedBody generates block which reads message into v
// func (f *FieldCopyToGenerator) genNestedBody(m *MessageCopyToGenerator, fieldName string, g *Group) {
// 	copyObj := func(g *Group) {
// 		g.BlockFunc(func(g *Group) {
// 			g.Id("obj").Op(":=").Id(fieldName)
// 			g.Id("tf").Op(":=").Id("&v")
// 			m.GenerateFields(g)
// 		})
// 	}

// 	// v := types.Object{Attrs: make(map[string]attr.Value, len(o.AttrTypes)), AttrTypes: o.AttrTypes}
// 	g.Id("v").Op(":=").Add(f.c.GetElemAttrValueType(f.Field)).Block(Dict{
// 		Id("Attrs"):     Make(Map(String()).Add(f.c.AttrValue()), Len(Id("o.AttrTypes"))),
// 		Id("AttrTypes"): Id("o.AttrTypes"),
// 	})

// 	if f.IsNullable {
// 		// if obj.Nested == nil
// 		g.If(Id(fieldName).Op("==").Nil()).Block(
// 			Id("v.Null").Op("=").True(),
// 		).Else().BlockFunc(
// 			copyObj,
// 		)
// 	} else {
// 		copyObj(g)
// 	}
// }

// // getValueField returns list/map value field
// func (f *FieldCopyToGenerator) getValueField() *desc.Field {
// 	if f.IsMap {
// 		return f.MapValueField
// 	}

// 	return f.Field
// }

// // genPrimitive generates CopyTo statement for a primitive type
// func (f *FieldCopyToGenerator) genPrimitive() *Statement {
// 	fieldName := "obj." + f.Name

// 	return f.nextField("_", func(g *Group) {
// 		f.genPrimitiveBody(fieldName, g)
// 		g.Id("tf.Attrs").Index(Lit(f.NameSnake)).Op("=").Id("v")
// 	})
// }

// // genNested generates CopyTo statement for a nested message
// func (f *FieldCopyToGenerator) genNested() *Statement {
// 	m := NewMessageCopyToGenerator(f.Message, f.c)
// 	typ := f.c.GetElemAttrTypeType(f.Field)
// 	fieldName := "obj." + f.Name

// 	copyObj := func(g *Group) {
// 		g.BlockFunc(func(g *Group) {
// 			g.Id("obj").Op(":=").Id(fieldName)
// 			g.Id("tf").Op(":=").Id("&v")
// 			m.GenerateFields(g)
// 		})
// 	}

// 	// TODO: assertTo -> DRY
// 	return f.nextField("a", func(g *Group) {
// 		// v, ok := a.(types.ObjectType)
// 		g.List(Id("o"), Id("ok")).Op(":=").Id("a").Assert(typ)
// 		g.If(Id("!ok")).Block(
// 			Id("diags.AddError").Call(
// 				Lit(errWriting),
// 				Lit(fmt.Sprintf("A type for %v can not be converted to %v", f.Path, typ.GoString())),
// 			),
// 		).Else().BlockFunc(func(g *Group) {
// 			// v := types.Object{Attrs: make(map[string]attr.Value, len(o.AttrTypes)), AttrTypes: o.AttrTypes}
// 			g.Id("v").Op(":=").Add(f.c.GetAttrValueType(f.Field)).Block(Dict{
// 				Id("Attrs"):     Make(Map(String()).Add(f.c.AttrValue()), Len(Id("o.AttrTypes"))),
// 				Id("AttrTypes"): Id("o.AttrTypes"),
// 			})

// 			if f.IsNullable {
// 				// if obj.Nested == nil
// 				g.If(Id(fieldName).Op("==").Nil()).Block(
// 					Id("v.Null").Op("=").True(),
// 				).Else().BlockFunc(
// 					copyObj,
// 				)
// 			} else {
// 				copyObj(g)
// 			}

// 			g.Id("tf.Attrs").Index(Lit(f.NameSnake)).Op("=").Id("v")
// 		})
// 	})
// }

// func (f *FieldCopyToGenerator) genPrimitiveListOrMap() *Statement {
// 	typ := f.c.GetAttrTypeType(f.Field)
// 	fieldName := "obj." + f.Name

// 	return f.nextField("a", func(g *Group) {
// 		// v, ok := a.(types.ListType)
// 		g.List(Id("o"), Id("ok")).Op(":=").Id("a").Assert(typ)
// 		g.If(Id("!ok")).Block(
// 			Id("diags.AddError").Call(
// 				Lit(errWriting),
// 				Lit(fmt.Sprintf("A type for %v can not be converted to %v", f.Path, typ.GoString())),
// 			),
// 		).Else().BlockFunc(func(g *Group) {
// 			var mk Code

// 			if f.IsMap {
// 				// make(map[string]attr.Value, len(obj.Map))
// 				mk = Make(Map(String()).Add(f.c.AttrValue()), Len(Id(fieldName)))
// 			}
// 			if f.IsRepeated {
// 				// make(map[string]attr.Value, len(obj.List))
// 				mk = Make(Index().Add(f.c.AttrValue()), Len(Id(fieldName)))
// 			}

// 			// v := types.Object{Elems: make([]attr.Value, ElemType: o.ElemType}
// 			g.Id("c").Op(":=").Add(f.c.GetAttrValueType(f.Field)).Block(Dict{
// 				Id("Elems"):    mk,
// 				Id("ElemType"): Id("o.ElemType"),
// 			})

// 			g.If(Id(fieldName)).Op("==").Nil().Block(
// 				Id("c.Null").Op("=").True(),
// 			).Else().BlockFunc(func(g *Group) {
// 				if (f.Kind != desc.PrimitiveList) && (f.Kind != desc.PrimitiveMap) {
// 					g.Id("o").Op(":=").Id("o.ElemType").Assert(f.c.Types("ObjectType"))
// 				}
// 				// for k, a := range obj.List
// 				g.For(List(Id("k"), Id("a"))).Op(":=").Range().Id(fieldName).BlockFunc(func(g *Group) {
// 					if (f.Kind == desc.PrimitiveList) || (f.Kind == desc.PrimitiveMap) {
// 						f.genPrimitiveBody("a", g)
// 					} else {
// 						m := NewMessageCopyToGenerator(f.getValueField().Message, f.c)
// 						f.genNestedBody(m, "a", g)
// 					}
// 					g.Id("c.Elems").Index(Id("k")).Op("=").Id("v")
// 				})
// 			})

// 			g.Id("tf.Attrs").Index(Lit(f.NameSnake)).Op("=").Id("c")
// 		})
// 	})
// }

// // genCustom generates statement representing custom type
// func (f *FieldCopyToGenerator) genCustom() *Statement {
// 	return f.nextField("_", func(g *Group) {
// 		g.Id("v").Op(":=").Id("CopyTo"+f.Suffix).Params(Id("diags"), Id("obj."+f.Name))
// 		g.Id("tf.Attrs").Index(Lit(f.NameSnake)).Op("=").Id("v")
// 	})
// }
