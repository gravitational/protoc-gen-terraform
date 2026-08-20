package test

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// MockValidator ensures that a time is in the future
type MockValidator struct{}

// UseMockValidator returns TimeValueInFutureValidator
func UseMockValidator() validator.String {
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
func (v MockValidator) ValidateString(_ context.Context, req validator.StringRequest, resp *validator.StringResponse) {

}
