resource "example_objects" "test" {
  primitives = {
    string_value = ""
    int32_value  = 0
    float_value  = 0
    bool_value   = false
    enum_value   = 0
  }

  string_map = {
    key1 = ""
    key2 = ""
  }

  int_map = {
    key1 = 0
    key2 = 0
  }

  bool_map = {
    key1 = false
    key2 = false
  }

  nested_value = {
    leaf = {
      value = ""
    }
  }

  nested_nullable = null

  nested_list = [
    { leaf = { value = "" } },
    { leaf = { value = "" } },
  ]

  nested_nullable_list = null

  nested_map = {
    key1 = { leaf = { value = "" } }
    key2 = { leaf = { value = "" } }
  }

  nested_nullable_map = null

  branch1 = {
    leaf = {
      value = ""
    }
  }

  leaf = {
    value = ""
  }
}
