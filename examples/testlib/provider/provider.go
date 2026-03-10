package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	types "github.com/gravitational/protoc-gen-terraform/v3/examples/types"
)

var _ tfsdk.Provider = &exampleProvider{}

type exampleProvider struct {
	primitives map[string]*types.Primitives
	time       map[string]*types.Time
	objects    map[string]*types.Objects
	custom     map[string]*types.Custom
}

func New() tfsdk.Provider {
	return &exampleProvider{
		primitives: make(map[string]*types.Primitives),
		time:       make(map[string]*types.Time),
		objects:    make(map[string]*types.Objects),
		custom:     make(map[string]*types.Custom),
	}
}

// GetSchema satisfies the provider.Provider interface for exampleProvider.
func (p *exampleProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{}, nil
}

// Configure satisfies the provider.Provider interface for exampleProvider.
func (p *exampleProvider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Nothing to configure
}

// GetDataSources satisfies the provider.Provider interface for exampleProvider.
func (p *exampleProvider) GetDataSources(ctx context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		// TODO: Add example data source types
	}, nil
}

// GetResources satisfies the provider.Provider interface for exampleProvider.
func (p *exampleProvider) GetResources(ctx context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{
		"example_primitives": primitivesResourceType{},
		"example_time":       timeResourceType{},
		"example_objects":    objectsResourceType{},
		"example_custom":     customResourceType{},
	}, nil
}
