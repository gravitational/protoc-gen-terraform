package testlib

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

func (s *TerraformSuite) TestTime() {
	t := s.T()
	name := "example_time.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("time.tf"),
				Check:  s.testCheckTimeResource(name),
			},
		},
	})
}

func (s *TerraformSuite) TestTimeZeroValues() {
	t := s.T()
	name := "example_time.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("time_zero_values.tf"),
				Check:  s.testCheckTimeZeroValuesResource(name),
			},
		},
	})
}

func (s *TerraformSuite) TestTimeUpdate() {
	t := s.T()
	name := "example_time.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("time.tf"),
				Check:  s.testCheckTimeResource(name),
			},
			{
				Config:   s.getFixture("time.tf"),
				PlanOnly: true,
			},
			{
				Config: s.getFixture("time_zero_values.tf"),
				Check:  s.testCheckTimeZeroValuesResource(name),
			},
			{
				Config:   s.getFixture("time_zero_values.tf"),
				PlanOnly: true,
			},
			{
				Config: s.getFixture("time.tf"),
				Check:  s.testCheckTimeResource(name),
			},
			{
				Config:   s.getFixture("time.tf"),
				PlanOnly: true,
			},
		},
	})
}

func (s *TerraformSuite) TestTimeNullValues() {
	t := s.T()
	name := "example_time.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("time_null_values.tf"),
				Check:  s.testCheckTimeNullValuesResource(name),
			},
		},
	})
}

func (s *TerraformSuite) testCheckTimeResource(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(name, "timestamp_value", "2026-01-02T03:04:05Z"),
		resource.TestCheckResourceAttr(name, "timestamp_list.0", "2026-01-02T03:04:05Z"),
		resource.TestCheckResourceAttr(name, "timestamp_list.1", "2026-01-02T03:04:06Z"),

		resource.TestCheckResourceAttr(name, "duration_standard", "5m0s"),
		resource.TestCheckResourceAttr(name, "duration_list.0", "5m0s"),
		resource.TestCheckResourceAttr(name, "duration_list.1", "10m0s"),

		resource.TestCheckResourceAttr(name, "duration_custom", "5m0s"),
		resource.TestCheckResourceAttr(name, "duration_custom_list.0", "5m0s"),
		resource.TestCheckResourceAttr(name, "duration_custom_list.1", "10m0s"),

		resource.TestCheckResourceAttr(name, "nullable_timestamp", "2026-01-02T03:04:05Z"),
		resource.TestCheckResourceAttr(name, "nullable_duration", "5m0s"),
	)
}

func (s *TerraformSuite) testCheckTimeZeroValuesResource(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(name, "timestamp_value", "0001-01-01T00:00:00Z"),
		resource.TestCheckResourceAttr(name, "timestamp_list.0", "0001-01-01T00:00:00Z"),
		resource.TestCheckResourceAttr(name, "timestamp_list.1", "0001-01-01T00:00:00Z"),

		resource.TestCheckResourceAttr(name, "duration_standard", "0s"),
		resource.TestCheckResourceAttr(name, "duration_list.0", "0s"),
		resource.TestCheckResourceAttr(name, "duration_list.1", "0s"),

		resource.TestCheckResourceAttr(name, "duration_custom", "0s"),
		resource.TestCheckResourceAttr(name, "duration_custom_list.0", "0s"),
		resource.TestCheckResourceAttr(name, "duration_custom_list.1", "0s"),

		resource.TestCheckResourceAttr(name, "nullable_timestamp", "0001-01-01T00:00:00Z"),
		resource.TestCheckResourceAttr(name, "nullable_duration", "0s"),
	)
}

func (s *TerraformSuite) testCheckTimeNullValuesResource(name string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttr(name, "timestamp_value", "0001-01-01T00:00:00Z"),
		resource.TestCheckResourceAttr(name, "duration_standard", "0s"),
		resource.TestCheckResourceAttr(name, "duration_custom", "0s"),

		resource.TestCheckNoResourceAttr(name, "timestamp_list.0"),
		resource.TestCheckNoResourceAttr(name, "duration_list.0"),
		resource.TestCheckNoResourceAttr(name, "duration_custom_list.0"),

		resource.TestCheckNoResourceAttr(name, "nullable_timestamp"),
		resource.TestCheckNoResourceAttr(name, "nullable_duration"),
	)
}
