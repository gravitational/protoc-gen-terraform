package provider

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v4/examples/tfschema/custom/v1"
	extypes "github.com/gravitational/protoc-gen-terraform/v4/examples/types"
)

var _ resource.Resource = &customResource{}

type customResource struct {
	p *exampleProvider
}

func (r customResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom"
}

func (r customResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema, diags := schemav1.GenSchemaCustomResource(ctx)
	resp.Schema = schema
	resp.Diagnostics.Append(diags...)
}

func (r customResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("unable to generate uuid", err.Error()))
	}

	plan.Attributes()["id"] = types.StringValue(id)
	plan.Attributes()["computed"] = types.StringValue("computed")
	plan.Attributes()["injected"] = types.StringValue("injected")

	custom := &extypes.Custom{}
	resp.Diagnostics.Append(schemav1.CopyCustomFromTerraform(ctx, plan, custom)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.Lock()
	r.p.custom[id] = custom
	r.p.Unlock()

	result, diags := schemav1.CopyCustomToTerraform(ctx, custom, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r customResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	r.p.RLock()
	custom, ok := r.p.custom[id.ValueString()]
	r.p.RUnlock()
	if !ok {
		resp.State.RemoveResource(ctx)
		return
	}

	result, diags := schemav1.CopyCustomToTerraform(ctx, custom, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r customResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	r.p.Lock()
	r.p.custom[custom.Id] = custom
	r.p.Unlock()

	result, diags := schemav1.CopyCustomToTerraform(ctx, custom, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r customResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

	r.p.Lock()
	delete(r.p.custom, id.ValueString())
	r.p.Unlock()
}

func (r customResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the entire plan is null, the resource is planned for destruction.
	if req.Plan.Raw.IsNull() {
		return
	}

	if req.State.Raw.IsNull() {
		return
	}

	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the provider-managed ID, but rewrite all other fields from
	// config so omitted or null values become explicit zero values in the plan.
	id, hasID := plan.Attributes()["id"]

	custom := &extypes.Custom{}
	resp.Diagnostics.Append(schemav1.CopyCustomFromTerraform(ctx, plan, custom)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, diags := schemav1.CopyCustomToTerraformPreserveUnknown(ctx, custom, &plan, true)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if hasID {
		result.Attributes()["id"] = id
	}

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &result)...)
}
