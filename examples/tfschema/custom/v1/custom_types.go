package v1

import "github.com/gravitational/protoc-gen-terraform/v3/examples/types"

type BoolCustom = types.BoolCustom
type OverrideCastType = types.OverrideCastType

var (
	GenSchemaBoolSpecial = types.GenSchemaBoolSpecial
	CopyToBoolSpecial    = types.CopyToBoolSpecial
	CopyFromBoolSpecial  = types.CopyFromBoolSpecial

	GenSchemaStringCustom = types.GenSchemaStringCustom
	CopyToStringCustom    = types.CopyToStringCustom
	CopyFromStringCustom  = types.CopyFromStringCustom
)
