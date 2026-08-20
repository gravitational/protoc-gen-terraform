package types

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MockPlanModifier struct{}

// DefaultRoleOptions returns the default implementation of the DefaultRoleOptionsModifier
func UseMockPlanModifier() planmodifier.String {
	return MockPlanModifier{}
}

func (m MockPlanModifier) Description(_ context.Context) string {
	return "Mock plan modifier"
}

func (m MockPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Mock plan modifier"
}

func (m MockPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = types.StringValue("modified_value")
		return
	}
}
