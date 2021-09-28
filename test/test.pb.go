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
	// Bytes byte[] field
	Bytes []byte `protobuf:"bytes,7,opt,name=Bytes,proto3" json:"Bytes,omitempty"`
	// Timestamp time.Time field
	Timestamp time.Time `protobuf:"bytes,8,opt,name=Timestamp,proto3,stdtime" json:"Timestamp"`
	// Timestamp time.Time field
	TimestampMissing time.Time `protobuf:"bytes,9,opt,name=TimestampMissing,proto3,stdtime" json:"TimestampMissing"`
	// TimestampNullable *time.Time field
	TimestampNullable *time.Time `protobuf:"bytes,10,opt,name=TimestampNullable,proto3,stdtime" json:"TimestampNullable,omitempty"`
	// TimestampNullableWithNilValue *time.Time field
	TimestampNullableWithNilValue *time.Time `protobuf:"bytes,11,opt,name=TimestampNullableWithNilValue,proto3,stdtime" json:"TimestampNullableWithNilValue,omitempty"`
	// DurationStandard time.Duration field (standard)
	DurationStandard time.Duration `protobuf:"varint,12,opt,name=DurationStandard,proto3,stdduration" json:"DurationStandard,omitempty"`
	// DurationStandardMissing time.Duration field (standard) missing in input data
	DurationStandardMissing time.Duration `protobuf:"varint,13,opt,name=DurationStandardMissing,proto3,stdduration" json:"DurationStandardMissing,omitempty"`
	// DurationCustom time.Duration field (with casttype)
	DurationCustom Duration `protobuf:"varint,14,opt,name=DurationCustom,proto3,casttype=Duration" json:"DurationCustom,omitempty"`
	// DurationCustomMissing time.Duration field (with casttype) missing in input data
	DurationCustomMissing Duration `protobuf:"varint,15,opt,name=DurationCustomMissing,proto3,casttype=Duration" json:"DurationCustomMissing,omitempty"`
	// StringList []string field
	StringList []string `protobuf:"bytes,16,rep,name=StringList,proto3" json:"StringList,omitempty"`
	// StringListEmpty []string field
	StringListEmpty []string `protobuf:"bytes,17,rep,name=StringListEmpty,proto3" json:"StringListEmpty,omitempty"`
	// BoolCustomList []bool field
	BoolCustomList []BoolCustom `protobuf:"varint,18,rep,packed,name=BoolCustomList,proto3,customtype=BoolCustom" json:"BoolCustomList,omitempty"`
	// BytesList [][]byte field
	BytesList [][]byte `protobuf:"bytes,19,rep,name=BytesList,proto3" json:"BytesList,omitempty"`
	// TimestampList []time.Time field
	TimestampList []*time.Time `protobuf:"bytes,20,rep,name=TimestampList,proto3,stdtime" json:"TimestampList,omitempty"`
	// DurationCustomList []time.Duration field
	DurationCustomList []Duration `protobuf:"varint,21,rep,packed,name=DurationCustomList,proto3,casttype=Duration" json:"DurationCustomList,omitempty"`
	// Nested nested message field, non-nullable
	Nested Nested `protobuf:"bytes,22,opt,name=Nested,proto3" json:"Nested"`
	// NestedNullable nested message field, nullabel
	NestedNullable *Nested `protobuf:"bytes,23,opt,name=NestedNullable,proto3" json:"NestedNullable,omitempty"`
	// NestedNullableWithNilValue nested message field, with no value set
	NestedNullableWithNilValue *Nested `protobuf:"bytes,24,opt,name=NestedNullableWithNilValue,proto3" json:"NestedNullableWithNilValue,omitempty"`
	// NestedList nested message array
	NestedList []Nested `protobuf:"bytes,25,rep,name=NestedList,proto3" json:"NestedList"`
	// NestedListNullable nested message array
	NestedListNullable []*Nested `protobuf:"bytes,26,rep,name=NestedListNullable,proto3" json:"NestedListNullable,omitempty"`
	// Map normal map
	Map map[string]string `protobuf:"bytes,27,rep,name=Map,proto3" json:"Map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// MapObject is the object map
	MapObject map[string]Nested `protobuf:"bytes,29,rep,name=MapObject,proto3" json:"MapObject" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// MapObjectNullable is the object map with nullable values
	MapObjectNullable    map[string]*Nested `protobuf:"bytes,30,rep,name=MapObjectNullable,proto3" json:"MapObjectNullable,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
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
	NestedList []*OtherNested `protobuf:"bytes,2,rep,name=NestedList,proto3" json:"NestedList,omitempty"`
	// Nested map repeated nested messages
	Map map[string]string `protobuf:"bytes,3,rep,name=Map,proto3" json:"Map,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// MapObjectNested nested object map
	MapObjectNested      map[string]OtherNested `protobuf:"bytes,4,rep,name=MapObjectNested,proto3" json:"MapObjectNested" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
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

// OtherNested message nested into nested message
type OtherNested struct {
	// Str string field
	Str                  string   `protobuf:"bytes,1,opt,name=Str,proto3" json:"Str,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OtherNested) Reset()         { *m = OtherNested{} }
func (m *OtherNested) String() string { return proto.CompactTextString(m) }
func (*OtherNested) ProtoMessage()    {}
func (*OtherNested) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{2}
}
func (m *OtherNested) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OtherNested.Unmarshal(m, b)
}
func (m *OtherNested) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OtherNested.Marshal(b, m, deterministic)
}
func (m *OtherNested) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OtherNested.Merge(m, src)
}
func (m *OtherNested) XXX_Size() int {
	return xxx_messageInfo_OtherNested.Size(m)
}
func (m *OtherNested) XXX_DiscardUnknown() {
	xxx_messageInfo_OtherNested.DiscardUnknown(m)
}

var xxx_messageInfo_OtherNested proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Test)(nil), "test.Test")
	proto.RegisterMapType((map[string]string)(nil), "test.Test.MapEntry")
	proto.RegisterMapType((map[string]Nested)(nil), "test.Test.MapObjectEntry")
	proto.RegisterMapType((map[string]*Nested)(nil), "test.Test.MapObjectNullableEntry")
	proto.RegisterType((*Nested)(nil), "test.Nested")
	proto.RegisterMapType((map[string]string)(nil), "test.Nested.MapEntry")
	proto.RegisterMapType((map[string]OtherNested)(nil), "test.Nested.MapObjectNestedEntry")
	proto.RegisterType((*OtherNested)(nil), "test.OtherNested")
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_c161fcfdc0c3ff1e) }

var fileDescriptor_c161fcfdc0c3ff1e = []byte{
	// 788 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x55, 0x4d, 0x4f, 0xdb, 0x48,
	0x18, 0xce, 0xc4, 0x26, 0x9b, 0xbc, 0x84, 0x90, 0x0c, 0x01, 0x86, 0xec, 0x92, 0x78, 0x91, 0x56,
	0x58, 0x7b, 0x08, 0xbb, 0x80, 0x50, 0x85, 0xda, 0x1e, 0x5c, 0x40, 0x6a, 0xd5, 0x40, 0x3b, 0x40,
	0x7b, 0x76, 0x8a, 0x1b, 0x0c, 0x4e, 0x6c, 0xd9, 0x93, 0x56, 0xf9, 0x17, 0x3d, 0xf6, 0x27, 0x71,
	0xa9, 0xd4, 0x63, 0xd5, 0x03, 0x55, 0xf9, 0x19, 0x3d, 0x55, 0x9e, 0x89, 0x3f, 0xe3, 0xd0, 0xc2,
	0x6d, 0xde, 0xe7, 0x7d, 0x9e, 0xc7, 0xf3, 0x7e, 0x64, 0x02, 0xc0, 0x0c, 0x8f, 0xb5, 0x1d, 0xd7,
	0x66, 0x36, 0x96, 0xfd, 0x73, 0xa3, 0xde, 0xb3, 0x7b, 0x36, 0x07, 0x36, 0xfc, 0x93, 0xc8, 0x35,
	0x5a, 0x3d, 0xdb, 0xee, 0x59, 0xc6, 0x06, 0x8f, 0xba, 0xc3, 0xb7, 0x1b, 0xcc, 0xec, 0x1b, 0x1e,
	0xd3, 0xfb, 0xce, 0x98, 0xd0, 0x4c, 0x13, 0xde, 0xbb, 0xba, 0xe3, 0x18, 0xae, 0x27, 0xf2, 0x6b,
	0x9f, 0xe6, 0x40, 0x3e, 0x31, 0x3c, 0x86, 0xab, 0x20, 0x1d, 0x33, 0x97, 0x20, 0x05, 0xa9, 0x25,
	0xea, 0x1f, 0x71, 0x1d, 0x66, 0x9e, 0x0e, 0xd8, 0xd6, 0x26, 0xc9, 0x2b, 0x48, 0x9d, 0xa1, 0x22,
	0x18, 0xa3, 0x3b, 0xdb, 0x44, 0x52, 0x90, 0x2a, 0x51, 0x11, 0xf8, 0xe8, 0x81, 0x65, 0xeb, 0x8c,
	0xc8, 0x0a, 0x52, 0xf3, 0x54, 0x04, 0x78, 0x09, 0x0a, 0x7b, 0xf6, 0xb0, 0x6b, 0x19, 0x64, 0x46,
	0x41, 0x2a, 0xa2, 0xe3, 0x08, 0x63, 0x90, 0x35, 0xdb, 0xb6, 0x48, 0x41, 0x41, 0x6a, 0x91, 0xf2,
	0xb3, 0xef, 0xa0, 0x8d, 0x98, 0xe1, 0x91, 0x3f, 0x14, 0xa4, 0x96, 0xa9, 0x08, 0xb0, 0x06, 0xa5,
	0x93, 0xa0, 0x22, 0x52, 0x54, 0x90, 0x3a, 0xbb, 0xd9, 0x68, 0x8b, 0x92, 0xda, 0x41, 0x49, 0xed,
	0x90, 0xa1, 0x15, 0xaf, 0xae, 0x5b, 0xb9, 0x0f, 0xdf, 0x5a, 0x88, 0x46, 0x32, 0xfc, 0x02, 0xaa,
	0x61, 0xd0, 0x31, 0x3d, 0xcf, 0x1c, 0xf4, 0x48, 0xe9, 0x0e, 0x56, 0x13, 0x6a, 0x4c, 0xa1, 0x16,
	0x62, 0x87, 0x43, 0xcb, 0xd2, 0xfd, 0x12, 0xe1, 0xb7, 0x2c, 0x11, 0xb7, 0x9c, 0x94, 0xe3, 0x0b,
	0x58, 0x9d, 0x00, 0x5f, 0x9b, 0xec, 0xfc, 0xd0, 0xb4, 0x5e, 0xe9, 0xd6, 0xd0, 0x20, 0xb3, 0x77,
	0xf0, 0xbf, 0xdd, 0x0a, 0xff, 0x07, 0xd5, 0xbd, 0xa1, 0xab, 0x33, 0xd3, 0x1e, 0x1c, 0x33, 0x7d,
	0x70, 0xa6, 0xbb, 0x67, 0xa4, 0xec, 0x8f, 0x53, 0x93, 0x3f, 0xf2, 0x8a, 0xd3, 0x59, 0xfc, 0x18,
	0x96, 0xd3, 0x58, 0xd0, 0xca, 0xb9, 0x98, 0x70, 0x1a, 0x09, 0x6f, 0x43, 0x25, 0x48, 0x3d, 0x19,
	0x7a, 0xcc, 0xee, 0x93, 0x0a, 0x97, 0x95, 0x7f, 0x5c, 0xb7, 0x8a, 0x41, 0x86, 0xa6, 0x38, 0x58,
	0x83, 0xc5, 0x24, 0x12, 0x7c, 0x73, 0x3e, 0x43, 0x9c, 0x4d, 0xc5, 0x4d, 0x80, 0x63, 0xe6, 0x9a,
	0x83, 0xde, 0x73, 0xd3, 0x63, 0xa4, 0xaa, 0x48, 0x6a, 0x89, 0xc6, 0x10, 0xac, 0xc2, 0x7c, 0x14,
	0xed, 0xf7, 0x1d, 0x36, 0x22, 0x35, 0x4e, 0x4a, 0xc3, 0x78, 0x07, 0x2a, 0xfe, 0xa6, 0x0a, 0x7b,
	0xee, 0x86, 0x15, 0x49, 0x2d, 0x6a, 0x95, 0xaf, 0xd7, 0x2d, 0x88, 0x32, 0x34, 0xc5, 0xc2, 0x7f,
	0x41, 0x89, 0x2f, 0x33, 0x97, 0x2c, 0x28, 0x92, 0x5a, 0xa6, 0x11, 0x80, 0x0f, 0x60, 0x2e, 0x1c,
	0x16, 0x67, 0xd4, 0x15, 0xe9, 0x17, 0x73, 0x96, 0xf9, 0x8c, 0x93, 0x32, 0xfc, 0x10, 0x70, 0xb2,
	0x01, 0xdc, 0x6c, 0x51, 0x91, 0x26, 0x1a, 0x95, 0xc1, 0xc3, 0xff, 0x42, 0xe1, 0xd0, 0xf0, 0x98,
	0x71, 0x46, 0x96, 0xf8, 0x9a, 0x95, 0xdb, 0xfc, 0x01, 0x12, 0x98, 0x26, 0xfb, 0xbf, 0x05, 0x3a,
	0x66, 0xe0, 0x5d, 0xa8, 0x88, 0x53, 0xb8, 0xfa, 0xcb, 0x53, 0x34, 0x88, 0xa6, 0x98, 0x98, 0x42,
	0x23, 0x89, 0x24, 0x56, 0x9c, 0x4c, 0xf5, 0xb9, 0x45, 0x85, 0x37, 0x01, 0x44, 0x96, 0x57, 0xbc,
	0xc2, 0xdb, 0x97, 0x75, 0xff, 0x18, 0x0b, 0x6b, 0x80, 0xa3, 0x28, 0xac, 0xa3, 0x31, 0x45, 0x8b,
	0x68, 0x06, 0x1b, 0xff, 0x03, 0x52, 0x47, 0x77, 0xc8, 0x9f, 0x5c, 0xb4, 0x20, 0x44, 0xfe, 0x53,
	0xda, 0xee, 0xe8, 0xce, 0xfe, 0x80, 0xb9, 0x23, 0xea, 0xe7, 0xf1, 0x23, 0x28, 0x75, 0x74, 0xe7,
	0xa8, 0x7b, 0x61, 0xbc, 0x61, 0x64, 0x95, 0x93, 0x57, 0x92, 0x64, 0x91, 0xe3, 0x92, 0xf1, 0x55,
	0x23, 0x05, 0x3e, 0x85, 0x5a, 0x18, 0x84, 0x17, 0x6d, 0x72, 0x9b, 0xbf, 0xb3, 0x6c, 0x02, 0x4e,
	0x64, 0x87, 0xe8, 0xa4, 0x43, 0x63, 0x07, 0x8a, 0xc1, 0x35, 0xfd, 0xa7, 0xff, 0xd2, 0x18, 0x05,
	0x4f, 0xff, 0xa5, 0x31, 0xf2, 0x1f, 0xe3, 0x77, 0x7c, 0x22, 0x79, 0x8e, 0x89, 0x60, 0x37, 0xff,
	0x00, 0x35, 0x9e, 0x41, 0x25, 0x79, 0xe3, 0x0c, 0xf5, 0x5a, 0x5c, 0x9d, 0xea, 0x67, 0xdc, 0x8b,
	0xc2, 0x52, 0xf6, 0xb5, 0xef, 0xef, 0xb9, 0xf6, 0x25, 0x1f, 0x6c, 0x72, 0xc6, 0x3f, 0xda, 0xff,
	0x89, 0x4d, 0xc9, 0xf3, 0x26, 0xd6, 0x84, 0xd3, 0x11, 0x3b, 0x37, 0xdc, 0xb1, 0x5d, 0x7c, 0x51,
	0xd6, 0xc5, 0x90, 0x25, 0xce, 0x5d, 0x8c, 0x7f, 0x35, 0x35, 0xe6, 0x97, 0x30, 0x1f, 0x15, 0x23,
	0x7e, 0x4a, 0x72, 0x7c, 0x4a, 0x91, 0x28, 0xce, 0x89, 0x0f, 0x3d, 0xad, 0xbf, 0xf7, 0x8c, 0x4e,
	0xa1, 0x9e, 0xf5, 0x99, 0x0c, 0x8f, 0xf5, 0x64, 0x57, 0x33, 0x7a, 0x11, 0x6b, 0x6d, 0x0b, 0x66,
	0x63, 0x99, 0xc9, 0xf6, 0x6a, 0xe5, 0xab, 0xef, 0xcd, 0xdc, 0xd5, 0x4d, 0x33, 0xf7, 0xf9, 0xa6,
	0x99, 0xeb, 0x16, 0xf8, 0xcb, 0xb5, 0xf5, 0x33, 0x00, 0x00, 0xff, 0xff, 0x72, 0xef, 0x65, 0x76,
	0xcb, 0x08, 0x00, 0x00,
}
