resource "example_computed" "test" {
  string_value = "string"
  int32_value  = 123
  int64_value  = 456
  float_value  = 0.75
  double_value = 1.25
  bool_value   = true
  bytes_value  = "bytes"
  enum_value   = 1

  timestamp_value = "2026-01-02T03:04:05Z"
  duration_value  = "5m0s"

  primitives_value = {
    string_value = "string"
    int32_value  = 123
    int64_value  = 456
    float_value  = 0.75
    double_value = 0.75
    bool_value   = true
    bytes_value  = "bytes"
    enum_value   = 1
  }

  nested_value = {
    leaf = {
      value = "nested-value"
    }
  }

  nested_nullable = {
    leaf = {
      value = "nested-nullable"
    }
  }
}
