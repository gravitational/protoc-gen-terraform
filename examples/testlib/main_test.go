package testlib

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/suite"

	"github.com/gravitational/protoc-gen-terraform/v3/examples/testlib/provider"
)

func TestTerraform(t *testing.T) {
	suite.Run(t, &TerraformSuite{})
}

type TerraformSuite struct {
	suite.Suite

	terraformProvider tfsdk.Provider

	terraformProviders map[string]func() (tfprotov6.ProviderServer, error)
}

func (s *TerraformSuite) SetupSuite() {
	os.Setenv("TF_ACC", "true")
	s.terraformProvider = provider.New()
	s.terraformProviders = make(map[string]func() (tfprotov6.ProviderServer, error))
	s.terraformProviders["example"] = func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6(s.terraformProvider)(), nil
	}
}

//go:embed fixtures/*
var fixtures embed.FS

// getFixture loads fixture and returns it as string or <error> if failed
func (s *TerraformSuite) getFixture(name string, formatArgs ...any) string {
	b, err := fixtures.ReadFile(filepath.Join("fixtures", name))
	if err != nil {
		return fmt.Sprintf("<error: %v fixture not found>", name)
	}
	return fmt.Sprintf(string(b), formatArgs...)
}

func (s *TerraformSuite) TestPrimitives() {
	t := s.T()
	name := "example_primitives.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("primitives.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "string_value", "string"),
					resource.TestCheckResourceAttr(name, "int32_value", "123"),
					resource.TestCheckResourceAttr(name, "int64_value", "456"),
					resource.TestCheckResourceAttr(name, "float_value", "0.75"),
					resource.TestCheckResourceAttr(name, "double_value", "0.75"),
					resource.TestCheckResourceAttr(name, "bool_value", "true"),
					resource.TestCheckResourceAttr(name, "bytes_value", "bytes"),
					resource.TestCheckResourceAttr(name, "enum_value", "1"),
					resource.TestCheckResourceAttr(name, "string_list.0", "el1"),
					resource.TestCheckResourceAttr(name, "string_list.1", "el2"),
					resource.TestCheckResourceAttr(name, "int32_list.0", "123"),
					resource.TestCheckResourceAttr(name, "int32_list.1", "456"),
					resource.TestCheckResourceAttr(name, "int64_list.0", "234"),
					resource.TestCheckResourceAttr(name, "int64_list.1", "567"),
					resource.TestCheckResourceAttr(name, "float_list.0", "0.75"),
					resource.TestCheckResourceAttr(name, "float_list.1", "1.25"),
					resource.TestCheckResourceAttr(name, "double_list.0", "0.75"),
					resource.TestCheckResourceAttr(name, "double_list.1", "1.25"),
					resource.TestCheckResourceAttr(name, "bool_list.0", "true"),
					// TODO: Bool false value is treated as null.
					// This should not be the case...
					// resource.TestCheckResourceAttr(name, "bool_list.1", "false"),
					resource.TestCheckResourceAttr(name, "bytes_list.0", "bytes1"),
					resource.TestCheckResourceAttr(name, "bytes_list.1", "bytes2"),
					resource.TestCheckResourceAttr(name, "enum_list.0", "1"),
					resource.TestCheckResourceAttr(name, "enum_list.1", "2"),
					resource.TestCheckNoResourceAttr(name, "nullable_value"),
				),
			},
		},
	})
}

func (s *TerraformSuite) TestTime() {
	t := s.T()
	name := "example_time.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("time.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "timestamp_value", "2026-01-02T03:04:05Z"),
					resource.TestCheckResourceAttr(name, "timestamp_list.0", "2026-01-02T03:04:05Z"),
					resource.TestCheckResourceAttr(name, "timestamp_list.1", "2026-01-02T03:04:06Z"),

					resource.TestCheckResourceAttr(name, "duration_standard", "5m0s"),
					resource.TestCheckResourceAttr(name, "duration_list.0", "5m0s"),
					resource.TestCheckResourceAttr(name, "duration_list.1", "10m0s"),

					resource.TestCheckResourceAttr(name, "duration_custom", "5m0s"),
					resource.TestCheckResourceAttr(name, "duration_custom_list.0", "5m0s"),
					resource.TestCheckResourceAttr(name, "duration_custom_list.1", "10m0s"),
				),
			},
		},
	})
}

func (s *TerraformSuite) TestObjects() {
	t := s.T()
	name := "example_objects.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("objects.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "primitives.string_value", "string"),
					resource.TestCheckResourceAttr(name, "primitives.int32_value", "123"),
					resource.TestCheckResourceAttr(name, "primitives.float_value", "0.75"),
					resource.TestCheckResourceAttr(name, "primitives.bool_value", "true"),
					resource.TestCheckResourceAttr(name, "primitives.enum_value", "1"),
					resource.TestCheckNoResourceAttr(name, "primitives.nullable_value"),

					resource.TestCheckResourceAttr(name, "string_map.key1", "value1"),
					resource.TestCheckResourceAttr(name, "string_map.key2", "value2"),
					resource.TestCheckResourceAttr(name, "int_map.one", "1"),
					resource.TestCheckResourceAttr(name, "int_map.two", "2"),
					resource.TestCheckResourceAttr(name, "bool_map.enabled", "true"),
					resource.TestCheckResourceAttr(name, "bool_map.disabled", "false"),

					resource.TestCheckResourceAttr(name, "nested_value.leaf.value", "nested-value"),
					resource.TestCheckNoResourceAttr(name, "nested_nullable"),

					resource.TestCheckResourceAttr(name, "nested_list.0.leaf.value", "list-1"),
					resource.TestCheckResourceAttr(name, "nested_list.1.leaf.value", "list-2"),
					resource.TestCheckNoResourceAttr(name, "nested_nullable_list"),

					resource.TestCheckResourceAttr(name, "nested_map.first.leaf.value", "map-1"),
					resource.TestCheckResourceAttr(name, "nested_map.second.leaf.value", "map-2"),
					resource.TestCheckNoResourceAttr(name, "nested_nullable_map"),

					resource.TestCheckResourceAttr(name, "branch1.leaf.value", "branch-1"),
					resource.TestCheckNoResourceAttr(name, "branch2"),

					resource.TestCheckResourceAttr(name, "leaf.value", "embedded-leaf"),
					// TODO: Unepxected behavior with embedded fields.
					// This embedded value is overwrites the leaf.value field...
					// resource.TestCheckResourceAttr(name, "value", "embedded-nullable-value"),
				),
			},
		},
	})
}

func (s *TerraformSuite) TestCustom() {
	t := s.T()
	name := "example_custom.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("custom.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "computed", "computed"),
					resource.TestCheckResourceAttr(name, "injected", "injected"),
					resource.TestCheckResourceAttr(name, "required", "required"),
					resource.TestCheckResourceAttr(name, "sensitive", "sensitive"),
					resource.TestCheckNoResourceAttr(name, "excluded"),

					resource.TestCheckResourceAttr(name, "bool_custom_list.0", "true"),
					resource.TestCheckResourceAttr(name, "bool_custom_list.1", "false"),
					resource.TestCheckResourceAttr(name, "bool_custom_list.2", "true"),

					resource.TestCheckResourceAttr(name, "string_override.0", "foo"),
					resource.TestCheckResourceAttr(name, "string_override.1", "bar"),
					resource.TestCheckResourceAttr(name, "schema_override", "schema-override"),
				),
			},
		},
	})
}
