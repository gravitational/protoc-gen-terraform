package testlib

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

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
					resource.TestCheckResourceAttr(name, "bool_list.1", "false"),
					resource.TestCheckResourceAttr(name, "bytes_list.0", "bytes1"),
					resource.TestCheckResourceAttr(name, "bytes_list.1", "bytes2"),
					resource.TestCheckResourceAttr(name, "enum_list.0", "1"),
					resource.TestCheckResourceAttr(name, "enum_list.1", "2"),
				),
			},
		},
	})
}

func (s *TerraformSuite) TestPrimitivesZeroValues() {
	t := s.T()
	name := "example_primitives.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("primitives_zero_values.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "string_value", ""),
					resource.TestCheckResourceAttr(name, "int32_value", "0"),
					resource.TestCheckResourceAttr(name, "int64_value", "0"),
					resource.TestCheckResourceAttr(name, "float_value", "0"),
					resource.TestCheckResourceAttr(name, "double_value", "0"),
					resource.TestCheckResourceAttr(name, "bool_value", "false"),
					resource.TestCheckResourceAttr(name, "bytes_value", ""),
					resource.TestCheckResourceAttr(name, "enum_value", "0"),
					resource.TestCheckResourceAttr(name, "string_list.0", ""),
					resource.TestCheckResourceAttr(name, "string_list.1", ""),
					resource.TestCheckResourceAttr(name, "int32_list.0", "0"),
					resource.TestCheckResourceAttr(name, "int32_list.1", "0"),
					resource.TestCheckResourceAttr(name, "int64_list.0", "0"),
					resource.TestCheckResourceAttr(name, "int64_list.1", "0"),
					resource.TestCheckResourceAttr(name, "float_list.0", "0"),
					resource.TestCheckResourceAttr(name, "float_list.1", "0"),
					resource.TestCheckResourceAttr(name, "double_list.0", "0"),
					resource.TestCheckResourceAttr(name, "double_list.1", "0"),
					resource.TestCheckResourceAttr(name, "bool_list.0", "false"),
					resource.TestCheckResourceAttr(name, "bool_list.1", "false"),
					resource.TestCheckResourceAttr(name, "bytes_list.0", ""),
					resource.TestCheckResourceAttr(name, "bytes_list.1", ""),
					resource.TestCheckResourceAttr(name, "enum_list.0", "0"),
					resource.TestCheckResourceAttr(name, "enum_list.1", "0"),
				),
			},
		},
	})
}

func (s *TerraformSuite) TestPrimitivesUpdate() {
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
					resource.TestCheckResourceAttr(name, "bool_list.1", "false"),
					resource.TestCheckResourceAttr(name, "bytes_list.0", "bytes1"),
					resource.TestCheckResourceAttr(name, "bytes_list.1", "bytes2"),
					resource.TestCheckResourceAttr(name, "enum_list.0", "1"),
					resource.TestCheckResourceAttr(name, "enum_list.1", "2"),
				),
			},
			{
				Config:   s.getFixture("primitives.tf"),
				PlanOnly: true,
			},
			{
				Config: s.getFixture("primitives_zero_values.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "string_value", ""),
					resource.TestCheckResourceAttr(name, "int32_value", "0"),
					resource.TestCheckResourceAttr(name, "int64_value", "0"),
					resource.TestCheckResourceAttr(name, "float_value", "0"),
					resource.TestCheckResourceAttr(name, "double_value", "0"),
					resource.TestCheckResourceAttr(name, "bool_value", "false"),
					resource.TestCheckResourceAttr(name, "bytes_value", ""),
					resource.TestCheckResourceAttr(name, "enum_value", "0"),
					resource.TestCheckResourceAttr(name, "string_list.0", ""),
					resource.TestCheckResourceAttr(name, "string_list.1", ""),
					resource.TestCheckResourceAttr(name, "int32_list.0", "0"),
					resource.TestCheckResourceAttr(name, "int32_list.1", "0"),
					resource.TestCheckResourceAttr(name, "int64_list.0", "0"),
					resource.TestCheckResourceAttr(name, "int64_list.1", "0"),
					resource.TestCheckResourceAttr(name, "float_list.0", "0"),
					resource.TestCheckResourceAttr(name, "float_list.1", "0"),
					resource.TestCheckResourceAttr(name, "double_list.0", "0"),
					resource.TestCheckResourceAttr(name, "double_list.1", "0"),
					resource.TestCheckResourceAttr(name, "bool_list.0", "false"),
					resource.TestCheckResourceAttr(name, "bool_list.1", "false"),
					resource.TestCheckResourceAttr(name, "bytes_list.0", ""),
					resource.TestCheckResourceAttr(name, "bytes_list.1", ""),
					resource.TestCheckResourceAttr(name, "enum_list.0", "0"),
					resource.TestCheckResourceAttr(name, "enum_list.1", "0"),
				),
			},
			{
				Config:   s.getFixture("primitives_zero_values.tf"),
				PlanOnly: true,
			},
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
					resource.TestCheckResourceAttr(name, "bool_list.1", "false"),
					resource.TestCheckResourceAttr(name, "bytes_list.0", "bytes1"),
					resource.TestCheckResourceAttr(name, "bytes_list.1", "bytes2"),
					resource.TestCheckResourceAttr(name, "enum_list.0", "1"),
					resource.TestCheckResourceAttr(name, "enum_list.1", "2"),
				),
			},
			{
				Config:   s.getFixture("primitives.tf"),
				PlanOnly: true,
			},
		},
	})
}
