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

This will generate `types_terraform.go` in _out directory. This file will contain `GetUserV2FromResourceData` and `GetRoleV3FromResourceData` along with `SchemaUserV2` and `SchemaRolesV3` methods.

Schema method should have the following prototype:

```
func SchemawrappersLabelValues() *schema.Schema
```

Getter method should look like this:
```
func GetwrappersLabelValuesFromResourceData(string path, data *schema.ResourceData, target *types.Traits) error
```

Where types.Traits, and wrappers.LabelValues is the type name in proto file.

See [Makefile](Makefile) for details.

Options:

* `types` - the list of top level types to export (with namespace).
* `exclude_fields` - list of a fields to exclude from export including type name (with namespace, ex: 'types.UserV2.Name`).
* `pkg` - default package name to prepend to type names with no package reference. This option is required if the target package of Terraform generated code is different from package of original protobuf generated code.
* `target_pkg` - the name of the target package
* `custom_duration` - the name of custom Duration type, if used.
* `custom_imports` - comma-separated package list to add into target file

# Testing

Run:

```make test```

# Note on maps of messages

Terraform does not support map of resources. If a field in protoc object is map of messages, it could not be defined in Terraform. This case is emulated via generating list with `key` and `value` fields instead. See `NestedMObj` field of [`Test`](test/custom_types.go) message for example.

Map of arrays of elementary types are not supported as well.

# Note on gogoproto.customtype

If a field has `gogoproto.casttype` flag, it can not be automatically unmarshalled from `ResourceData`. You need to define your own custom `Schema<type>` and `Get<type>FromResourceData` methods. See [test/custom_types.go](test/custom_types.go).

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
- [ ] Extract comments from original protoc file
- [x] Add argument to provide custom duration type
- [x] Add argument to provide custom imports for target file
- [ ] Add argument which will represent specific []byte fields as byte lists on Terraform side
- [x] Manually replace target package name
- [x] Run goimports to remove unused packages