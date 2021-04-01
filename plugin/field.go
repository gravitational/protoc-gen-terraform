/*
Copyright 2015-2021 Gravitational, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package plugin

import (
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/stoewer/go-strcase"
)

// Field represents metadata about protobuf message field descriptor.
// This struct and the following methods know about both schema details and target field details.
type Field struct {
	// Name type name
	Name string

	// NameSnake type name in snake case (Terraform schema field name)
	NameSnake string

	// SchemaRawType Terraform schema raw value go type (float64 for types.Float), used in schema generation
	SchemaRawType string

	// SchemaGoType Go type to convert schema raw type to (uint32, []byte, time.Time, time.Duration)
	SchemaGoType string

	// GoType Go type without [] and *, but with package name prepended
	GoType string

	// GoTypeIsSlice specifies whether Go type is a slice
	GoTypeIsSlice bool

	// GoTypeIsPtr specifies whether Go type is a pointer
	GoTypeIsPtr bool

	// GoTypeFull Go type with [] and * and package name prepended
	GoTypeFull string

	// Kind Field kind (resulting of combination of flags below, see setKind method)
	Kind string

	// IsRepeated field is list?
	IsRepeated bool

	// IsMap field is map?
	IsMap bool

	// IsMessage field is message? (might be map or list at the same time)
	IsMessage bool

	// IsRequired field is required
	IsRequired bool

	// IsComputed field is computed
	IsComputed bool

	// IsTime field contains time? (needs to be parsed from string)
	IsTime bool

	// IsDuration field contains duration? (needs to be parsed from string)
	IsDuration bool

	// IsCustomType field has gogo.customtype flag?
	IsCustomType bool

	// CustomTypeMethodInfix custom type schema and unmarshal method name infix
	CustomTypeMethodInfix string

	// Message reference to nested message
	Message *Message

	// MapValueField reference to map value field reflection
	MapValueField *Field

	// RawComment is field comment in proto file without // prepended
	RawComment string

	// Comment is field comment in proto file with // prepended
	Comment string

	// Default is field default value
	Default string
}

// BuildFields builds []*Field from descriptors of specified message
func BuildFields(m *Message, g *generator.Generator, d *generator.Descriptor) error {
	for i, fd := range d.GetField() {
		fd := &FieldDescriptorProtoExt{fd}

		c, err := NewFieldBuildContext(m, g, d, fd, i)
		if err != nil {
			return err
		}

		f, err := BuildField(c)
		if err != nil {
			return err
		}

		if f != nil {
			m.Fields = append(m.Fields, f)
		}
	}

	return nil
}

// BuildField builds Field structure
func BuildField(c *FieldBuildContext) (*Field, error) {
	var err error

	if c.IsExcluded() {
		return nil, nil
	}

	n := c.GetName()

	f := &Field{
		Name:      n,
		NameSnake: strcase.SnakeCase(n),
	}

	f.RawComment, f.Comment = c.GetComment()
	f.SchemaRawType, f.SchemaGoType, f.IsMessage, err = c.GetTypeAndIsMessage()
	if err != nil {
		return nil, err
	}

	f.IsTime = c.IsTime()
	f.IsDuration = c.IsDuration()

	// Byte slice is an exception, should not be treated as normal array
	if c.IsByteSlice() {
		f.GoType = c.GetRawGoType() // []byte
		f.GoTypeFull = c.GetRawGoType()
	} else {
		f.GoTypeIsPtr = c.GetGoTypeIsPtr()
		f.GoTypeIsSlice = c.GetGoTypeIsSlice()

		if f.IsMessage {
			d, err := c.GetMessageDescriptor()
			if err != nil {
				return nil, err
			}

			m, err := BuildMessage(c.g, d, false, c.path)
			if err != nil {
				return nil, err
			}
			if m == nil {
				return nil, nil
			}

			f.Message = m
			f.GoType = m.GoTypeName
		} else {
			f.GoType = c.GetGoType()
		}

		f.setGoTypeFull()
	}

	f.setKind()

	return f, nil
}

// setGoTypeFull constructs and sets full go type name for the field
func (f *Field) setGoTypeFull() {
	r := f.GoType

	if f.GoTypeIsPtr {
		r = "*" + r
	}
	if f.GoTypeIsSlice {
		r = "[]" + r
	}

	f.GoTypeFull = r
}

// setKind resolves and sets kind for current field
func (f *Field) setKind() {
	switch {
	case f.IsCustomType:
		f.Kind = "CUSTOM_TYPE"
	case f.IsMap && f.MapValueField.IsMessage:
		// Terraform does not support map of objects. We replace such field with list of objects having key and value fields.
		f.Kind = "OBJECT_MAP" // ex: map[string]struct, requires additional steps to unmarshal
	case f.IsMap:
		f.Kind = "MAP" // ex: map[string]string
	case f.IsRepeated && f.IsMessage:
		f.Kind = "REPEATED_MESSAGE" // ex: []struct
	case f.IsRepeated:
		f.Kind = "REPEATED_ELEMENTARY" // ex: []string
	case f.IsMessage:
		f.Kind = "SINGULAR_MESSAGE" // ex: struct
	default:
		f.Kind = "SINGULAR_ELEMENTARY" // ex: string
	}
}

// err := f.setMessage(ctx)
// if err != nil {
// 	return trace.Wrap(err)
// }

// f.setGoType(ctx)
// f.setAggregate(ctx)
// // // if err != nil {
// // // 	return nil, trace.Wrap(err)
// // // }

// f.setCustomType(ctx.f)

// f.setComment(ctx.d)
// f.setRequired()
// f.setComputed()
// f.setDefault()
// f.setKind()

// f, err := BuildField(g, d, f, path, i)
// if err != nil {
// 	invErr, ok := trace.Unwrap(err).(*invalidFieldError)

// 	// invalidFieldError is not considered fatal, we just need to report it to the user and skip it
// 	if !ok {
// 		return err
// 	}

// 	logrus.Warning(invErr.Error())
// 	continue
// }
//}

// setMessage sets reference to a nested message
// func (f *Field) setMessage(ctx *context) error {
// 	// Resolve underlying message via protobuf
// 	x := ctx.g.ObjectNamed(ctx.f.GetTypeName())
// 	desc, ok := x.(*generator.Descriptor)
// 	if desc == nil || !ok {
// 		return nil
// 	}

// 	// Try to analyse it
// 	m, err := BuildMessage(g, d, false, f.Path)
// 	if err != nil {
// 		// If underlying message is invalid, we must consider current field as invalid and not stop
// 		_, ok := trace.Unwrap(err).(*invalidMessageError)
// 		if ok {
// 			return trace.Wrap(
// 				newInvalidFieldError(b, fmt.Sprintf("failed to reflect message type information: %v", err.Error())),
// 			)
// 		}

// 		return trace.Wrap(err)
// 	} else if m == nil {
// 		return trace.Wrap(newInvalidFieldError(b, "field marked as skipped"))
// 	}

// 	// Nested message schema, or nil if message is not whitelisted
// 	f.Message = m

// 	return nil
// }

// // setAggregate detects and sets IsList and IsMap flags.
// func (f *Field) setAggregate(ctx *context) error {
// 	if ctx.g.IsMap(ctx.f.FieldDescriptorProto) {
// 		f.IsMap = true
// 		err := f.setMap(ctx)
// 		if err != nil {
// 			return trace.Wrap(err)
// 		}
// 	} else if ctx.f.IsRepeated() {
// 		f.IsRepeated = true
// 	}

// 	return nil
// }

// // setCustomType sets detects and set IsCustomType flag and custom type method infix
// func (f *Field) setCustomType(fd *FieldDescriptorProtoExt) {
// 	if !fd.IsCustomType() {
// 		return
// 	}

// 	f.IsCustomType = true
// 	f.CustomTypeMethodInfix = strings.ReplaceAll(strings.ReplaceAll(f.GoType, "/", ""), ".", "")
// }

// // setMap sets map value properties
// func (f *Field) setMap(ctx *context) error {
// 	m := ctx.g.GoMapType(nil, ctx.f.FieldDescriptorProto)

// 	keyGoType, _ := ctx.g.GoType(ctx.d, m.KeyField)
// 	if keyGoType != "string" {
// 		return trace.Wrap(newInvalidFieldError(b, "non-string map keys are not supported"))
// 	}

// 	c := &context{m: ctx.m, g: ctx.g, d: ctx.d, f: &FieldDescriptorProtoExt{m.ValueField}}

// 	valueField, err := BuildField(c, 0)
// 	if err != nil {
// 		return err
// 	}
// 	f.MapValueField = valueField

// 	return nil
// }

// // setSchemaTypes sets SchemaRawType and SchemaGoType
// func (f *Field) setSchemaTypes(schemaRawType string, goTypeCast string) {
// 	f.SchemaRawType = schemaRawType
// 	f.SchemaGoType = goTypeCast
// }

// // isInSet returns true if that field name or path is in set
// func isInSet(f *Field, m map[string]struct{}) bool {
// 	_, ok := m[f.TypeName]
// 	if ok {
// 		return true
// 	}

// 	_, ok = m[f.Path]
// 	return ok
// }

// // setComment resolves leading comment for this field
// func (f *Field) setComment(d *generator.Descriptor) {
// 	p := d.Path() + ",2," + strconv.Itoa(f.Index)

// 	for _, l := range d.File().GetSourceCodeInfo().GetLocation() {
// 		if getLocationPath(l) == p {
// 			c := strings.Trim(l.GetLeadingComments(), "\n")

// 			f.RawComment = commentToSingleLine(strings.TrimSpace(c))
// 			f.Comment = appendSlashSlash(c, false)
// 		}
// 	}
// }

// // setRequired sets IsRequired flag
// func (f *Field) setRequired() {
// 	if f.isInSet(config.RequiredFields) {
// 		f.IsRequired = true
// 	}
// }

// // setComputed sets IsComputed flag
// func (f *Field) setComputed() {
// 	if f.isInSet(config.ComputedFields) {
// 		f.IsComputed = true
// 	}
// }

// // setDefault sets default value
// func (f *Field) setDefault() {
// 	v, ok := config.Defaults[f.TypeName]
// 	if ok {
// 		f.Default = fmt.Sprintf("%s", v)
// 	}
// }

// // setKind sets field kind which represents field meta type for generation

// // getFieldTypeName returns field type name with package
// func getFieldTypeName(d *generator.Descriptor, f *FieldDescriptorProtoExt) string {
// 	return getMessageTypeName(d) + "." + f.GetName()
// }

// // // build fills in a Field structure
// // func (b *fieldBuilder) build() (*Field, error) {
// // 	// b.resolveName()

// // 	// if hasFieldInSet(config.ExcludeFields, b.field) {
// // 	// 	return nil, nil
// // 	// }

// // 	// err := b.resolveType()
// // 	// if err != nil {
// // 	// 	return nil, trace.Wrap(err)
// // 	// }

// // 	// err = b.setGoType()
// // 	// if err != nil {
// // 	// 	return nil, trace.Wrap(err)
// // 	// }

// // 	// // err = b.setAggregate()
// // 	// // if err != nil {
// // 	// // 	return nil, trace.Wrap(err)
// // 	// // }

// // 	// b.setCustomType()
// // 	//b.resolveKind()
// // 	//b.setComment()
// // 	// b.setRequired()
// // 	// b.setComputed()
// // 	// b.setDefault()

// // 	return b.field, nil
// // }

// // // setMessage sets reference to a nested message
// // func (b *fieldBuilder) setMessage() error {
// // 	// Resolve underlying message via protobuf
// // 	x := b.generator.ObjectNamed(b.fieldDescriptor.GetTypeName())
// // 	desc, ok := x.(*generator.Descriptor)
// // 	if desc == nil || !ok {
// // 		return nil
// // 	}

// // 	// Try to analyse it
// // 	m, err := BuildMessage(b.generator, desc, false, b.field.path)
// // 	if err != nil {
// // 		// If underlying message is invalid, we must consider current field as invalid and not stop
// // 		_, ok := trace.Unwrap(err).(*invalidMessageError)
// // 		if ok {
// // 			return trace.Wrap(
// // 				newInvalidFieldError(b, fmt.Sprintf("failed to reflect message type information: %v", err.Error())),
// // 			)
// // 		}

// // 		return trace.Wrap(err)
// // 	} else if m == nil {
// // 		return trace.Wrap(newInvalidFieldError(b, "field marked as skipped"))
// // 	}

// // 	// Nested message schema, or nil if message is not whitelisted
// // 	b.field.Message = m

// // 	return nil
// // }

// // // setAggregate detects and sets IsList and IsMap flags.
// // func (b *fieldBuilder) setAggregate() error {
// // 	f := b.field

// // 	if b.generator.IsMap(b.fieldDescriptor) {
// // 		f.IsMap = true
// // 		err := b.setMap()
// // 		if err != nil {
// // 			return trace.Wrap(err)
// // 		}
// // 	} else if b.fieldDescriptor.IsRepeated() {
// // 		f.IsRepeated = true
// // 	}

// // 	return nil
// // }

// // // setCustomType sets detects and set IsCustomType flag and custom type method infix
// // func (b *fieldBuilder) setCustomType() {
// // 	if !gogoproto.IsCustomType(b.fieldDescriptor) {
// // 		return
// // 	}

// // 	b.field.IsCustomType = true
// // 	b.field.CustomTypeMethodInfix = strings.ReplaceAll(strings.ReplaceAll(b.field.GoType, "/", ""), ".", "")
// // }

// // // setMap sets map value properties
// // func (b *fieldBuilder) setMap() error {
// // 	m := b.generator.GoMapType(nil, b.fieldDescriptor)

// // 	keyGoType, _ := b.generator.GoType(b.descriptor, m.KeyField)
// // 	if keyGoType != "string" {
// // 		return trace.Wrap(newInvalidFieldError(b, "non-string map keys are not supported"))
// // 	}

// // 	valueField, err := BuildField(b.generator, b.descriptor, m.ValueField, b.field.path, 0)
// // 	if err != nil {
// // 		return err
// // 	}
// // 	b.field.MapValueField = valueField

// // 	return nil
// // }
