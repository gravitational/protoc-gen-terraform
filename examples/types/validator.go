package types

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MockValidator struct{}

func UseMockValidator() tfsdk.AttributeValidator {
	return MockValidator{}
}

func (v MockValidator) Description(_ context.Context) string {
	return "Mock validator"
}

func (v MockValidator) MarkdownDescription(_ context.Context) string {
	return "Mock validator"
}

func (v MockValidator) Validate(_ context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
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
		return
	}

	if value.Value != "valid" {
		resp.Diagnostics.AddError("mock error", fmt.Sprintf(
			`Attribute %q value must be "valid"`,
			req.AttributePath.String()))
		return
	}
}
