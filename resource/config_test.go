package resource

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	cfg, err := ReadConfig(map[string]string{"config": "fixtures/config.yaml", "types": "foo+bar"})
	require.NoError(t, err)

	require.Equal(t, cfg.Types, flagMap{"foo": struct{}{}, "bar": struct{}{}})
	require.Equal(t, cfg.DurationCustomType, "Duration")
	require.Equal(t, cfg.UseStateForUnknownByDefault, true)
	require.Equal(t, cfg.Sort, true)
	require.Equal(t, cfg.TargetPackageName, "test")

	require.Equal(t, cfg.ExcludeFields, flagMap{"Test.Excluded": struct{}{}})
	require.Equal(t, cfg.ComputedFields, flagMap{"Test.Str": struct{}{}})
	require.Equal(t, cfg.RequiredFields, flagMap{"Test.RequiredStr": struct{}{}})
	require.Equal(t, cfg.SensitiveFields, flagMap{"Test.Str": struct{}{}})

	require.Equal(t, cfg.Suffixes, map[string]string{"BoolCustom": "BoolSpecial"})
	require.Equal(t, cfg.NameOverrides, map[string]string{"Test.Str": "str"})

	require.Equal(t, cfg.PlanModifiers, map[string][]string{"Test.Str": {"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier.UseStateForUnknown()"}})
	require.Equal(t, cfg.Validators, map[string][]string{"Test.Str": {"UseMockValidator()"}})

	require.Equal(t, cfg.TimeType, &SchemaType{
		Type:               "TimeType",
		ValueType:          "TimeValue",
		ValueFromMethod:    "ValueTime",
		ValueToMethod:      "NewTime",
		NullValueMethod:    "NullTime",
		UnknownValueMethod: "UnknownTime",
		CastToType:         "time.Time",
		CastFromType:       "time.Time",
		TypeConstructor:    "UseRFC3339Time()",
	})

	require.Equal(t, cfg.DurationType, &SchemaType{
		Type:               "DurationType",
		ValueType:          "DurationValue",
		ValueFromMethod:    "ValueDuration",
		ValueToMethod:      "NewDuration",
		NullValueMethod:    "NullDuration",
		UnknownValueMethod: "UnknownDuration",
		CastToType:         "time.Duration",
		CastFromType:       "time.Duration",
		TypeConstructor:    "DurationType{}",
	})

	require.Equal(t, cfg.InjectedFields, map[string][]InjectedField{
		"Test": {{
			Name:               "id",
			Type:               "github.com/hashicorp/terraform-plugin-framework/types.StringType",
			Computed:           true,
			DefaultValueMethod: "github.com/hashicorp/terraform-plugin-framework/types.StringNull",
		}},
	})
	require.Equal(t, cfg.SchemaTypes, map[string]SchemaType{
		"Test.SchemaOverride": {
			Type:               "github.com/hashicorp/terraform-plugin-framework/types.StringType",
			ValueType:          "github.com/hashicorp/terraform-plugin-framework/types.String",
			ValueFromMethod:    "ValueString",
			ValueToMethod:      "github.com/hashicorp/terraform-plugin-framework/types.StringValue",
			NullValueMethod:    "github.com/hashicorp/terraform-plugin-framework/types.StringNull",
			UnknownValueMethod: "github.com/hashicorp/terraform-plugin-framework/types.StringUnknown",
			CastToType:         "string",
			CastFromType:       "OverrideCastType",
		},
	})

}

func TestConfigSchemaTypeValidation(t *testing.T) {
	path := writeConfig(t, `
types:
  - Test
time_type:
  type: TimeType
  value_type: TimeValue
  cast_to_type: time.Time
  cast_from_type: time.Time
`)

	_, err := ReadConfig(map[string]string{"config": path})
	require.Error(t, err)
	require.Contains(t, err.Error(), "time_type.value_from_method is required")
}

func TestConfigInjectedFieldValidation(t *testing.T) {
	path := writeConfig(t, `
types:
  - Test
injected_fields:
  Test:
    - name: id
      type: github.com/hashicorp/terraform-plugin-framework/types.StringType
      computed: true
`)

	_, err := ReadConfig(map[string]string{"config": path})
	require.Error(t, err)
	require.Contains(t, err.Error(), "injected_fields.Test.id.default_value_method is required")
}

func writeConfig(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))

	return path
}
