resource "example_custom" "test" {
  # fields omitted:
  # - computed
  # - injected
  # - excluded

  required  = "required"
  sensitive = "sensitive"

  bool_custom_list = [true, false, true]

  # Stored in Go as a single string joined with "/".
  string_override = ["foo", "bar"]

  schema_override = "schema-override"
}
