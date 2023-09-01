package main

import (
	"io"

	j "github.com/dave/jennifer/jen"
)

// MessageCopyFromGenerator is the visitor struct to generate tfsdk.Schema of a message
type MessageCopyFromGenerator struct {
	*Message
	i *Imports
}

// FieldCopyFromGenerator is a visitor for a field
type FieldCopyFromGenerator struct {
	*Field
	i *Imports
}

// NewMessageCopyFromGenerator returns new MessageCopyFromGenerator struct
func NewMessageCopyFromGenerator(m *Message, i *Imports) *MessageCopyFromGenerator {
	return &MessageCopyFromGenerator{m, i}
}

// NewFieldCopyFromGenerator returns new FieldCopyFromGenerator struct
func NewFieldCopyFromGenerator(f *Field, i *Imports) *FieldCopyFromGenerator {
	return &FieldCopyFromGenerator{f, i}
}

// Generate generates Copy<Name>FromTerraform method
func (m *MessageCopyFromGenerator) Generate(writer io.Writer) (int, error) {
	methodName := "Copy" + m.Name + "FromTerraform"
	tf := j.Id("tf").Id(m.i.WithPackage(Types, "Object"))
	obj := j.Id("obj").Op("*").Id(m.i.WithType(m.GoType))
	diags := j.Var().Id("diags").Id(m.i.WithPackage(Diag, "Diagnostics"))
	ctx := j.Id("_").Id(m.i.WithPackage("context", "Context"))

	// func Copy<name>FromTerraform(_ context.Context, tf types.Object, obj *apitypes.<name>) diag.Diagnostics
	// ... statements for a fields
	method :=
		j.Commentf("// %v copies contents of the source Terraform object into a target struct\n", methodName).
			Func().Id(methodName).
			Params(ctx, tf, obj).
			Id(m.i.WithPackage(Diag, "Diagnostics")).
			BlockFunc(func(g *j.Group) {
				g.Add(diags)
				m.GenerateFields(g)
				g.Return(j.Id("diags"))
			})

	return writer.Write([]byte(method.GoString() + "\n"))
}

// GenerateFields generates specific statements for CopyToTF<name> methods
func (m *MessageCopyFromGenerator) GenerateFields(g *j.Group) {
	// Reset all oneOf fields in advance, otherwise if all oneOf branches would be null in the passed
	// object, the oneOf field won't be nil
	for _, m := range m.OneOfNames {
		g.Add(j.Id("obj." + m).Op("=").Nil())
	}

	for _, f := range m.Fields {
		g.Add(NewFieldCopyFromGenerator(f, m.i).Generate())
	}
}

// Generate generates CopyFrom fragment for a field of different kind
func (f *FieldCopyFromGenerator) Generate() *j.Statement {
	switch f.Kind {
	case PrimitiveKind:
		return f.genPrimitive()
	case ObjectKind:
		return f.genObject()
	case PrimitiveListKind, PrimitiveMapKind:
		return f.genPrimitiveListOrMap()
	case ObjectListKind, ObjectMapKind:
		return f.genObjectListOrMap()
	case CustomKind:
		return f.genCustom()
	}
	return nil
}

// errMissingDiag diags.Append(attrMissingDiag{path})
func (f *FieldCopyFromGenerator) errAttrMissingDiag(g *j.Group) {
	g.Id("diags.Append").Call(
		j.Id("attrReadMissingDiag").Values(j.Lit(f.Path)),
	)
}

// errAttrConversionFailure diags.Append(attrConversionFailureDiag{path, typ})
func (f *FieldCopyFromGenerator) errAttrConversionFailure(path string, typ string) func(g *j.Group) {
	return func(g *j.Group) {
		g.Id("diags.Append").Call(
			j.Id("attrReadConversionFailureDiag").Values(j.Lit(path), j.Lit(typ)),
		)
	}
}

// nextField reads current field value from Terraform object and asserts it's type against expected
func (f *FieldCopyFromGenerator) nextField(g func(g *j.Group)) *j.Statement {
	return j.Block(
		// a, ok := ft.Attrs["key"]
		j.List(j.Id("a"), j.Id("ok")).Op(":=").Id("tf.Attrs").Index(j.Lit(f.NameSnake)),
		j.If(j.Id("!ok")).BlockFunc(f.errAttrMissingDiag).Else().Block(
			// v, ok := a.(types.Int64)
			j.List(j.Id("v"), j.Id("ok")).Op(":=").Id("a").Assert(j.Id(f.i.WithType(f.ValueType))),
			j.If(j.Id("!ok")).BlockFunc(
				f.errAttrConversionFailure(f.Path, f.ValueType),
			).Else().BlockFunc(g),
		),
	)
}

// genPrimitiveBody generates fragment which converts attr.Value v to go variable t
func (f *FieldCopyFromGenerator) genPrimitiveBody(g *j.Group) {
	// var t float32 || *float32, acts as zero value if needed
	g.Var().Id("t").Id(f.i.WithType(f.GoElemType))
	// if !v.Null {
	g.If(j.Id("!v.Null && !v.Unknown")).BlockFunc(func(g *j.Group) {
		if !f.IsNullable {
			// obj.Float = float32(v.Value)
			g.Id("t").Op("=").Id(f.i.WithType(f.ValueCastFromType)).Parens(j.Id("v.Value"))
		} else {
			// c := float32(v.Value)
			g.Id("c").Op(":=").Id(f.i.WithType(f.ValueCastFromType)).Parens(j.Id("v.Value"))
			// obj.Float = &c
			g.Id("t").Op("=&").Id("c")
		}
	})
}

// genListOrMapBody generates base iterator statement, which iterates over a list or map
func (f *FieldCopyFromGenerator) genListOrMapIterator(g *j.Group, typ *j.Statement, els func(g *j.Group)) {
	objFieldName := "obj." + f.Name

	// obj.List = make([]string, len(v.Elems)) - same for maps
	g.Id(objFieldName).Op("=").Make(j.Id(f.i.WithType(f.GoType)), j.Len(j.Id("v.Elems")))

	// if !v.Null
	g.If(j.Id("!v.Null && !v.Unknown")).BlockFunc(func(g *j.Group) {
		// for k, el := range v.Elems - where k is either index or map key
		g.For(j.List(j.Id("k"), j.Id("a"))).Op(":=").Range().Id("v.Elems").BlockFunc(func(g *j.Group) {
			// v, ok := a.(types.String)
			g.List(j.Id("v"), j.Id("ok")).Op(":=").Id("a").Assert(typ)
			g.If(j.Id("!ok")).BlockFunc(
				f.errAttrConversionFailure(f.Path, typ.GoString()),
			).Else().BlockFunc(els)
		})
	})
}

// genPrimitive generates CopyFrom fragment for a primitive field, wrapped by oneOf extraction
func (f *FieldCopyFromGenerator) genPrimitive() *j.Statement {
	return f.nextField(func(g *j.Group) {
		f.genPrimitiveBody(g)

		if f.OneOfName != "" {
			// Do not set empty oneOf value to not override values possibly set by other branches
			g.If(j.Id("!v.Null && !v.Unknown")).BlockFunc(func(g *j.Group) {
				g.Id("obj." + f.OneOfName).Op("=").Id("&" + f.i.WithType(f.OneOfType)).Values(j.Dict{
					j.Id(f.Name): j.Id("t"),
				})
			})
			return
		}

		if f.ParentIsOptionalEmbed {
			// If the current value is Null or Unknown, we should not set the parent field, otherwise we will get the default values for all the inner fields.
			g.If(j.Id("!v.Null && !v.Unknown")).BlockFunc(func(g *j.Group) {
				g.If(j.Id("obj." + f.ParentIsOptionalEmbedFieldName).Op("==").Nil()).Block(
					j.Id("obj." + f.ParentIsOptionalEmbedFieldName).Op("=").Id("&" + f.ParentIsOptionalEmbedFullType + "{}"),
				)
				g.Id("obj." + f.Name).Op("=").Id("t")
			})
			return
		}

		g.Id("obj." + f.Name).Op("=").Id("t")
	})
}

// genObject generates CopyFrom fragment for a nested object
func (f *FieldCopyFromGenerator) genObject() *j.Statement {
	m := NewMessageCopyFromGenerator(f.Message, f.i)
	objFieldName := "obj." + f.Name

	return f.nextField(func(g *j.Group) {
		if f.OneOfName == "" {
			if f.IsNullable {
				// obj.Nested = nil
				g.Id(objFieldName).Op("=").Nil()
			} else {
				// obj.Nested = Nested{}
				g.Id(objFieldName).Op("=").Id(f.i.WithType(f.GoElemType)).Values()
			}
			// if !v.Null
			g.If(j.Id("!v.Null && !v.Unknown")).BlockFunc(func(g *j.Group) {
				if !m.IsEmpty {
					// tf := v
					g.Id("tf").Op(":=").Id("v")

					if f.IsNullable {
						// obj.Nested = &Nested{}
						g.Id(objFieldName).Op("=&").Id(f.i.WithType(f.GoElemTypeIndirect)).Values()
						// obj := obj.Nested
						g.Id("obj").Op(":=").Id(objFieldName)
					} else {
						// obj := &obj.Nested
						g.Id("obj").Op(":=&").Id(objFieldName)
					}

					m.GenerateFields(g)
				}
			})
		} else {
			// We do not need nullable checks because all oneOf branches are nullable by design
			// We do not need to assign OneOf explicitly to not overrite other OneOf branch values
			g.If(j.Id("!v.Null && !v.Unknown")).BlockFunc(func(g *j.Group) {
				g.Id("b").Op(":=&").Id(f.i.WithType(f.GoElemTypeIndirect)).Values()

				g.Id("obj." + f.OneOfName).Op("=").Id("&" + f.i.WithType(f.OneOfType)).Values(j.Dict{
					j.Id(f.Name): j.Id("b"),
				})

				if !m.IsEmpty {
					g.Id("obj").Op(":=").Id("b")
					g.Id("tf").Op(":=").Id("v")
					m.GenerateFields(g)
				}
			})
		}
	})
}

// getValueField returns list/map value field
func (f *FieldCopyFromGenerator) getValueField() *Field {
	if f.IsMap {
		return f.MapValueField
	}

	return f.Field
}

// genPrimitiveListOrMap generates CopyFrom fragment for a list or a map of a primitive values
func (f *FieldCopyFromGenerator) genPrimitiveListOrMap() *j.Statement {
	objFieldName := "obj." + f.Name

	field := f.getValueField()

	return f.nextField(func(g *j.Group) {
		f.genListOrMapIterator(g, j.Id(f.i.WithType(field.ElemValueType)), func(g *j.Group) {
			f.genPrimitiveBody(g)

			// obj.List[k] = t
			g.Id(objFieldName).Index(j.Id("k")).Op("=").Id("t")
		})
	})
}

// genNestedListOrMap generates CopyFrom fragment for a list or a map of a nested objects
func (f *FieldCopyFromGenerator) genObjectListOrMap() *j.Statement {
	objFieldName := "obj." + f.Name

	field := f.getValueField()
	m := NewMessageCopyFromGenerator(field.Message, f.i)

	return f.nextField(func(g *j.Group) {
		f.genListOrMapIterator(g, j.Id(f.i.WithType(field.ElemValueType)), func(g *j.Group) {
			// var t Nested || *Nested
			g.Var().Id("t").Id(f.i.WithType(f.GoElemType))

			g.If(j.Id("!v.Null && !v.Unknown")).BlockFunc(func(g *j.Group) {
				// tf := v
				g.Id("tf").Op(":=").Id("v")

				if f.IsNullable {
					// t = &Nested{}
					g.Id("t").Op("=&").Id(f.i.WithType(f.GoElemTypeIndirect)).Values()
					// obj := t - obj is just an alias to reuse field generator code
					g.Id("obj").Op(":=").Id("t")
				} else {
					// obj := &t
					g.Id("obj").Op(":=&").Id("t")
				}
				m.GenerateFields(g)
			})

			// obj.List[k] = t
			g.Id(objFieldName).Index(j.Id("k")).Op("=").Id("t")
		})
	})
}

// genCustom generates statement representing custom type
func (f *FieldCopyFromGenerator) genCustom() *j.Statement {
	return j.Block(
		// a, ok := ft.Attrs["key"]
		j.List(j.Id("a"), j.Id("ok")).Op(":=").Id("tf.Attrs").Index(j.Lit(f.NameSnake)),
		j.If(j.Id("!ok")).BlockFunc(f.errAttrMissingDiag),
		j.Id("CopyFrom"+f.Suffix).Params(j.Id("diags"), j.Id("a"), j.Id("&obj."+f.Name)),
	)
}
