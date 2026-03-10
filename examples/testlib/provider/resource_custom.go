package provider

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v3/examples/tfschema/custom/v1"
	extypes "github.com/gravitational/protoc-gen-terraform/v3/examples/types"
)

var _ tfsdk.ResourceType = &customResourceType{}

type customResourceType struct{}

func (c customResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemav1.GenSchemaCustom(ctx)
}

func (c customResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return customResource{
		p: p.(*exampleProvider),
	}, nil
}

var _ tfsdk.Resource = &customResource{}

type customResource struct {
	p *exampleProvider
}

func (r customResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("unable to generate uuid", err.Error()))
	}

	plan.Attrs["id"] = types.String{Value: id}
	plan.Attrs["computed"] = types.String{Value: "computed"}
	plan.Attrs["injected"] = types.String{Value: "injected"}

	custom := &extypes.Custom{}
	resp.Diagnostics.Append(schemav1.CopyCustomFromTerraform(ctx, plan, custom)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.custom[id] = custom

	resp.Diagnostics.Append(schemav1.CopyCustomToTerraform(ctx, custom, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r customResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	custom := r.p.custom[id.Value]

	resp.Diagnostics.Append(schemav1.CopyCustomToTerraform(ctx, custom, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read resource using 3rd party API.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r customResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	custom := &extypes.Custom{}
	resp.Diagnostics.Append(schemav1.CopyCustomFromTerraform(ctx, plan, custom)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.custom[custom.Id] = custom

	resp.Diagnostics.Append(schemav1.CopyCustomToTerraform(ctx, custom, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r customResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state types.Object
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var id types.String
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	delete(r.p.custom, id.Value)
}
