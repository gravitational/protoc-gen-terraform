/*
Copyright 2015-2021 Gravitational, Inc.

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
	"os"
	"regexp"
	"strings"

	plugin_go "github.com/gogo/protobuf/protoc-gen-gogo/plugin"
	"github.com/gogo/protobuf/vanity/command"
	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/imports"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"

	_ "embed"
)

var (
	// packageReplacementRegexp is used to replace package name in a target file
	packageReplacementRegexp = regexp.MustCompile("package (.+)\n")

	//go:embed license.txt
	license string
)

func main() {
	log.Infof("protoc-gen-terraform %s", Version)

	if len(os.Args) > 1 && os.Args[1] == "version" {
		return
	}

	p := NewPlugin()

	req := command.Read()
	resp := command.GeneratePlugin(req, p, "_terraform.go")

	err := runGoImports(p, resp)
	if err != nil {
		p.Fail(err.Error())
	}

	// Convert the gogo response to a regular protobuf response. This allows us
	// to pack in the SupportedFeatures field, which indicates that the optional
	// field is supported.
	response := &pluginpb.CodeGeneratorResponse{}
	response.Error = resp.Error
	response.File = make([]*pluginpb.CodeGeneratorResponse_File, 0, len(resp.File))
	for _, file := range resp.File {
		response.File = append(response.File, &pluginpb.CodeGeneratorResponse_File{
			Name:           file.Name,
			InsertionPoint: file.InsertionPoint,
			Content:        file.Content,
		})
	}
	features := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	response.SupportedFeatures = &features

	// Send back the results. The code below was taken from the vanity command,
	// but it now uses the regular response instead of the gogo specific one.
	data, err := proto.Marshal(response)
	if err != nil {
		p.Fail(err.Error(), "failed to marshal output proto")
	}
	_, err = os.Stdout.Write(data)
	if err != nil {
		p.Fail(err.Error(), "failed to write output proto")
	}
}

// runGoImports formats code and removes unused imports from the resulting code using goimports tool
func runGoImports(p *Plugin, resp *plugin_go.CodeGeneratorResponse) error {
	opts := imports.Options{
		FormatOnly: false,
		Comments:   true,
	}

	for _, file := range resp.GetFile() {
		if file.Content == nil {
			continue
		}

		result, err := imports.Process("", []byte(*file.Content), &opts)
		if err != nil {
			return trace.Wrap(err)
		}

		s := string(result)

		s, err = prependLicense(p, s)
		if err != nil {
			return trace.Wrap(err)
		}

		if p.Config.TargetPackageName != "" {
			s = replacePackageName(s, p.Config.TargetPackageName)
		}

		file.Content = &s
	}

	return nil
}

// prependLicense prepends license information
func prependLicense(p *Plugin, s string) (string, error) {
	return license + s, nil
}

// replacePackageName replaces package name in target file with provided from cli
func replacePackageName(s string, target string) string {
	// Replace one string
	pkg := packageReplacementRegexp.FindString(s)
	if pkg == "" {
		log.Warning("Package directive not found in target file, can't replace package name, skipping")
		return s
	}

	return strings.Replace(s, pkg, "package "+target+"\n", 1)
}
