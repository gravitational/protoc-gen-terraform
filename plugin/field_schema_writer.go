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
		w.pAggregateType()
		w.pRequired()
		w.pMaxItems()
		if w.field.HasNestedMessage() {
			w.pNestedMessageElem()
		} else {
			w.pSimpleElem()
		}
	} else {
		w.pType()
		w.pRequired()
	}

	w.p(`},`)
	return w.buf.String()
}

func (w *fieldSchemaWriter) pAggregateType() {
	w.p(`Type: schema.`, w.field.TFSchemaAggregateType, `,`)
}

func (w *fieldSchemaWriter) pMaxItems() {
	max := w.field.TFSchemaMaxItems
	if max > 0 {
		w.p(`MaxItems: `, max, `,`)
	}
}

func (w *fieldSchemaWriter) pType() {
	w.p(`Type: schema.`, w.field.TFSchemaType, `,`)
}

func (w *fieldSchemaWriter) pRequired() {
	w.p(`Optional: true,`)
}

func (w *fieldSchemaWriter) pNestedMessageElem() {
	w.p(`Elem: &schema.Resource {`)
	w.p(`  Schema: `, newMessageSchemaWriter(w.field.Message).methodName(), `,`)
	w.p(`},`)

}

func (w *fieldSchemaWriter) pSimpleElem() {
	w.p(`Elem: &schema.Schema {`)
	w.pType()
	w.p(`},`)
}
