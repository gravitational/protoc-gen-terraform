package v1

import "github.com/gravitational/protoc-gen-terraform/v4/examples/types"

type BoolCustom = types.BoolCustom
type BoolCustomList = types.BoolCustomList
type OverrideCastType = types.OverrideCastType

var (
	GenSchemaBoolSpecialResource   = types.GenSchemaBoolSpecialResource
	GenSchemaBoolSpecialDataSource = types.GenSchemaBoolSpecialDataSource
	CopyToBoolSpecial              = types.CopyToBoolSpecial
	CopyFromBoolSpecial            = types.CopyFromBoolSpecial

	GenSchemaBoolSpecialListResource   = types.GenSchemaBoolSpecialListResource
	GenSchemaBoolSpecialListDataSource = types.GenSchemaBoolSpecialListDataSource
	CopyToBoolSpecialList              = types.CopyToBoolSpecialList
	CopyFromBoolSpecialList            = types.CopyFromBoolSpecialList

	GenSchemaStringCustomResource   = types.GenSchemaStringCustomResource
	GenSchemaStringCustomDataSource = types.GenSchemaStringCustomDataSource
	CopyToStringCustom              = types.CopyToStringCustom
	CopyFromStringCustom            = types.CopyFromStringCustom

	UseMockValidator    = types.UseMockValidator
	UseMockPlanModifier = types.UseMockPlanModifier
)
