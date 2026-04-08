package testlib

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

func (s *TerraformSuite) TestComputed() {
	t := s.T()
	name := "example_computed.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("computed.tf"),
				Check:  testCheckComputedResource(name),
			},
		},
	})
}

func (s *TerraformSuite) TestComputedNullValues() {
	t := s.T()
	name := "example_computed.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("computed_null_values.tf"),
				Check:  testCheckComputedResourceZeroValue(name),
			},
		},
	})
}

func (s *TerraformSuite) TestComputedUpdate() {
	t := s.T()
	name := "example_computed.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("computed.tf"),
				Check:  testCheckComputedResource(name),
			},
			{
				Config:   s.getFixture("computed.tf"),
				PlanOnly: true,
			},
			{
				Config: s.getFixture("computed_null_values.tf"),
				Check:  testCheckComputedResourceZeroValue(name),
			},
			{
				Config:   s.getFixture("computed_null_values.tf"),
				PlanOnly: true,
			},
			{
				Config: s.getFixture("computed.tf"),
				Check:  testCheckComputedResource(name),
			},
			{
				Config:   s.getFixture("computed.tf"),
				PlanOnly: true,
			},
		},
	})
}

func testCheckComputedResource(name string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttr(name, "string_value", "string"),
		resource.TestCheckResourceAttr(name, "int32_value", "123"),
		resource.TestCheckResourceAttr(name, "int64_value", "456"),
		resource.TestCheckResourceAttr(name, "float_value", "0.75"),
		resource.TestCheckResourceAttr(name, "double_value", "1.25"),
		resource.TestCheckResourceAttr(name, "bool_value", "true"),
		resource.TestCheckResourceAttr(name, "bytes_value", "bytes"),
		resource.TestCheckResourceAttr(name, "enum_value", "1"),
		resource.TestCheckResourceAttr(name, "timestamp_value", "2026-01-02T03:04:05Z"),
		resource.TestCheckResourceAttr(name, "duration_value", "5m0s"),
		resource.TestCheckResourceAttr(name, "primitives_value.string_value", "string"),
		resource.TestCheckResourceAttr(name, "primitives_value.int32_value", "123"),
		resource.TestCheckResourceAttr(name, "primitives_value.int64_value", "456"),
		resource.TestCheckResourceAttr(name, "primitives_value.float_value", "0.75"),
		resource.TestCheckResourceAttr(name, "primitives_value.double_value", "0.75"),
		resource.TestCheckResourceAttr(name, "primitives_value.bool_value", "true"),
		resource.TestCheckResourceAttr(name, "primitives_value.bytes_value", "bytes"),
		resource.TestCheckResourceAttr(name, "primitives_value.enum_value", "1"),
		resource.TestCheckResourceAttr(name, "nested_value.leaf.value", "nested-value"),
		resource.TestCheckResourceAttr(name, "nested_nullable.leaf.value", "nested-nullable"),
	)
}

func testCheckComputedResourceZeroValue(name string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(name, "id"),
		resource.TestCheckResourceAttr(name, "string_value", ""),
		resource.TestCheckResourceAttr(name, "int32_value", "0"),
		resource.TestCheckResourceAttr(name, "int64_value", "0"),
		resource.TestCheckResourceAttr(name, "float_value", "0"),
		resource.TestCheckResourceAttr(name, "double_value", "0"),
		resource.TestCheckResourceAttr(name, "bool_value", "false"),
		resource.TestCheckResourceAttr(name, "bytes_value", ""),
		resource.TestCheckResourceAttr(name, "enum_value", "0"),
		resource.TestCheckResourceAttr(name, "timestamp_value", "0001-01-01T00:00:00Z"),
		resource.TestCheckResourceAttr(name, "duration_value", "0s"),

		resource.TestCheckResourceAttr(name, "primitives_value.string_value", ""),
		resource.TestCheckResourceAttr(name, "primitives_value.int32_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives_value.int64_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives_value.float_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives_value.double_value", "0"),
		resource.TestCheckResourceAttr(name, "primitives_value.bool_value", "false"),
		resource.TestCheckResourceAttr(name, "primitives_value.bytes_value", ""),
		resource.TestCheckResourceAttr(name, "primitives_value.enum_value", "0"),

		resource.TestCheckResourceAttr(name, "nested_value.leaf.value", ""),
		resource.TestCheckNoResourceAttr(name, "nested_nullable.leaf.value"),
	)
}
