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
	go build -o $(BINFILE) -ldflags "-X main.Sha=`git rev-parse HEAD` -X main.Version=$(PACKAGE_VERSION)"

.PHONY: install
install: build
	go install -ldflags "-X main.Sha=`git rev-parse HEAD` -X main.Version=$(PACKAGE_VERSION)"

.PHONY: test
test: build
	@protoc \
		-I$(PWD) \
		-I$(PWD)/test \
		-I$(shell go list -m -f {{.Dir}} github.com/gogo/protobuf) \
		-I$(SRCPATH) \
		--gogo_out=test \
		test.proto

	@protoc \
		-I$(PWD) \
		-I$(PWD)/test \
		-I./vendor/github.com/gogo/protobuf \
		-I$(SRCPATH) \
		--plugin=$(BINFILE) \
		--terraform_out=config=test/config.yaml:test \
		test.proto

	@go test ./...
