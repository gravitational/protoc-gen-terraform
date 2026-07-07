package types

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestCopyToBoolSpecialPreserveUnknownWithNil(t *testing.T) {
	var diags diag.Diagnostics

	result := CopyToBoolSpecial(diags, BoolCustom(true), types.BoolType, nil, true)

	boolValue, ok := result.(types.Bool)
	require.True(t, ok)
	require.Equal(t, types.Bool{
		Null:    false,
		Unknown: false,
		Value:   true,
	}, boolValue)
}

func TestCopyToBoolSpecialPreserveUnknown(t *testing.T) {
	var diags diag.Diagnostics

	val := types.Bool{Unknown: true, Value: false}
	result := CopyToBoolSpecial(diags, BoolCustom(true), types.BoolType, val, true)

	boolValue, ok := result.(types.Bool)
	require.True(t, ok)
	require.Equal(t, types.Bool{
		Null:    false,
		Unknown: true,
		Value:   true,
	}, boolValue)
}

func TestCopyToBoolSpecialListPreserveUnknownWithEmptySlice(t *testing.T) {
	var diags diag.Diagnostics

	val := types.List{
		Unknown:  true,
		Elems:    []attr.Value{},
		ElemType: types.BoolType,
	}
	result := CopyToBoolSpecialList(diags, []BoolCustomList{true}, types.ListType{ElemType: types.BoolType}, val, true)

	listValue, ok := result.(types.List)
	require.True(t, ok)
	require.Equal(t, types.List{
		Null:     false,
		Unknown:  true,
		Elems:    []attr.Value{types.Bool{Null: false, Unknown: false, Value: true}},
		ElemType: types.BoolType,
	}, listValue)
}
