package provider

import (
	"context"
	"sync"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/gravitational/protoc-gen-terraform/v5/examples/types"
)

var _ provider.Provider = &exampleProvider{}

type exampleProvider struct {
	*exampleStore
}

type exampleStore struct {
	sync.RWMutex

	primitives map[string]*types.Primitives
	time       map[string]*types.Time
	objects    map[string]*types.Objects
	custom     map[string]*types.Custom
}

func New() provider.Provider {
	return &exampleProvider{
		exampleStore: &exampleStore{
			primitives: make(map[string]*types.Primitives),
			time:       make(map[string]*types.Time),
			objects:    make(map[string]*types.Objects),
			custom:     make(map[string]*types.Custom),
		},
	}
}

func (p *exampleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "example"
}

func (p *exampleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
}

// Configure satisfies the provider.Provider interface for exampleProvider.
func (p *exampleProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Nothing to configure
}

func (p *exampleProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource { return primitivesDataSource{p: p} },
		func() datasource.DataSource { return timeDataSource{p: p} },
		func() datasource.DataSource { return objectsDataSource{p: p} },
		func() datasource.DataSource { return customDataSource{p: p} },
	}
}

func (p *exampleProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource { return primitivesResource{p: p} },
		func() resource.Resource { return timeResource{p: p} },
		func() resource.Resource { return objectsResource{p: p} },
		func() resource.Resource { return customResource{p: p} },
	}
}
