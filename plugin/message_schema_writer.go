package plugin

// MessageSchemaWriter represents message writer instance
type messageSchemaWriter struct {
	writer
	message *Message
}

// NewMessageSchemaWriter returns message writer instance
func newMessageSchemaWriter(m *Message) *messageSchemaWriter {
	return &messageSchemaWriter{message: m}
}

// Name of generated schema method
func (w *messageSchemaWriter) methodName() string {
	return `Schema` + w.message.Name + `()`
}

// Writer writes terraform.Schema for a message
func (w *messageSchemaWriter) write() string {
	w.p(`// `, w.message.Name)
	w.p(`func `, w.methodName(), ` map[string]*schema.Schema {`)
	w.p(`  return map[string]*schema.Schema {`)

	for _, field := range w.message.Fields {
		w.p(newFieldSchemaWriter(field).write())
	}

	w.p(`  }`)
	w.p(`}`)

	return w.buf.String()
}
