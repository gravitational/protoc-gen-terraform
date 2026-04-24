package testlib

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

func (s *TerraformSuite) TestObjects() {
	t := s.T()
	name := "example_objects.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("objects.tf"),
				Check:  s.testCheckObjectResource(name),
			},
		},
	})
}

func (s *TerraformSuite) TestObjectsZeroValues() {
	t := s.T()
	name := "example_objects.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("objects_zero_values.tf"),
				Check:  s.testCheckObjectZeroValuesResource(name),
			},
		},
	})
}

func (s *TerraformSuite) TestObjectsUpdate() {
	t := s.T()
	name := "example_objects.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("objects.tf"),
				Check:  s.testCheckObjectResource(name),
			},
			{
				Config:   s.getFixture("objects.tf"),
				PlanOnly: true,
			},
			{
				Config: s.getFixture("objects_zero_values.tf"),
				Check:  s.testCheckObjectZeroValuesResource(name),
			},
			{
				Config:   s.getFixture("objects_zero_values.tf"),
				PlanOnly: true,
			},
			{
				Config: s.getFixture("objects.tf"),
				Check:  s.testCheckObjectResource(name),
			},
			{
				Config:   s.getFixture("objects.tf"),
				PlanOnly: true,
			},
		},
	})
}

func (s *TerraformSuite) TestObjectsNullValues() {
	t := s.T()
	name := "example_objects.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("objects_null_values.tf"),
				Check:  s.testCheckObjectNullValuesResource(name),
			},
		},
	})
}

func (s *TerraformSuite) testCheckObjectResource(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(name, "primitives.string_value", "string"),
		resource.TestCheckResourceAttr(name, "primitives.int32_value", "123"),
		resource.TestCheckResourceAttr(name, "primitives.float_value", "0.75"),
		resource.TestCheckResourceAttr(name, "primitives.bool_value", "true"),
		resource.TestCheckResourceAttr(name, "primitives.enum_value", "1"),

		resource.TestCheckResourceAttr(name, "string_map.key1", "value1"),
		resource.TestCheckResourceAttr(name, "string_map.key2", "value2"),
		resource.TestCheckResourceAttr(name, "int_map.key1", "1"),
		resource.TestCheckResourceAttr(name, "int_map.key2", "2"),
		resource.TestCheckResourceAttr(name, "bool_map.key1", "true"),
		resource.TestCheckResourceAttr(name, "bool_map.key2", "false"),

		resource.TestCheckResourceAttr(name, "nested_value.leaf.value", "nested-value"),
		resource.TestCheckResourceAttr(name, "nested_nullable.leaf.value", "nested-value"),

		resource.TestCheckResourceAttr(name, "nested_list.0.leaf.value", "list-1"),
		resource.TestCheckResourceAttr(name, "nested_list.1.leaf.value", "list-2"),
		resource.TestCheckResourceAttr(name, "nested_nullable_list.0.leaf.value", "list-1"),
		resource.TestCheckResourceAttr(name, "nested_nullable_list.1.leaf.value", "list-2"),

		resource.TestCheckResourceAttr(name, "nested_map.key1.leaf.value", "map-1"),
		resource.TestCheckResourceAttr(name, "nested_map.key2.leaf.value", "map-2"),
		resource.TestCheckResourceAttr(name, "nested_nullable_map.key1.leaf.value", "map-1"),
		resource.TestCheckResourceAttr(name, "nested_nullable_map.key2.leaf.value", "map-2"),

		resource.TestCheckResourceAttr(name, "branch1.leaf.value", "branch-1"),
		resource.TestCheckNoResourceAttr(name, "branch2"),

		resource.TestCheckResourceAttr(name, "branch_string", "branch-string"),
		resource.TestCheckNoResourceAttr(name, "branch_bool"),
		resource.TestCheckNoResourceAttr(name, "branch_int"),

		resource.TestCheckResourceAttr(name, "leaf.value", "embedded-leaf"),
		resource.TestCheckResourceAttr(name, "embedded_value", "embedded-nullable-value"),
	)
}

func (s *TerraformSuite) testCheckObjectZeroValuesResource(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(name, "primitives.string_value", ""),
		resource.TestCheckResourceAttr(name, "primitives.int32_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives.float_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives.bool_value", "false"),
		resource.TestCheckResourceAttr(name, "primitives.enum_value", "0"),

		resource.TestCheckResourceAttr(name, "string_map.key1", ""),
		resource.TestCheckResourceAttr(name, "string_map.key2", ""),
		resource.TestCheckResourceAttr(name, "int_map.key1", "0"),
		resource.TestCheckResourceAttr(name, "int_map.key2", "0"),
		resource.TestCheckResourceAttr(name, "bool_map.key1", "false"),
		resource.TestCheckResourceAttr(name, "bool_map.key2", "false"),

		resource.TestCheckResourceAttr(name, "nested_value.leaf.value", ""),
		resource.TestCheckNoResourceAttr(name, "nested_nullable"),

		resource.TestCheckResourceAttr(name, "nested_list.0.leaf.value", ""),
		resource.TestCheckResourceAttr(name, "nested_list.1.leaf.value", ""),
		resource.TestCheckNoResourceAttr(name, "nested_nullable_list"),

		resource.TestCheckResourceAttr(name, "nested_map.key1.leaf.value", ""),
		resource.TestCheckResourceAttr(name, "nested_map.key2.leaf.value", ""),
		resource.TestCheckNoResourceAttr(name, "nested_nullable_map"),

		resource.TestCheckResourceAttr(name, "branch1.leaf.value", ""),
		resource.TestCheckNoResourceAttr(name, "branch2"),

		resource.TestCheckResourceAttr(name, "branch_string", ""),
		resource.TestCheckNoResourceAttr(name, "branch_bool"),
		resource.TestCheckNoResourceAttr(name, "branch_int"),

		resource.TestCheckResourceAttr(name, "leaf.value", ""),
		resource.TestCheckResourceAttr(name, "embedded_value", ""),
	)
}

func (s *TerraformSuite) testCheckObjectNullValuesResource(name string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttr(name, "primitives.string_value", ""),
		resource.TestCheckResourceAttr(name, "primitives.int32_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives.int64_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives.float_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives.double_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives.bool_value", "false"),
		resource.TestCheckResourceAttr(name, "primitives.bytes_value", ""),
		resource.TestCheckResourceAttr(name, "primitives.enum_value", "0"),

		resource.TestCheckNoResourceAttr(name, "string_map.key1"),
		resource.TestCheckNoResourceAttr(name, "string_map.key2"),
		resource.TestCheckNoResourceAttr(name, "int_map.key1"),
		resource.TestCheckNoResourceAttr(name, "int_map.key2"),
		resource.TestCheckNoResourceAttr(name, "bool_map.key1"),
		resource.TestCheckNoResourceAttr(name, "bool_map.key2"),

		resource.TestCheckResourceAttr(name, "nested_value.leaf.value", ""),
		resource.TestCheckNoResourceAttr(name, "nested_nullable.leaf.value"),

		resource.TestCheckNoResourceAttr(name, "nested_list.0.leaf.value"),
		resource.TestCheckNoResourceAttr(name, "nested_list.1.leaf.value"),
		resource.TestCheckNoResourceAttr(name, "nested_nullable_list.0.leaf.value"),
		resource.TestCheckNoResourceAttr(name, "nested_nullable_list.1.leaf.value"),

		resource.TestCheckNoResourceAttr(name, "nested_map.key1.leaf.value"),
		resource.TestCheckNoResourceAttr(name, "nested_map.key2.leaf.value"),
		resource.TestCheckNoResourceAttr(name, "nested_nullable_map.key1.leaf.value"),
		resource.TestCheckNoResourceAttr(name, "nested_nullable_map.key2.leaf.value"),

		resource.TestCheckNoResourceAttr(name, "branch1.leaf.value"),
		resource.TestCheckNoResourceAttr(name, "branch2.leaf.value"),

		resource.TestCheckNoResourceAttr(name, "branch_bool"),
		resource.TestCheckNoResourceAttr(name, "branch_int"),
		resource.TestCheckNoResourceAttr(name, "branch_string"),

		resource.TestCheckResourceAttr(name, "leaf.value", ""),
		resource.TestCheckNoResourceAttr(name, "embedded_value"),
	)
}
