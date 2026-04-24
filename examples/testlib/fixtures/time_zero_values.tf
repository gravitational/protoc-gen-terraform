resource "example_time" "test" {
  timestamp_value = "0001-01-01T00:00:00Z"
  timestamp_list  = ["0001-01-01T00:00:00Z", "0001-01-01T00:00:00Z"]

  duration_standard = "0s"
  duration_list     = ["0s", "0s"]

  duration_custom      = "0s"
  duration_custom_list = ["0s", "0s"]
}
