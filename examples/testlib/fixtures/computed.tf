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

  computed_list = {
    string_list = ["s1", "s2"]
    int64_list  = [123, 456]
    float_list  = [0.5, 1.5]
    bool_list   = [true, false]
  }

  computed_map = {
    string_map = ["s1", "s2"]
    int64_map  = [123, 456]
    float_map  = [0.5, 1.5]
    bool_map   = [true, false]
  }
}
