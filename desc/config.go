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

package desc

import (
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// schemaType represents a struct used for the schema type overrides
type schemaType struct {
	// Type is a Go attr.Type struct name
	Type string `yaml:"type,omitempty"`
	// ValueType is a Go attr.Value struct name
	ValueType string `yaml:"value_type,omitempty"`
	// CastType is a Go attr.Value .Value member type
	CastType string `yaml:"cast_type,omitempty"`
}

// Config represents the plugin config
type Config struct {
	// Types is the list of top level types to export. This list must be explicit.
	//
	// Passed from command line (--terraform_out=types=types.UserV2:./_out)
	Types map[string]struct{} `yaml:"-"`
	// DurationCustomType this type name will be treated as a custom extendee of time.Duration
	DurationCustomType string `yaml:"duration_custom_type,omitempty"`
	// ExcludeFields is the list of fields to ignore.
	//
	// Passed from command line (--terraform_out=excludeFields=types.UserV2.Expires:./_out)
	ExcludeFields map[string]struct{} `yaml:"-"`
	// TargetPackageName sets the name of the target package
	TargetPackageName string `yaml:"target_package_name,omitempty"`
	// DefaultPackageName default package name, gets appended to type name if its import
	// path is ".", but the type itself is located in another package
	DefaultPackageName string `yaml:"default_package_name,omitempty"`
	// ComputedFields is the list of fields to mark as 'Computed: true'
	//
	// Passed from command line (--terraform_out=computed=types.UserV2.Kind:./_out)
	ComputedFields map[string]struct{} `yaml:"-"`
	// RequiredFields is the list of fields to mark as 'Required: true'
	//
	// Passed from command line (--terraform_out=required=types.Metadata.Name:./_out)
	RequiredFields map[string]struct{} `yaml:"-"`
	// SensitiveFields is the list of fields to mark as 'Sensitive: true'
	//
	// Passed from command line (--terraform_out=sensitive=types.Token.Name:./_out)
	SensitiveFields map[string]struct{} `yaml:"-"`
	// CustomImports adds imports required in target file
	CustomImports []string `yaml:"custom_imports,omitempty"`
	// Suffixes represents map of suffixes for custom types
	Suffixes map[string]string `yaml:"suffixes,omitempty"`
	// NameOverrides represents map of CamelCased field names to under_score field names if needs replacement
	NameOverrides map[string]string `yaml:"name_overrides,omitempty"`
	// Validators represents map of validators for a fields
	Validators map[string][]string `yaml:"validators,omitempty"`
	// SchemaTypes represents a map of a schema field type overrides
	SchemaTypes map[string]schemaType `yaml:"schema_types,omitempty"`
	// Sort sort fields and messages by name (otherwise, will keep the order as it was in .proto file)
	Sort bool `yaml:"sort,omitempty"`
	// TimeType represents time.Time type for the Terraform Framework if set in SchemaTypes
	TimeType *schemaType `yaml:"time_type,omitempty"`
	// DurationType represents time.Duration type for the Terraform Framework if set in SchemaTypes
	DurationType *schemaType `yaml:"duration_type,omitempty"`

	// TypesRaw types loaded from a yaml file as is
	TypesRaw []string `yaml:"types,omitempty"`
	// ComputedFieldsRaw computed fields loaded from a yaml file as is
	ComputedFieldsRaw []string `yaml:"computed_fields,omitempty"`
	// RequiredFieldsRaw required fields loaded from a yaml file as is
	RequiredFieldsRaw []string `yaml:"required_fields,omitempty"`
	// SensitiveFieldsRaw sensitive fields loaded from a yaml file as is
	SensitiveFieldsRaw []string `yaml:"sensitive_fields,omitempty"`
	// ForceNewFieldsRaw force new fields loaded from a yaml file as is
	ForceNewFieldsRaw []string `yaml:"force_new_fields,omitempty"`
	// ExcludeFieldsRaw exclude fields loaded from a yaml file as is
	ExcludeFieldsRaw []string `yaml:"exclude_fields,omitempty"`

	// params represents CLI params passed from the plugin
	params map[string]string `yaml:"-"`
}

const (
	paramDelimiter = "+" // Delimiter for arrays in CLI param
)

// ReadConfig reads creates configuration instance from the CLI params
func ReadConfig(params map[string]string) (*Config, error) {
	c := &Config{params: params}

	err := c.readFromYaml()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	err = c.readFromCLI()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	err = c.validate()
	if err != nil {
		return nil, trace.Wrap(err)
	}

	c.dump()

	return c, nil
}

// getStringParam trims and returns string CLI param value
func (c *Config) getStringParam(name string, d string) string {
	p := strings.TrimSpace(c.params[name])
	if p == "" {
		return d
	}
	return p
}

// getSliceParam trims, splits to elements and returns []string CLI param value
func (c *Config) getSliceParam(name string, d []string) []string {
	v := c.getStringParam(name, "")

	// Prevents returning slice with an empty single element
	if v == "" {
		return d
	}

	return strings.Split(v, paramDelimiter)
}

// getBoolParam returns bool CLI param value, false by default
func (c *Config) getBoolParam(name string, d bool) bool {
	a := strings.ToLower(c.getStringParam(name, ""))
	if a == "" {
		return d
	}

	b, err := strconv.ParseBool(a)
	if err != nil {
		log.Printf("Failed to parse %v bool param from the %v, defaulting to %v.", name, a, d)
		return d
	}

	return b
}

// decodeRawValue converts []string (usual type for a field list) to map[string]struct{}
func (c *Config) decodeRawValue(v []string) map[string]struct{} {
	r := make(map[string]struct{})
	for _, n := range v {
		r[n] = struct{}{}
	}
	return r
}

// readFromYaml reads configuration from a yaml file passed in the config parameter
func (c *Config) readFromYaml() error {
	p := c.getStringParam("config", "")

	if p == "" {
		return nil
	}

	cfg, err := ioutil.ReadFile(p)
	if err != nil {
		return trace.Wrap(err)
	}

	err = yaml.Unmarshal(cfg, &c)
	if err != nil {
		return trace.Wrap(err)
	}

	return nil
}

// readFromCLI reads configuration from the CLI param
func (c *Config) readFromCLI() error {
	c.TypesRaw = c.getSliceParam("types", c.TypesRaw)
	c.ExcludeFieldsRaw = c.getSliceParam("exclude_fields", c.ExcludeFieldsRaw)
	c.ComputedFieldsRaw = c.getSliceParam("computed_fields", c.ComputedFieldsRaw)
	c.RequiredFieldsRaw = c.getSliceParam("required_fields", c.RequiredFieldsRaw)
	c.ForceNewFieldsRaw = c.getSliceParam("force_new", c.ForceNewFieldsRaw)
	c.SensitiveFieldsRaw = c.getSliceParam("sensitive", c.SensitiveFieldsRaw)

	c.CustomImports = c.getSliceParam("custom_imports", c.CustomImports)
	c.DefaultPackageName = c.getStringParam("default_package_name", c.DefaultPackageName)
	c.TargetPackageName = c.getStringParam("target_package_name", c.TargetPackageName)
	c.DurationCustomType = c.getStringParam("custom_duration", c.DurationCustomType)
	c.Sort = c.getBoolParam("sort", c.Sort)

	return nil
}

// validate validates the configuration
func (c *Config) validate() error {
	c.Types = c.decodeRawValue(c.TypesRaw)
	c.ExcludeFields = c.decodeRawValue(c.ExcludeFieldsRaw)
	c.ComputedFields = c.decodeRawValue(c.ComputedFieldsRaw)
	c.RequiredFields = c.decodeRawValue(c.RequiredFieldsRaw)
	c.SensitiveFields = c.decodeRawValue(c.SensitiveFieldsRaw)

	if len(c.Types) == 0 {
		return trace.Errorf("Please, specify explicit top level type list, e.g. --terraform-out=types=UserV2+UserSpecV2:./_out")
	}

	return nil
}

// dump prints configuration
func (c *Config) dump() {
	c.logMap("Types: %s", c.Types)
	c.logMap("Excluded fields: %s", c.ExcludeFields)

	if c.DefaultPackageName != "" {
		log.Printf("Default package name: %v", c.DefaultPackageName)
	}

	if c.TargetPackageName != "" {
		log.Printf("Target package name: %v", c.TargetPackageName)
	}

	if len(c.CustomImports) > 0 {
		log.Printf("Custom imports: %s", c.CustomImports)
	}

	if c.DurationCustomType != "" {
		log.Printf("Duration custom type: %s", c.DurationCustomType)
	}

	if c.TimeType != nil {
		log.Printf("Custom time type set: %v", c.TimeType.Type)
	}

	if c.DurationType != nil {
		log.Printf("Custom duration type set: %v", c.DurationType.Type)
	}

	c.logMap("Computed fields: %s", c.ComputedFields)
	c.logMap("Required fields: %s", c.RequiredFields)
	c.logMap("Suffixes set for: %v", c.Suffixes)
	c.logMap("Field name replacements set for: %v", c.NameOverrides)
	c.logMap("Sensitive flags set for: %v", c.SensitiveFields)
	c.logMap("Validatos set for: %v", c.Validators)
	c.logMap("Schema types set for: %v", c.SchemaTypes)

	if c.Sort {
		log.Printf("Sorting is enabled")
	}
}

// logMap outputs keys of a map[string]struct{} if any
func (c *Config) logMap(f string, m interface{}) {
	k := reflect.ValueOf(m).MapKeys()
	if len(k) > 0 {
		log.Printf(f, k)
	}
}
