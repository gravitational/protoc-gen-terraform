resource "example_custom" "test" {
  required  = "required"
  sensitive = "sensitive"
  validated = "valid"

  custom_name_override = "name-override"

  bool_custom      = true
  bool_custom_list = [true, false]

  string_override = ["foo", "bar"]

  schema_override = "schema-override"
}

data "example_custom" "test" {
  id       = example_custom.test.id
  required = example_custom.test.required
}
