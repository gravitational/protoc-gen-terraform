package testlib

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

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
					resource.TestCheckResourceAttr(name, "validated", "valid"),
					resource.TestCheckNoResourceAttr(name, "excluded"),

					resource.TestCheckResourceAttr(name, "custom_name_override", "name-override"),
					resource.TestCheckResourceAttr(name, "plan_modifier", "modified_value"),

					resource.TestCheckResourceAttr(name, "bool_custom", "true"),
					resource.TestCheckResourceAttr(name, "bool_custom_list.0", "true"),
					resource.TestCheckResourceAttr(name, "bool_custom_list.1", "false"),

					resource.TestCheckResourceAttr(name, "string_override.0", "foo"),
					resource.TestCheckResourceAttr(name, "string_override.1", "bar"),
					resource.TestCheckResourceAttr(name, "schema_override", "schema-override"),
				),
			},
		},
	})
}
