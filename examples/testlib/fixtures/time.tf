resource "example_time" "test" {
  timestamp_value = "2026-01-02T03:04:05Z"
  timestamp_list  = ["2026-01-02T03:04:05Z", "2026-01-02T03:04:06Z"]

  duration_standard = "5m0s"
  duration_list     = ["5m0s", "10m0s"]

  # Casttype duration fields are int64 in this generated schema (nanoseconds).
  duration_custom      = "5m0s"
  duration_custom_list = ["5m0s", "10m0s"]
}
