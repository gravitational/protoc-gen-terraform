package provider

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v3/examples/tfschema/objects/v1"
	extypes "github.com/gravitational/protoc-gen-terraform/v3/examples/types"
)

var _ tfsdk.ResourceType = &objectsResourceType{}

type objectsResourceType struct{}

func (c objectsResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemav1.GenSchemaObjects(ctx)
}

func (c objectsResourceType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return objectsResource{
		p: p.(*exampleProvider),
	}, nil
}

var _ tfsdk.Resource = &objectsResource{}

type objectsResource struct {
	p *exampleProvider
}

func (r objectsResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
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

	objects := &extypes.Objects{}
	resp.Diagnostics.Append(schemav1.CopyObjectsFromTerraform(ctx, plan, objects)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.objects[id] = objects

	resp.Diagnostics.Append(schemav1.CopyObjectsToTerraform(ctx, objects, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r objectsResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
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

	objects := r.p.objects[id.Value]

	resp.Diagnostics.Append(schemav1.CopyObjectsToTerraform(ctx, objects, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r objectsResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan types.Object
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	objects := &extypes.Objects{}
	resp.Diagnostics.Append(schemav1.CopyObjectsFromTerraform(ctx, plan, objects)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.p.objects[objects.Id] = objects

	resp.Diagnostics.Append(schemav1.CopyObjectsToTerraform(ctx, objects, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r objectsResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
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

	delete(r.p.objects, id.Value)
}
