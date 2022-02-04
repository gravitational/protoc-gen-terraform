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

package desc

import (
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
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
func (f *FieldDescriptorProtoExt) IsDuration(durationCustomType string) bool {
	ct := f.GetCastType()
	t := f.TypeName

	isStdDuration := gogoproto.IsStdDuration(f.FieldDescriptorProto)
	isGoogleDuration := (t != nil && strings.HasSuffix(*t, protobufDurationTypeName))
	isCastToCustomDuration := durationCustomType != "" && ct == durationCustomType
	isCastToDuration := ct == "time.Duration"

	return isStdDuration || isGoogleDuration || isCastToDuration || isCastToCustomDuration
}

// IsMessage returns true if field is a message
func (f *FieldDescriptorProtoExt) IsMessage() bool {
	return f.IsTypeEq(descriptor.FieldDescriptorProto_TYPE_MESSAGE)
}

// IsCastType returns true if field has gogoproto.casttype flag
func (f *FieldDescriptorProtoExt) IsCastType() bool {
	return gogoproto.IsCastType(f.FieldDescriptorProto)
}

// IsCustomType returns true if field has gogoproto.customtype flag
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

// GetJSONName returns JSON name set in JSON tag
func (f *FieldDescriptorProtoExt) GetJSONName() string {
	t := gogoproto.GetJsonTag(f.FieldDescriptorProto)
	if t != nil {
		j := strings.Split(*t, ",")
		if j[0] != "-" {
			return j[0]
		}
	}

	return ""
}
