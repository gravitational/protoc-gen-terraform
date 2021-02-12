package plugin

// FieldSchemaWriter represents message writer instance
type fieldSchemaWriter struct {
	writer
	field *Field
}

// NewField returns message writer instance
func newFieldSchemaWriter(f *Field) *fieldSchemaWriter {
	return &fieldSchemaWriter{field: f}
}

// write returns schema definition
func (w *fieldSchemaWriter) write() string {
	w.p(`"`, w.field.NameSnake, `": {`)

	if w.field.IsAggregate() || w.field.HasNestedMessage() {
		w.writeAggregateType()
		w.writeRequired()
		w.writeMaxItems()
		if w.field.HasNestedMessage() {
			w.writeNestedMessageElem()
		} else {
			w.writeSimpleElem()
		}
	} else {
		w.writeType()
		w.writeRequired()
	}

	w.p(`},`)
	return w.buf.String()
}

func (w *fieldSchemaWriter) writeAggregateType() {
	w.p(`Type: schema.`, w.field.TFSchemaAggregateType, `,`)
}

func (w *fieldSchemaWriter) writeMaxItems() {
	max := w.field.TFSchemaMaxItems
	if max > 0 {
		w.p(`MaxItems: `, max, `,`)
	}
}

func (w *fieldSchemaWriter) writeType() {
	w.p(`Type: schema.`, w.field.TFSchemaType, `,`)
}

func (w *fieldSchemaWriter) writeRequired() {
	w.p(`Optional: true,`)
}

func (w *fieldSchemaWriter) writeNestedMessageElem() {
	w.p(`Elem: &schema.Resource {`)
	w.p(`  Schema: `, newMessageSchemaWriter(w.field.Message).methodName(), `,`)
	w.p(`},`)

}

func (w *fieldSchemaWriter) writeSimpleElem() {
	w.p(`Elem: &schema.Schema {`)
	w.writeType()
	w.p(`},`)
}
