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

func (s *TerraformSuite) TestTimeZeroValues() {
	t := s.T()
	name := "example_time.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: s.terraformProviders,
		IsUnitTest:               true,
		Steps: []resource.TestStep{
			{
				Config: s.getFixture("time_zero_values.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(name, "duration_standard", "0s"),
					resource.TestCheckResourceAttr(name, "duration_list.0", "0s"),
					resource.TestCheckResourceAttr(name, "duration_list.1", "0s"),

					resource.TestCheckResourceAttr(name, "duration_custom", "0s"),
					resource.TestCheckResourceAttr(name, "duration_custom_list.0", "0s"),
					resource.TestCheckResourceAttr(name, "duration_custom_list.1", "0s"),
				),
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
			{
				Config:   s.getFixture("time.tf"),
				PlanOnly: true,
			},
			{
				Config: s.getFixture("time_zero_values.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr(name, "timestamp_value"),
					resource.TestCheckNoResourceAttr(name, "timestamp_list.0"),
					resource.TestCheckNoResourceAttr(name, "timestamp_list.1"),

					resource.TestCheckResourceAttr(name, "duration_standard", "0s"),
					resource.TestCheckResourceAttr(name, "duration_list.0", "0s"),
					resource.TestCheckResourceAttr(name, "duration_list.1", "0s"),

					resource.TestCheckResourceAttr(name, "duration_custom", "0s"),
					resource.TestCheckResourceAttr(name, "duration_custom_list.0", "0s"),
					resource.TestCheckResourceAttr(name, "duration_custom_list.1", "0s"),
				),
			},
			{
				Config:   s.getFixture("time_zero_values.tf"),
				PlanOnly: true,
			},
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
			{
				Config:   s.getFixture("time.tf"),
				PlanOnly: true,
			},
		},
	})
}
