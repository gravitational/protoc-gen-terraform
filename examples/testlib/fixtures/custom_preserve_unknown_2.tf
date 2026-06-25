resource "example_custom" "reference" {
  required = "required"
}

resource "example_custom" "preserve_unknown" {
  required = example_custom.reference.computed
}
