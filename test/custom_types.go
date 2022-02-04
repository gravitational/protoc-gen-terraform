package test

import (
	"context"
	fmt "fmt"
	time "time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	diag "github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Duration custom duration type
type Duration int64

// String returns duration string representation, must be implemented for custom duration type
func (d Duration) String() string {
	return time.Duration(d).String()
}

// BoolCustom custom bool array
type BoolCustom bool

// GenSchemaBoolCustom generates custom field schema (bool list)
func GenSchemaBoolCustom(_ context.Context) tfsdk.Attribute {
	return tfsdk.Attribute{
		Type: types.ListType{
			ElemType: types.BoolType,
		},
	}
}

// CopyFromBoolCustom copies target value to the source
func CopyFromBoolCustom(diags diag.Diagnostics, tf attr.Value, obj *[]BoolCustom) {
	v, ok := tf.(types.List)
	if !ok {
		diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.List", tf))
		return
	}

	arr := make([]BoolCustom, len(v.Elems))
	for i, raw := range v.Elems {
		el, ok := raw.(types.Bool)
		if !ok {
			diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.Bool", raw))
			return
		}

		if !el.Null && !el.Unknown {
			arr[i] = BoolCustom(el.Value)
		}
	}

	*obj = arr
}

// CopyToBoolCustom copies source value to the target
func CopyToBoolCustom(diags diag.Diagnostics, obj []BoolCustom) attr.Value {
	v := types.List{
		Null:     true,
		Unknown:  false,
		ElemType: types.BoolType,
	}

	if len(obj) > 0 {
		v.Null = false
		v.Elems = make([]attr.Value, len(obj))

		for i, b := range obj {
			v.Elems[i] = types.Bool{Null: false, Unknown: false, Value: bool(b)}
		}
	}

	return v
}
