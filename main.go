package main

import (
	"strings"

	plugin_go "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
	"github.com/gogo/protobuf/vanity/command"
	"github.com/gravitational/protoc-gen-terraform/config"
	"github.com/gravitational/protoc-gen-terraform/plugin"
	"github.com/sirupsen/logrus"
	"golang.org/x/tools/imports"
)

func main() {
	logrus.Infof("protoc-gen-terraform %s", Version)
	logrus.Infof("protoc-gen-terraform build hash: %s", Sha)

	p := plugin.NewPlugin()

	req := command.Read()
	resp := command.GeneratePlugin(req, p, "_terraform.go")

	err := fmt(resp)
	if err != nil {
		p.Fail(err.Error())
	}

	command.Write(resp)
}

// fmt removes unused imports from the resulting code
func fmt(resp *plugin_go.CodeGeneratorResponse) error {
	opts := imports.Options{
		FormatOnly: false,
		Comments:   true,
	}

	for _, file := range resp.GetFile() {
		result, err := imports.Process("", []byte(*file.Content), &opts)
		if err != nil {
			return err
		}

		s := string(result)
		s = replacePackageName(s)
		file.Content = &s
	}

	return nil
}

// replacePackageName replaces package name in target file with provided from cli
func replacePackageName(s string) string {
	if config.TargetPkgName == "" {
		return s
	}

	n1 := strings.Index(s, "package")
	n2 := strings.Index(s[n1:], "\n")

	r := s[0:n1+8] + config.TargetPkgName + s[n1+n2:]

	return r
}
