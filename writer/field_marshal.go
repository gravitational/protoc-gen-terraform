package writer

import "github.com/gzigzigzeo/protoc-gen-terraform/plugin"

// FieldWriter represents logic required to generate field read go code
type FieldWriter struct {
	field *plugin.Field
}

// NewFieldWriter returns an instance of FieldWriter
func NewFieldWriter(f *plugin.Field) *FieldWriter {
	return &FieldWriter{field: f}
}

// Write generates code required for field
func (w *FieldWriter) Write() {
	// // 1. Read field value
	// _value_before_type_cast, ok := d.GetOk(path)
	// if (ok) {
	// 	// scalar
	// 	_value_after_type_cast = _value_before_type_cast.()

	// 	// array
	// 	var _value_after_type_cast make([]string, len(_valu_value_before_type_cast)
	// 	for _n, _item := range _value_before_type_cast {
	// 		_value_after_type_cast[_n] = _value_before_type_cast[n].() || parse
	// 	}
	// }
	// 1.1. List
	// 1.2. Array
	//
	// 2. Cast field value to schema type
	// 3. Cast schema type to struct type
	// 4. Assign

	// m.valueFmt = "time.Parse(time.RFC3339, %v)"
}

func (w *FieldWriter) writeReadRaw() {

}

func (w *FieldWriter) writeReadRawList() {

}

func (w *FieldWriter) writeCastRawToSchemaType() {

}

func (w *FieldWriter) writeCastRawListToSchemaType() {

}

func (w *FieldWriter) writeCastSchemaTypeToStructType() {

}

func (w *FieldWriter) writeCastSchemaTypeToStructTypeList() {

}
