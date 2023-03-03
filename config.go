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
	"io/ioutil"
	"reflect"
	"strconv"
	"strings"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// flagMap represents flag map
type flagMap map[string]struct{}

// UnmarshalYAML unmarshals values from yaml
func (lm *flagMap) UnmarshalYAML(unmarshal func(interface{}) error) error {
	//Unmarshal time to string then convert to time.Time manually
	var list []string
	if err := unmarshal(&list); err != nil {
		return err
	}

	*lm = flagMapFromArray(list)
	return nil
}

// listMapFrom array converts array of strings to flag map
func flagMapFromArray(v []string) flagMap {
	r := make(flagMap)
	for _, n := range v {
		r[n] = struct{}{}
	}
	return r
}

// SchemaType represents a struct used for the schema type overrides
type SchemaType struct {
	// Type is a Go attr.Type struct name
	Type string `yaml:"type,omitempty"`
	// ValueType is a Go attr.Value struct name
	ValueType string `yaml:"value_type,omitempty"`
	// CastToType is a Go attr.Value .Value member type
	CastToType string `yaml:"cast_to_type,omitempty"`
	// CastToType is a go type of the object field
	CastFromType string `yaml:"cast_from_type,omitempty"`
	// TypeConstructor represents statement used to produce empty type value in schema
	TypeConstructor string `yaml:"type_constructor,omitempty"`
}

// InjectedField represents custom injected field descriptor
type InjectedField struct {
	// Name represents field schema name
	Name string `yaml:"name,omitempty"`
	// Type represents field type
	Type string `yaml:"type,omitempty"`
	// Required is the required flag
	Required bool `yaml:"required,omitempty"`
	// Computed is the computed flag
	Computed bool `yaml:"computed,omitempty"`
	// Optional is the optional flag
	Optional bool `yaml:"optional,omitempty"`
	// PlanModifiers is the array of PlanModifiers
	PlanModifiers []string `yaml:"plan_modifiers,omitempty"`
	// PlanModifiers is the array of Validators
	Validators []string `yaml:"validators,omitempty"`
}

// Config represents the plugin config
type Config struct {
	// Types is the list of top level types to export. This list must be explicit.
	//
	// Passed from command line (--terraform_out=types=types.UserV2:./_out)
	Types flagMap `yaml:"types"`
	// DurationCustomType this type name will be treated as a custom extendee of time.Duration
	DurationCustomType string `yaml:"duration_custom_type,omitempty"`
	// ExcludeFields is the list of fields to ignore.
	//
	// Passed from command line (--terraform_out=excludeFields=types.UserV2.Expires:./_out)
	ExcludeFields flagMap `yaml:"exclude_fields"`
	// TargetPackageName sets the name of the target package
	TargetPackageName string `yaml:"target_package_name,omitempty"`
	// DefaultPackageName represents the package name of the proto-generated code
	DefaultPackageName string `yaml:"default_package_name,omitempty"`
	// Sort sort fields and messages by name (otherwise, will keep the order as it was in .proto file)
	Sort bool `yaml:"sort,omitempty"`
	// UseStateForUnknownByDefault represents flag, if true - appends UseStateForUnknown to all computed fields
	UseStateForUnknownByDefault bool `yaml:"use_state_for_unknown_by_default,omitempty"`
	// ComputedFields is the list of fields to mark as 'Computed: true'
	//
	// Passed from command line (--terraform_out=computed=types.UserV2.Kind:./_out)
	ComputedFields flagMap `yaml:"computed_fields"`
	// RequiredFields is the list of fields to mark as 'Required: true'
	//
	// Passed from command line (--terraform_out=required=types.Metadata.Name:./_out)
	RequiredFields flagMap `yaml:"required_fields"`
	// SensitiveFields is the list of fields to mark as 'Sensitive: true'
	//
	// Passed from command line (--terraform_out=sensitive=types.Token.Name:./_out)
	SensitiveFields flagMap `yaml:"sensitive_fields"`
	// Suffixes represents map of suffixes for custom types
	Suffixes map[string]string `yaml:"suffixes,omitempty"`
	// NameOverrides represents map of CamelCased field names to under_score field names
	NameOverrides map[string]string `yaml:"name_overrides,omitempty"`
	// Validators represents the map of validators for a fields
	Validators map[string][]string `yaml:"validators,omitempty"`
	// PlanModifiers represents the map of plan modifiers for a fields
	PlanModifiers map[string][]string `yaml:"plan_modifiers,omitempty"`
	// SchemaTypes represents a map of a schema field type overrides
	SchemaTypes map[string]SchemaType `yaml:"schema_types,omitempty"`
	// TimeType represents time.Time type for the Terraform Framework if set in SchemaTypes
	TimeType *SchemaType `yaml:"time_type,omitempty"`
	// DurationType represents time.Duration type for the Terraform Framework if set in SchemaTypes
	DurationType *SchemaType `yaml:"duration_type,omitempty"`
	// InjectedFields represents array of fields which are missing in object, but must be injected in the schema
	InjectedFields map[string][]InjectedField `yaml:"injected_fields,omitempty"`
	// ImportPathOverrides represents fully qualified go import paths which
	// should be used for the given package names, it can be used when the
	// correct import path is not found automatically.
	ImportPathOverrides map[string]string `yaml:"import_path_overrides,omitempty"`

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

	if len(c.Types) == 0 {
		return nil, trace.Errorf("Please, specify explicit top level type list, e.g. --terraform-out=types=UserV2+UserSpecV2:./_out")
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
func (c *Config) getSliceParam(name string, d flagMap) flagMap {
	v := c.getStringParam(name, "")

	// Prevents returning slice with an empty single element
	if v == "" {
		return d
	}

	return flagMapFromArray(strings.Split(v, paramDelimiter))
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
	c.Types = c.getSliceParam("types", c.Types)
	c.ExcludeFields = c.getSliceParam("exclude_fields", c.ExcludeFields)
	c.ComputedFields = c.getSliceParam("computed_fields", c.ComputedFields)
	c.RequiredFields = c.getSliceParam("required_fields", c.RequiredFields)
	c.SensitiveFields = c.getSliceParam("sensitive", c.SensitiveFields)

	c.DefaultPackageName = c.getStringParam("default_package_name", c.DefaultPackageName)
	c.TargetPackageName = c.getStringParam("target_package_name", c.TargetPackageName)
	c.DurationCustomType = c.getStringParam("custom_duration", c.DurationCustomType)
	c.Sort = c.getBoolParam("sort", c.Sort)

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
	c.logMap("Validators set for: %v", c.Validators)
	c.logMap("PlanModifiers set for: %v", c.PlanModifiers)
	c.logMap("Schema types set for: %v", c.SchemaTypes)

	if c.Sort {
		log.Printf("Sorting is enabled")
	}

	if c.UseStateForUnknownByDefault {
		log.Printf("StateForUnknown used by default")
	}

	if len(c.InjectedFields) > 0 {
		c.logMap("Fields are injected to: %v", c.InjectedFields)
	}
}

// logMap outputs keys of a map[string]struct{} if any
func (c *Config) logMap(f string, m interface{}) {
	k := reflect.ValueOf(m).MapKeys()
	if len(k) > 0 {
		log.Printf(f, k)
	}
}
