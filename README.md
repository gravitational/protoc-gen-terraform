# protoc-gen-terraform

Generates Terraform schema and unmarshalling methods from gogo/protobuf .protoc file.

# Note on maps of messages

Terraform does not support map of resources. If a field in protoc object is map of messages, it could not be defined in Terraform. This case is emulated via generating list with `key` and `value` fields instead. See `NestedMObj` field of [`Test`](test/custom_types.go) message for example.

# Note on gogoproto.customtype

If a field has `gogoproto.casttype` flag, it can not be automatically unmarshalled from `ResourceData`. You need to define your own custom `Schema<type>` and `Unmarshal<type>` methods. See [test/custom_types.go](test/custom_types.go).

# TODO

1. Oneof is not supported yet
2. Extract comments from original protoc file