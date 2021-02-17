package config

import (
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
)

const (
	paramDelimiter = "+" // Delimiter for types and ignoreFields
)

// ParseTypes parses and sets Types
func ParseTypes(arg string) {
	Types = strings.Split(arg, paramDelimiter)

	if len(Types) == 0 {
		logrus.Fatal("Please, specify explicit top level type list, eg. --terraform-out=types=UserV2+UserSpecV2:./_out")
	}

	logrus.Printf("Types: %s", Types)
}

// ParseExcludeFields parses and sets ExcludeFields
func ParseExcludeFields(arg string) {
	ExcludeFields = strings.Split(arg, paramDelimiter)

	logrus.Printf("Excluded fields: %s", ExcludeFields)
}
