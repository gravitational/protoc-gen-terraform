package gen

import (
	"github.com/dave/jennifer/jen"
	. "github.com/dave/jennifer/jen"
	"github.com/gravitational/protoc-gen-terraform/desc"
)

// MessageCopyToGenerator is the visitor struct to generate tfsdk.Schema of a message
type MessageCopyToGenerator struct {
	*desc.Message
	c GeneratorContext
}

// NewMessageCopyToGenerator returns new MessageCopyToGenerator struct
func NewMessageCopyToGenerator(m *desc.Message, c GeneratorContext) *MessageCopyToGenerator {
	return &MessageCopyToGenerator{Message: m, c: c}
}

// Generate generates CopyToTF<Name> method
func (m *MessageCopyToGenerator) Generate() []byte {
	f := m.GenFields()
	f.Add(Return(Nil()))

	// func Copy<name>ToTerraform(tf types.Object, obj apitypes.<name>)
	j := Commentf("// Copy"+m.Name+"ToTerraform copies the contents of the source Terraform object into target struct\n").
		Func().Id("Copy"+m.Name+"ToTerraform").
		Params(
			Id("tf").Op("*").Add(m.c.TFType("Object")),
			Id("obj").Id(m.GoType)).
		Error().
		Block(f...)

	return []byte(j.GoString() + "\n")
}

// GenFields generates specific statements for CopyToTF<name> methods
func (m *MessageCopyToGenerator) GenFields() Statement {
	l := Statement{}

	for _, f := range m.Fields {
		f := NewFieldCopyToGenerator(f, m.c)
		l.Add(f.Generate())
	}

	return l
}

// FieldCopyToGenerator is a visitor for a field
type FieldCopyToGenerator struct {
	*desc.Field
	c GeneratorContext
}

// NewFieldCopyToGenerator returns new FieldCopyToGenerator struct
func NewFieldCopyToGenerator(f *desc.Field, c GeneratorContext) *FieldCopyToGenerator {
	return &FieldCopyToGenerator{Field: f, c: c}
}

// Generate generates CopyTo fragment for a field of different kind
func (f *FieldCopyToGenerator) Generate() *Statement {
	switch f.Kind {
	case desc.Primitive:
		return f.genPrimitive()
	case desc.Nested:
		return f.genNested()
	case desc.NestedList, desc.NestedMap, desc.PrimitiveList, desc.PrimitiveMap:
		return f.genNestedListOrMap()
	case desc.Custom:
		return f.genCustom()
	}
	return nil
}

// genPrimitive generates CopyTo statement for a primitive type
func (f *FieldCopyToGenerator) genPrimitive() *Statement {
	j := Statement{}

	j.
		Add(f.getTFAttr()...).
		Add(f.c.AssertTFAttrToPrimitiveType(f.Field)...)

	name := Id("obj." + f.Name)
	nullEq := Id("v.Null").Op("=")
	unknownEq := Id("v.Unknown").Op("=")
	valueEq := Id("v.Value").Op("=")

	if f.IsNullable {
		j.Add(
			If(name.Clone().Op("==").Nil()).Block(
				nullEq.Clone().True(),
				unknownEq.Clone().False(),
				valueEq.Clone().Id(f.GoElemTypeZeroValue),
			).Else().Block(
				nullEq.Clone().False(),
				unknownEq.Clone().False(),
				valueEq.Clone().Id(f.GoElemType).Parens(Op("*").Add(name.Clone())),
			),
		)
	} else {
		j.
			Add(
				valueEq.Clone().Id(f.TFSchemaValueCastType).Parens(name.Clone()),
				If(Parens(name.Clone().Op("!=").Id(f.GoElemTypeZeroValue))).Block(
					nullEq.Clone().False(),
				),
				unknownEq.Clone().False(),
			)
	}

	j.Add(Id("tf.Attrs").Index(Lit(f.NameSnake)).Op("=").Id("v"))

	return Block(j...)
}

// genNested generates CopyTo statement for a nested message
func (f *FieldCopyToGenerator) genNested() *Statement {
	m := NewMessageCopyToGenerator(f.Message, f.c)

	j := f.getTFAttr()
	j.Add(f.c.AssertTFAttrToObject(f.Field)...)

	source := Id("obj." + f.Name)

	if f.IsNullable {
		b := Statement{}

		b.Add(
			Id("obj").Op(":=").Add(source),
			Id("tf").Op(":=").Id("v"),
		).
			Add(m.GenFields()...)

		j.Add(
			If(Parens(Add(source).Op("!=").Nil())).
				Block(b...).
				Else().Block(
				Id("v.Null").Op("=").True(),
				Id("v.Unknown").Op("=").False(),
			),
		)
	} else {
		j.Add(
			Id("obj").Op(":=").Add(source),
			Id("tf").Op(":=").Id("v"),
		)
		j.Add(m.GenFields()...)
	}

	return Block(j...)
}

// genNestedListOrMap generates nested statements for Map and List
func (f *FieldCopyToGenerator) genNestedListOrMap() *Statement {
	var mk *jen.Statement
	var iterateEl jen.Statement
	var vf *desc.Field

	x := Id("x")
	objField := Id("obj." + f.Name)

	if f.IsRepeated {
		vf = f.Field
		mk = Make(Index().Add(f.c.AttrValue()), Len(objField.Clone()))
	}

	if f.IsMap {
		vf = f.MapValueField
		mk = Make(Map(String()).Add(f.c.AttrValue()))
	}

	if f.Kind == desc.PrimitiveList || f.Kind == desc.PrimitiveMap {
		if f.IsNullable {
			x = Op("*").Add(x.Clone())
		}
		iterateEl.Add(
			Id("e").Index(Id("k")).Op("=").Add(f.c.GetPrimitiveType(vf)).Values(Dict{
				Id("Value"): Id(vf.TFSchemaValueCastType).Parens(x.Clone()),
			}),
		)
	} else if f.Kind == desc.NestedList || f.Kind == desc.NestedMap {
		var m *MessageCopyToGenerator
		if f.IsRepeated {
			m = NewMessageCopyToGenerator(f.Message, f.c)
		} else if f.IsMap {
			m = NewMessageCopyToGenerator(f.MapValueField.Message, f.c)
		}

		iterateEl.Add(
			Id("tf").Op(":=").Add(f.c.TFType("Object")).Values(Dict{
				Id("AttrTypes"): Id("v.ElemType.").Parens(f.c.TFType("ObjectType")).Id(".AttrTypes"),
			}),
			Id("obj").Op(":=").Add(x),
		)
		iterateEl.Add(m.GenFields()...)
		iterateEl.Add(Id("e").Index(Id("k")).Op("=").Id("tf"))
	}

	j := f.getTFAttr()
	j.
		Add(f.c.AssertTFAttrToNestedType(f.Field)...).
		Add(
			If(objField.Clone().Op("==").Nil()).Block(
				Id("v.Null").Op("=").True(),
			).Else().Block(
				Id("e").Op(":=").Add(mk),
				For(List(Id("k"), Id("x")).Op(":=").Range().Add(objField.Clone())).Block(
					iterateEl...,
				),
				Id("v.Elems").Op("=").Id("e"),
				Id("v.Null").Op("=").False(),
			),
			Id("v.Unknown").Op("=").False(),
		)

	return Block(j...)
}

// genCustom generates statement representing custom type
func (f *FieldCopyToGenerator) genCustom() *Statement {
	var j Statement
	j.
		Add(f.getTFAttr()...).
		Add(
			Err().Op(":=").Id("CopyTo"+f.Suffix).Params(Id("a"), Id("obj."+f.Name)),
			If(Err().Op("!=").Nil()).Block(
				Return(Err()),
			),
		)

	return Block(j...)
}

// getTFAttr is a shortcut to f.c.getTFAttr
func (f *FieldCopyToGenerator) getTFAttr() Statement {
	return f.c.getTFAttr(f.NameSnake, f.Path)
}
