package main

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/gogo/protobuf/vanity/command"
	"github.com/sirupsen/logrus"
)

type stringer struct {
	*generator.Generator
	generator.PluginImports
}

func (s *stringer) Init(g *generator.Generator) {
	logrus.Info("Init() called")

	s.Generator = g
}

func (p *stringer) Name() string {
	return "tfschema"
}

func (s *stringer) Generate(file *generator.FileDescriptor) {
	s.PluginImports = generator.NewPluginImports(s.Generator)

	logrus.Println("Generate called")

	for _, message := range file.Messages() {
		if message != nil && message.Name != nil {
			logrus.Println("Chponk")
		}

		// s.P("package tfschema")
		// s.Out()
	}
}

func main() {
	logrus.Info("Starting protoc-gen-terraform")

	p := &stringer{}
	req := command.Read()
	//files := req.GetProtoFile()
	//files = vanity.FilterFiles(files, vanity.NotGoogleProtobufDescriptorProto)

	resp := command.GeneratePlugin(req, p, "_terraform.pb.go")

	command.Write(resp)
}
