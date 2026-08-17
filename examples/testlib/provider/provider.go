package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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

func (p *exampleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *exampleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return primitivesResource{p: p} },
		func() resource.Resource { return timeResource{p: p} },
		func() resource.Resource { return objectsResource{p: p} },
		func() resource.Resource { return customResource{p: p} },
	}
}
