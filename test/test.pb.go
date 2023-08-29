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

type Mode int32

const (
	Mode_UNKNOWN Mode = 0
	Mode_ON      Mode = 1
	Mode_OFF     Mode = 2
)

var Mode_name = map[int32]string{
	0: "UNKNOWN",
	1: "ON",
	2: "OFF",
}

var Mode_value = map[string]int32{
	"UNKNOWN": 0,
	"ON":      1,
	"OFF":     2,
}

func (x Mode) String() string {
	return proto.EnumName(Mode_name, int32(x))
}

func (Mode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{0}
}

// Test message definition.
type Test struct {
	// Str string field
	Str string `protobuf:"bytes,1,opt,name=Str,proto3" json:"str1"`
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
	// bytes byte[] field
	Bytes []byte `protobuf:"bytes,7,opt,name=bytes,proto3" json:"bytes,omitempty"`
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
	// NestedNullable nested message field, nullable
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
	MapObjectNullable map[string]*Nested `protobuf:"bytes,30,rep,name=MapObjectNullable,proto3" json:"MapObjectNullable,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// Mode is the enum value
	Mode Mode `protobuf:"varint,31,opt,name=Mode,proto3,enum=test.Mode" json:"Mode,omitempty"`
	// Excluded is the excluded field
	Excluded bool `protobuf:"varint,32,opt,name=Excluded,proto3" json:"Excluded,omitempty"`
	// Types that are valid to be assigned to OneOf:
	//
	//	*Test_Branch1
	//	*Test_Branch2
	//	*Test_Branch3
	OneOf isTest_OneOf `protobuf_oneof:"OneOf"`
	// Types that are valid to be assigned to OneOfWithEmptyMessage:
	//
	//	*Test_EmptyMessageBranch
	//	*Test_StringBranch
	OneOfWithEmptyMessage isTest_OneOfWithEmptyMessage `protobuf_oneof:"OneOfWithEmptyMessage"`
	// EmbeddedField encapsulates fields which can be shared among various types
	EmbeddedField        `protobuf:"bytes,38,opt,name=EmbeddedField,proto3,embedded=EmbeddedField" json:""`
	*MaxAgeDuration      `protobuf:"bytes,39,opt,name=EmbedNullable,proto3,embedded=EmbedNullable" json:""`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
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

type isTest_OneOf interface {
	isTest_OneOf()
}
type isTest_OneOfWithEmptyMessage interface {
	isTest_OneOfWithEmptyMessage()
}

type Test_Branch1 struct {
	Branch1 *Branch1 `protobuf:"bytes,33,opt,name=Branch1,proto3,oneof" json:"Branch1,omitempty"`
}
type Test_Branch2 struct {
	Branch2 *Branch2 `protobuf:"bytes,34,opt,name=Branch2,proto3,oneof" json:"Branch2,omitempty"`
}
type Test_Branch3 struct {
	Branch3 string `protobuf:"bytes,35,opt,name=Branch3,proto3,oneof" json:"Branch3,omitempty"`
}
type Test_EmptyMessageBranch struct {
	EmptyMessageBranch *EmptyMessageBranch `protobuf:"bytes,36,opt,name=EmptyMessageBranch,proto3,oneof" json:"EmptyMessageBranch,omitempty"`
}
type Test_StringBranch struct {
	StringBranch string `protobuf:"bytes,37,opt,name=StringBranch,proto3,oneof" json:"StringBranch,omitempty"`
}

func (*Test_Branch1) isTest_OneOf()                            {}
func (*Test_Branch2) isTest_OneOf()                            {}
func (*Test_Branch3) isTest_OneOf()                            {}
func (*Test_EmptyMessageBranch) isTest_OneOfWithEmptyMessage() {}
func (*Test_StringBranch) isTest_OneOfWithEmptyMessage()       {}

func (m *Test) GetOneOf() isTest_OneOf {
	if m != nil {
		return m.OneOf
	}
	return nil
}
func (m *Test) GetOneOfWithEmptyMessage() isTest_OneOfWithEmptyMessage {
	if m != nil {
		return m.OneOfWithEmptyMessage
	}
	return nil
}

func (m *Test) GetBranch1() *Branch1 {
	if x, ok := m.GetOneOf().(*Test_Branch1); ok {
		return x.Branch1
	}
	return nil
}

func (m *Test) GetBranch2() *Branch2 {
	if x, ok := m.GetOneOf().(*Test_Branch2); ok {
		return x.Branch2
	}
	return nil
}

func (m *Test) GetBranch3() string {
	if x, ok := m.GetOneOf().(*Test_Branch3); ok {
		return x.Branch3
	}
	return ""
}

func (m *Test) GetEmptyMessageBranch() *EmptyMessageBranch {
	if x, ok := m.GetOneOfWithEmptyMessage().(*Test_EmptyMessageBranch); ok {
		return x.EmptyMessageBranch
	}
	return nil
}

func (m *Test) GetStringBranch() string {
	if x, ok := m.GetOneOfWithEmptyMessage().(*Test_StringBranch); ok {
		return x.StringBranch
	}
	return ""
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*Test) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*Test_Branch1)(nil),
		(*Test_Branch2)(nil),
		(*Test_Branch3)(nil),
		(*Test_EmptyMessageBranch)(nil),
		(*Test_StringBranch)(nil),
	}
}

type MaxAgeDuration struct {
	Value                Duration `protobuf:"varint,1,opt,name=Value,proto3,casttype=Duration" json:"max_age"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MaxAgeDuration) Reset()         { *m = MaxAgeDuration{} }
func (m *MaxAgeDuration) String() string { return proto.CompactTextString(m) }
func (*MaxAgeDuration) ProtoMessage()    {}
func (*MaxAgeDuration) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{1}
}
func (m *MaxAgeDuration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MaxAgeDuration.Unmarshal(m, b)
}
func (m *MaxAgeDuration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MaxAgeDuration.Marshal(b, m, deterministic)
}
func (m *MaxAgeDuration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MaxAgeDuration.Merge(m, src)
}
func (m *MaxAgeDuration) XXX_Size() int {
	return xxx_messageInfo_MaxAgeDuration.Size(m)
}
func (m *MaxAgeDuration) XXX_DiscardUnknown() {
	xxx_messageInfo_MaxAgeDuration.DiscardUnknown(m)
}

var xxx_messageInfo_MaxAgeDuration proto.InternalMessageInfo

// EmptyMessageBranch message for empty oneof branch
type EmptyMessageBranch struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmptyMessageBranch) Reset()         { *m = EmptyMessageBranch{} }
func (m *EmptyMessageBranch) String() string { return proto.CompactTextString(m) }
func (*EmptyMessageBranch) ProtoMessage()    {}
func (*EmptyMessageBranch) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{2}
}
func (m *EmptyMessageBranch) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmptyMessageBranch.Unmarshal(m, b)
}
func (m *EmptyMessageBranch) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmptyMessageBranch.Marshal(b, m, deterministic)
}
func (m *EmptyMessageBranch) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmptyMessageBranch.Merge(m, src)
}
func (m *EmptyMessageBranch) XXX_Size() int {
	return xxx_messageInfo_EmptyMessageBranch.Size(m)
}
func (m *EmptyMessageBranch) XXX_DiscardUnknown() {
	xxx_messageInfo_EmptyMessageBranch.DiscardUnknown(m)
}

var xxx_messageInfo_EmptyMessageBranch proto.InternalMessageInfo

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
	return fileDescriptor_c161fcfdc0c3ff1e, []int{3}
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
	return fileDescriptor_c161fcfdc0c3ff1e, []int{4}
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

// Branch1 message is OneOf branch 1
type Branch1 struct {
	// Str string field
	Str                  string   `protobuf:"bytes,1,opt,name=Str,proto3" json:"Str,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Branch1) Reset()         { *m = Branch1{} }
func (m *Branch1) String() string { return proto.CompactTextString(m) }
func (*Branch1) ProtoMessage()    {}
func (*Branch1) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{5}
}
func (m *Branch1) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Branch1.Unmarshal(m, b)
}
func (m *Branch1) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Branch1.Marshal(b, m, deterministic)
}
func (m *Branch1) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Branch1.Merge(m, src)
}
func (m *Branch1) XXX_Size() int {
	return xxx_messageInfo_Branch1.Size(m)
}
func (m *Branch1) XXX_DiscardUnknown() {
	xxx_messageInfo_Branch1.DiscardUnknown(m)
}

var xxx_messageInfo_Branch1 proto.InternalMessageInfo

// Branch2 message is OneOf branch 2
type Branch2 struct {
	// Int32 int field
	Int32                int32    `protobuf:"varint,1,opt,name=Int32,proto3" json:"Int32,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Branch2) Reset()         { *m = Branch2{} }
func (m *Branch2) String() string { return proto.CompactTextString(m) }
func (*Branch2) ProtoMessage()    {}
func (*Branch2) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{6}
}
func (m *Branch2) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Branch2.Unmarshal(m, b)
}
func (m *Branch2) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Branch2.Marshal(b, m, deterministic)
}
func (m *Branch2) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Branch2.Merge(m, src)
}
func (m *Branch2) XXX_Size() int {
	return xxx_messageInfo_Branch2.Size(m)
}
func (m *Branch2) XXX_DiscardUnknown() {
	xxx_messageInfo_Branch2.DiscardUnknown(m)
}

var xxx_messageInfo_Branch2 proto.InternalMessageInfo

// EmbeddedField encapsulates fields which can be shared among various types
type EmbeddedField struct {
	// EmbeddedString string field
	EmbeddedString string `protobuf:"bytes,1,opt,name=EmbeddedString,proto3" json:"embedded_string"`
	// Nested EmbeddedNestedField field
	EmbeddedNestedField  *EmbeddedNestedField `protobuf:"bytes,2,opt,name=EmbeddedNestedField,proto3" json:"EmbeddedNestedField,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *EmbeddedField) Reset()         { *m = EmbeddedField{} }
func (m *EmbeddedField) String() string { return proto.CompactTextString(m) }
func (*EmbeddedField) ProtoMessage()    {}
func (*EmbeddedField) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{7}
}
func (m *EmbeddedField) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmbeddedField.Unmarshal(m, b)
}
func (m *EmbeddedField) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmbeddedField.Marshal(b, m, deterministic)
}
func (m *EmbeddedField) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmbeddedField.Merge(m, src)
}
func (m *EmbeddedField) XXX_Size() int {
	return xxx_messageInfo_EmbeddedField.Size(m)
}
func (m *EmbeddedField) XXX_DiscardUnknown() {
	xxx_messageInfo_EmbeddedField.DiscardUnknown(m)
}

var xxx_messageInfo_EmbeddedField proto.InternalMessageInfo

type EmbeddedNestedField struct {
	// EmbeddedNestedString string field
	EmbeddedNestedString string   `protobuf:"bytes,1,opt,name=EmbeddedNestedString,proto3" json:"embedded_nested_string"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EmbeddedNestedField) Reset()         { *m = EmbeddedNestedField{} }
func (m *EmbeddedNestedField) String() string { return proto.CompactTextString(m) }
func (*EmbeddedNestedField) ProtoMessage()    {}
func (*EmbeddedNestedField) Descriptor() ([]byte, []int) {
	return fileDescriptor_c161fcfdc0c3ff1e, []int{8}
}
func (m *EmbeddedNestedField) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EmbeddedNestedField.Unmarshal(m, b)
}
func (m *EmbeddedNestedField) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EmbeddedNestedField.Marshal(b, m, deterministic)
}
func (m *EmbeddedNestedField) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EmbeddedNestedField.Merge(m, src)
}
func (m *EmbeddedNestedField) XXX_Size() int {
	return xxx_messageInfo_EmbeddedNestedField.Size(m)
}
func (m *EmbeddedNestedField) XXX_DiscardUnknown() {
	xxx_messageInfo_EmbeddedNestedField.DiscardUnknown(m)
}

var xxx_messageInfo_EmbeddedNestedField proto.InternalMessageInfo

func init() {
	proto.RegisterEnum("test.Mode", Mode_name, Mode_value)
	proto.RegisterType((*Test)(nil), "test.Test")
	proto.RegisterMapType((map[string]string)(nil), "test.Test.MapEntry")
	proto.RegisterMapType((map[string]Nested)(nil), "test.Test.MapObjectEntry")
	proto.RegisterMapType((map[string]*Nested)(nil), "test.Test.MapObjectNullableEntry")
	proto.RegisterType((*MaxAgeDuration)(nil), "test.MaxAgeDuration")
	proto.RegisterType((*EmptyMessageBranch)(nil), "test.EmptyMessageBranch")
	proto.RegisterType((*Nested)(nil), "test.Nested")
	proto.RegisterMapType((map[string]string)(nil), "test.Nested.MapEntry")
	proto.RegisterMapType((map[string]OtherNested)(nil), "test.Nested.MapObjectNestedEntry")
	proto.RegisterType((*OtherNested)(nil), "test.OtherNested")
	proto.RegisterType((*Branch1)(nil), "test.Branch1")
	proto.RegisterType((*Branch2)(nil), "test.Branch2")
	proto.RegisterType((*EmbeddedField)(nil), "test.EmbeddedField")
	proto.RegisterType((*EmbeddedNestedField)(nil), "test.EmbeddedNestedField")
}

func init() { proto.RegisterFile("test.proto", fileDescriptor_c161fcfdc0c3ff1e) }

var fileDescriptor_c161fcfdc0c3ff1e = []byte{
	// 1141 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x56, 0xcb, 0x72, 0x1a, 0x47,
	0x14, 0xa5, 0x19, 0x24, 0xe0, 0x0a, 0x21, 0x74, 0x85, 0xa4, 0x36, 0x8e, 0x99, 0x31, 0xb1, 0xe3,
	0x89, 0xab, 0x82, 0x22, 0xe4, 0x52, 0xa5, 0x9c, 0x57, 0x65, 0x22, 0x29, 0x8e, 0x1d, 0x20, 0x69,
	0x59, 0xf1, 0xd2, 0x35, 0x88, 0x36, 0xc2, 0x01, 0x86, 0x62, 0x86, 0x94, 0xf4, 0x17, 0x59, 0x26,
	0x7f, 0x90, 0x4f, 0xd1, 0xd2, 0xcb, 0x54, 0x16, 0xa4, 0xe2, 0xa5, 0x3e, 0x21, 0xab, 0xd4, 0x74,
	0x33, 0x4f, 0x46, 0x4e, 0xec, 0x5d, 0xdf, 0x73, 0xcf, 0x39, 0xd3, 0xaf, 0x7b, 0xa7, 0x01, 0x1c,
	0x6e, 0x3b, 0xf5, 0xf1, 0xc4, 0x72, 0x2c, 0xcc, 0xb8, 0xe3, 0x4a, 0xb9, 0x67, 0xf5, 0x2c, 0x01,
	0xec, 0xb8, 0x23, 0x99, 0xab, 0xa8, 0x3d, 0xcb, 0xea, 0x0d, 0xf8, 0x8e, 0x88, 0x3a, 0xd3, 0x17,
	0x3b, 0x4e, 0x7f, 0xc8, 0x6d, 0xc7, 0x1c, 0x8e, 0x25, 0xa1, 0xf6, 0xfb, 0x3a, 0x64, 0x9e, 0x72,
	0xdb, 0xc1, 0x0a, 0x28, 0xc7, 0xce, 0x84, 0x12, 0x8d, 0xe8, 0x79, 0x23, 0x77, 0x35, 0x53, 0x33,
	0xb6, 0x33, 0xd9, 0x65, 0x2e, 0x88, 0x65, 0x58, 0xfa, 0x76, 0xe4, 0xec, 0x35, 0x68, 0x5a, 0x23,
	0xfa, 0x12, 0x93, 0xc1, 0x1c, 0xdd, 0x7f, 0x40, 0x15, 0x8d, 0xe8, 0x0a, 0x93, 0x81, 0x8b, 0x1e,
	0x0d, 0x2c, 0xd3, 0xa1, 0x19, 0x8d, 0xe8, 0x69, 0x26, 0x03, 0xdc, 0x82, 0xe5, 0x03, 0x6b, 0xda,
	0x19, 0x70, 0xba, 0xa4, 0x11, 0x9d, 0xb0, 0x79, 0x84, 0x08, 0x19, 0xc3, 0xb2, 0x06, 0x74, 0x59,
	0x23, 0x7a, 0x8e, 0x89, 0xb1, 0xeb, 0xd0, 0xb9, 0x70, 0xb8, 0x4d, 0xb3, 0x1a, 0xd1, 0x0b, 0x4c,
	0x06, 0x68, 0x40, 0xfe, 0xa9, 0x37, 0x77, 0x9a, 0xd3, 0x88, 0xbe, 0xd2, 0xa8, 0xd4, 0xe5, 0xea,
	0xea, 0xde, 0xea, 0xea, 0x3e, 0xc3, 0xc8, 0x5d, 0xce, 0xd4, 0xd4, 0x2f, 0x7f, 0xa9, 0x84, 0x05,
	0x32, 0xfc, 0x1e, 0x4a, 0x7e, 0xd0, 0xec, 0xdb, 0x76, 0x7f, 0xd4, 0xa3, 0xf9, 0xb7, 0xb0, 0x5a,
	0x50, 0x23, 0x83, 0x75, 0x1f, 0x6b, 0x4d, 0x07, 0x03, 0xd3, 0x5d, 0x22, 0xfc, 0x2f, 0x4b, 0x22,
	0x2c, 0x17, 0xe5, 0xf8, 0x12, 0x6e, 0x2d, 0x80, 0xcf, 0xfa, 0xce, 0x59, 0xab, 0x3f, 0xf8, 0xd1,
	0x1c, 0x4c, 0x39, 0x5d, 0x79, 0x0b, 0xff, 0x37, 0x5b, 0xe1, 0xc7, 0x50, 0x3a, 0x98, 0x4e, 0x4c,
	0xa7, 0x6f, 0x8d, 0x8e, 0x1d, 0x73, 0xd4, 0x35, 0x27, 0x5d, 0x5a, 0x70, 0x8f, 0xd3, 0xc8, 0xfc,
	0x2a, 0x56, 0x1c, 0xcf, 0xe2, 0x17, 0xb0, 0x1d, 0xc7, 0xbc, 0xad, 0x5c, 0x0d, 0x09, 0xaf, 0x23,
	0xe1, 0x03, 0x28, 0x7a, 0xa9, 0xaf, 0xa7, 0xb6, 0x63, 0x0d, 0x69, 0x51, 0xc8, 0x0a, 0xff, 0xcc,
	0xd4, 0x9c, 0x97, 0x61, 0x31, 0x0e, 0x1a, 0xb0, 0x19, 0x45, 0xbc, 0x6f, 0xae, 0x25, 0x88, 0x93,
	0xa9, 0x58, 0x05, 0x38, 0x76, 0x26, 0xfd, 0x51, 0xef, 0xbb, 0xbe, 0xed, 0xd0, 0x92, 0xa6, 0xe8,
	0x79, 0x16, 0x42, 0x50, 0x87, 0xb5, 0x20, 0x3a, 0x1c, 0x8e, 0x9d, 0x0b, 0xba, 0x2e, 0x48, 0x71,
	0x18, 0xf7, 0xa1, 0xe8, 0xde, 0x54, 0x69, 0x2f, 0xdc, 0x50, 0x53, 0xf4, 0x9c, 0x51, 0xfc, 0x73,
	0xa6, 0x42, 0x90, 0x61, 0x31, 0x16, 0xbe, 0x07, 0x79, 0xc3, 0xbd, 0xcc, 0x42, 0xb2, 0xa1, 0x29,
	0x7a, 0x81, 0x05, 0x00, 0x1e, 0xc1, 0xaa, 0x7f, 0x58, 0x82, 0x51, 0xd6, 0x94, 0xff, 0x38, 0xe7,
	0x8c, 0x38, 0xe3, 0xa8, 0x0c, 0x3f, 0x03, 0x8c, 0x6e, 0x80, 0x30, 0xdb, 0xd4, 0x94, 0x85, 0x8d,
	0x4a, 0xe0, 0xe1, 0x7d, 0x58, 0x6e, 0x71, 0xdb, 0xe1, 0x5d, 0xba, 0x25, 0xae, 0x59, 0xa1, 0x2e,
	0x5a, 0x8d, 0xc4, 0x8c, 0x8c, 0x5b, 0x0b, 0x6c, 0xce, 0xc0, 0x87, 0x50, 0x94, 0x23, 0xff, 0xea,
	0x6f, 0x5f, 0xa3, 0x21, 0x2c, 0xc6, 0x44, 0x06, 0x95, 0x28, 0x12, 0xb9, 0xe2, 0xf4, 0x5a, 0x9f,
	0x37, 0xa8, 0xb0, 0x01, 0x20, 0xb3, 0x62, 0xc5, 0x37, 0xc4, 0xf6, 0x25, 0xcd, 0x3f, 0xc4, 0x42,
	0x03, 0x30, 0x88, 0xfc, 0x75, 0x54, 0xae, 0xd1, 0x12, 0x96, 0xc0, 0xc6, 0xbb, 0xa0, 0x34, 0xcd,
	0x31, 0xbd, 0x29, 0x44, 0x1b, 0x52, 0xe4, 0x36, 0xd5, 0x7a, 0xd3, 0x1c, 0x1f, 0x8e, 0x9c, 0xc9,
	0x05, 0x73, 0xf3, 0xf8, 0x39, 0xe4, 0x9b, 0xe6, 0xb8, 0xdd, 0x79, 0xc9, 0x4f, 0x1d, 0x7a, 0x4b,
	0x90, 0x6f, 0x44, 0xc9, 0x32, 0x27, 0x24, 0xf3, 0xa9, 0x06, 0x0a, 0x3c, 0x81, 0x75, 0x3f, 0xf0,
	0x27, 0x5a, 0x15, 0x36, 0xb7, 0x93, 0x6c, 0x3c, 0x4e, 0x60, 0x47, 0xd8, 0xa2, 0x03, 0x56, 0x21,
	0xd3, 0xb4, 0xba, 0x9c, 0xaa, 0x1a, 0xd1, 0x8b, 0x0d, 0x90, 0x4e, 0x2e, 0xc2, 0x04, 0x8e, 0x15,
	0xc8, 0x1d, 0x9e, 0x9f, 0x0e, 0xa6, 0x5d, 0xde, 0xa5, 0x9a, 0x68, 0xd3, 0x7e, 0x8c, 0x1f, 0x42,
	0xd6, 0x98, 0x98, 0xa3, 0xd3, 0xb3, 0x5d, 0x7a, 0x5b, 0x9c, 0xd8, 0xaa, 0x94, 0xcf, 0xc1, 0x47,
	0x29, 0xe6, 0xe5, 0x03, 0x6a, 0x83, 0xd6, 0x16, 0xa9, 0x8d, 0x80, 0xda, 0xc0, 0x8a, 0x47, 0xdd,
	0xa3, 0xef, 0xbb, 0xbf, 0xa3, 0x20, 0xb7, 0x87, 0x8f, 0x01, 0x45, 0x0d, 0x36, 0xb9, 0x6d, 0x9b,
	0x3d, 0x2e, 0x61, 0x7a, 0x47, 0x38, 0x52, 0xe9, 0xb8, 0x98, 0x7f, 0x44, 0x58, 0x82, 0x0a, 0xef,
	0x40, 0x41, 0x56, 0xf6, 0xdc, 0xe5, 0xae, 0xf8, 0x18, 0x61, 0x11, 0x14, 0xbf, 0x81, 0xd5, 0xc3,
	0x61, 0x87, 0x77, 0xbb, 0xbc, 0x7b, 0xd4, 0xe7, 0x83, 0x2e, 0xfd, 0x40, 0x7c, 0x6c, 0xc3, 0xfb,
	0x58, 0x28, 0x65, 0x14, 0xdc, 0x33, 0x7b, 0x35, 0x53, 0xc9, 0x95, 0x7b, 0x76, 0x51, 0x1d, 0x1e,
	0xcc, 0x8d, 0xfc, 0xb3, 0xbb, 0x27, 0x8c, 0xca, 0xf3, 0x1d, 0x37, 0xcf, 0xbf, 0xea, 0x71, 0xaf,
	0x20, 0x8d, 0x5c, 0xcc, 0xc5, 0x13, 0x55, 0xf6, 0x21, 0xe7, 0xdd, 0x2a, 0x2c, 0x81, 0xf2, 0x13,
	0xbf, 0x90, 0xff, 0x6c, 0xe6, 0x0e, 0xdd, 0x7f, 0xe7, 0xcf, 0xa2, 0x80, 0xd2, 0x02, 0x93, 0xc1,
	0xc3, 0xf4, 0x27, 0xa4, 0xf2, 0x18, 0x8a, 0xd1, 0x0b, 0x96, 0xa0, 0xae, 0x85, 0xd5, 0xb1, 0xeb,
	0x1f, 0xf6, 0x62, 0xb0, 0x95, 0x7c, 0xcb, 0xde, 0xdd, 0xd3, 0xc8, 0xc2, 0x52, 0x7b, 0xc4, 0xdb,
	0x2f, 0x8c, 0x6d, 0xd8, 0x14, 0x03, 0xb7, 0xb2, 0xc3, 0x87, 0x56, 0xfb, 0xd2, 0x5d, 0x41, 0x78,
	0x93, 0xf0, 0x23, 0x58, 0x92, 0xed, 0x82, 0x88, 0xbf, 0xc0, 0xf6, 0xd5, 0x4c, 0xcd, 0x0e, 0xcd,
	0xf3, 0xe7, 0x66, 0x8f, 0x47, 0xfa, 0x9c, 0x64, 0xd5, 0xca, 0x49, 0x77, 0xa7, 0xf6, 0x47, 0xda,
	0xeb, 0x78, 0xee, 0xec, 0xfd, 0x37, 0x90, 0x7c, 0xf9, 0xec, 0x46, 0x3a, 0x4a, 0x5a, 0x14, 0xdb,
	0xba, 0x5c, 0x42, 0xdb, 0x39, 0xe3, 0x93, 0xf9, 0x3a, 0xc2, 0x0d, 0xe5, 0x9e, 0x6c, 0x06, 0x8a,
	0xe0, 0x6e, 0x86, 0x97, 0x1b, 0x6b, 0x07, 0x3f, 0xc0, 0x5a, 0xb0, 0x8b, 0xb2, 0xe5, 0x66, 0xc2,
	0xd5, 0x1c, 0x88, 0xc2, 0x9c, 0x70, 0x73, 0x88, 0xeb, 0xdf, 0xf9, 0x72, 0x9c, 0x40, 0x39, 0xe9,
	0x33, 0x09, 0x1e, 0xf7, 0xa2, 0xc7, 0x99, 0xb0, 0x17, 0x81, 0x6d, 0x4d, 0x85, 0x95, 0x50, 0x66,
	0x71, 0x7b, 0x6b, 0x37, 0xfd, 0xfe, 0x91, 0x90, 0x54, 0xfd, 0x8e, 0x11, 0x3c, 0x40, 0x49, 0xe8,
	0x01, 0x5a, 0xfb, 0x8d, 0xc4, 0x4a, 0x13, 0x3f, 0x85, 0xa2, 0x07, 0xc8, 0x1a, 0x9e, 0xbf, 0x67,
	0x37, 0xae, 0x66, 0xea, 0x1a, 0x9f, 0x67, 0x9e, 0xdb, 0x22, 0xc5, 0x62, 0x54, 0x7c, 0x02, 0x1b,
	0x1e, 0x22, 0x27, 0x2c, 0xcb, 0x5d, 0x2e, 0xf4, 0x46, 0xb4, 0xdc, 0x43, 0x04, 0x96, 0xa4, 0xaa,
	0xf1, 0x44, 0x33, 0x6c, 0x41, 0x39, 0x0a, 0x47, 0xa6, 0x59, 0xb9, 0x9a, 0xa9, 0x5b, 0xfe, 0x34,
	0x47, 0x82, 0xe0, 0xcd, 0x36, 0x51, 0x77, 0xff, 0x8e, 0x6c, 0xde, 0xb8, 0x02, 0xd9, 0x93, 0xd6,
	0x93, 0x56, 0xfb, 0x59, 0xab, 0x94, 0xc2, 0x65, 0x48, 0xb7, 0x5b, 0x25, 0x82, 0x59, 0x50, 0xda,
	0x47, 0x47, 0xa5, 0xb4, 0x51, 0xb8, 0xfc, 0xbb, 0x9a, 0xba, 0x7c, 0x5d, 0x4d, 0xbd, 0x7a, 0x5d,
	0x4d, 0x75, 0x96, 0xc5, 0x43, 0x62, 0xef, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4e, 0xfe, 0x42,
	0x84, 0x44, 0x0c, 0x00, 0x00,
}
