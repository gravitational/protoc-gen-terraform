package plugin

import (
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gravitational/protoc-gen-terraform/config"
)

const (
	// protobufTimestampTypeName gogo protobuf timestamp type name
	protobufTimestampTypeName = "google.protobuf.Timestamp"

	// protobufDurationTypeName gogo protobuf duration type name
	protobufDurationTypeName = "google.protobuf.Duration"
)

// FieldDescriptorProtoExt adds some useful extension methods to gogo protobuf field descriptor
type FieldDescriptorProtoExt struct {
	*descriptor.FieldDescriptorProto
}

// IsTypeEq returns true if type equals current field descriptor type
func (f *FieldDescriptorProtoExt) IsTypeEq(t descriptor.FieldDescriptorProto_Type) bool {
	return *f.Type == t
}

// IsTime returns true if field stores a time value (protobuf or standard library)
func (f *FieldDescriptorProtoExt) IsTime() bool {
	t := f.TypeName

	isStdTime := gogoproto.IsStdTime(f.FieldDescriptorProto)
	isGoogleTime := (t != nil && strings.HasSuffix(*t, protobufTimestampTypeName))
	isCastToTime := f.GetCastType() == "time.Time"

	return isStdTime || isGoogleTime || isCastToTime
}

// IsDuration returns true if field stores a duration value (protobuf or cast to a standard library type)
func (f *FieldDescriptorProtoExt) IsDuration() bool {
	ct := f.GetCastType()
	t := f.TypeName

	isStdDuration := gogoproto.IsStdDuration(f.FieldDescriptorProto)
	isGoogleDuration := (t != nil && strings.HasSuffix(*t, protobufDurationTypeName))
	isCastToCustomDuration := ct == config.DurationCustomType
	isCastToDuration := ct == "time.Duration"

	return isStdDuration || isGoogleDuration || isCastToDuration || isCastToCustomDuration
}

// IsMessage returns true if field is a message
func (f *FieldDescriptorProtoExt) IsMessage() bool {
	return f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_MESSAGE)
}

// IsCastType returns true if field has gogoprotobuf.casttype flag
func (f *FieldDescriptorProtoExt) IsCastType() bool {
	return gogoproto.IsCastType(f.FieldDescriptorProto)
}

// IsCustomType returns true if field has gogoprotobuf.customtype flag
func (f *FieldDescriptorProtoExt) IsCustomType() bool {
	return gogoproto.IsCustomType(f.FieldDescriptorProto)
}

// GetCastType returns field cast type name
func (f *FieldDescriptorProtoExt) GetCastType() string {
	return gogoproto.GetCastType(f.FieldDescriptorProto)
}

// GetCustomType returns field custom type name
func (f *FieldDescriptorProtoExt) GetCustomType() string {
	return gogoproto.GetCustomType(f.FieldDescriptorProto)
}
