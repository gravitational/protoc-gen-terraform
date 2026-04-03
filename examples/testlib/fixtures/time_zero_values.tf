resource "example_time" "test_zero_values" {
  duration_standard = "0s"
  duration_list     = ["0s", "0s"]

  duration_custom      = "0s"
  duration_custom_list = ["0s", "0s"]
}
