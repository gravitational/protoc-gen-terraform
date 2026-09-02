package provider

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v4/examples/tfschema/primitives/v1"
	extypes "github.com/gravitational/protoc-gen-terraform/v4/examples/types"
)

var _ resource.Resource = &primitivesResource{}

type primitivesResource struct {
	p *exampleProvider
}

func (r primitivesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "example_primitives"
}

func (r primitivesResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return schemav1.GenSchemaPrimitivesResource(ctx)
}

func (r primitivesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	primitives := &extypes.Primitives{}
	resp.Diagnostics.Append(schemav1.CopyPrimitivesFromTerraform(ctx, plan, primitives)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// We proto marshall and unmarshall the resource to test stability through proto round-trip.
	// This is required because protouf does not preserve the distinction between nil and empty
	// for some types (lists for examples).
	wireMsg, err := proto.Marshal(primitives)
	if err != nil {
		resp.Diagnostics.AddError("failed to proto marshal", err.Error())
		return
	}

	var newPrimitives extypes.Primitives
	if err := proto.Unmarshal(wireMsg, &newPrimitives); err != nil {
		resp.Diagnostics.AddError("failed to proto unmarshal", err.Error())
		return
	}

	r.p.primitives[id] = &newPrimitives

	result, diags := schemav1.CopyPrimitivesToTerraform(ctx, &newPrimitives, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r primitivesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	primitives := r.p.primitives[id.ValueString()]

	result, diags := schemav1.CopyPrimitivesToTerraform(ctx, primitives, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r primitivesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	// We proto marshall and unmarshall the resource to test stability through proto round-trip.
	// This is required because protouf does not preserve the distinction between nil and empty
	// for some types (lists for examples).
	wireMsg, err := proto.Marshal(primitives)
	if err != nil {
		resp.Diagnostics.AddError("failed to proto marshal", err.Error())
		return
	}

	var newPrimitives extypes.Primitives
	if err := proto.Unmarshal(wireMsg, &newPrimitives); err != nil {
		resp.Diagnostics.AddError("failed to proto unmarshal", err.Error())
		return
	}

	r.p.primitives[newPrimitives.Id] = &newPrimitives

	result, diags := schemav1.CopyPrimitivesToTerraform(ctx, &newPrimitives, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r primitivesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

	delete(r.p.primitives, id.ValueString())
}
