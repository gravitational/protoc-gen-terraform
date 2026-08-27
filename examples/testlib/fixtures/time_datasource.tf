resource "example_time" "test" {
  timestamp_value = "2026-01-02T03:04:05Z"
  timestamp_list  = ["2026-01-02T03:04:05Z", "2026-01-02T03:04:06Z"]

  duration_standard = "5m0s"
  duration_list     = ["5m0s", "10m0s"]

  duration_custom      = "5m0s"
  duration_custom_list = ["5m0s", "10m0s"]

  nullable_timestamp = "2026-01-02T03:04:05Z"
  nullable_duration  = "5m0s"
}

data "example_time" "test" {
  id = example_time.test.id

  nullable_timestamp = example_time.test.nullable_timestamp
  nullable_duration  = example_time.test.nullable_duration
}
