# Examples

The `examples/` directory contains a set of sample protobuf schemas, protoc
configuration, generated Terraform code, and test fixtures used to test against
realistic inputs.

## Layout

- `config/`: protoc-gen-terraform YAML config files.
- `proto/`: source `.proto` files for each example.
- `types/`: generated protobuf Go types plus custom types used by the examples.
- `tfschema`: generated Terraform schema and conversion code.
- `testlib/fixtures/`: Terraform configurations used by the example tests.
- `testlib/provider/`: an in-memory provider used for testing.

## Generated artifacts

The following files are generated and should not be edited by hand:

- `examples/types/*.pb.go`
- `examples/tfschema/**/*_terraform.go`

The following files are handwritten support code and may need updates when an
example changes:

- `examples/types/custom_types.go`
- `examples/types/time_duration.go`
- `examples/types/plan_modifier.go`
- `examples/types/validator.go`
- `examples/tfschema/**/custom_types.go`
- `examples/testlib/provider/*.go`
- `examples/testlib/fixtures/*.tf`

## Running tests

Run the following target to regenerate artifacts:

```sh
make gen
```

Run the following target to run tests:

```sh
make test
```

## How to add or update an existing example

For a change to an existing example, the usual edit sequence is:

1. Update source schema in `examples/proto/`.
2. Update Terraform generator config in `examples/config/`.
3. Update `gen` target in `Makefile` if needed.
4. Update custom types or custom Terraform conversions if needed.
5. Update Terraform fixtures in `examples/testlib/fixtures/`.
6. Update provider in `examples/testlib/provider/` if needed.
7. Run `make gen` and `make test`.
