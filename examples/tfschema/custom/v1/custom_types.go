package v1

import (
	"github.com/gravitational/protoc-gen-terraform/v5/examples/types"
)

type BoolCustom = types.BoolCustom
type BoolCustomList = types.BoolCustomList
type OverrideCastType = types.OverrideCastType

var (
	GenSchemaBoolSpecialDataSource = types.GenSchemaBoolSpecialDataSource
	GenSchemaBoolSpecialResource   = types.GenSchemaBoolSpecialResource
	CopyToBoolSpecial              = types.CopyToBoolSpecial
	CopyFromBoolSpecial            = types.CopyFromBoolSpecial

	GenSchemaBoolSpecialListDataSource = types.GenSchemaBoolSpecialListDataSource
	GenSchemaBoolSpecialListResource   = types.GenSchemaBoolSpecialListResource
	CopyToBoolSpecialList              = types.CopyToBoolSpecialList
	CopyFromBoolSpecialList            = types.CopyFromBoolSpecialList

	GenSchemaStringCustomDataSource = types.GenSchemaStringCustomDataSource
	GenSchemaStringCustomResource   = types.GenSchemaStringCustomResource
	CopyToStringCustom              = types.CopyToStringCustom
	CopyFromStringCustom            = types.CopyFromStringCustom

	GenSchemaBoolOptionDataSource = types.GenSchemaBoolOptionDataSource
	GenSchemaBoolOptionResource   = types.GenSchemaBoolOptionResource
	CopyToBoolOption              = types.CopyToBoolOption
	CopyFromBoolOption            = types.CopyFromBoolOption

	UseMockValidator    = types.UseMockValidator
	UseMockPlanModifier = types.UseMockPlanModifier
)
