// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: test.proto

package test

import (
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	schema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	validation "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	math "math"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

func UnmarshalTest(d *schema.ResourceData, t *Test) error {
	meta := []struct {
		Name          string
		SchemaName    string
		SchemaRawType string
		SchemaGoType  string
	}{

		{
			Name:          "Str",
			SchemaName:    "str",
			SchemaRawType: "string",
			SchemaGoType:  "string",
		},

		{
			Name:          "Int32",
			SchemaName:    "int32",
			SchemaRawType: "int",
			SchemaGoType:  "int32",
		},

		{
			Name:          "Int64",
			SchemaName:    "int64",
			SchemaRawType: "int",
			SchemaGoType:  "int64",
		},

		{
			Name:          "Float",
			SchemaName:    "float",
			SchemaRawType: "float64",
			SchemaGoType:  "float32",
		},

		{
			Name:          "Double",
			SchemaName:    "double",
			SchemaRawType: "float64",
			SchemaGoType:  "double64",
		},

		{
			Name:          "Bool",
			SchemaName:    "bool",
			SchemaRawType: "bool",
			SchemaGoType:  "bool",
		},

		{
			Name:          "Bytes",
			SchemaName:    "bytes",
			SchemaRawType: "string",
			SchemaGoType:  "[]byte",
		},

		{
			Name:          "Timestamp",
			SchemaName:    "timestamp",
			SchemaRawType: "string",
			SchemaGoType:  "time.Time",
		},

		{
			Name:          "DurationStd",
			SchemaName:    "duration_std",
			SchemaRawType: "string",
			SchemaGoType:  "time.Duration",
		},

		{
			Name:          "DurationCustom",
			SchemaName:    "duration_custom",
			SchemaRawType: "string",
			SchemaGoType:  "time.Duration",
		},

		{
			Name:          "BoolN",
			SchemaName:    "bool_n",
			SchemaRawType: "bool",
			SchemaGoType:  "bool",
		},

		{
			Name:          "BytesN",
			SchemaName:    "bytes_n",
			SchemaRawType: "string",
			SchemaGoType:  "[]byte",
		},

		{
			Name:          "TimestampN",
			SchemaName:    "timestamp_n",
			SchemaRawType: "string",
			SchemaGoType:  "time.Time",
		},

		{
			Name:          "DurationN",
			SchemaName:    "duration_n",
			SchemaRawType: "string",
			SchemaGoType:  "time.Duration",
		},

		{
			Name:          "StringA",
			SchemaName:    "string_a",
			SchemaRawType: "string",
			SchemaGoType:  "string",
		},

		{
			Name:          "BoolA",
			SchemaName:    "bool_a",
			SchemaRawType: "bool",
			SchemaGoType:  "bool",
		},

		{
			Name:          "BytesA",
			SchemaName:    "bytes_a",
			SchemaRawType: "string",
			SchemaGoType:  "[]byte",
		},

		{
			Name:          "TimestampA",
			SchemaName:    "timestamp_a",
			SchemaRawType: "string",
			SchemaGoType:  "time.Time",
		},

		{
			Name:          "GracePeriodA",
			SchemaName:    "grace_period_a",
			SchemaRawType: "string",
			SchemaGoType:  "time.Duration",
		},

		{
			Name:          "Nested",
			SchemaName:    "nested",
			SchemaRawType: "",
			SchemaGoType:  "",
		},

		{
			Name:          "NestedA",
			SchemaName:    "nested_a",
			SchemaRawType: "",
			SchemaGoType:  "",
		},
	}
}
