# protoc-gen-terraform

protoc plugin to generate Terraform Framework schema definitions and getter/setter methods from gogo/protobuf .proto files.

# Installation

Install the generator binary.

```
go install github.com/gravitational/protoc-gen-terraform
```

# Usage

Given that you have `gogo/protobuf` and `gravitational/teleport/api` in your $GOSRC dir:

```
mkdir -p ./tfschema
protoc \
    -I$(go env GOPATH)/src/github.com/gravitational/teleport/api/types \
    -I$(go env GOPATH)/src/github.com/gogo/protobuf \
    -I$(go env GOPATH)/src \
    --plugin=./build/protoc-gen-terraform \
    --terraform_out=types=RoleSpecV4,pkg=types:tfschema \
    types.proto
```

This command will generate `types_terraform.go` in `tfschema` folder. 

See [Makefile](Makefile) for details.

# Options

Options can be set using either CLI args or [YAML](test/config.yaml). The path to the config file can be specified with `config` argument. Be advised that some options can only be set via the config file

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
	"Role.Options":
		- "github.com/hashicorp/terraform-plugin-framework/tfsdk.RequiresReplace()"
```

## UseStateForUnknown by default

The following setting:

```
use_state_for_unknown_by_default: true
```

will add `tfsdk.UseStateForUnknown()` PlanModifier to all computed fields.

## Injecting fields into schema

There are cases when you need to add fields not existing in the object to schema. For example, artificial id field is required for Terraform acceptance tests to work. You can achieve it using `injected_fields` option:

```yaml
injected_fields:
  Test: # Path to inject
    -
      name: id
      type: github.com/hashicorp/terraform-plugin-framework/types.StringType
      computed: true
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
    type: "TimeType"                    # attr.Type
    value_type: "TimeValue"             # attr.Value
    cast_to_type: "time.Time"           # TimeValue.Value type
	cast_from_type: "time.Time"         # Go object field type
	type_constructor: UseRFC3339Time()  # Function to put into schema definition Type, will generate TimeType{} if missing
```

## Custom duration value

If you schema uses the following definition for the duration fields:

```golang
int64 MaxSessionTTL = 2 [ (gogoproto.casttype) = "Duration" ];
```

you can set the `duration_custom_type` to make such fields act as duration custom type:

```
duration_custom_type=Duration
```

## Generated methods

`Copy*ToTerraform` and `Copy*FromTerraform` methods are generated for every .proto message. They convert Terraform state object to go proto type and vice versa using normal go assignment operations (no reflect).

### CopyFrom

Copies Terraform data to an object.

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

The following rules apply:
1. Source Terraform object must contain values for all target object fields.
2. Unknown values are treated as nulls. Target object value would be set to either nil or zero value.

So, the source Terraform object might be Plan, State or Object.

### CopyTo

Copies object to Terraform object.

```go
func CopyTestToTerraform(obj Test, tf *types.Object, updateOnly bool) error
```

Target Terraform object must have AttrTypes for all fields of Object.

The following rules apply:
1. All target attributes are marked as known.
2. In case an attribute is present in AttrTypes, but is missing in AttrValues, it is created.

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
