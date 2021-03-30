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

// Package config contains global configuration variables and methods to parse them
package config

import (
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	// Types is the list of top level types to export. This list must be explicit.
	//
	// Passed from command line (--terraform_out=types=types.UserV2:./_out)
	Types map[string]struct{} = make(map[string]struct{})

	// ExcludeFields is the list of fields to ignore.
	//
	// Passed from command line (--terraform_out=excludeFields=types.UserV2.Expires:./_out)
	ExcludeFields map[string]struct{} = make(map[string]struct{})

	// DurationCustomType this type name will be treated as a custom extendee of time.Duration
	DurationCustomType = ""

	// DefaultPackageName default package name, gets appended to type name if its import
	// path is ".", but the type itself is located in another package
	DefaultPackageName string

	// CustomImports adds imports required in target file
	CustomImports []string

	// TargetPackageName sets the name of the target package
	TargetPackageName string

	// ComputedFields is the list of fields to mark as 'Computed: true'
	//
	// Passed from command line (--terraform_out=computed=types.UserV2.Kind:./_out)
	ComputedFields map[string]struct{} = make(map[string]struct{})

	// RequiredFields is the list of fields to mark as 'Required: true'
	//
	// Passed from command line (--terraform_out=required=types.Metadata.Name:./_out)
	RequiredFields map[string]struct{} = make(map[string]struct{})
)

const (
	paramDelimiter = "+" // Delimiter for types and ignoreFields
)

// MustSet sets config variables from command line
func MustSet(p map[string]string) {
	mustSetTypes(p["types"])
	setExcludeFields(p["exclude_fields"])
	setComputedFields(p["computed"])
	setRequiredFields(p["required"])
	setDefaultPackageName(p["pkg"])
	setDurationType(p["custom_duration"])
	setCustomImports(p["custom_imports"])
	setTargetPackageName(p["target_pkg"])
}

// mustSetTypes parses and sets Types.
// Panics if the argument is not a valid type list
func mustSetTypes(arg string) {
	t := strings.Split(arg, paramDelimiter)

	if len(t) == 0 {
		logrus.Fatal("Please, specify explicit top level type list, e.g. --terraform-out=types=UserV2+UserSpecV2:./_out")
	}

	for _, n := range t {
		Types[n] = struct{}{}
	}

	logrus.Printf("Types: %s", t)
}

// setExcludeFields parses and sets ExcludeFields
func setExcludeFields(arg string) {
	f := getSetFromArg(arg, ExcludeFields)

	if len(f) > 0 {
		logrus.Printf("Excluded fields: %s", f)
	}
}

// setDefaultPackageName sets the default package name
func setDefaultPackageName(arg string) {
	if trimArg(arg) == "" {
		return
	}

	_, name := filepath.Split(arg)
	DefaultPackageName = name

	logrus.Printf("Default package name: %v", DefaultPackageName)
}

// setDurationType sets the custom duration type
func setDurationType(arg string) {
	if trimArg(arg) != "" {
		DurationCustomType = arg
	}

	logrus.Printf("Duration custom type: %s", DurationCustomType)
}

// setCustomImports parses custom import packages
func setCustomImports(arg string) {
	if trimArg(arg) == "" {
		return
	}

	CustomImports = strings.Split(arg, paramDelimiter)

	logrus.Printf("Custom imports: %s", CustomImports)
}

// setTargetPackageName sets the target package name
func setTargetPackageName(arg string) {
	if trimArg(arg) == "" {
		return
	}

	_, name := filepath.Split(arg)
	TargetPackageName = name

	logrus.Printf("Target package name: %v", TargetPackageName)
}

// setComputedFields parses and sets ExcludeFields
func setComputedFields(arg string) {
	f := getSetFromArg(arg, ComputedFields)

	if len(f) > 0 {
		logrus.Printf("Computed fields: %s", f)
	}
}

// setRequiredFields parses and sets ExcludeFields
func setRequiredFields(arg string) {
	f := getSetFromArg(arg, RequiredFields)

	if len(f) > 0 {
		logrus.Printf("Required fields: %s", f)
	}
}

// trimArg returns argument value without spaces and line breaks
func trimArg(s string) string {
	return strings.Trim(s, " \n\r")
}

// getSetFromArg parses array argument and converts it to set
func getSetFromArg(arg string, r map[string]struct{}) []string {
	if trimArg(arg) == "" {
		return nil
	}

	f := strings.Split(arg, paramDelimiter)

	for _, n := range f {
		r[n] = struct{}{}
	}

	return f
}
