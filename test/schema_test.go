package test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestSchema(t *testing.T) {
	s, diags := GenSchemaTest(t.Context())
	require.False(t, diags.HasError())

	t.Run("string flags", func(t *testing.T) {
		str := stringAttribute(t, s, path.Root("str"))
		require.True(t, str.IsComputed())
		require.False(t, str.IsRequired())
		require.True(t, str.IsSensitive())
		require.Len(t, str.StringPlanModifiers(), 1)
		require.Len(t, str.StringValidators(), 1)
	})

	t.Run("required string", func(t *testing.T) {
		requiredStr := stringAttribute(t, s, path.Root("required_str"))
		require.False(t, requiredStr.IsComputed())
		require.True(t, requiredStr.IsRequired())
		require.False(t, requiredStr.IsSensitive())
	})

	t.Run("id is computed", func(t *testing.T) {
		require.True(t, stringAttribute(t, s, path.Root("id")).IsComputed())
	})

	t.Run("repeated custom scalar has element type", func(t *testing.T) {
		boolCustomList := listAttribute(t, s, path.Root("bool_custom_list"))
		require.Equal(t, "BoolCustomList []bool field", boolCustomList.GetDescription())
		require.True(t, boolCustomList.IsOptional())
		require.True(t, boolCustomList.ElementType.Equal(types.BoolType))
	})

	t.Run("custom schema override can change Terraform shape", func(t *testing.T) {
		stringOverride := listAttribute(t, s, path.Root("string_override"))
		require.True(t, stringOverride.IsComputed())
		require.True(t, stringOverride.IsOptional())
		require.True(t, stringOverride.ElementType.Equal(types.StringType))
		require.Len(t, stringOverride.ListPlanModifiers(), 1)
	})
}

func stringAttribute(t *testing.T, s schema.Schema, attrPath path.Path) schema.StringAttribute {
	t.Helper()

	attr := attributeAtPath(t, s, attrPath)
	stringAttr, ok := attr.(schema.StringAttribute)
	require.True(t, ok, "schema attribute %q should be schema.StringAttribute, got %T", attrPath.String(), attr)

	return stringAttr
}

func listAttribute(t *testing.T, s schema.Schema, attrPath path.Path) schema.ListAttribute {
	t.Helper()

	attr := attributeAtPath(t, s, attrPath)
	listAttr, ok := attr.(schema.ListAttribute)
	require.True(t, ok, "schema attribute %q should be schema.ListAttribute, got %T", attrPath.String(), attr)

	return listAttr
}

func attributeAtPath(t *testing.T, s schema.Schema, attrPath path.Path) schema.Attribute {
	t.Helper()

	attr, diags := s.AttributeAtPath(t.Context(), attrPath)
	require.False(t, diags.HasError())
	require.NotNil(t, attr, "schema attribute %q should exist", attrPath.String())

	return attr
}
