// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: test.proto

package test

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
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

type Nested struct {
	String_              string        `protobuf:"bytes,1,opt,name=String,proto3" json:"String,omitempty"`
	Nested               *NestedLevel2 `protobuf:"bytes,2,opt,name=Nested,proto3" json:"Nested,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *Nested) Reset()         { *m = Nested{} }
func (m *Nested) String() string { return proto.CompactTextString(m) }
func (*Nested) ProtoMessage()    {}
func (*Nested) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{0}
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

type NestedLevel2 struct {
	String_              string   `protobuf:"bytes,1,opt,name=String,proto3" json:"String,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NestedLevel2) Reset()         { *m = NestedLevel2{} }
func (m *NestedLevel2) String() string { return proto.CompactTextString(m) }
func (*NestedLevel2) ProtoMessage()    {}
func (*NestedLevel2) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{1}
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

type Test struct {
	// string field
	Str string `protobuf:"bytes,1,opt,name=Str,proto3" json:"Str,omitempty"`
	// int32 field
	Int32 int32 `protobuf:"varint,2,opt,name=Int32,proto3" json:"Int32,omitempty"`
	// int64 field
	Int64 int64 `protobuf:"varint,3,opt,name=Int64,proto3" json:"Int64,omitempty"`
	// float field
	Float float32 `protobuf:"fixed32,4,opt,name=Float,proto3" json:"Float,omitempty"`
	// double field
	Double float64 `protobuf:"fixed64,5,opt,name=Double,proto3" json:"Double,omitempty"`
	// bool field
	Bool bool `protobuf:"varint,6,opt,name=Bool,proto3" json:"Bool,omitempty"`
	// byte[] field
	Bytes []byte `protobuf:"bytes,7,opt,name=Bytes,proto3" json:"Bytes,omitempty"`
	// time.Time field
	Timestamp time.Time `protobuf:"bytes,8,opt,name=Timestamp,proto3,stdtime" json:"Timestamp"`
	// time.Duration field (standard)
	DurationStd time.Duration `protobuf:"varint,9,opt,name=DurationStd,proto3,stdduration" json:"DurationStd,omitempty"`
	// time.Duration field (custom)
	DurationCustom Duration `protobuf:"varint,10,opt,name=DurationCustom,proto3,casttype=Duration" json:"DurationCustom,omitempty"`
	// *bool field
	BoolN bool `protobuf:"varint,11,opt,name=BoolN,proto3" json:"BoolN,omitempty"`
	// *byte[] field
	BytesN []byte `protobuf:"bytes,12,opt,name=BytesN,proto3" json:"BytesN,omitempty"`
	// *time.Time field
	TimestampN *time.Time `protobuf:"bytes,13,opt,name=TimestampN,proto3,stdtime" json:"TimestampN,omitempty"`
	// *time.Duration field
	DurationN Duration `protobuf:"varint,14,opt,name=DurationN,proto3,casttype=Duration" json:"DurationN,omitempty"`
	// []string field
	StringA []string `protobuf:"bytes,15,rep,name=StringA,proto3" json:"StringA,omitempty"`
	// []bool field
	BoolA []BoolCustomArray `protobuf:"varint,16,rep,packed,name=BoolA,proto3,customtype=BoolCustomArray" json:"BoolA,omitempty"`
	// [][]byte field
	BytesA [][]byte `protobuf:"bytes,17,rep,name=BytesA,proto3" json:"BytesA,omitempty"`
	// []time.Time field
	TimestampA []*time.Time `protobuf:"bytes,18,rep,name=TimestampA,proto3,stdtime" json:"TimestampA,omitempty"`
	// []time.Duration field
	GracePeriodA []Duration `protobuf:"varint,19,rep,packed,name=GracePeriodA,proto3,casttype=Duration" json:"GracePeriodA,omitempty"`
	// Nested field
	Nested *Nested `protobuf:"bytes,20,opt,name=Nested,proto3" json:"Nested,omitempty"`
	// Nested array
	NestedA              []*Nested `protobuf:"bytes,21,rep,name=NestedA,proto3" json:"NestedA,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Test) Reset()         { *m = Test{} }
func (m *Test) String() string { return proto.CompactTextString(m) }
func (*Test) ProtoMessage()    {}
func (*Test) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{2}
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

func init() {
	proto.RegisterType((*Nested)(nil), "test.Nested")
	proto.RegisterType((*NestedLevel2)(nil), "test.NestedLevel2")
	proto.RegisterType((*Test)(nil), "test.Test")
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_c161fcfdc0c3ff1e) }

var fileDescriptor_c161fcfdc0c3ff1e = []byte{
	// 519 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x93, 0x41, 0x6f, 0x12, 0x41,
	0x14, 0xc7, 0x99, 0xee, 0xb2, 0xc0, 0xb0, 0xb6, 0xf5, 0x15, 0xcd, 0x0b, 0x31, 0xec, 0x84, 0x18,
	0xb2, 0xf6, 0x00, 0x86, 0x36, 0x9e, 0x9d, 0x95, 0x68, 0x4c, 0x1a, 0x62, 0x86, 0x7e, 0x01, 0x90,
	0x91, 0x90, 0x00, 0xd3, 0xec, 0x0e, 0x26, 0xfd, 0x08, 0xde, 0x3c, 0xfa, 0x91, 0x38, 0x7a, 0xf6,
	0x80, 0xb1, 0x1f, 0xc3, 0x93, 0x99, 0x99, 0x5d, 0x4a, 0xab, 0x46, 0x6f, 0xef, 0xff, 0x7f, 0xff,
	0xdd, 0x79, 0xbf, 0x37, 0x19, 0x4a, 0xb5, 0xcc, 0x74, 0xf7, 0x2a, 0x55, 0x5a, 0x81, 0x6f, 0xea,
	0x66, 0x63, 0xa6, 0x66, 0xca, 0x1a, 0x3d, 0x53, 0xb9, 0x5e, 0x33, 0x9a, 0x29, 0x35, 0x5b, 0xc8,
	0x9e, 0x55, 0x93, 0xf5, 0x87, 0x9e, 0x9e, 0x2f, 0x65, 0xa6, 0xc7, 0xcb, 0x2b, 0x17, 0x68, 0x5f,
	0xd0, 0x60, 0x28, 0x33, 0x2d, 0xa7, 0xf0, 0x98, 0x06, 0x23, 0x9d, 0xce, 0x57, 0x33, 0x24, 0x8c,
	0xc4, 0x35, 0x91, 0x2b, 0x38, 0x2d, 0x12, 0x78, 0xc0, 0x48, 0x5c, 0xef, 0x43, 0xd7, 0x9e, 0xed,
	0xbc, 0x0b, 0xf9, 0x51, 0x2e, 0xfa, 0x22, 0x4f, 0xb4, 0x3b, 0x34, 0xdc, 0xf7, 0xff, 0xf6, 0xcf,
	0xf6, 0xa7, 0x80, 0xfa, 0x97, 0x32, 0xd3, 0x70, 0x4c, 0xbd, 0x91, 0x4e, 0xf3, 0xae, 0x29, 0xa1,
	0x41, 0xcb, 0x6f, 0x57, 0xfa, 0xac, 0x6f, 0x4f, 0x2b, 0x0b, 0x27, 0x72, 0xf7, 0xc5, 0x39, 0x7a,
	0x8c, 0xc4, 0x9e, 0x70, 0xc2, 0xb8, 0xaf, 0x17, 0x6a, 0xac, 0xd1, 0x67, 0x24, 0x3e, 0x10, 0x4e,
	0x98, 0x43, 0x07, 0x6a, 0x3d, 0x59, 0x48, 0x2c, 0x33, 0x12, 0x13, 0x91, 0x2b, 0x00, 0xea, 0x27,
	0x4a, 0x2d, 0x30, 0x60, 0x24, 0xae, 0x0a, 0x5b, 0x9b, 0x3f, 0x24, 0xd7, 0x5a, 0x66, 0x58, 0x61,
	0x24, 0x0e, 0x85, 0x13, 0x90, 0xd0, 0xda, 0x65, 0xb1, 0x27, 0xac, 0x5a, 0xea, 0x66, 0xd7, 0x6d,
	0xb2, 0x5b, 0x6c, 0xb2, 0xbb, 0x4b, 0x24, 0xd5, 0xcd, 0x36, 0x2a, 0x7d, 0xfe, 0x1e, 0x11, 0x71,
	0xfb, 0x19, 0x74, 0x68, 0x7d, 0xb0, 0x4e, 0xc7, 0x7a, 0xae, 0x56, 0x23, 0x3d, 0xc5, 0x9a, 0x99,
	0x3b, 0xf1, 0xbf, 0x98, 0xd4, 0x7e, 0x03, 0xce, 0xe9, 0x61, 0x21, 0x5f, 0xad, 0x33, 0xad, 0x96,
	0x48, 0x6d, 0x34, 0xfc, 0xb9, 0x8d, 0xaa, 0x45, 0x47, 0xdc, 0xcb, 0x40, 0x93, 0x96, 0xcd, 0xfc,
	0x43, 0xac, 0x1b, 0x98, 0xc4, 0xdf, 0x6c, 0x23, 0x22, 0x9c, 0x05, 0x4f, 0x68, 0x60, 0x31, 0x86,
	0x18, 0x1a, 0xa8, 0xbc, 0x99, 0x7b, 0x30, 0xa0, 0x74, 0x37, 0xe4, 0x10, 0x1f, 0xfc, 0x17, 0x1c,
	0xb1, 0x70, 0x7b, 0xdf, 0xc1, 0x29, 0xad, 0x15, 0x13, 0x0d, 0xf1, 0xf0, 0x0f, 0x03, 0xdf, 0xb6,
	0x01, 0x69, 0xc5, 0x5d, 0x3b, 0xc7, 0x23, 0xe6, 0xc5, 0x35, 0x51, 0x48, 0x78, 0xe6, 0x28, 0x38,
	0x1e, 0x33, 0x2f, 0xae, 0x26, 0x27, 0xdf, 0xb6, 0xd1, 0x91, 0x31, 0x1c, 0x24, 0x4f, 0xd3, 0xf1,
	0xb5, 0x83, 0xe2, 0xe6, 0x52, 0x2d, 0x00, 0xc7, 0x87, 0xcc, 0x8b, 0xc3, 0x1c, 0x87, 0xc3, 0xcb,
	0x3d, 0x1c, 0x8e, 0xc0, 0xbc, 0x7f, 0xe0, 0xf8, 0xf7, 0x50, 0x38, 0x3c, 0xa7, 0xe1, 0x9b, 0x74,
	0xfc, 0x5e, 0xbe, 0x93, 0xe9, 0x5c, 0x4d, 0x39, 0x9e, 0x30, 0xef, 0x37, 0x9a, 0x3b, 0x09, 0x78,
	0xba, 0x7b, 0x11, 0x0d, 0xbb, 0xbe, 0x70, 0xff, 0x45, 0x14, 0x6f, 0x01, 0x3a, 0xb4, 0xe2, 0x2a,
	0x8e, 0x8f, 0xec, 0x58, 0x77, 0x63, 0x45, 0x33, 0x09, 0x37, 0x3f, 0x5a, 0xa5, 0xcd, 0x4d, 0xab,
	0xf4, 0xf5, 0xa6, 0x55, 0x9a, 0x04, 0x76, 0xe6, 0xb3, 0x5f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x0c,
	0xf4, 0xfd, 0x85, 0xe1, 0x03, 0x00, 0x00,
}
