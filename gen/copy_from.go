package gen

import (
	"fmt"

	. "github.com/dave/jennifer/jen"
	"github.com/gravitational/protoc-gen-terraform/desc"
)

const (
	errReading = "Error reading value from Terraform"
)

// MessageCopyFromGenerator is the visitor struct to generate tfsdk.Schema of a message
type MessageCopyFromGenerator struct {
	*desc.Message
	i *desc.Imports
}

// NewMessageCopyFromGenerator returns new MessageCopyFromGenerator struct
func NewMessageCopyFromGenerator(m *desc.Message, i *desc.Imports) *MessageCopyFromGenerator {
	return &MessageCopyFromGenerator{Message: m, i: i}
}

// Generate generates Copy<Name>FromTerraform method
func (m *MessageCopyFromGenerator) Generate() []byte {
	methodName := "Copy" + m.Name + "FromTerraform"
	tf := Id("tf").Id(m.i.WithPackage(Types, "Object"))
	obj := Id("obj").Op("*").Id(m.GoType)
	diags := Var().Id("diags").Id(m.i.WithPackage(Diag, "Diagnostics"))

	// func Copy<name>FromTerraform(tf types.Object, obj *apitypes.<name>) diag.Diagnostics
	// ... statements for a fields
	method :=
		Commentf("// %v copies contents of the source Terraform object into a target struct\n", methodName).
			Func().Id(methodName).
			Params(tf, obj).
			Id(m.i.WithPackage(Diag, "Diagnostics")).
			BlockFunc(func(g *Group) {
				g.Add(diags)
				m.GenerateFields(g)
				g.Return(Id("diags"))
			})

	return []byte(method.GoString() + "\n")
}

// GenerateFields generates specific statements for CopyToTF<name> methods
func (m *MessageCopyFromGenerator) GenerateFields(g *Group) {
	for _, f := range m.Fields {
		g.Add(NewFieldCopyFromGenerator(f, m.i).Generate())
	}
}

// FieldCopyFromGenerator is a visitor for a field
type FieldCopyFromGenerator struct {
	*desc.Field
	i *desc.Imports
}

// NewFieldCopyFromGenerator returns new FieldCopyFromGenerator struct
func NewFieldCopyFromGenerator(f *desc.Field, i *desc.Imports) *FieldCopyFromGenerator {
	return &FieldCopyFromGenerator{Field: f, i: i}
}

// Generate generates CopyFrom fragment for a field of different kind
func (f *FieldCopyFromGenerator) Generate() *Statement {
	switch f.Kind {
	case desc.Primitive:
		return f.genPrimitive()
	case desc.Nested:
		return f.genNested()
	case desc.PrimitiveList, desc.PrimitiveMap:
		return f.genPrimitiveListOrMap()
	case desc.NestedList, desc.NestedMap:
		return f.genNestedListOrMap()
	case desc.Custom:
		return f.genCustom()
	}
	return nil
}

// nextField reads current field value from Terraform object and asserts it's type against expected
func (f *FieldCopyFromGenerator) nextField(g func(g *Group)) *Statement {
	return Block(
		// a, ok := ft.Attrs["key"]
		List(Id("a"), Id("ok")).Op(":=").Id("tf.Attrs").Index(Lit(f.NameSnake)),
		If(Id("!ok")).Block(
			Id("diags.AddError").Call(
				Lit(errReading),
				Lit(fmt.Sprintf("A value for %v is missing in the source Terraform object Attrs", f.Path)),
			),
		).Else().Block(
			// v, ok := a.(types.Int64)
			List(Id("v"), Id("ok")).Op(":=").Id("a").Assert(Id(f.ValueType)),
			If(Id("!ok")).Block(
				Id("diags.AddError").Call(
					Lit(errReading),
					Lit(fmt.Sprintf("A value for %v can not be converted to %v", f.Path, f.ValueType)),
				),
			).Else().BlockFunc(g),
		),
	)
}

// genPrimitiveBody generates fragment which converts attr.Value v to go variable t
func (f *FieldCopyFromGenerator) genPrimitiveBody(g *Group) {
	// var t float32 || *float32, acts as zero value if needed
	g.Var().Id("t").Id(f.GoElemType)
	// if !v.Null {
	g.If(Id("!v.Null && !v.Unknown")).BlockFunc(func(g *Group) {
		if !f.IsNullable {
			// obj.Float = float32(v.Value)
			g.Id("t").Op("=").Id(f.ValueCastFromType).Parens(Id("v.Value"))
		} else {
			// c := float32(v.Value)
			g.Id("c").Op(":=").Id(f.ValueCastFromType).Parens(Id("v.Value"))
			// obj.Float = &c
			g.Id("t").Op("=&").Id("c")
		}
	})
}

// genListOrMapBody generates base iterator statement, which iterates over a list or map
func (f *FieldCopyFromGenerator) genListOrMapIterator(g *Group, typ *Statement, els func(g *Group)) {
	objFieldName := "obj." + f.Name

	// obj.List = make([]string, len(v.Elems)) - same for maps
	g.Id(objFieldName).Op("=").Make(Id(f.GoType), Len(Id("v.Elems")))

	// if !v.Null
	g.If(Id("!v.Null && !v.Unknown")).BlockFunc(func(g *Group) {
		// for k, el := range v.Elems - where k is either index or map key
		g.For(List(Id("k"), Id("a"))).Op(":=").Range().Id("v.Elems").BlockFunc(func(g *Group) {
			// v, ok := a.(types.String)
			g.List(Id("v"), Id("ok")).Op(":=").Id("a").Assert(typ)
			g.If(Id("!ok")).Block(
				Id("diags.AddError").Call(
					Lit(errReading),
					Lit(fmt.Sprintf("An element value for %v can not be converted to %v", f.Path, typ.GoString())),
				),
			).Else().BlockFunc(els)
		})
	})
}

// genPrimitive generates CopyFrom fragment for a primitive field
func (f *FieldCopyFromGenerator) genPrimitive() *Statement {
	return f.nextField(func(g *Group) {
		f.genPrimitiveBody(g)
		g.Id("obj." + f.Name).Op("=").Id("t")
	})
}

// genNested generates CopyFrom fragment for a nested object
func (f *FieldCopyFromGenerator) genNested() *Statement {
	m := NewMessageCopyFromGenerator(f.Message, f.i)
	objFieldName := "obj." + f.Name

	return f.nextField(func(g *Group) {
		if f.IsNullable {
			// obj.Nested = nil
			g.Id(objFieldName).Op("=").Nil()
		} else {
			// obj.Nested = Nested{}
			g.Id(objFieldName).Op("=").Id(f.GoElemType).Values()
		}
		// if !v.Null
		g.If(Id("!v.Null && !v.Unknown")).BlockFunc(func(g *Group) {
			// tf := v
			g.Id("tf").Op(":=").Id("v")

			if f.IsNullable {
				// obj.Nested = &Nested{}
				g.Id(objFieldName).Op("=&").Id(f.GoElemTypeIndirect).Values()
				// obj := obj.Nested
				g.Id("obj").Op(":=").Id(objFieldName)
			} else {
				// obj := &obj.Nested
				g.Id("obj").Op(":=&").Id(objFieldName)
			}

			m.GenerateFields(g)
		})
	})
}

// getValueField returns list/map value field
func (f *FieldCopyFromGenerator) getValueField() *desc.Field {
	if f.IsMap {
		return f.MapValueField
	}

	return f.Field
}

// genPrimitiveListOrMap generates CopyFrom fragment for a list or a map of a primitive values
func (f *FieldCopyFromGenerator) genPrimitiveListOrMap() *Statement {
	objFieldName := "obj." + f.Name

	field := f.getValueField()

	return f.nextField(func(g *Group) {
		f.genListOrMapIterator(g, Id(field.ElemValueType), func(g *Group) {
			f.genPrimitiveBody(g)

			// obj.List[k] = t
			g.Id(objFieldName).Index(Id("k")).Op("=").Id("t")
		})
	})
}

// genNestedListOrMap generates CopyFrom fragment for a list or a map of a nested objects
func (f *FieldCopyFromGenerator) genNestedListOrMap() *Statement {
	objFieldName := "obj." + f.Name

	field := f.getValueField()
	m := NewMessageCopyFromGenerator(field.Message, f.i)

	return f.nextField(func(g *Group) {
		f.genListOrMapIterator(g, Id(field.ElemValueType), func(g *Group) {
			// var t Nested || *Nested
			g.Var().Id("t").Id(f.GoElemType)

			g.If(Id("!v.Null && !v.Unknown")).BlockFunc(func(g *Group) {
				// tf := v
				g.Id("tf").Op(":=").Id("v")

				if f.IsNullable {
					// t = &Nested{}
					g.Id("t").Op("=&").Id(f.GoElemTypeIndirect).Values()
					// obj := t - obj is just an alias to reuse field generator code
					g.Id("obj").Op(":=").Id("t")
				} else {
					// obj := &t
					g.Id("obj").Op(":=&").Id("t")
				}
				m.GenerateFields(g)
			})

			// obj.List[k] = t
			g.Id(objFieldName).Index(Id("k")).Op("=").Id("t")
		})
	})
}

// genCustom generates statement representing custom type
func (f *FieldCopyFromGenerator) genCustom() *Statement {
	return Block(
		// a, ok := ft.Attrs["key"]
		List(Id("a"), Id("ok")).Op(":=").Id("tf.Attrs").Index(Lit(f.NameSnake)),
		If(Id("!ok")).Block(
			Id("diags.AddError").Call(
				Lit(errReading),
				Lit(fmt.Sprintf("A value for %v is missing in Terraform object Attrs", f.Path)),
			)),
		Id("CopyFrom"+f.Suffix).Params(Id("diags"), Id("a"), Id("&obj."+f.Name)),
	)
}
