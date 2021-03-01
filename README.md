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

This will generate `types_terraform.go` in _out directory. This file will contain `UnmarshalUserV2` and `UnmarshalRolesV3` along with `SchemaUserV2` and `SchemaRolesV3` methods. Target package name should be changed manually. It also might contain some unused imports.

Schema method should have the following prototype:

```
func SchemawrappersLabelValues() *schema.Schema
```

Unmarshalling method should look like this:
```
func SchemawrappersLabelValues(string path, data *schema.ResourceData, target *types.Traits) error
```

Where types.Traits changes according to populated field type.

See [Makefile](Makefile) for details.

Options:

* `types` - the list of top level types to export (with namespace).
* `exclude_fields` - list of a fields to exclude from export including type name (with namespace, ex: 'types.UserV2.Name`).
* `pkg` - default package name to prepend to type names with no package reference. This option is required if the target package of Terraform generated code is different from package of original protobuf generated code.

# Testing

Run:

```make test```

# Note on maps of messages

Terraform does not support map of resources. If a field in protoc object is map of messages, it could not be defined in Terraform. This case is emulated via generating list with `key` and `value` fields instead. See `NestedMObj` field of [`Test`](test/custom_types.go) message for example.

Map of arrays of elementary types are not supported as well.

# Note on gogoproto.customtype

If a field has `gogoproto.casttype` flag, it can not be automatically unmarshalled from `ResourceData`. You need to define your own custom `Schema<type>` and `Unmarshal<type>` methods. See [test/custom_types.go](test/custom_types.go).

# TODO

- [ ] Oneof is not supported yet
- [ ] Extract comments from original protoc file
- [x] Add argument to provide custom duration type
- [x] Add argument to provide custom imports for target file
- [ ] Add argument which will represent specific []byte fields as byte lists on Terraform side
- [ ] Manually replace target package name
- [ ] Run goimports to remove unused packages