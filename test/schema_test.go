package test

// func TestSchema(t *testing.T) {
// 	schema, diags := GenSchemaTest(context.Background())
// 	require.False(t, diags.HasError())
// 	require.True(t, schema.Attributes["str"].Computed)
// 	require.False(t, schema.Attributes["str"].Required)
// 	require.True(t, schema.Attributes["str"].Sensitive)
// 	require.False(t, schema.Attributes["required_str"].Computed)
// 	require.True(t, schema.Attributes["required_str"].Required)
// 	require.False(t, schema.Attributes["required_str"].Sensitive)
// 	require.True(t, schema.Attributes["id"].Computed)
// 	require.Len(t, schema.Attributes["str"].PlanModifiers, 1)
// 	require.Len(t, schema.Attributes["str"].Validators, 1)
// 	require.Equal(t, "BoolCustomList []bool field", schema.Attributes["bool_custom_list"].Description)
// 	require.True(t, schema.Attributes["bool_custom_list"].Optional)
// }
