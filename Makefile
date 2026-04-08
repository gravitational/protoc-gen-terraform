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

	@protoc \
		-I$(PWD) \
		-I$(EXAMPLES_PROTO_DIR) \
		-I$(PWD)/vendor/github.com/gogo/protobuf \
		--plugin=$(BINFILE) \
		--terraform_out=config=examples/config/computed.yaml:./examples/tfschema \
		computed.proto

	mv ./examples/tfschema/primitives_terraform.go ./examples/tfschema/primitives/v1/
	mv ./examples/tfschema/time_terraform.go ./examples/tfschema/time/v1/
	mv ./examples/tfschema/objects_terraform.go ./examples/tfschema/objects/v1/
	mv ./examples/tfschema/custom_terraform.go ./examples/tfschema/custom/v1/
	mv ./examples/tfschema/computed_terraform.go ./examples/tfschema/computed/v1/
