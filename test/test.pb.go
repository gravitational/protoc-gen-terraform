// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: test.proto

package test

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb"
	math "math"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Test message definition.
type Test struct {
	// Str string field
	Str string `protobuf:"bytes,1,opt,name=Str,proto3" json:"Str,omitempty"`
	// Int32 int32 field
	Int32 int32 `protobuf:"varint,2,opt,name=Int32,proto3" json:"Int32,omitempty"`
	// Int64 int64 field
	Int64 int64 `protobuf:"varint,3,opt,name=Int64,proto3" json:"Int64,omitempty"`
	// Float float field
	Float float32 `protobuf:"fixed32,4,opt,name=Float,proto3" json:"Float,omitempty"`
	// Double double field
	Double float64 `protobuf:"fixed64,5,opt,name=Double,proto3" json:"Double,omitempty"`
	// Bool bool field
	Bool bool `protobuf:"varint,6,opt,name=Bool,proto3" json:"Bool,omitempty"`
	// Bytest byte[] field
	Bytes []byte `protobuf:"bytes,7,opt,name=Bytes,proto3" json:"Bytes,omitempty"`
	// Timestamp time.Time field
	Timestamp time.Time `protobuf:"bytes,8,opt,name=Timestamp,proto3,stdtime" json:"Timestamp"`
	// TimestampN *time.Time field
	TimestampNullable *time.Time `protobuf:"bytes,9,opt,name=TimestampNullable,proto3,stdtime" json:"TimestampNullable,omitempty"`
	// TimestampN *time.Time field
	TimestampNullableWithNilValue *time.Time `protobuf:"bytes,10,opt,name=TimestampNullableWithNilValue,proto3,stdtime" json:"TimestampNullableWithNilValue,omitempty"`
	// DurationStd time.Duration field (standard)
	DurationStd time.Duration `protobuf:"varint,110,opt,name=DurationStd,proto3,stdduration" json:"DurationStd,omitempty"`
	// DurationCustom time.Duration field (custom)
	DurationCustom Duration `protobuf:"varint,120,opt,name=DurationCustom,proto3,casttype=Duration" json:"DurationCustom,omitempty"`
	// StringA []string field
	StringA []string `protobuf:"bytes,12,rep,name=StringA,proto3" json:"StringA,omitempty"`
	// BoolA []bool field
	BoolA []BoolCustom `protobuf:"varint,13,rep,packed,name=BoolA,proto3,customtype=BoolCustom" json:"BoolA,omitempty"`
	// BytesA [][]byte field
	BytesA [][]byte `protobuf:"bytes,14,rep,name=BytesA,proto3" json:"BytesA,omitempty"`
	// TimestampA []time.Time field
	TimestampA []*time.Time `protobuf:"bytes,15,rep,name=TimestampA,proto3,stdtime" json:"TimestampA,omitempty"`
	// DurationCustomA []time.Duration field
	DurationCustomA []Duration `protobuf:"varint,16,rep,packed,name=DurationCustomA,proto3,casttype=Duration" json:"DurationCustomA,omitempty"`
	// Nested nested message field
	Nested *Nested `protobuf:"bytes,17,opt,name=Nested,proto3" json:"Nested,omitempty"`
	// NestedA nested message array
	NestedA []*Nested `protobuf:"bytes,18,rep,name=NestedA,proto3" json:"NestedA,omitempty"`
	// NestedM normal map
	NestedM map[string]string `protobuf:"bytes,19,rep,name=NestedM,proto3" json:"NestedM,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// NestedMObj object map
	NestedMObj           map[string]*Nested `protobuf:"bytes,20,rep,name=NestedMObj,proto3" json:"NestedMObj,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *Test) Reset()         { *m = Test{} }
func (m *Test) String() string { return proto.CompactTextString(m) }
func (*Test) ProtoMessage()    {}
func (*Test) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{0}
}
func (m *Test) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Test.Unmarshal(m, b)
}
func (m *Test) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Test.Marshal(b, m, deterministic)
}
func (m *Test) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Test.Merge(m, src)
}
func (m *Test) XXX_Size() int {
	return xxx_messageInfo_Test.Size(m)
}
func (m *Test) XXX_DiscardUnknown() {
	xxx_messageInfo_Test.DiscardUnknown(m)
}

var xxx_messageInfo_Test proto.InternalMessageInfo

// Nested message definition
type Nested struct {
	// Str string field
	Str string `protobuf:"bytes,1,opt,name=Str,proto3" json:"Str,omitempty"`
	// Nested repeated nested messages
	Nested []*NestedLevel2 `protobuf:"bytes,2,rep,name=Nested,proto3" json:"Nested,omitempty"`
	// Nested map repeated nested messages
	NestedM map[string]string `protobuf:"bytes,3,rep,name=NestedM,proto3" json:"NestedM,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// NestedMObj nested object map
	NestedMObj           map[string]*NestedLevel2 `protobuf:"bytes,4,rep,name=NestedMObj,proto3" json:"NestedMObj,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *Nested) Reset()         { *m = Nested{} }
func (m *Nested) String() string { return proto.CompactTextString(m) }
func (*Nested) ProtoMessage()    {}
func (*Nested) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{1}
}
func (m *Nested) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Nested.Unmarshal(m, b)
}
func (m *Nested) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Nested.Marshal(b, m, deterministic)
}
func (m *Nested) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Nested.Merge(m, src)
}
func (m *Nested) XXX_Size() int {
	return xxx_messageInfo_Nested.Size(m)
}
func (m *Nested) XXX_DiscardUnknown() {
	xxx_messageInfo_Nested.DiscardUnknown(m)
}

var xxx_messageInfo_Nested proto.InternalMessageInfo

// Message nested into nested message
type NestedLevel2 struct {
	// Str string field
	Str                  string   `protobuf:"bytes,1,opt,name=Str,proto3" json:"Str,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NestedLevel2) Reset()         { *m = NestedLevel2{} }
func (m *NestedLevel2) String() string { return proto.CompactTextString(m) }
func (*NestedLevel2) ProtoMessage()    {}
func (*NestedLevel2) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{2}
}
func (m *NestedLevel2) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NestedLevel2.Unmarshal(m, b)
}
func (m *NestedLevel2) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NestedLevel2.Marshal(b, m, deterministic)
}
func (m *NestedLevel2) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NestedLevel2.Merge(m, src)
}
func (m *NestedLevel2) XXX_Size() int {
	return xxx_messageInfo_NestedLevel2.Size(m)
}
func (m *NestedLevel2) XXX_DiscardUnknown() {
	xxx_messageInfo_NestedLevel2.DiscardUnknown(m)
}

var xxx_messageInfo_NestedLevel2 proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Test)(nil), "test.Test")
	proto.RegisterMapType((map[string]string)(nil), "test.Test.NestedMEntry")
	proto.RegisterMapType((map[string]*Nested)(nil), "test.Test.NestedMObjEntry")
	proto.RegisterType((*Nested)(nil), "test.Nested")
	proto.RegisterMapType((map[string]string)(nil), "test.Nested.NestedMEntry")
	proto.RegisterMapType((map[string]*NestedLevel2)(nil), "test.Nested.NestedMObjEntry")
	proto.RegisterType((*NestedLevel2)(nil), "test.NestedLevel2")
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_c161fcfdc0c3ff1e) }

var fileDescriptor_c161fcfdc0c3ff1e = []byte{
	// 629 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x54, 0xcd, 0x6e, 0x13, 0x3d,
	0x14, 0x8d, 0x33, 0x93, 0x34, 0xb9, 0x9d, 0xaf, 0x3f, 0xfe, 0x2a, 0x30, 0x11, 0x24, 0x56, 0x55,
	0x55, 0x16, 0x8b, 0x54, 0xb4, 0x55, 0x85, 0x2a, 0x16, 0x64, 0x28, 0x48, 0x08, 0xda, 0x0a, 0xb7,
	0x82, 0x75, 0xa2, 0x9a, 0x90, 0x76, 0x1a, 0x47, 0x33, 0x9e, 0x42, 0xde, 0x82, 0x15, 0xe2, 0x91,
	0xb2, 0x64, 0xcd, 0x22, 0x88, 0x3e, 0x06, 0x2b, 0x64, 0x7b, 0x26, 0x99, 0xfc, 0x00, 0x45, 0x62,
	0x77, 0xcf, 0xbd, 0xe7, 0x5c, 0xf9, 0x1c, 0x8f, 0x07, 0x40, 0x89, 0x48, 0xd5, 0x7b, 0xa1, 0x54,
	0x12, 0xbb, 0xba, 0xae, 0xac, 0xb5, 0x65, 0x5b, 0x9a, 0xc6, 0x96, 0xae, 0xec, 0xac, 0x52, 0x6b,
	0x4b, 0xd9, 0x0e, 0xc4, 0x96, 0x41, 0xad, 0xf8, 0xed, 0x96, 0xea, 0x5c, 0x8a, 0x48, 0x35, 0x2f,
	0x7b, 0x09, 0xa1, 0x3a, 0x4d, 0x78, 0x1f, 0x36, 0x7b, 0x3d, 0x11, 0x46, 0x76, 0xbe, 0xfe, 0xa9,
	0x04, 0xee, 0xa9, 0x88, 0x14, 0x5e, 0x01, 0xe7, 0x44, 0x85, 0x04, 0x51, 0xc4, 0xca, 0x5c, 0x97,
	0x78, 0x0d, 0x0a, 0xcf, 0xbb, 0x6a, 0x67, 0x9b, 0xe4, 0x29, 0x62, 0x05, 0x6e, 0x41, 0xd2, 0xdd,
	0xdb, 0x25, 0x0e, 0x45, 0xcc, 0xe1, 0x16, 0xe8, 0xee, 0xb3, 0x40, 0x36, 0x15, 0x71, 0x29, 0x62,
	0x79, 0x6e, 0x01, 0xbe, 0x05, 0xc5, 0x03, 0x19, 0xb7, 0x02, 0x41, 0x0a, 0x14, 0x31, 0xc4, 0x13,
	0x84, 0x31, 0xb8, 0xbe, 0x94, 0x01, 0x29, 0x52, 0xc4, 0x4a, 0xdc, 0xd4, 0x7a, 0x83, 0xdf, 0x57,
	0x22, 0x22, 0x0b, 0x14, 0x31, 0x8f, 0x5b, 0x80, 0x7d, 0x28, 0x9f, 0xa6, 0x8e, 0x48, 0x89, 0x22,
	0xb6, 0xb8, 0x5d, 0xa9, 0x5b, 0x4b, 0xf5, 0xd4, 0x52, 0x7d, 0xc4, 0xf0, 0x4b, 0x83, 0x61, 0x2d,
	0xf7, 0xf1, 0x5b, 0x0d, 0xf1, 0xb1, 0x0c, 0x73, 0x58, 0x1d, 0x81, 0xa3, 0x38, 0x08, 0x9a, 0xfa,
	0x40, 0xe5, 0x1b, 0xed, 0x42, 0x66, 0xd7, 0xac, 0x1c, 0x9f, 0xc3, 0xbd, 0x99, 0xe6, 0x9b, 0x8e,
	0x7a, 0x77, 0xd4, 0x09, 0x5e, 0x37, 0x83, 0x58, 0x10, 0xf8, 0x8b, 0xfd, 0xbf, 0x5f, 0x85, 0x37,
	0x61, 0xf1, 0x20, 0x0e, 0x9b, 0xaa, 0x23, 0xbb, 0x27, 0xea, 0x8c, 0x74, 0x75, 0xee, 0xbe, 0xfb,
	0x59, 0x2b, 0xb3, 0x03, 0xbc, 0x0b, 0x4b, 0x29, 0x7c, 0x12, 0x47, 0x4a, 0x5e, 0x92, 0x0f, 0x86,
	0xea, 0xfd, 0x18, 0xd6, 0x4a, 0xe9, 0x84, 0x4f, 0x71, 0x30, 0x81, 0x85, 0x13, 0x15, 0x76, 0xba,
	0xed, 0x06, 0xf1, 0xa8, 0xc3, 0xca, 0x3c, 0x85, 0x78, 0x03, 0x0a, 0xfa, 0x66, 0x1a, 0xe4, 0x3f,
	0xea, 0xb0, 0x92, 0xbf, 0xf4, 0x75, 0x58, 0x03, 0xdd, 0xb0, 0x42, 0x6e, 0x87, 0xfa, 0x8e, 0xcd,
	0x55, 0x35, 0xc8, 0x12, 0x75, 0x98, 0xc7, 0x13, 0x84, 0x1f, 0x03, 0x8c, 0x6c, 0x35, 0xc8, 0x32,
	0x75, 0xfe, 0x10, 0x87, 0x6b, 0xa2, 0xc8, 0x68, 0xf0, 0x1e, 0x2c, 0x4f, 0x9e, 0xb5, 0x41, 0x56,
	0xa8, 0x33, 0x63, 0x68, 0x9a, 0x84, 0x37, 0xa0, 0x78, 0x24, 0x22, 0x25, 0xce, 0xc8, 0xaa, 0xb9,
	0x04, 0xaf, 0x6e, 0x1e, 0x93, 0xed, 0xf1, 0x64, 0x86, 0x37, 0x61, 0xc1, 0x56, 0x0d, 0x82, 0xcd,
	0xe1, 0x26, 0x69, 0xe9, 0x10, 0x3f, 0x48, 0x79, 0x87, 0xe4, 0x7f, 0xc3, 0xbb, 0x6d, 0x79, 0xa7,
	0x63, 0xf2, 0xe1, 0xd3, 0xae, 0x0a, 0xfb, 0xa9, 0xe4, 0x10, 0xef, 0x03, 0x24, 0xe5, 0x71, 0xeb,
	0x9c, 0xac, 0x25, 0xd6, 0x67, 0x54, 0xc7, 0xad, 0x73, 0x2b, 0xcc, 0xb0, 0x2b, 0xfb, 0xe0, 0x65,
	0x97, 0xea, 0x67, 0x79, 0x21, 0xfa, 0xe9, 0xb3, 0xbc, 0x10, 0x7d, 0xfd, 0x50, 0xae, 0xcc, 0x27,
	0x96, 0x37, 0x3d, 0x0b, 0xf6, 0xf3, 0x0f, 0x51, 0xe5, 0x05, 0x2c, 0x4f, 0xad, 0x9e, 0x23, 0x5f,
	0xcf, 0xca, 0xa7, 0x5d, 0x8f, 0x97, 0xad, 0x0f, 0xf2, 0x69, 0x8c, 0x73, 0x7e, 0x0d, 0xf7, 0x47,
	0x11, 0xe7, 0x8d, 0x3b, 0x9c, 0xdd, 0xf2, 0x52, 0x5c, 0x89, 0x60, 0x7b, 0x14, 0xf4, 0xce, 0x38,
	0x40, 0xc7, 0x90, 0xef, 0x64, 0xc9, 0xbf, 0x88, 0xf0, 0xd1, 0x44, 0x84, 0xae, 0xd1, 0xdd, 0x9d,
	0xa7, 0xfb, 0xe7, 0x21, 0xbe, 0xba, 0x49, 0x88, 0x6c, 0x32, 0xc4, 0x79, 0xf6, 0x33, 0x51, 0xd2,
	0xf4, 0x38, 0x76, 0x34, 0x9b, 0xa7, 0xef, 0x0d, 0xbe, 0x57, 0x73, 0x83, 0xeb, 0x6a, 0xee, 0xcb,
	0x75, 0x35, 0xd7, 0x2a, 0x9a, 0xe7, 0xb1, 0xf3, 0x33, 0x00, 0x00, 0xff, 0xff, 0x55, 0x1f, 0x1f,
	0x00, 0x05, 0x06, 0x00, 0x00,
}
