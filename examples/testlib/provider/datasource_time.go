package provider

import (
	"context"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v5/examples/tfschema/time/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &timeDataSource{}

type timeDataSource struct {
	p *exampleProvider
}

func (d timeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "example_time"
}

func (d timeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema, diags := schemav1.GenSchemaTimeDataSource(ctx)
	resp.Schema = schema
	resp.Diagnostics.Append(diags...)
}

func (d timeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
	time, ok := d.p.time[id.ValueString()]
	d.p.RUnlock()
	if !ok {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("time not found", "no example_time resource exists with the provided id"))
		return
	}

	result, diags := schemav1.CopyTimeToTerraform(ctx, time, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}
