package provider

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v3/examples/tfschema/time/v1"
	extypes "github.com/gravitational/protoc-gen-terraform/v3/examples/types"
)

var _ tfsdk.ResourceType = &timeResourceType{}

type timeResourceType struct{}

func (c timeResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemav1.GenSchemaTime(ctx)
}

func (c timeResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return timeResource{
		p: p.(*exampleProvider),
	}, nil
}

var _ tfsdk.Resource = &timeResource{}

type timeResource struct {
	p *exampleProvider
}

func (r timeResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
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

	time := &extypes.Time{}
	resp.Diagnostics.Append(schemav1.CopyTimeFromTerraform(ctx, plan, time)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.time[id] = time

	resp.Diagnostics.Append(schemav1.CopyTimeToTerraform(ctx, time, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r timeResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
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

	time := r.p.time[id.Value]

	resp.Diagnostics.Append(schemav1.CopyTimeToTerraform(ctx, time, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read resource using 3rd party API.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r timeResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	time := &extypes.Time{}
	resp.Diagnostics.Append(schemav1.CopyTimeFromTerraform(ctx, plan, time)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.time[time.Id] = time

	resp.Diagnostics.Append(schemav1.CopyTimeToTerraform(ctx, time, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r timeResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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

	delete(r.p.time, id.Value)
}
