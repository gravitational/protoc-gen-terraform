package types

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCopyToBoolSpecialPreserveUnknownWithNil(t *testing.T) {
	var diags diag.Diagnostics

	result := CopyToBoolSpecial(diags, BoolCustom(true), types.BoolType, nil, true)

	boolValue, ok := result.(types.Bool)
	require.True(t, ok)
	require.Equal(t, types.BoolValue(true), boolValue)
}

func TestCopyToBoolSpecialPreserveUnknown(t *testing.T) {
	var diags diag.Diagnostics

	result := CopyToBoolSpecial(diags, BoolCustom(true), types.BoolType, types.BoolUnknown(), true)

	boolValue, ok := result.(types.Bool)
	require.True(t, ok)
	require.Equal(t, types.BoolUnknown(), boolValue)
}

func TestCopyToBoolSpecialListPreserveUnknown(t *testing.T) {
	var diags diag.Diagnostics

	val := types.ListUnknown(types.BoolType)
	result := CopyToBoolSpecialList(diags, []BoolCustomList{true}, types.ListType{ElemType: types.BoolType}, val, true)

	listValue, ok := result.(types.List)
	require.True(t, ok)
	require.Equal(t, types.ListUnknown(types.BoolType), listValue)
}
