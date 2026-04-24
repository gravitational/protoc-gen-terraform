# Get the current directory relative to Makefile location
PWD = $(realpath $(dir $(CURDIR)/$(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))

PACKAGE_VERSION = $(shell cat $(PWD)/VERSION)
BUILD_DIR = build
BINFILE = $(PWD)/$(BUILD_DIR)/protoc-gen-terraform

GOPATH = $(shell go env GOPATH)
SRCPATH = $(GOPATH)/src

.PHONY: clean
clean:
	@mkdir -p ./$(BUILD_DIR)
	@rm -rf ./$(BUILD_DIR)/*
	go clean

.PHONY: build
build: clean
	go build -o $(BINFILE)

.PHONY: install
install: build
	go install

.PHONY: test
test: build gen
	@protoc \
		-I$(PWD) \
		-I$(PWD)/test \
		-I$(PWD)/vendor/github.com/gogo/protobuf \
		--gogo_out=test \
		test.proto

	@protoc \
		-I$(PWD) \
		-I$(PWD)/test \
		-I$(PWD)/vendor/github.com/gogo/protobuf \
		--plugin=$(BINFILE) \
		--terraform_out=config=test/config.yaml:test \
		test.proto

# Generate Go types for optional test using standard protoc-gen-go (not gogo)
	@protoc \
		-I$(PWD)/test/optional \
		--go_out=test/optional \
		--go_opt=paths=source_relative \
		optional.proto

# Generate Terraform code for optional test
	@protoc \
		-I$(PWD)/test/optional \
		--plugin=$(BINFILE) \
		--terraform_out=config=test/optional/config.yaml:. \
		optional.proto
	mv ./github.com/gravitational/protoc-gen-terraform/v3/test/optional/optional_terraform.go ./test/optional/
	rm -rf ./github.com

	@go test ./...

# EXAMPLES_PROTO_DIR and EXAMPLES_PROTO_FILES specifies the directory and files
# containing example proto config for testing.
EXAMPLES_PROTO_DIR := examples/proto
EXAMPLES_PROTO_FILES := $(EXAMPLES_PROTO_DIR)/*proto

.PHONY: gen
gen: build
# generate protobuf go code
	@protoc \
		-I$(PWD) \
		-I$(EXAMPLES_PROTO_DIR) \
		-I$(PWD)/vendor/github.com/gogo/protobuf \
		--gogo_out=examples/types \
		$(EXAMPLES_PROTO_FILES)

# generate terraform code
	@protoc \
		-I$(PWD) \
		-I$(EXAMPLES_PROTO_DIR) \
		-I$(PWD)/vendor/github.com/gogo/protobuf \
		--plugin=$(BINFILE) \
		--terraform_out=config=examples/config/primitives.yaml:examples/tfschema \
		primitives.proto

	@protoc \
		-I$(PWD) \
		-I$(EXAMPLES_PROTO_DIR) \
		-I$(PWD)/vendor/github.com/gogo/protobuf \
		--plugin=$(BINFILE) \
		--terraform_out=config=examples/config/time.yaml:examples/tfschema \
		time.proto

	@protoc \
		-I$(PWD) \
		-I$(EXAMPLES_PROTO_DIR) \
		-I$(PWD)/vendor/github.com/gogo/protobuf \
		--plugin=$(BINFILE) \
		--terraform_out=config=examples/config/objects.yaml:examples/tfschema \
		objects.proto

	@protoc \
		-I$(PWD) \
		-I$(EXAMPLES_PROTO_DIR) \
		-I$(PWD)/vendor/github.com/gogo/protobuf \
		--plugin=$(BINFILE) \
		--terraform_out=config=examples/config/custom.yaml:./examples/tfschema \
		custom.proto

	mv ./examples/tfschema/primitives_terraform.go ./examples/tfschema/primitives/v1/
	mv ./examples/tfschema/time_terraform.go ./examples/tfschema/time/v1/
	mv ./examples/tfschema/objects_terraform.go ./examples/tfschema/objects/v1/
	mv ./examples/tfschema/custom_terraform.go ./examples/tfschema/custom/v1/
