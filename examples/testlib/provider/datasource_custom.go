package provider

import (
	"context"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v5/examples/tfschema/custom/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &customDataSource{}

type customDataSource struct {
	p *exampleProvider
}

func (d customDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "example_custom"
}

func (d customDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schema, diags := schemav1.GenSchemaCustomDataSource(ctx)
	resp.Schema = schema
	resp.Diagnostics.Append(diags...)
}

func (d customDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
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

	custom, ok := d.p.custom[id.ValueString()]
	if !ok {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("custom not found", "no example_custom resource exists with the provided id"))
		return
	}

	config.Attributes()["injected"] = types.StringValue("injected")

	result, diags := schemav1.CopyCustomToTerraform(ctx, custom, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}
