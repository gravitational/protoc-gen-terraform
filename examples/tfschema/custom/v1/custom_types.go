package v1

import "github.com/gravitational/protoc-gen-terraform/v5/examples/types"

type BoolCustom = types.BoolCustom
type BoolCustomList = types.BoolCustomList
type OverrideCastType = types.OverrideCastType

var (
	CopyToBoolSpecial   = types.CopyToBoolSpecial
	CopyFromBoolSpecial = types.CopyFromBoolSpecial

	CopyToBoolSpecialList   = types.CopyToBoolSpecialList
	CopyFromBoolSpecialList = types.CopyFromBoolSpecialList

	CopyToStringCustom   = types.CopyToStringCustom
	CopyFromStringCustom = types.CopyFromStringCustom

	UseMockValidator    = types.UseMockValidator
	UseMockPlanModifier = types.UseMockPlanModifier
)
