package test

import (
	"context"
	fmt "fmt"
	"strings"
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

// GenSchemaBoolSpecialResource generates custom field schema (bool list)
func GenSchemaBoolSpecialResource(_ context.Context, attr tfsdk.Attribute) tfsdk.Attribute {
	attr.Type = types.ListType{
		ElemType: types.BoolType,
	}
	return attr
}

// GenSchemaBoolSpecialDataSource generates custom field schema (bool list)
func GenSchemaBoolSpecialDataSource(_ context.Context, attr tfsdk.Attribute) tfsdk.Attribute {
	attr.Type = types.ListType{
		ElemType: types.BoolType,
	}
	return attr
}

// CopyFromBoolSpecial copies target value to the source
func CopyFromBoolSpecial(diags diag.Diagnostics, tf attr.Value, obj *[]BoolCustom) {
	v, ok := tf.(types.List)
	if !ok {
		diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.List", tf))
		return
	}

	arr := make([]BoolCustom, len(v.Elements()))
	for i, raw := range v.Elements() {
		el, ok := raw.(types.Bool)
		if !ok {
			diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.Bool", raw))
			return
		}

		if !el.IsNull() && !el.IsUnknown() {
			arr[i] = BoolCustom(el.ValueBool())
		}
	}

	*obj = arr
}

// CopyToBoolSpecial copies source value to the target
func CopyToBoolSpecial(diags diag.Diagnostics, obj []BoolCustom, t attr.Type, v attr.Value, preserveUnknown bool) attr.Value {
	value, ok := v.(types.List)
	if !ok {
		value = types.ListNull(types.BoolType)
	}

	if preserveUnknown && value.IsUnknown() {
		return types.ListUnknown(types.BoolType)
	}

	elems := value.Elements()
	if len(elems) != len(obj) {
		resized := make([]attr.Value, len(obj))
		copy(resized, elems)
		elems = resized
	}

	for i, b := range obj {
		if preserveUnknown && elems[i] != nil && elems[i].IsUnknown() {
			elems[i] = types.BoolUnknown()
			continue
		}

		elems[i] = types.BoolValue(bool(b))
	}

	return types.ListValueMust(types.BoolType, elems)
}

// StringCustom is a custom type that maps a Terraform List of string, onto a
// single go string by joining all elements with "/".

// GenSchemaStringCustomResource returns the StringCustom schema.
func GenSchemaStringCustomResource(_ context.Context, attr tfsdk.Attribute) tfsdk.Attribute {
	attr.Type = types.ListType{
		ElemType: types.StringType,
	}
	return attr
}

// GenSchemaStringCustomDataSource returns the StringCustom schema.
func GenSchemaStringCustomDataSource(_ context.Context, attr tfsdk.Attribute) tfsdk.Attribute {
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
	for _, raw := range v.Elements() {
		el, ok := raw.(types.String)
		if !ok {
			diags.AddError("Error reading value from Terraform", fmt.Sprintf("Failed to cast %T to types.Bool", raw))
			return
		}

		if !el.IsNull() && !el.IsUnknown() {
			items = append(items, el.ValueString())
		}
	}

	*obj = strings.Join(items, "/")
}

// CopyToStringCustom copies a source value (single string) into the Terraform
// value (a list of strings) by splitting the string on every "/".
func CopyToStringCustom(diags diag.Diagnostics, obj string, t attr.Type, v attr.Value, preserveUnknown bool) attr.Value {
	value, ok := v.(types.List)
	if !ok {
		value = types.ListNull(types.StringType)
	}

	if preserveUnknown && value.IsUnknown() {
		return types.ListUnknown(types.StringType)
	}

	parts := strings.Split(obj, "/")
	elems := value.Elements()
	if len(elems) != len(parts) {
		resized := make([]attr.Value, len(parts))
		copy(resized, elems)
		elems = resized
	}

	for i, b := range parts {
		if preserveUnknown && elems[i] != nil && elems[i].IsUnknown() {
			elems[i] = types.StringUnknown()
			continue
		}

		elems[i] = types.StringValue(b)
	}

	return types.ListValueMust(types.StringType, elems)
}

type OverrideCastType string
