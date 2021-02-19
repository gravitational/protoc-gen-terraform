package test

import (
	"log"
	"testing"
)

// Custom duration type
type Duration int64

// Custom bool array
type BoolCustomArray []bool

// https://pkg.go.dev/github.com/hashicorp/terraform-plugin-sdk@v1.16.0/helper/resource
// https://www.terraform.io/docs/extend/best-practices/testing.html
// schema.TestResourceDataRaw

func TestAll(*testing.T) {
	log.Println("AKA")
}
