resource "example_objects" "test" {
  primitives = {
    string_value   = "string"
    int32_value    = 123
    float_value    = 0.75
    bool_value     = true
    enum_value     = 1
    nullable_value = null
  }

  string_map = {
    key1 = "value1"
    key2 = "value2"
  }

  int_map = {
    one = 1
    two = 2
  }

  bool_map = {
    enabled  = true
    disabled = false
  }

  nested_value = {
    leaf = {
      value = "nested-value"
    }
  }

  nested_nullable = null

  nested_list = [
    { leaf = { value = "list-1" } },
    { leaf = { value = "list-2" } },
  ]

  nested_nullable_list = null

  nested_map = {
    first  = { leaf = { value = "map-1" } }
    second = { leaf = { value = "map-2" } }
  }

  nested_nullable_map = null

  # oneof: set only one branch.
  branch1 = {
    leaf = {
      value = "branch-1"
    }
  }

  leaf = {
    value = "embedded-leaf"
  }
  # TODO: Unepxected behavior with embedded fields.
  # This embedded value overwrites the embedded leaf.value field.
  # value = "embedded-nullable-value"
}
