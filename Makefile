pwd = $(realpath $(dir $(CURDIR)/$(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))

package_version = $(shell cat $(pwd)/VERSION)
build_dir = build
outfile = $(pwd)/$(build_dir)/protoc-gen-terraform

gopath = $(shell go env GOPATH)
srcpath = $(gopath)/src
teleport_url = github.com/gravitational/teleport
teleport_repo = https://$(teleport_url)
teleport_dir = $(srcpath)/$(teleport_url)
out_dir := "$(pwd)/out"

.PHONY: clean
clean:
	@mkdir -p ./$(build_dir)
	@rm -rf ./$(build_dir)/*
	go clean

.PHONY: build
build: clean
	go build -o $(outfile) -ldflags "-X main.Sha=`git rev-parse HEAD` -X main.Version=$(package_version)"

.PHONY: install
install: build
	go install -ldflags "-X main.Sha=`git rev-parse HEAD` -X main.Version=$(package_version)"

.PHONY: teleport
teleport: build
ifeq ("$(wildcard $(teleport_dir))", "")
	$(warning Teleport source code is required to build this example!)
	$(warning git clone ${teleport_repo} ${teleport_dir} to proceed)
	$(error Teleport source code is required to build this example)
endif
	@mkdir -p ./_out
	@protoc \
		-I$(teleport_dir)/api/types \
		-I$(teleport_dir)/vendor/github.com/gogo/protobuf \
		-I$(srcpath) \
		--plugin=$(outfile) \
		--terraform_out=config=example/teleport.yaml:${out_dir} \
		types.proto

.PHONY: test
test: build
	@protoc \
		-I$(pwd) \
		-I$(shell go list -m -f {{.Dir}} github.com/gogo/protobuf) \
		-I$(srcpath) \
		--gogo_out=test \
		test.proto

	@protoc \
		-I$(pwd) \
		-I$(pwd)/test \
		-I./vendor/github.com/gogo/protobuf \
		-I$(srcpath) \
		--plugin=$(outfile) \
		--terraform_out=config=test/config.yaml:test \
		test.proto

	@go test ./...
