# protoc-gen-terraform

protoc plugin to generate Terraform Framework schema from gogo/protobuf .proto file.

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

This will generate `types_terraform.go` in `tfschema` folder. 

See [Makefile](Makefile) for details.

# Options

`protoc-gen-terraform` supports two output formats: `sdkv2` for the legacy Terraform SDK, `tfsdk` for the new [Terraform Plugin Framework](http://github.com/hashicorp/terraform-plugin-framework/tree/v0.4.1/tfsdk).

Set output format using `foramt` option: `sdkv2` or `tfsdk`.

All config variables could be set in [YAML](example/teleport.yaml). Path to config file can be specified using `config` variable. 

Schema field names are extracted from `json` tag by default. If a `json` tag is missing, a snake case of a field name is used. 

## Common options

* `types` - the list of top level types to export (with namespaces).
* `exclude_fields` - list of a fields to exclude from export including type name (with namespace, ex: `types.UserV2.Name`). Fields could also be addressed via full path in a structure (ex: `types.UserV2.Spec.Allow.Logins`). This distinction is needed if a field needs to have a different behaviour in different parent structs even if it belongs to a same type. For example, `Metadata.Expires` must be excluded for `User` and is required for `Token`.
* `default_package name` - default package name to prepend to type names with no package reference. This option is required if the target package of Terraform generated code is different from the package of original protobuf generated code.
* `target_package_name` - the name of the target package
* `custom_imports` - comma-separated package list to add into a generated file
* `required` - list of a generated Terraform schema fields to mark as `Required: true`
* `computed` - list of a generated Terraform schema fields to mark as `Computed: true`
* `name_overrides` - map of schema field name overrides (`"a_ws_arns" -> "aws_arns"`)
* `sort` - if true, sort types and fields by name.
* `suffixes` - map of overrides of method names generated for `gogo.customtype` fields (available in config file only)
* `custom_duration` - the name of a custom Duration type, if used. Fields of this type will be treated as `time.Duration` fields.

## SDK V2 output options

_Please, note that current development is focused on Terraform framework_.

* `force_new` - list of a generated Terraform schema fields to mark as `ForceNew: true`
* `defaults` - default values for a fields in generated Terraform schema (available in config file only). Note that default value type in YAML file is taken into account.
* `state_func` - state func names to set (available in config file only).

## Terraform Framework output options

* `sensitive` - list of a generated Terraform schema fields to mark as `Sensitive: true`
* `validators` - map of arrays of attribute validators set to a fields (available in config file only).
* `schema_types` - map of type overrides (available in config file only). There are two special keys: `time` and `duration` for custom time and duration types respectively.

# Testing

Run:

```make test```

# Notes on SDK V2

## Accessors package

Converts golang structures to the Terraform objects and vice versa using reflect. See [get](example/get_test.go) and [set](example/get_test.go) tests for details. By default, they use `time.RFC3339` format for time fields with microsecond truncation.

## Note on maps of messages

Terraform does not support map of resources in SDK V2. If a field in protoc object is map of messages, it could not be defined in Terraform. This case is emulated via generating list with `key` and `value` fields instead. See `MapObject` field of [`Test`](test/custom_types.go) struct for example.

Map of arrays of elementary types are not supported as well.

## Note on gogoproto.customtype

If a field has `gogoproto.customtype` flag, it can not be automatically unmarshalled from `ResourceData`. You need to define your own custom `Schema<type>`, `Get<type>` and `Set<type>` methods. See [test/custom_types.go](test/custom_types.go) for example.

# Notes on Terraform SDK

## CopyTo and CopyFrom methods

`CopyTo` and `CopyFrom` functions are generated for every tfsdk type. They convert Terraform state object to protobuf and vice versa using standard go assignment operations. By default, they use `time.RFC3339` format for time fields with microsecond truncation.

## Note on gogoproto.customtype

If a field has `gogoproto.customtype` flag, schema and converters for this field can not be generated automatically. You need to define `Gen<type>Schema` method.

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
