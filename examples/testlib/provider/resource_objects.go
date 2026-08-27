package provider

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	schemav1 "github.com/gravitational/protoc-gen-terraform/v5/examples/tfschema/objects/v1"
	extypes "github.com/gravitational/protoc-gen-terraform/v5/examples/types"
)

var _ resource.Resource = &objectsResource{}

type objectsResource struct {
	p *exampleProvider
}

func (r objectsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "example_objects"
}

func (r objectsResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema, diags := schemav1.GenSchemaObjectsResource(ctx)
	resp.Schema = schema
	resp.Diagnostics.Append(diags...)
}

func (r objectsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
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

	objects := &extypes.Objects{}
	resp.Diagnostics.Append(schemav1.CopyObjectsFromTerraform(ctx, plan, objects)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// We proto marshall and unmarshall the resource to test stability through proto round-trip.
	// This is required because protouf does not preserve the distinction between nil and empty
	// for some types (lists for examples).
	wireMsg, err := proto.Marshal(objects)
	if err != nil {
		resp.Diagnostics.AddError("failed to proto marshal", err.Error())
		return
	}

	var newObjects extypes.Objects
	if err := proto.Unmarshal(wireMsg, &newObjects); err != nil {
		resp.Diagnostics.AddError("failed to proto unmarshal", err.Error())
		return
	}

	r.p.objects[id] = &newObjects

	result, diags := schemav1.CopyObjectsToTerraform(ctx, &newObjects, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r objectsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
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

	objects := r.p.objects[id.ValueString()]

	result, diags := schemav1.CopyObjectsToTerraform(ctx, objects, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r objectsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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

	// We proto marshall and unmarshall the resource to test stability through proto round-trip.
	// This is required because protouf does not preserve the distinction between nil and empty
	// for some types (lists for examples).
	wireMsg, err := proto.Marshal(objects)
	if err != nil {
		resp.Diagnostics.AddError("failed to proto marshal", err.Error())
		return
	}

	var newObjects extypes.Objects
	if err := proto.Unmarshal(wireMsg, &newObjects); err != nil {
		resp.Diagnostics.AddError("failed to proto unmarshal", err.Error())
		return
	}

	r.p.objects[newObjects.Id] = &newObjects

	result, diags := schemav1.CopyObjectsToTerraform(ctx, &newObjects, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}

func (r objectsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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

	delete(r.p.objects, id.ValueString())
}
