package types

import (
	"context"
	fmt "fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MockPlanModifier struct{}

// DefaultRoleOptions returns the default implementation of the DefaultRoleOptionsModifier
func UseMockPlanModifier() tfsdk.AttributePlanModifier {
	return MockPlanModifier{}
}

func (m MockPlanModifier) Description(_ context.Context) string {
	return "Mock plan modifier"
}

func (m MockPlanModifier) MarkdownDescription(_ context.Context) string {
	return "Mock plan modifier"
}

func (m MockPlanModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeConfig == nil {
		return
	}

	value, ok := req.AttributeConfig.(types.String)
	if !ok {
		resp.Diagnostics.AddError("mock error", fmt.Sprintf(
			"Attribute %q can not be converted to StringValue",
			req.AttributePath.String()))
		return
	}

	if value.Null || value.Unknown {
		resp.AttributePlan = types.String{Value: "modified_value"}
		return
	}

	resp.AttributePlan = value
}
