# protoc-gen-terraform

protoc plugin to generate Terraform Framework schema definitions and getter/setter methods from gogo/protobuf .proto files.

# Installation

Install the generator binary.

```
go install github.com/gravitational/protoc-gen-terraform/v4@v4.0.0
```

# Dependencies

- `protoc-gen-go`: `protoc-gen-gogo` does not support all proto3 fields/keywords such as `optional`; if your `.proto` files uses features that are not supported in gogo, you must use `protoc-gen-go` to generate the `.pb.go` files.

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

If package import paths are not being correctly found automatically, use the
`import_path_overrides` field in the yaml config to override the import path for
specific package names.

```yaml
import_path_overrides:
  "types": "github.com/gravitational/teleport/api/types"
```

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
    - "github.com/hashicorp/terraform-plugin-framework/resource.RequiresReplace()"
```

## UseStateForUnknown by default

The following setting:

```
use_state_for_unknown_by_default: true
```

will add `resource.UseStateForUnknown()` PlanModifier to all computed fields.

## Injecting fields into schema

There are cases when you need to add fields not existing in the object to schema. For example, artificial id field is required for Terraform acceptance tests to work. You can achieve it using `injected_fields` option:

```yaml
injected_fields:
  Test: # Path to inject
    - name: id
      type: github.com/hashicorp/terraform-plugin-framework/types.StringType
      computed: true
      default_value_method: github.com/hashicorp/terraform-plugin-framework/types.StringNull
```

## Schema field naming

Schema field names are extracted from `json` tag by default. If a `json` tag is missing, a snake case of a field name is used.

If you need to rename field in schema, use `name_overrides` option:

```yaml
name_overrides:
  "Role.Spec.AWSRoleARNs": aws_arns
```

## Custom fields

If your proto generated objects use type alias for duration fields, you can set `duration_custom_type` to the name of a custom duration type.

`time_type`, `duration_type` and `schema_types` options are used to override Terraform types.

Available configuration options include:
- `type`: Terraform attr.Type used in schema, for example types.StringType.
- `value_type`: Terraform attr.Value implementation expected when reading, for example types.String.
- `value_from_method`: method called on the Terraform value to extract the Go value, for example ValueString().
- `value_to_method`: constructor used to create a known Terraform value from Go, for example types.StringValue(...).
- `null_value_method`: constructor used for null values, for example types.StringNull().
- `unknown_value_method`: constructor used when preserving unknown values, for example types.StringUnknown().
- `cast_to_type`: Go type passed into value_to_method.
- `cast_from_type`: Go field type when reading back from Terraform.

```yaml
# time_type overrides the protobuf time fields.
time_type:
  type: "TimeType"
  value_type: "TimeValue"
  value_from_method: "ValueTime"
  value_to_method: "NewTime"
  null_value_method: "NullTime"
  unknown_value_method: "UnknownTime"
  cast_to_type: "time.Time"
  cast_from_type: "time.Time"
  type_constructor: UseRFC3339Time() # Function to put into schema definition Type, will generate TimeType{} if missing

# duration_type overrides the protobuf duration fields.
duration_type:
  type: "DurationType"
  value_type: "DurationValue"
  value_from_method: "ValueDuration"
  value_to_method: "NewDuration"
  null_value_method: "NullDuration"
  unknown_value_method: "UnknownDuration"
  cast_to_type: "time.Duration"
  cast_from_type: "time.Duration"

# schema_types lets you override the Terraform type used for a specific proto field.
schema_types:
  "Custom.schema_override":
    type: "github.com/hashicorp/terraform-plugin-framework/types.StringType"
    value_type: "github.com/hashicorp/terraform-plugin-framework/types.String"
    value_from_method: "ValueString"
    value_to_method: "github.com/hashicorp/terraform-plugin-framework/types.StringValue"
    null_value_method: "github.com/hashicorp/terraform-plugin-framework/types.StringNull"
    unknown_value_method: "github.com/hashicorp/terraform-plugin-framework/types.StringUnknown"
    cast_to_type: "string"
    cast_from_type: "OverrideCastType"
```

If every time or duration field is overridden with `schema_types`, you can use
placeholder definitions for `time_type` and `duration_type`. The generator only
needs these global entries to recognize time and duration fields; the
field-specific `schema_types` entries provide the real Terraform type and copy
methods.

```yaml
# This must be defined for the generator to be happy, but in reality all time
# fields are overridden because protobuf timestamps contain locks and the
# linter gets mad if raw structs are used instead of pointers.
time_type:
  type: "PlaceholderType"
duration_type:
  type: "PlaceholderType"
```

## Custom types

`custom_types` lets you replace the generator's normal schema and copy logic
for a specific field with handwritten functions.

```yaml
custom_types:
  "Custom.string_override": StringCustom
```

The map key is the field path. The value is the suffix used to find the custom
functions. For the example above, the generator expects these functions to be
available:

```go
func GenSchemaStringCustom(ctx context.Context, attr tfsdk.Attribute) tfsdk.Attribute
func CopyFromStringCustom(diags diag.Diagnostics, tf attr.Value, obj *string)
func CopyToStringCustom(diags diag.Diagnostics, obj string, t attr.Type, v attr.Value, preserveUnknown bool) attr.Value
```

Use `custom_types` when the Terraform representation is shaped differently from
the Go field. For example, a Go `string` field can be represented in Terraform
as a list of strings, with custom copy functions joining and splitting the
value.

Use `schema_types` when the generator can still perform the conversion and only
needs different Terraform type/value constructors. Use `custom_types` when you
need to own the schema and both copy directions for the field.

## Custom duration value

If your schema uses the following definition for the duration fields:

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
func CopyTestFromTerraform(ctx context.Context, tf types.Object, obj *Test) diag.Diagnostics
```

They can be used as following:

```go
// Create template resource create method
func (r resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan types.Object
    diags := req.Plan.Get(ctx, &plan)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    obj := &Test{}
    diags = tfschema.CopyTestFromTerraform(ctx, plan, obj)
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

Copies object to Terraform object returning a new Terraform object.

```go
func CopyTestToTerraform(ctx context.Context, obj *Test, tf *types.Object) (types.Object, diag.Diagnostics)
```

Target Terraform object is used as the source of existing values and attribute
types. The returned object contains the updated attributes.

The following rules apply:

1. The source Terraform object must have attribute types for all generated fields.
2. Missing attribute values are created in the returned object.
3. Unknown values are replaced by known values by default.
4. Use `Copy*ToTerraformPreserveUnknown(..., true)` to preserve unknown values already present in the source object.

## Note on gogoproto.customtype

If a field has `gogoproto.customtype` flag, schema and converters for this field can not be generated automatically. You need to define `Gen<type>Schema`, `Copy<type>FromTerraform`, `Copy<type>ToTerraform` methods.

`suffixes` option can be used to control method names:

```yaml
suffixes:
  "github.com/gravitational/teleport/api/types/wrappers.Traits": "Traits"
```

In the example above, `GenTraitsSchema` method will be called. Without this option, method name would be `GenGithubComGravitationalTeleportApiTypesWrappersTraits`.

# Note on empty messages

Protobuf allows to define messages with no fields. Terraform treats such objects as errors. If a message has no fields, generator defines an artificial `active` field in the schema. It will always be null.

# Testing

Run:

`make test`

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

# Printing version

`protoc-gen-terraform version`

will print version number and quit.

# Releasing the new version

Current version number resides in the VERSION file. To update the file contents from current git tag, run:

```
git tag v1.1.1
go generate
```

# Debugging

The plugin can be configured to dump the `protoc` generation request by setting the `PROTOC_GEN_TERRAFORM_DUMP` env var.
The dumped request can then be replayed by sending it to the plugin's stdin.
For example:

```
export PROTOC_GEN_TERRAFORM_DUMP="$(mktemp)"
echo "$PROTOC_GEN_TERRAFORM_DUMP"

# note that make gen calls the plugin several times, only the last request is logged, which is usually the one that fails.
make gen

# invoke the plugin again, but replay the saved request
./build/protoc-gen-terraform < "$PROTOC_GEN_TERRAFORM_DUMP"
```

Debuggers can be configured to read stdin from the file:

- in GoLand: on a `Go Build` target, set `Redirect input from` to the dump file path
- in VSCode: in the `Launch.json` file, set `stdinFrom` to the dump file path
- with `dlv`: `dlv exec ./build/protoc-gen-terraform --listen=:2345 --headless=true --api-version=2 --accept-multiclient --redirect stdin:"$PROTOC_GEN_TERRAFORM_DUMP"
