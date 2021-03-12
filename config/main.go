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
	Types []string

	// ExcludeFields is the list of fields to ignore.
	//
	// Passed from command line (--terraform_out=excludeFields=types.UserV2.Expires:./_out)
	ExcludeFields []string

	// DurationCustomType this type name will be treated as a custom extendee of time.Duration
	DurationCustomType = ""

	// DefaultPkgName default package name, gets appended to type name if its import
	// path is ".", but the type itself is located in another package
	DefaultPkgName string

	// CustomImports adds imports required in target file
	CustomImports []string

	// TargetPkgName sets the name of the target package
	TargetPkgName string
)

const (
	paramDelimiter = "+" // Delimiter for types and ignoreFields
)

// MustParseTypes parses and sets Types.
// Panics if the argument is not a valid type list
func MustSetTypes(arg string) {
	Types = strings.Split(arg, paramDelimiter)

	if len(Types) == 0 {
		logrus.Fatal("Please, specify explicit top level type list, e.g. --terraform-out=types=UserV2+UserSpecV2:./_out")
	}

	logrus.Printf("Types: %s", Types)
}

// ParseExcludeFields parses and sets ExcludeFields
func SetExcludeFields(arg string) {
	if arg == "" {
		return
	}

	ExcludeFields = strings.Split(arg, paramDelimiter)

	logrus.Printf("Excluded fields: %s", ExcludeFields)
}

// SetDefaultPkgName sets the default package name
func SetDefaultPkgName(arg string) {
	if arg == "" {
		return
	}

	_, name := filepath.Split(arg)
	DefaultPkgName = name

	logrus.Printf("Default package name: %v", DefaultPkgName)
}

// SetDuration sets the custom duration type
func SetDuration(arg string) {
	if arg != "" {
		DurationCustomType = arg
	}

	logrus.Printf("Duration custom type: %s", DurationCustomType)
}

// ParseCustomImports parses custom import packages
func SetCustomImports(arg string) {
	if arg == "" {
		return
	}

	CustomImports = strings.Split(arg, paramDelimiter)

	logrus.Printf("Custom imports: %s", CustomImports)
}

// SetTargetPkgName sets the target package name
func SetTargetPkgName(arg string) {
	if arg == "" {
		return
	}

	_, name := filepath.Split(arg)
	TargetPkgName = name

	logrus.Printf("Target package name: %v", TargetPkgName)
}
