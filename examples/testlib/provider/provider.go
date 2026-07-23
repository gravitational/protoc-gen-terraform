package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/gravitational/protoc-gen-terraform/v4/examples/types"
)

var _ provider.Provider = &exampleProvider{}

type exampleProvider struct {
	primitives map[string]*types.Primitives
	time       map[string]*types.Time
	objects    map[string]*types.Objects
	custom     map[string]*types.Custom
}

func New() provider.Provider {
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
func (p *exampleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Nothing to configure
}

// GetDataSources satisfies the provider.Provider interface for exampleProvider.
func (p *exampleProvider) GetDataSources(ctx context.Context) (map[string]provider.DataSourceType, diag.Diagnostics) {
	return map[string]provider.DataSourceType{
		// TODO: Add example data source types
	}, nil
}

// GetResources satisfies the provider.Provider interface for exampleProvider.
func (p *exampleProvider) GetResources(ctx context.Context) (map[string]provider.ResourceType, diag.Diagnostics) {
	return map[string]provider.ResourceType{
		"example_primitives": primitivesResourceType{},
		"example_time":       timeResourceType{},
		"example_objects":    objectsResourceType{},
		"example_custom":     customResourceType{},
	}, nil
}
