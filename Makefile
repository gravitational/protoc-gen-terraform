include version.mk

.PHONY: clean

clean:
	@mkdir -p ./_build
	rm -rf ./_build/*
	go clean

.PHONY: build
build: clean
	go build -o _build/protoc-gen-terraform -ldflags "-X main.Sha=`git rev-parse HEAD` -X main.Version=$(package_version)"

.PHONY: install
install: build
	go install -ldflags "-X main.Sha=`git rev-parse HEAD` -X main.Version=$(package_version)"

pwd = $(realpath $(dir $(CURDIR)/$(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST))))

gopath = $(shell go env GOPATH)
srcpath = $(gopath)/src
teleport_url = github.com/gravitational/teleport
teleport_repo = https://$(teleport_url)
teleport_dir = $(srcpath)/$(teleport_url)
out_dir := "$(pwd)/_out"

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
		--plugin=./_build/protoc-gen-terraform \
		--terraform_out=config=example/teleport.yaml:${out_dir} \
		types.proto

.PHONY: test
test: build
	@protoc \
		-I$(pwd)/test \
		-I$(pwd) \
		-I./vendor/github.com/gogo/protobuf \
		-I$(srcpath) \
		--plugin=./_build/protoc-gen-terraform \
		--terraform_out=target_pkg=test,types=Test,sort=true,custom_duration=Duration:test \
		--gogo_out=test \
		test.proto
	@go test -v ./test