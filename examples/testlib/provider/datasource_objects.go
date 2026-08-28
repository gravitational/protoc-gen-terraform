package provider

import (
	"context"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v5/examples/tfschema/objects/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &objectsDataSource{}

type objectsDataSource struct {
	p *exampleProvider
}

func (d objectsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "example_objects"
}

func (d objectsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema, diags := schemav1.GenSchemaObjectsDataSource(ctx)
	resp.Schema = schema
	resp.Diagnostics.Append(diags...)
}

func (d objectsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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
	objects, ok := d.p.objects[id.ValueString()]
	d.p.RUnlock()
	if !ok {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("objects not found", "no example_objects resource exists with the provided id"))
		return
	}

	result, diags := schemav1.CopyObjectsToTerraform(ctx, objects, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}
