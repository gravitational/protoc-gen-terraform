/*
Copyright 2015-2020 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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

	if n1 == -1 || n2 == -1 {
		logrus.Warning("Package directive not found in target file, can't replace package name, skipping")
		return s
	}

	r := s[0:n1+8] + config.TargetPkgName + s[n1+n2:]

	return r
}
