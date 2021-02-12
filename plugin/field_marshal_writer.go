package plugin

import "strings"

// fieldMarshalWriter represents logic required to generate field read go code
type fieldMarshalWriter struct {
	writer
	field *Field

	rawVarName  string
	castVarName string
	castGoName  string
}

// newFieldMarshalWriter returns an instance of FieldWriter
func newFieldMarshalWriter(f *Field) *fieldMarshalWriter {
	return &fieldMarshalWriter{
		field:       f,
		rawVarName:  `_` + f.NameSnake + `_raw`,
		castVarName: `_` + f.NameSnake + `_cast`,
		castGoName:  `_` + f.NameSnake + `_go_cast`,
	}
}

// Write generates code required for field
func (w *fieldMarshalWriter) write() string {
	if w.field.HasNestedMessage() {
		// w.pNested()
	} else {
		if w.field.IsAggregate() {
			if w.field.IsRepeated {
				// w.pList()
			}
		} else {
			w.pGetRawValue()
			w.pIfOk()

			w.pCastToSchemaGoType()
			w.pCastToTargetGoType()
			w.pAssign()

			w.pEndIfOk()
		}
	}

	return w.buf.String()
}

func (w *fieldMarshalWriter) pGetRawValue() {
	w.p(w.rawVarName, `, ok := d.GetOk(prefix + "`, w.field.NameSnake, `")`)
}

func (w *fieldMarshalWriter) pCastToSchemaGoType() {
	w.p(w.castVarName, ` := `, w.rawVarName, `.(`, w.field.TFSchemaTypeCast, `)`)
}

func (w *fieldMarshalWriter) pCastToTargetGoType() {
	if w.isTime() {
		// TODO: Handle error
		w.p(w.castGoName, `, ok := time.Parse(time.RFC3339, `, w.castVarName, `)`)
		w.p(`if !ok {`)
		w.p(`  return fmt.Errorf("Malformed time value for field `, w.field.Name, `")`)
		w.p(`}`)
	} else if w.isDuration() {
		// TODO: Handle error
		w.p(w.castGoName, `, ok := time.ParseDuration(`, w.castVarName, `)`)
		w.p(`if !ok {`)
		w.p(`  return fmt.Errorf("Malformed duration for field `, w.field.Name, `")`)
		w.p(`}`)
	} else {
		w.p(w.castGoName, ` := `, w.castVarName, `.(`, w.field.GoTypeCast, `)`)
	}
}

func (w *fieldMarshalWriter) pIfOk() {
	w.p(`if (ok) {`)
}

func (w *fieldMarshalWriter) pEndIfOk() {
	w.p(`}`)
}

func (w *fieldMarshalWriter) pAssign() {
	w.a(`t.`, w.field.Name, ` = `, w.castGoNameWithPtr())
}

func (w *fieldMarshalWriter) castGoNameWithPtr() string {
	s := ""
	if w.field.IsNullable {
		s = s + "&"
	}

	return s + w.castGoName
}

func (w *fieldMarshalWriter) isTime() bool {
	return strings.HasSuffix(w.field.GoTypeCast, "time.Time")
}

func (w *fieldMarshalWriter) isDuration() bool {
	return strings.HasSuffix(w.field.GoTypeCast, "time.Duration")
}

// func (w *FieldWriter) writeReadRaw() {

// }

// func (w *FieldWriter) writeReadRawList() {

// }

// func (w *FieldWriter) writeCastRawToSchemaType() {

// }

// func (w *FieldWriter) writeCastRawListToSchemaType() {

// }

// func (w *FieldWriter) writeCastSchemaTypeToStructType() {

// }

// func (w *FieldWriter) writeCastSchemaTypeToStructTypeList() {

// }

// // Message writes message descriptor
// func Message(g *generator.Generator, d *generator.Descriptor) {
// 	//logrus.Println("// ", d.File().GetPackage(), ".", d.GetName())
// 	//g.P("// ", d.File().GetPackage(), ".", d.GetName())

// 	//logrus.Println(d.GetName())
// 	//logrus.Println(d.File().GetPackage())

// 	g.P("func Unmarshal", d.GetName(), "(s *ResourceDataFacade, t *", d.GetName(), ", string prefix) {")

// 	logrus.Println(d.TypeName(), d.File().GoPackageName())
// 	logrus.Println(d.GoImportPath())

// 	for _, f := range d.GetField() {
// 		m := buildFieldMarshalMeta(g, d, f)

// 		if m.assignable {
// 			g.P(`_`, m.snakeName, ` := `, `s.Get("`+m.snakeName+`")`)
// 			//g.P(`t.`, m.name, ` = `, m.sprintfValue())
// 		} else {
// 			// g.P(`// t.`, m.Name, ` `, m.GoType)
// 		}

// 		// if f.IsMessage() && !f.IsRepeated() {
// 		// 	g.P(snakeFieldName, ` := NewFacadeFromNestedMap(s.Get("`, snakeFieldName, `").([]interface{})[0].(map[string]interface{}))`)
// 		// 	g.P(`Unmarshal`, fieldName, `(`, snakeFieldName, `,`, `t.`, fieldName, `)`)
// 		// }

// 		//

// 		// logrus.Println("  ", f.GetName(), " ", g.IsMap(f), " ", f.IsRepeated(), " ", gt)
// 		// if g.IsMap(f) {
// 		// 	gf, _ := g.GoType(d, g.GoMapType(nil, f).ValueField)
// 		// 	logrus.Println("      ", gf)
// 		// }

// 		//logrus.Println(f.GetTypeName())
// 		//logrus.Println(gogoproto.GetCastType(f))
// 		//logrus.Println(gogoproto.GetCustomType(f))

// 	}

// 	g.P("}")
// 	g.P()
// }
