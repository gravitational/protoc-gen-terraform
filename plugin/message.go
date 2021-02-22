package plugin

import (
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gzigzigzeo/protoc-gen-terraform/config"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/stew/slice"

	"fmt"
)

var (
	cache map[string]*Message = make(map[string]*Message)
)

// Message holds reflection information about message
type Message struct {
	Name       string   // Type name
	NameSnake  string   // Type name in snake case, schema field name
	GoTypeName string   // Go type name for this message with package name
	Fields     []*Field // Collection of fields
	// TODO: Comments
}

// BuildMessage builds Message from it's protobuf descriptor
// checkValiditiy is should be false for nested fields, otherwise we'll have to be over-explicit in allowed type
func BuildMessage(g *generator.Generator, d *generator.Descriptor, checkValidity bool) *Message {
	typeName := getMessageTypeName(d)

	if !slice.Contains(config.Types, typeName) && checkValidity {
		return nil
	}

	if cache[typeName] != nil {
		return cache[typeName]
	}

	for _, field := range d.GetField() {
		if field.OneofIndex != nil {
			logrus.Println("Oneof messages are not supported yet")
			return nil
		}
	}

	name := d.GetName()

	message := &Message{
		Name:       name,
		NameSnake:  strcase.SnakeCase(name),
		GoTypeName: typeName,
	}

	BuildFields(message, g, d)

	return message
}

// getMessageTypeName returns full message name, with prepended DefaultPkgName if needed
func getMessageTypeName(d *generator.Descriptor) string {
	if d.GoImportPath() == "." {
		if config.DefaultPkgName != "" {
			return config.DefaultPkgName + "." + d.GetName()
		} else {
			return d.GetName()
		}
	}
	return d.File().GetPackage() + "." + d.GetName()
}

// NOTE: temp comments render
func (m *Message) GoTypeMapString(prefixa string) string {
	b := strings.Builder{}

	//b.WriteString(fmt.Sprintf("//%-30v\n", prefixa+m.GoTypeName))

	for _, f := range m.Fields {
		s := fmt.Sprintf("// %-40v %-50v %-25v %-7v\n", prefixa+f.Name, f.GoType, f.Kind, f.IsMap)
		b.WriteString(s)

		if f.IsMessage {
			b.WriteString(f.Message.GoTypeMapString(prefixa + "  "))
		}
	}

	return b.String()
}
