package provider

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v3/examples/tfschema/computed/v1"
	extypes "github.com/gravitational/protoc-gen-terraform/v3/examples/types"
)

var _ tfsdk.ResourceType = &computedResourceType{}

type computedResourceType struct{}

func (c computedResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemav1.GenSchemaComputed(ctx)
}

func (c computedResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return computedResource{
		p: p.(*exampleProvider),
	}, nil
}

var _ tfsdk.ResourceWithModifyPlan = &computedResource{}

type computedResource struct {
	p *exampleProvider
}

func (r computedResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("unable to generate uuid", err.Error()))
	}

	plan.Attrs["id"] = types.String{Value: id}

	computed := &extypes.Computed{}
	resp.Diagnostics.Append(schemav1.CopyComputedFromTerraform(ctx, plan, computed)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.computed[id] = computed

	resp.Diagnostics.Append(schemav1.CopyComputedToTerraform(ctx, computed, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r computedResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
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

	computed := r.p.computed[id.Value]

	resp.Diagnostics.Append(schemav1.CopyComputedToTerraform(ctx, computed, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r computedResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	computed := &extypes.Computed{}
	resp.Diagnostics.Append(schemav1.CopyComputedFromTerraform(ctx, plan, computed)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.computed[computed.Id] = computed

	resp.Diagnostics.Append(schemav1.CopyComputedToTerraform(ctx, computed, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r computedResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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

	delete(r.p.computed, id.Value)
}

func (r computedResource) ModifyPlan(ctx context.Context, req tfsdk.ModifyResourcePlanRequest, resp *tfsdk.ModifyResourcePlanResponse) {
	// If the entire plan is null, the resource is planned for destruction.
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config types.Object
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the provider-managed ID, but rewrite all other fields from
	// config so omitted or null values become explicit zero values in the plan.
	id, hasID := plan.Attrs["id"]

	computed := &extypes.Computed{}
	resp.Diagnostics.Append(schemav1.CopyComputedFromTerraform(ctx, config, computed)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(schemav1.CopyComputedToTerraform(ctx, computed, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if hasID {
		plan.Attrs["id"] = id
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}
