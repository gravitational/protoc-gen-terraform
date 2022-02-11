# protoc-gen-terraform

protoc plugin to generate Terraform Framework schema definitions and getter/setter methods from gogo/protobuf .proto files.

# Installation

Install the generator binary.

```
go install github.com/gravitational/protoc-gen-terraform
```

# Usage

```
mkdir -p ./tfschema
protoc \
    -I$(teleport_dir)/api/types \
    -I$(teleport_dir)/vendor/github.com/gogo/protobuf \
    -I$(srcpath) \
    --plugin=./_build/protoc-gen-terraform \
    --terraform_out=types=types.UserV2+types.RoleV3,pkg=types:tfschema \
    types.proto
```

This command will generate `types_terraform.go` in `tfschema` folder. 

See [Makefile](Makefile) for details.

# Options

Options could be set either using cli args or [YAML](example/teleport.yaml). Path to config file can be specified with `config` argument. Be advised that some options could be set in config file only.

## Setting target and default package name

By default, generated code is assumed to reside in the same package as your go generated code.

Use `target_package_name` option to change the target package name:

```
target_package_name=tfschema
```

Please also specify the full name of the go package where your generated code is located:

```
default_package_name="github.com/gravitational/teleport/api/types"
```

If you use types from external packages, please specify them in `external_imports`:

```
external_imports="github.com/gravitational/teleport/api/wrappers"
```

and reference that types anywhere in config using full package name (`github.com/gravitational/teleport/api/wrappers.Wrapper`). Generator will handle the rest for you.

## Specifying types to export

List message names you want to export in `types` option:

```
types=UserV2+RoleV3
```

## Excluding fields

Let's consider we have the following proto definition:

```proto
message Metadata {
	string ID = 1;
}

message User {
	Metadata Metadata = 1;
}

message AuthPreference {
	Metadata Metadata = 1;
}
```

Specify `exclude_fields` option:

```
exclude_fields=Metadata.ID+AuthPreference.Metadata.Name
```

In this case, `Metadata.ID` would be omitted for both `User` and `AuthPreference`, and `Metadata.Name` would be omitted for `AuthPreference` only. `User.Metadata.Name` won't be affected.

## Terraform Schema flags

You can specify `Required: true` (`required_fields`), `Computed: true` (`computed_fields`) and `Sensitive: true` (`sensitive_fields`) flags for your Terraform schema:

```
required_fields=Metadata.Name
```

You also can set list of `Validators` and `PlanModifiers` using configuration file:

```yaml
validators:
	"Metadata.Expires":
		- rfc3339TimeValidator

plan_modifiers:
	"Role.Options"
		- "github.com/hashicorp/terraform-plugin-framework/tfsdk.UseStateForUnknown()"
```

## Schema field naming

Schema field names are extracted from `json` tag by default. If a `json` tag is missing, a snake case of a field name is used.

If you need to rename field in schema, use `name_overrides` option:

```yaml
name_overrides:
	"Role.Spec.AWSRoleARNs": aws_arns 
```

## Custom fields

If your proto generated objects use type alias for duration fields, you can set `custom_duration` to the name of a custom duration type.

`time_type`, `duration_type` and `schema_types` options are used to override Terraform types.

```yaml
time_type:
    type: "TimeType"
    value_type: "Time"
    cast_type: "time.Time"
```

`type` is a go type name for a field's `attr.Type`, `value` is a go type for a field's `attr.Value`, `cast_type` is a go type to cast `{value_type}.Value` from.

## Generated methods

`Copy*ToTerraform` and `Copy*FromTerraform` methods are generated for every .proto message. They convert Terraform state object to go proto type and vice versa using normal go assignment operations (no reflect).

### CopyFrom

The signatures for `Test` resource would be the following:

```go
// CopyTestFromTerraform copies Terraform object fields to obj
// tf must have all the object attrs present (including null and unknown). 
// Hence, tf must be the result of req.Plan.Get or similar Terraform method.
// Otherwise, error would be returned.
func CopyTestFromTerraform(tf types.Object, obj *Test) diag.Diagnostics
```

They can be used as following:

```go
// Create template resource create method
func (r resource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan types.Object
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	obj := types.Object{}
	diags := tfschema.CopyObjFromTerraform(plan, &obj)	
	resp.Diagnostics.Append(diags...)	
	if resp.Diagnostics.HasError() {
		return
	}
}
```

Source object can contain unknown values. Unknown values are treated as nulls: the corresponding target object field is set to zero value/nil.

Source object must have all target values referenced in `Attrs`, even if they are null or unknown.

### CopyTo

Copies object to Terraform object.

```go
func CopyTestToTerraform(obj Test, tf *types.Object, updateOnly bool) error
```

Target Terraform object must have AttrTypes for all fields of Object.

## Note on gogoproto.customtype

If a field has `gogoproto.customtype` flag, schema and converters for this field can not be generated automatically. You need to define `Gen<type>Schema`, `Copy<type>FromTerraform`, `Copy<type>ToTerraform` methods.

`suffixes` option can be used to control method names:

```yaml
suffixes:
    "github.com/gravitational/teleport/api/types/wrappers.Traits": "Traits"
```

In the example above, `GenTraitsSchema` method will be called. Without this option, method name would be `GenGithubComGravitationalTeleportApiTypesWrappersTraits`.

# Testing

Run:

```make test```

# Build and test using Docker

```sh
cd build.assets
make build test
```

On Mac M1 use:
```sh
cd build.assets
make build test PROTOC_PLATFORM=linux-aarch_64
```

# TODO

- [ ] Make time format customizable
- [ ] Ability to overwrite list and maps base types
- [ ] Check duplicate/unknown imports