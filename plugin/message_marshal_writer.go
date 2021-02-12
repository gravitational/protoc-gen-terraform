package plugin

// messageMarshalWriter represents message writer instance
type messageMarshalWriter struct {
	writer
	message *Message
}

// newMessageSchemaWriter returns message writer instance
func newMessageMarshalWriter(m *Message) *messageMarshalWriter {
	return &messageMarshalWriter{message: m}
}

// methodName returns name of generated schema method
func (w *messageMarshalWriter) methodName() string {
	// TODO GoTypeName needs improvement
	return `Unmarshal` + w.message.Name + `(d *schema.ResourceData, target *types.` + w.message.Name + `, prefix string) error`
}

// write writes terraform.Schema for a message
func (w *messageMarshalWriter) write() string {
	w.p(`// `, w.message.Name)
	w.p(`func `, w.methodName(), ` {`)

	for _, field := range w.message.Fields {
		w.p(newFieldMarshalWriter(field).write())
	}

	w.p(`return nil`)
	w.p(`}`)

	return w.buf.String()
}
