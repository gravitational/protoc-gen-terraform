package types

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type MockValidator struct{}

func UseMockValidator() validator.String {
	return MockValidator{}
}

func (v MockValidator) Description(_ context.Context) string {
	return "Mock validator"
}

func (v MockValidator) MarkdownDescription(_ context.Context) string {
	return "Mock validator"
}

func (v MockValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if req.ConfigValue.ValueString() != "valid" {
		resp.Diagnostics.AddError("mock error", fmt.Sprintf(
			`Attribute %q value must be "valid"`,
			req.Path.String()))
		return
	}
}
