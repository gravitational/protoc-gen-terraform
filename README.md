# protoc-gen-terraform

Generates Terraform schema and unmarshalling methods from gogo/protobuf .proto file.

# Installation

Install the generator package. 

```
go get github.com/gravitational/protoc-gen-terraform
```

# Usage

```
mkdir -p ./_out
protoc \
    -I$(teleport_dir)/api/types \
    -I$(teleport_dir)/vendor/github.com/gogo/protobuf \
    -I$(srcpath) \
    --plugin=./_build/protoc-gen-terraform \
    --terraform_out=types=types.UserV2+types.RoleV3,pkg=types:_out \
    types.proto
```

This will generate `types_terraform.go` in _out directory. This file will contain `SchemaUserV2`, `SchemaRolesV3` along with `SchemaMetaUserV2`, `SchemaMetaRolesV3`.

See [Makefile](Makefile) for details.

Options:

* `types` - the list of top level types to export (with namespace).
* `exclude_fields` - list of a fields to exclude from export including type name (with namespace, ex: 'types.UserV2.Name`). Fields could also be addressed via full path in a structure (ex: `types.UserV2.Spec.Allow.Logins`). This distinction is needed if a field needs to have a different behaviour in different parent structs even if it belongs to a same type. For example, `Metadata.Expires` must be excluded for `User` and is required for `Token`.
* `pkg` - default package name to prepend to type names with no package reference. This option is required if the target package of Terraform generated code is different from the package of original protobuf generated code.
* `target_pkg` - the name of the target package
* `custom_duration` - the name of custom Duration type, if used
* `custom_imports` - comma-separated package list to add into generated file
* `required` - list of a generated Terraform schema fields to mark as `Required: true`
* `computed` - list of a generated Terraform schema fields to mark as `Computed: true`
* `force_new` - list of a generated Terraform schema fields to mark as `ForceNew: true`
* `config_mode_attr` - list of a generated Terraform schema fields to mark as `SchemaConfigMode: schema.SchemaConfigModeAttr`
* `config_mode_block` - list of a generated Terraform schema fields to mark as `SchemaConfigMode: schema.SchemaConfigModeBlock`
* `suffixes` - map of overrides of method names generated for `gogo.customtype` fields (available in config file only)
* `defaults` - default values for a fields in generated Terraform schema (available in config file only). Note that default value type in YAML file is taken into account.

All config variables could be set in [YAML](example/teleport.yaml). Path to config file can be specified using `config` variable. 

# Usage

See [get](example/get_test.go) and [set](example/get_test.go) tests for details.

# Testing

Run:

```make test```

# Note on maps of messages

Terraform does not support map of resources. If a field in protoc object is map of messages, it could not be defined in Terraform. This case is emulated via generating list with `key` and `value` fields instead. See `MapObject` field of [`Test`](test/custom_types.go) struct for example.

Map of arrays of elementary types are not supported as well.

# Note on gogoproto.customtype

If a field has `gogoproto.customtype` flag, it can not be automatically unmarshalled from `ResourceData`. You need to define your own custom `Get<type>` and `Set<type>` methods. See [test/custom_types.go](test/custom_types.go) for example.

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

- [ ] Oneof is not supported yet
- [x] Extract comments from original protoc file
- [x] Add argument to provide custom duration type
- [x] Add argument to provide custom imports for target file
- [ ] Add argument which will represent specific []byte fields as byte lists on Terraform side
- [x] Manually replace target package name
- [x] Run goimports to remove unused packages
- [x] Separate config file