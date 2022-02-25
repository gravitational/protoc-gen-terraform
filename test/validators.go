package test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// MockValidator ensures that a time is in the future
type MockValidator struct{}

// UseMockValidator returns TimeValueInFutureValidator
func UseMockValidator() tfsdk.AttributeValidator {
	return MockValidator{}
}

// Description returns validator description
func (v MockValidator) Description(_ context.Context) string {
	return "Mock validator"
}

// MarkdownDescription returns validator markdown description
func (v MockValidator) MarkdownDescription(_ context.Context) string {
	return "Mock validator"
}

// Validate performs the validation.
func (v MockValidator) Validate(_ context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {

}
