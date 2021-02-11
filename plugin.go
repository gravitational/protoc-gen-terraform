package main

import (
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/stew/slice"
)

// Plugin is terraform generator plugin
type Plugin struct {
	*generator.Generator
	generator.PluginImports

	types         []string         // The list of types to export passed from command line (--terraform_out=types=types.UserV2:./_out)
	schemaPkg     generator.Single // Reference to terraform schema package
	validationPkg generator.Single // Reference to terraform validation package
	pkg           generator.Single // Reference to package with protoc types
}

// NewPlugin creates the new plugin
func NewPlugin() *Plugin {
	return &Plugin{}
}

// Init initializes plugin and sets the generator instance
func (p *Plugin) Init(g *generator.Generator) {
	p.Generator = g

	if g.Param["types"] == "" {
		logrus.Fatal("Please, specify explicit top level type list, eg. --terraform-out=types=UserV2+UserSpecV2:./_out")
	}

	p.types = strings.Split(g.Param["types"], "+")

	if len(p.types) == 0 {
		if g.Param["types"] == "" {
			logrus.Fatal("Types list is malformed!:./_out")
		}
	}

	logrus.Printf("Types: %s", p.types)
}

// Name returns the name of the plugin
func (p *Plugin) Name() string {
	return "terraform"
}

// Generate is the plugin body
func (p *Plugin) Generate(file *generator.FileDescriptor) {
	logrus.Printf("Processing: %s", *file.Name)

	p.setImports()
	p.setBaseImport(file)

	for _, message := range file.Messages() {
		if p.isMessageRequired(message) {
			r := p.newMessageReflect(message)
			p.writeMessage(r)
		}
	}

	//logrus.Println(file.GetDependency())
	//p.P("// FILE --", file.GetName())

	//imp := p.NewImport("basepackage.com/wrappers")
	//p.P("// ", imp.Use(), ".sometype")
	//writer.Message(p.Generator, message)
	//logrus.Println(message.GetName())
}

// setImports sets import definitions for current file
func (p *Plugin) setImports() {
	p.PluginImports = generator.NewPluginImports(p.Generator)
	p.schemaPkg = p.NewImport("github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema")
	p.validationPkg = p.NewImport("github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation")
}

// setBaseImport sets the reference to a package containing target types
func (p *Plugin) setBaseImport(file *generator.FileDescriptor) {
	p.pkg = p.NewImport(file.GoPackageName())
}

// isMessageRequired returns true if message was marked for export via command-line args
func (p *Plugin) isMessageRequired(d *generator.Descriptor) bool {
	typeName := d.File().GoPackageName() + "." + d.GetName()
	return slice.Contains(p.types, typeName)
}

// schemaRef returns type name with reference to terraform schema ns
func (p *Plugin) schemaRef(ref string) string {
	return p.schemaPkg.Use() + "." + ref
}

// writeMessage writes message go code to output buffer
func (p *Plugin) writeMessage(r *messageReflect) {
	p.P(`// `, r.name)
	p.P(`func Schema`, r.name, `() *`, p.schemaRef(`Schema {`))
	p.WriteString(`  return `)
	p.writeSchema(r)
	p.P(`}`)

	p.P()

	p.P(`func Unmarshal`, r.name, `(r *`, p.schemaRef(`ResourceData`), `, t *`, r.goType, `, prefix string) {`)
	p.writeUnmarshal(r)
	p.P(`}`)
	p.P()
}

// writeSchema writes schema definition
func (p *Plugin) writeSchema(r *messageReflect) {
	p.P(`map[string]*`, p.schemaRef(`Schema {`))

	for _, f := range r.fields {
		p.P(`// Cast to `, f.tfSchemaGoType)
		p.P(`"`, f.snakeName, `": {`)

		// This is scalar value
		if f.tfSchemaType != "" {
			p.P(`Type: `, p.schemaRef(f.tfSchemaType), `,`)
			p.P(`Optional: true,`)

			if f.tfSchemaValidate != "" {
				p.P(`Validate: `, p.validationPkg.Use(), `.`, f.tfSchemaValidate, `,`)
			}
		}

		p.P(`},`)
	}

	p.P(`}`)
}

// writeUnmarshal writes unmarshaling function
func (p *Plugin) writeUnmarshal(r *messageReflect) {
	for _, f := range r.fields {
		if f.tfSchemaType != "" {
			p.P(`_`, f.snakeName, `, ok := r.GetOk(prefix + "`, f.snakeName+`").(`, f.tfSchemaGoType, `)`)
			p.P(`if (ok) {`)
			p.P(`}`)
		}
	}
}
