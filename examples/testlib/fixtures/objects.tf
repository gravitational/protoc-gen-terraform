resource "example_objects" "test" {
  primitives = {
    string_value = "string"
    int32_value  = 123
    float_value  = 0.75
    bool_value   = true
    enum_value   = 1
  }

  string_map = {
    key1 = "value1"
    key2 = "value2"
  }

  int_map = {
    key1 = 1
    key2 = 2
  }

  bool_map = {
    key1 = true
    key2 = false
  }

  nested_value = {
    leaf = {
      value = "nested-value"
    }
  }

  nested_nullable = {
    leaf = {
      value = "nested-value"
    }
  }

  nested_list = [
    { leaf = { value = "list-1" } },
    { leaf = { value = "list-2" } },
  ]

  nested_nullable_list = [
    { leaf = { value = "list-1" } },
    { leaf = { value = "list-2" } },
  ]

  nested_map = {
    key1 = { leaf = { value = "map-1" } }
    key2 = { leaf = { value = "map-2" } }
  }

  nested_nullable_map = {
    key1 = { leaf = { value = "map-1" } }
    key2 = { leaf = { value = "map-2" } }
  }

  # oneof: set only one branch.
  branch1 = {
    leaf = {
      value = "branch-1"
    }
  }

  leaf = {
    value = "embedded-leaf"
  }

  embedded_value = "embedded-nullable-value"
}
