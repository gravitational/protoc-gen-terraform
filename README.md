# protoc-gen-terraform

Generates Terraform schema and unmarshalling methods from gogo/protobuf .protoc file.

# Installation

Install protobuf. 

```
go get github.com/gzigzigzeo/protoc-gen-terraform
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

# Note on maps of messages

Terraform does not support map of resources. If a field in protoc object is map of messages, it could not be defined in Terraform. This case is emulated via generating list with `key` and `value` fields instead. See `NestedMObj` field of [`Test`](test/custom_types.go) message for example.

# Note on gogoproto.customtype

If a field has `gogoproto.casttype` flag, it can not be automatically unmarshalled from `ResourceData`. You need to define your own custom `Schema<type>` and `Unmarshal<type>` methods. See [test/custom_types.go](test/custom_types.go).

# TODO

1. Oneof is not supported yet
2. Extract comments from original protoc file