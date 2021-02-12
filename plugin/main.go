package plugin

import (
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/stew/slice"
)

const (
	name = "terraform" // Plugin name
)

// Plugin is terraform generator plugin
type Plugin struct {
	*generator.Generator
	generator.PluginImports

	// The list of types to export. This list must be explicit. In case, a type is not listed, it would not be
	// present in output along with all fields which reference it. This is the way of keeping a structures private.
	//
	// Passed from command line (--terraform_out=types=types.UserV2:./_out)
	types []string

	messages map[string]*Message // Map of reflected messages

	// // NOTE: Replace with addImport
	// schemaPkg     generator.Single // Reference to terraform schema package
	// validationPkg generator.Single // Reference to terraform validation package
	// pkg           generator.Single // Reference to package with protoc types
}

// NewPlugin creates the new plugin
func NewPlugin() *Plugin {
	return &Plugin{
		messages: make(map[string]*Message),
	}
}

// Init initializes plugin and sets the generator instance
func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g
	p.fetchTypesFromCommandLine(p.Generator)
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return name
}

// Generate is the plugin body
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	logrus.Printf("Processing: %s", *file.Name)

	p.setImports()

	for _, message := range file.Messages() {
		if p.isRequiredMessage(message) {
			r := p.reflectMessage(message)
			p.registerMessage(r)
		}
	}

	for _, message := range p.messages {
		NewMessageSchemaWriter(message).Write()
	}
}

// fetchTypesFromCommandLine loads loads and parses type list from command line
func (p *Plugin) fetchTypesFromCommandLine(g *generator.Generator) {
	if g.Param["types"] == "" {
		logrus.Fatal("Please, specify explicit top level type list, eg. --terraform-out=types=UserV2+UserSpecV2:./_out")
	}

	p.types = strings.Split(g.Param["types"], "+")

	if len(p.types) == 0 {
		if g.Param["types"] == "" {
			logrus.Fatal("Types list is malformed or empty!")
		}
	}

	logrus.Printf("Types: %s", p.types)
}

// setImports sets import definitions for current file
func (p *Plugin) setImports() {
	p.PluginImports = generator.NewPluginImports(p.Generator)

	// TODO: Save names?
	p.AddImport("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
	p.AddImport("github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation")
}

// isMessageRequired returns true if message was marked for export via command-line args
func (p *Plugin) isRequiredMessage(d *generator.Descriptor) bool {
	typeName := d.File().GoPackageName() + "." + d.GetName()
	return slice.Contains(p.types, typeName)
}

// registerMessage puts message in the list of registered
func (p *Plugin) registerMessage(m *Message) {
	p.messages[m.Name] = m
}

// schemaRef returns type name with reference to terraform schema ns
// func (p *Plugin) schemaRef(ref string) string {
// 	return p.schemaPkg.Use() + "." + ref
// }

// // writeMessage writes message go code to output buffer
// func (p *Plugin) writeMessage(r *messageReflect) {
// 	p.P(`// `, r.name)
// 	p.P(`func Schema`, r.name, `() *`, p.schemaRef(`Schema {`))
// 	p.WriteString(`  return `)
// 	p.writeSchema(r, false)
// 	p.P(`}`)

// 	p.P()

// 	p.P(`func Unmarshal`, r.name, `(r *`, p.schemaRef(`ResourceData`), `, t *`, r.goType, `, prefix string) {`)
// 	p.writeUnmarshal(r)
// 	p.P(`}`)
// 	p.P()
// }

// // writeSchema writes schema definition
// func (p *Plugin) writeSchema(r *messageReflect, comma bool) {
// 	p.P(`map[string]*`, p.schemaRef(`Schema {`))

// 	for _, f := range r.fields {
// 		p.P(`"`, f.snakeName, `": {`)

// 		if f.hasNestedType {
// 			p.P(`Type: `, p.schemaRef("TypeList"), `,`)
// 			p.P(`Optional: true,`)
// 			p.P(`MaxItems: 1,`)
// 			p.P(`Elem: &`, p.schemaRef("Resource"), `{`)
// 			p.WriteString(`  Schema: `)
// 			p.writeSchema(f.message, true)
// 			p.P(`},`)

// 		} else {
// 			if f.tfSchemaCollectionType != "" {
// 				p.P(`Type: `, p.schemaRef(f.tfSchemaCollectionType), `,`)
// 				p.P(`Optional: true,`)
// 				p.P(`Elem: &`, p.schemaRef(`Schema {`))
// 				p.writeSchemaScalar(f)
// 				p.P(`},`)
// 			} else {
// 				p.writeSchemaScalar(f)
// 				p.P(`Optional: true,`)
// 			}
// 		}

// 		p.P(`},`)
// 	}

// 	if comma == true {
// 		p.P(`},`)
// 	} else {
// 		p.P(`}`)
// 	}
// }

// func (p *Plugin) writeSchemaScalar(f *fieldReflect) {
// 	// This is scalar value
// 	if f.tfSchemaType != "" {
// 		p.P(`Type: `, p.schemaRef(f.tfSchemaType), `,`)

// 		if f.tfSchemaValidate != "" {
// 			p.P(`Validate: `, p.validationPkg.Use(), `.`, f.tfSchemaValidate, `,`)
// 		}
// 	}
// }

// // writeUnmarshal writes unmarshaling function
// func (p *Plugin) writeUnmarshal(r *messageReflect) {
// 	for _, f := range r.fields {
// 		if f.tfSchemaType != "" {
// 			p.P(`_`, f.snakeName, `, ok := r.GetOk(prefix + "`, f.snakeName+`").(`, f.tfSchemaGoType, `)`)
// 			p.P(`if (ok) {`)
// 			p.P(`}`)
// 		}
// 	}
// }
