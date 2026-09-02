package testhelpers

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CustomNullValueFunc func(attr.Type) (attr.Value, bool)

func NullAttrs(attrTypes map[string]attr.Type, customNullValue CustomNullValueFunc) map[string]attr.Value {
	attrs := make(map[string]attr.Value, len(attrTypes))

	for name, attrType := range attrTypes {
		attrs[name] = NullValue(attrType, customNullValue)
	}

	return attrs
}

func NullValue(attrType attr.Type, customNullValue CustomNullValueFunc) attr.Value {
	if customNullValue != nil {
		if value, ok := customNullValue(attrType); ok {
			return value
		}
	}

	if attrType.Equal(types.StringType) {
		return types.StringNull()
	}
	if attrType.Equal(types.Int64Type) {
		return types.Int64Null()
	}
	if attrType.Equal(types.Float64Type) {
		return types.Float64Null()
	}
	if attrType.Equal(types.BoolType) {
		return types.BoolNull()
	}
	switch t := attrType.(type) {
	case types.ListType:
		return types.ListNull(t.ElementType())
	case types.MapType:
		return types.MapNull(t.ElementType())
	case types.ObjectType:
		return types.ObjectNull(t.AttributeTypes())
	default:
		panic(fmt.Sprintf("unsupported fixture attr type %T", attrType))
	}
}
