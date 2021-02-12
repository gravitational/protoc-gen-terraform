package writer

import "github.com/gzigzigzeo/protoc-gen-terraform/plugin"

// FieldSchema represents message writer instance
type FieldSchema struct {
	field *plugin.Field
}

// NewField returns message writer instance
func NewField(f *plugin.Field) *FieldSchema {
	return &FieldSchema{field: f}
}
