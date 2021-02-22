package test

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk@v1.16.0/helper/resource
// https://www.terraform.io/docs/extend/best-practices/testing.html
// schema.TestResourceDataRaw

const (
	defaultTimestamp = "2022-10-12T07:20:50.52Z"
)

var (
	fixutre map[string]interface{} = map[string]interface{}{
		"str":               "TestString",
		"int32":             999,
		"int64":             998,
		"float":             18.1,
		"double":            18.4,
		"bool":              true,
		"bytes":             "TestBytes",
		"timestamp":         defaultTimestamp,
		"duration_std":      "1h",
		"duration_custom":   "1m",
		"timestamp_n":       defaultTimestamp,
		"string_a":          []interface{}{"TestString1", "TestString2"},
		"bool_a":            []interface{}{false, true, false},
		"bytes_a":           []interface{}{"TestBytes1", "TestBytes2"},
		"timestamp_a":       []interface{}{defaultTimestamp},
		"duration_custom_a": []interface{}{"1m"},

		"nested": []interface{}{
			map[string]interface{}{
				"str": "TestString",
			},
		},
	}
)

func buildSubject(t *testing.T) (*Test, error) {
	subject := &Test{}
	data := schema.TestResourceDataRaw(t, SchemaTest(), fixutre)
	err := UnmarshalTest(data, subject)
	return subject, err
}

func TestUnmarshal(t *testing.T) {
	_, err := buildSubject(t)
	if err != nil {
		panic(err)
	}
}

func TestElementaries(t *testing.T) {
	subject, _ := buildSubject(t)

	assert.Equal(t, subject.Str, "TestString", "Test.Str")
	assert.Equal(t, subject.Int32, int32(999), "Test.Int32")
	assert.Equal(t, subject.Int64, int64(998), "Test.Int64")
	assert.Equal(t, subject.Float, float32(18.1), "Test.Float")
	assert.Equal(t, subject.Double, float64(18.4), "Test.Dobule")
	assert.Equal(t, subject.Bool, true, "Test.Bool")
	assert.Equal(t, subject.Bytes, []byte("TestBytes"), "Test.Bytes")
}

func TestTimes(t *testing.T) {
	subject, _ := buildSubject(t)

	timestamp, _ := time.Parse(time.RFC3339, defaultTimestamp)
	durationStd, _ := time.ParseDuration("1h")
	durationCustom, _ := time.ParseDuration("1m")

	assert.Equal(t, subject.Timestamp, timestamp, "Test.Timestamp")
	assert.Equal(t, subject.DurationStd, durationStd, "Test.DurationStd")
	assert.Equal(t, subject.DurationCustom, Duration(durationCustom), "Test.DurationCustom")
	assert.Equal(t, *(subject.TimestampN), timestamp, "Test.TimestampN")
}

func TestArrays(t *testing.T) {
	subject, _ := buildSubject(t)

	timestamp, _ := time.Parse(time.RFC3339, defaultTimestamp)
	duration, _ := time.ParseDuration("1m")

	assert.Equal(t, subject.StringA, []string{"TestString1", "TestString2"}, "Test.StringA[0]")
	assert.Equal(t, subject.BoolA, []BoolCustom{false, true, false}, "Test.BoolA")
	assert.Equal(t, subject.BytesA, [][]byte{[]byte("TestBytes1"), []byte("TestBytes2")}, "Test.BytesA")
	assert.Equal(t, subject.TimestampA, []*time.Time{&timestamp}, "Test.TimestampA")
	assert.Equal(t, subject.DurationCustomA, []Duration{Duration(duration)}, "Test.DurationCustomA")
}

func TestNestedMessage(t *testing.T) {
	subject, _ := buildSubject(t)

	assert.Equal(t, subject.Nested.Str, "TestString", "Test.Nested.Str")
}
