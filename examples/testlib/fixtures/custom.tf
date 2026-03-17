resource "example_custom" "test" {
  # fields omitted:
  # - computed is a computed field.
  # - injected is an injected field.
  # - excluded is an exluced field.
  # - plan_modifier sets a default value using a plan modifer.

  required  = "required"
  sensitive = "sensitive"
  validated = "valid"

  custom_name_override = "name-override"

  bool_custom      = true
  bool_custom_list = [true, false]

  # Stored in Go as a single string joined with "/".
  string_override = ["foo", "bar"]

  schema_override = "schema-override"
}
