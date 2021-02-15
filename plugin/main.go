package plugin

import (
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gravitational/trace"
	"github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"github.com/stretchr/stew/slice"
)

const (
	name          = "terraform"                                                      // Plugin name
	schemaPkg     = "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"     // Terraform schema package
	validationPkg = "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation" // Terraform validation package
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

	// Map of reflected messages, public just in case some post analysis is required
	Messages map[string]*Message

	// // NOTE: Replace with addImport
	// pkg           generator.Single // Reference to package with protoc types
	// referencePackages map[string]string
}

// NewPlugin creates the new plugin
func NewPlugin() *Plugin {
	return &Plugin{
		Messages: make(map[string]*Message),
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
		p.reflectMessage(message)
	}

	for _, message := range p.Messages {
		buf, err := message.GoString()
		if err != nil {
			p.Generator.Fail(trace.Wrap(err).Error())
		}
		p.P(buf.String())
		//p.P(newMessageMarshalWriter(message).write())
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

	// So those could be referenced via schema. and validation.
	p.AddImport(schemaPkg)
	p.AddImport(validationPkg)
}

// isMessageRequired returns true if message was marked for export via command-line args
func (p *Plugin) isMessageRequired(d *generator.Descriptor) bool {
	typeName := d.File().GetPackage() + "." + d.GetName()
	required := slice.Contains(p.types, typeName)

	if !required {
		logrus.Println("Skipping type:", typeName)
	}

	return required
}

// reflectMessage reflects message type
// todo: builder function?
func (p *Plugin) reflectMessage(d *generator.Descriptor) *Message {
	if !p.isMessageRequired(d) {
		return nil
	}

	name := d.GetName()

	if p.Messages[name] != nil {
		return p.Messages[name]
	}

	message := &Message{}

	message.Name = name
	message.NameSnake = strcase.SnakeCase(name)
	message.GoTypeName = d.File().GetPackage() + "." + name

	p.reflectFields(message, d)

	p.Messages[name] = message

	return message
}

// reflectFields builds array of message.Fields
func (p *Plugin) reflectFields(m *Message, d *generator.Descriptor) {
	for _, f := range d.GetField() {
		f := p.reflectField(d, f)
		if f != nil {
			m.Fields = append(m.Fields, f)
		}
	}
}
