package provider

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v3/examples/tfschema/primitives/v1"
	extypes "github.com/gravitational/protoc-gen-terraform/v3/examples/types"
)

var _ tfsdk.ResourceType = &primitivesResourceType{}

type primitivesResourceType struct{}

func (c primitivesResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemav1.GenSchemaPrimitives(ctx)
}

func (c primitivesResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return primitivesResource{
		p: p.(*exampleProvider),
	}, nil
}

var _ tfsdk.Resource = &primitivesResource{}

type primitivesResource struct {
	p *exampleProvider
}

func (r primitivesResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
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

	primitives := &extypes.Primitives{}
	resp.Diagnostics.Append(schemav1.CopyPrimitivesFromTerraform(ctx, plan, primitives)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.primitives[id] = primitives

	resp.Diagnostics.Append(schemav1.CopyPrimitivesToTerraform(ctx, primitives, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r primitivesResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
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

	primitives := r.p.primitives[id.Value]

	resp.Diagnostics.Append(schemav1.CopyPrimitivesToTerraform(ctx, primitives, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read resource using 3rd party API.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r primitivesResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	primitives := &extypes.Primitives{}
	resp.Diagnostics.Append(schemav1.CopyPrimitivesFromTerraform(ctx, plan, primitives)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.primitives[primitives.Id] = primitives

	resp.Diagnostics.Append(schemav1.CopyPrimitivesToTerraform(ctx, primitives, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r primitivesResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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

	delete(r.p.primitives, id.Value)
}
