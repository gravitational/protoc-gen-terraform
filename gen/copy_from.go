package gen

import (
	"strings"

	. "github.com/dave/jennifer/jen"
	"github.com/gravitational/protoc-gen-terraform/desc"
)

// MessageCopyFromGenerator is the visitor struct to generate tfsdk.Schema of a message
type MessageCopyFromGenerator struct {
	*desc.Message
	c GeneratorContext
}

// NewMessageCopyFromGenerator returns new MessageCopyFromGenerator struct
func NewMessageCopyFromGenerator(m *desc.Message, c GeneratorContext) *MessageCopyFromGenerator {
	return &MessageCopyFromGenerator{Message: m, c: c}
}

// Generate generates CopyFromTF<Name> method
func (m *MessageCopyFromGenerator) Generate() []byte {
	f := m.GenFields()
	f.Add(Return(Nil()))

	// func Copy<name>FromTerraform(tf types.Object, obj *apitypes.<name>)
	j := Commentf("// Copy"+m.Name+"FromTerraform copies the contents of the source Terraform object into target struct\n").
		Func().Id("Copy"+m.Name+"FromTerraform").
		Params(
			Id("tf").Add(m.c.TFType("Object")),
			Id("obj").Op("*").Id(m.GoType)).
		Error().
		Block(f...)

	return []byte(j.GoString() + "\n")
}

// GenFields generates specific statements for CopyFromTF<name> methods
func (m *MessageCopyFromGenerator) GenFields() Statement {
	l := Statement{}

	for _, f := range m.Fields {
		f := NewFieldCopyFromGenerator(f, m.c)
		l.Add(f.Generate())
	}

	return l
}

// FieldCopyFromGenerator is a visitor for a field
type FieldCopyFromGenerator struct {
	*desc.Field
	c GeneratorContext
}

// NewFieldCopyFromGenerator returns new FieldCopyFromGenerator struct
func NewFieldCopyFromGenerator(f *desc.Field, c GeneratorContext) *FieldCopyFromGenerator {
	return &FieldCopyFromGenerator{Field: f, c: c}
}

// Generate generates CopyFrom fragment for a field of different kind
func (f *FieldCopyFromGenerator) Generate() *Statement {
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

// genPrimitive generates CopyFrom fragment for a primitive field
func (f *FieldCopyFromGenerator) genPrimitive() *Statement {
	j := Statement{}

	j.
		Add(f.getTFAttr()...).                            // Read attribute
		Add(f.c.AssertTFAttrToPrimitiveType(f.Field)...). // Convert to Terraform type
		Add(f.genPrimitiveMain("obj." + f.Name))          // Set obj value

	return Block(j...)
}

// genPrimitiveMain generates CopyFrom fragment for a primitive field
func (f *FieldCopyFromGenerator) genPrimitiveMain(target string) *Statement {
	convert := Id("convert")
	zero := Id(f.GoElemTypeZeroValue)
	if f.IsNullable {
		convert = Op("&").Add(convert)
		zero = Nil()
	}

	return If(Id("v.Unknown || v.Null")).
		Block(
			Id(target).Op("=").Add(zero), // Set obj value to zero if a value is unknown or null
		).
		Else().
		Block(
			Id("convert").Op(":=").Add(f.c.CastValueToGoElemType(f.Field)), // Cast Terraform value to obj value and set
			Id(target).Op("=").Add(convert),
		)
}

// genNested generates main part for a CopyFrom fragment for a nested object
func (f *FieldCopyFromGenerator) genNested() *Statement {
	m := NewMessageCopyFromGenerator(f.Message, f.c)

	j := Statement{}

	j.
		Add(f.getTFAttr()...).                         // Read attribute
		Add(f.c.AssertTFAttrToNestedType(f.Field)...). // Cast to List or Map
		Add(f.genNestedMain(m, "obj."+f.Name))         // Generate nested statement

	return Block(j...)
}

// genNestedMain generates main part of a CopyFrom fragment for a nested object
func (f *FieldCopyFromGenerator) genNestedMain(m *MessageCopyFromGenerator, to string) *Statement {
	if !f.IsNullable {
		n := Statement{
			Id("tf").Op(":=").Id("v"),
			Id("obj").Op(":=").Id("&" + to),
		}
		n.Add(m.GenFields()...)

		return If(Id("v.Null || v.Unknown")).Block(
			Id(to).Op("=").Id(f.GoElemType).Values(),
		).Else().Block(n...)
	}

	n := Statement{
		Id(to).Op("=").Id(strings.ReplaceAll(f.GoElemType, "*", "&")).Values(),
		Id("tf").Op(":=").Id("v"),
		Id("obj").Op(":=").Id(to),
	}
	n.Add(m.GenFields()...)

	return If(Id("v.Null || v.Unknown")).Block(
		Id(to).Op("=").Nil(),
	).Else().Block(n...)
}

// genNestedListOrMap generates CopyFrom fragment for a nested/primitive list or a map
func (f *FieldCopyFromGenerator) genNestedListOrMap() *Statement {
	var construct *Statement
	var zero *Statement

	vf := f
	j := Statement{}

	target := Id("obj." + f.Name)
	zero = target.Clone().Op("=").Nil()

	if f.IsRepeated {
		construct = Make(Id(vf.GoType), Len(Id("v.Elems")))
	}

	if f.IsMap {
		vf = NewFieldCopyFromGenerator(f.MapValueField, f.c)
		construct = Make(Id(f.GoType))
	}

	iterate := &Statement{}

	if vf.IsMessage {
		iterate.
			Add(f.c.AssertTFAttrToObject(f.Field)...).
			Add(Var().Id("el").Id(vf.GoElemType))

		iterate.Add(f.genNestedMain(NewMessageCopyFromGenerator(vf.Message, vf.c), "el"))
	} else {
		var star string
		if vf.Field.IsNullable {
			star = "*"
		}

		iterate.
			Add(f.c.AssertTFAttrToPrimitiveType(vf.Field)...).
			Add(Var().Id("el").Op(star).Id(vf.GoElemType)).
			Add(vf.genPrimitiveMain("el"))
	}

	iterate.
		Add(
			target.Clone().Index(Id("k")).Op("=").Id("el"),
		)

	j.
		Add(f.getTFAttr()...).
		Add(f.c.AssertTFAttrToNestedType(f.Field)...)

	j.Add(
		If(Id("v.Null || v.Unknown")).Block(
			zero,
		).Else().Block(
			target.Clone().Op("=").Add(construct),
			For(List(Id("k"), Id("a")).Op(":=").Range().Id("v.Elems")).Block(*iterate...),
		),
	)

	return Block(j...)
}

// genCustom generates statement representing custom type
func (f *FieldCopyFromGenerator) genCustom() *Statement {
	var j Statement
	j.
		Add(f.getTFAttr()...).
		Add(
			Err().Op(":=").Id("CopyFrom"+f.Suffix).Params(Id("a"), Op("&").Id("obj."+f.Name)),
			If(Err().Op("!=").Nil()).Block(
				Return(Err()),
			),
		)

	return Block(j...)
}

// getTFAttr is a shortcut to f.c.getTFAttr
func (f *FieldCopyFromGenerator) getTFAttr() Statement {
	return f.c.getTFAttr(f.NameSnake, f.Path)
}
