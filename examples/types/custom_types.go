package types

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// GenSchemaBoolSpecial generates custom field schema (bool list)
func GenSchemaBoolSpecial(_ context.Context, attr tfsdk.Attribute) tfsdk.Attribute {
	attr.Type = types.BoolType
	return attr
}

// CopyFromBoolSpecial copies target value to the source
func CopyFromBoolSpecial(diags diag.Diagnostics, tf attr.Value, obj *BoolCustom) {
	v, ok := tf.(types.Bool)
	if !ok {
		diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.Bool", tf))
		return
	}

	if !v.Null && !v.Unknown {
		*obj = BoolCustom(v.Value)
	}
}

// CopyToBoolSpecial copies source value to the target
func CopyToBoolSpecial(diags diag.Diagnostics, obj BoolCustom, t attr.Type, v attr.Value, preserveUnknown bool) attr.Value {
	unknown := false
	if v != nil {
		unknown = preserveUnknown && v.IsUnknown()
	}

	value := types.Bool{
		Value:   bool(obj),
		Unknown: unknown,
	}

	return value
}

type BoolCustomList bool

// GenSchemaBoolSpecialList generates custom field schema (bool list)
func GenSchemaBoolSpecialList(_ context.Context, attr tfsdk.Attribute) tfsdk.Attribute {
	attr.Type = types.ListType{
		ElemType: types.BoolType,
	}
	return attr
}

// CopyFromBoolSpecialList copies target value to the source
func CopyFromBoolSpecialList(diags diag.Diagnostics, tf attr.Value, obj *[]BoolCustomList) {
	v, ok := tf.(types.List)
	if !ok {
		diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.List", tf))
		return
	}

	arr := make([]BoolCustomList, len(v.Elems))
	for i, raw := range v.Elems {
		el, ok := raw.(types.Bool)
		if !ok {
			diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.Bool", raw))
			return
		}

		if !el.Null && !el.Unknown {
			arr[i] = BoolCustomList(el.Value)
		}
	}

	*obj = arr
}

// CopyToBoolSpecialList copies source value to the target
func CopyToBoolSpecialList(diags diag.Diagnostics, obj []BoolCustomList, t attr.Type, v attr.Value, preserveUnknown bool) attr.Value {
	value, ok := v.(types.List)
	if !ok {
		value = types.List{
			Null:     true,
			ElemType: types.BoolType,
		}
	}
	value.Unknown = preserveUnknown && value.IsUnknown()

	if len(value.Elems) != len(obj) {
		newElems := make([]attr.Value, len(obj))
		copy(newElems, value.Elems)
		value.Elems = newElems
	}

	for i, b := range obj {
		elemUnknown := false
		if value.Elems[i] != nil {
			elemUnknown = preserveUnknown && value.Elems[i].IsUnknown()
		}

		value.Elems[i] = types.Bool{
			Value:   bool(b),
			Unknown: elemUnknown,
		}
	}

	return value
}

// StringCustom is a custom type that maps a Terraform List of string, onto a
// single go string by joining all elements with "/".

// GenSchemaStringCustom returns the StringCustom schema.
func GenSchemaStringCustom(_ context.Context, attr tfsdk.Attribute) tfsdk.Attribute {
	attr.Type = types.ListType{
		ElemType: types.StringType,
	}
	return attr
}

// CopyFromStringCustom copies the value from Terraform (a list of strings) into
// the source (a single string) by joining all elements with "/".
func CopyFromStringCustom(diags diag.Diagnostics, tf attr.Value, obj *string) {
	v, ok := tf.(types.List)
	if !ok {
		diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.List", tf))
		return
	}

	items := make([]string, 0)
	for _, raw := range v.Elems {
		el, ok := raw.(types.String)
		if !ok {
			diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.Bool", raw))
			return
		}

		if !el.Null && !el.Unknown {
			items = append(items, el.Value)
		}
	}

	*obj = strings.Join(items, "/")
}

// CopyToStringCustom copies a source value (single string) into the Terraform
// value (a list of strings) by splitting the string on every "/".
func CopyToStringCustom(diags diag.Diagnostics, obj string, t attr.Type, v attr.Value, preserveUnknown bool) attr.Value {
	value, ok := v.(types.List)
	if !ok {
		value = types.List{
			Null:     true,
			ElemType: types.StringType,
		}
	}
	value.Unknown = preserveUnknown && value.IsUnknown()

	parts := strings.Split(obj, "/")
	if len(value.Elems) != len(parts) {
		newElems := make([]attr.Value, len(parts))
		copy(newElems, value.Elems)
		value.Elems = newElems
	}

	for i, b := range parts {
		elemUnknown := false
		if value.Elems[i] != nil {
			elemUnknown = preserveUnknown && value.Elems[i].IsUnknown()
		}

		value.Elems[i] = types.String{
			Value:   b,
			Unknown: elemUnknown,
		}
	}

	return value
}

type OverrideCastType string
