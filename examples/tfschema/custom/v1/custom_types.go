package v1

import "github.com/gravitational/protoc-gen-terraform/v3/examples/types"

type BoolCustom = types.BoolCustom
type BoolCustomList = types.BoolCustomList
type OverrideCastType = types.OverrideCastType

var (
	GenSchemaBoolSpecial = types.GenSchemaBoolSpecial
	CopyToBoolSpecial    = types.CopyToBoolSpecial
	CopyFromBoolSpecial  = types.CopyFromBoolSpecial

	GenSchemaBoolSpecialList = types.GenSchemaBoolSpecialList
	CopyToBoolSpecialList    = types.CopyToBoolSpecialList
	CopyFromBoolSpecialList  = types.CopyFromBoolSpecialList

	GenSchemaStringCustom = types.GenSchemaStringCustom
	CopyToStringCustom    = types.CopyToStringCustom
	CopyFromStringCustom  = types.CopyFromStringCustom

	UseMockValidator    = types.UseMockValidator
	UseMockPlanModifier = types.UseMockPlanModifier
)
