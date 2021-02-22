package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk@v1.16.0/helper/resource
// https://www.terraform.io/docs/extend/best-practices/testing.html
// schema.TestResourceDataRaw

var (
	fixutre map[string]interface{} = map[string]interface{}{
		"str":   "TestString",
		"int32": 999,
		"int64": 999,
	}
)

func buildSubject(t *testing.T) *Test {
	subject := &Test{}
	data := schema.TestResourceDataRaw(t, SchemaTest(), fixutre)
	UnmarshalTest(data, subject)
	return subject
}

func TestAll(t *testing.T) {
	subject := buildSubject(t)

	assert.Equal(t, subject.Str, "TestString", "Test.Str")
	assert.Equal(t, subject.Int32, int32(999), "Test.Int32")
	assert.Equal(t, subject.Int64, int64(999), "Test.Int64")
}
