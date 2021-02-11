package writer

import (
	"fmt"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
)

type fieldMarshalMeta struct {
	name              string // struct field name
	snakeName         string // struct field name, camel cased
	goType            string // go type this marshalled into
	assignable        bool   // can be directly assigned
	byReferencePrefix string // & prefix or empty
	castToSuffix      string // .() suffix or empty

	// than it can be map or list or just nested resource (list with 1 item)

	valueFmt string // constructor function, for time
}

func buildFieldMarshalMeta(
	g *generator.Generator,
	d *generator.Descriptor,
	f *descriptor.FieldDescriptorProto,
) *fieldMarshalMeta {
	m := &fieldMarshalMeta{}

	goType, _ := g.GoType(d, f)

	m.goType = goType
	m.name = g.GetFieldName(d, f)
	m.snakeName = strcase.SnakeCase(m.name)
	m.valueFmt = "%s"

	if !gogoproto.IsNullable(f) {
		m.byReferencePrefix = "&"
	}

	if gogoproto.IsStdTime(f) {
		m.valueFmt = "time.Parse(time.RFC3339, %v)"
		m.assignable = true
	} else if f.IsMessage() {

		// Here goes map
	} else {
		m.castToSuffix = `.(` + m.goType + `)`
		m.assignable = true
	}

	return m
}

// Formats value
func (m *fieldMarshalMeta) sprintfValue(value string) string {
	return m.byReferencePrefix + fmt.Sprintf(m.valueFmt, value) + m.castToSuffix
}

// Message writes message descriptor
func Message(g *generator.Generator, d *generator.Descriptor) {
	//logrus.Println("// ", d.File().GetPackage(), ".", d.GetName())
	//g.P("// ", d.File().GetPackage(), ".", d.GetName())

	//logrus.Println(d.GetName())
	//logrus.Println(d.File().GetPackage())

	g.P("func Unmarshal", d.GetName(), "(s *ResourceDataFacade, t *", d.GetName(), ", string prefix) {")

	logrus.Println(d.TypeName(), d.File().GoPackageName())
	logrus.Println(d.GoImportPath())

	for _, f := range d.GetField() {
		m := buildFieldMarshalMeta(g, d, f)

		if m.assignable {
			g.P(`_`, m.snakeName, ` := `, `s.Get("`+m.snakeName+`")`)
			//g.P(`t.`, m.name, ` = `, m.sprintfValue())
		} else {
			// g.P(`// t.`, m.Name, ` `, m.GoType)
		}

		// if f.IsMessage() && !f.IsRepeated() {
		// 	g.P(snakeFieldName, ` := NewFacadeFromNestedMap(s.Get("`, snakeFieldName, `").([]interface{})[0].(map[string]interface{}))`)
		// 	g.P(`Unmarshal`, fieldName, `(`, snakeFieldName, `,`, `t.`, fieldName, `)`)
		// }

		//

		// logrus.Println("  ", f.GetName(), " ", g.IsMap(f), " ", f.IsRepeated(), " ", gt)
		// if g.IsMap(f) {
		// 	gf, _ := g.GoType(d, g.GoMapType(nil, f).ValueField)
		// 	logrus.Println("      ", gf)
		// }

		//logrus.Println(f.GetTypeName())
		//logrus.Println(gogoproto.GetCastType(f))
		//logrus.Println(gogoproto.GetCustomType(f))

	}

	g.P("}")
	g.P()
}
