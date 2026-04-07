resource "example_primitives" "test" {
  string_value = "string"
  int32_value  = 123
  int64_value  = 456
  // TODO: Float64 type validation error in terraform-plugin-framework v0.10.0
  // 0.75 works, but 0.76 fails.
  // Verify both cases are correctly validated after updating to > v1.2.0
  // See: https://github.com/hashicorp/terraform-plugin-framework/issues/647
  float_value  = 0.75
  double_value = 0.75
  bool_value   = true
  bytes_value  = "bytes"
  enum_value   = 1
  string_list  = ["el1", "el2"]
  int32_list   = [123, 456]
  int64_list   = [234, 567]
  float_list   = [0.75, 1.25]
  double_list  = [0.75, 1.25]
  bool_list    = [true, false]
  bytes_list   = ["bytes1", "bytes2"]
  enum_list    = [1, 2]
}
