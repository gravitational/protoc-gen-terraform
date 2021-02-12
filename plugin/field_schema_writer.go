package plugin

// FieldSchemaWriter represents message writer instance
type FieldSchemaWriter struct {
	field *Field
}

// NewField returns message writer instance
func NewField(f *Field) *FieldSchemaWriter {
	return &FieldSchemaWriter{field: f}
}
