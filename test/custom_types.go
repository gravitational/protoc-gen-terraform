package test

import (
	"context"
	time "time"

	"github.com/gravitational/trace"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
func CopyFromBoolCustom(tf attr.Value, obj *[]BoolCustom) error {
	v, ok := tf.(types.List)
	if !ok {
		return trace.Errorf("Failed to cast %T to types.List", tf)
	}

	arr := make([]BoolCustom, len(v.Elems))
	for i, raw := range v.Elems {
		el, ok := raw.(types.Bool)
		if !ok {
			return trace.Errorf("Failed to cast %T to types.List", raw)
		}

		if !el.Null && !el.Unknown {
			arr[i] = BoolCustom(el.Value)
		}
	}

	*obj = arr

	return nil
}

// CopyToBoolCustom copies source value to the target
func CopyToBoolCustom(tf attr.Value, obj []BoolCustom) error {
	// *source = make([]BoolCustom, len(target))
	// if len(target) > 0 {
	// 	copy(*source, target)
	// }
	return nil
}
