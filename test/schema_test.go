package test

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestSchema(t *testing.T) {
	s, diags := GenSchemaTestResource(t.Context())
	require.False(t, diags.HasError())

	t.Run("string flags", func(t *testing.T) {
		attr, diags := s.AttributeAtPath(t.Context(), path.Root("str"))
		require.False(t, diags.HasError())

		str, ok := attr.(tfsdk.Attribute)
		require.True(t, ok)

		require.True(t, str.IsComputed())
		require.False(t, str.IsRequired())
		require.True(t, str.IsSensitive())
		require.Len(t, str.GetPlanModifiers(), 1)
		require.Len(t, str.GetValidators(), 1)
	})

	t.Run("required string", func(t *testing.T) {
		attr, diags := s.AttributeAtPath(t.Context(), path.Root("required_str"))
		require.False(t, diags.HasError())

		requiredStr, ok := attr.(tfsdk.Attribute)
		require.True(t, ok)

		require.False(t, requiredStr.IsComputed())
		require.True(t, requiredStr.IsRequired())
		require.False(t, requiredStr.IsSensitive())
	})

	t.Run("id is computed", func(t *testing.T) {
		attr, diags := s.AttributeAtPath(t.Context(), path.Root("id"))
		require.False(t, diags.HasError())

		id, ok := attr.(tfsdk.Attribute)
		require.True(t, ok)

		require.True(t, id.IsComputed())
	})

	t.Run("repeated custom scalar has element type", func(t *testing.T) {
		attr, diags := s.AttributeAtPath(t.Context(), path.Root("bool_custom_list"))
		require.False(t, diags.HasError())

		boolCustomList, ok := attr.(tfsdk.Attribute)
		require.True(t, ok)

		listType, ok := boolCustomList.GetType().(types.ListType)
		require.True(t, ok)

		require.Equal(t, "BoolCustomList []bool field", boolCustomList.GetDescription())
		require.True(t, boolCustomList.IsOptional())
		require.True(t, listType.ElementType().Equal(types.BoolType))
	})

	t.Run("custom schema override can change Terraform shape", func(t *testing.T) {
		attr, diags := s.AttributeAtPath(t.Context(), path.Root("string_override"))
		require.False(t, diags.HasError())

		stringOverride, ok := attr.(tfsdk.Attribute)
		require.True(t, ok)

		listType, ok := stringOverride.GetType().(types.ListType)
		require.True(t, ok)

		require.True(t, stringOverride.IsComputed())
		require.True(t, stringOverride.IsOptional())
		require.True(t, listType.ElementType().Equal(types.StringType))
		require.Len(t, stringOverride.GetPlanModifiers(), 1)
	})
}

func TestDataSourceSchema(t *testing.T) {
	s, diags := GenSchemaTestDataSource(context.Background())
	require.False(t, diags.HasError())

	t.Run("data source attributes do not support plan modifiers ", func(t *testing.T) {
		attr, diags := s.AttributeAtPath(t.Context(), path.Root("str"))
		require.False(t, diags.HasError())

		str, ok := attr.(tfsdk.Attribute)
		require.True(t, ok)

		require.Empty(t, str.GetPlanModifiers())
		require.Len(t, str.GetValidators(), 1)
	})
}
