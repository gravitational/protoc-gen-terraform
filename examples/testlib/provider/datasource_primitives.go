package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v4/examples/tfschema/primitives/v1"
)

var _ datasource.DataSource = &primitivesDataSource{}

type primitivesDataSource struct {
	p *exampleProvider
}

func (d primitivesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_primitives"
}

func (d primitivesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema, diags := schemav1.GenSchemaPrimitivesDataSource(ctx)
	resp.Schema = schema
	resp.Diagnostics.Append(diags...)
}

func (d primitivesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config types.Object
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	d.p.RLock()
	primitives, ok := d.p.primitives[id.ValueString()]
	d.p.RUnlock()
	if !ok {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("primitives not found", "no example_primitives resource exists with the provided id"))
		return
	}

	result, diags := schemav1.CopyPrimitivesToTerraform(ctx, primitives, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}
