package_version = "0.0.2"

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
	go install .

gopath = $(shell go env GOPATH)
srcpath = $(gopath)/src
teleport_url = github.com/gravitational/teleport
teleport_repo = https://$(teleport_url)
teleport_dir = $(srcpath)/$(teleport_url)
out_dir := "./_out"
types = "types.UserV2+types.RoleV3+types.GithubConnectorV3+types.SAMLConnectorV2+types.OIDCConnectorV2+types.TrustedClusterV2+types.ProvisionTokenV2"
exclude_fields = "types.UserSpecV2.LocalAuth+types.Metadata.ID+types.UserSpecV2.Expires+types.UserSpecV2.CreatedBy+types.UserSpecV2.Status+types.UserV2.Version+types.GithubConnectorV3.Version"
computed="types.UserV2.Kind+types.UserV2.SubKind+types.Metadata.Namespace"
required="types.Metadata.Name"

custom_duration = "Duration"
custom_imports = "github.com/gravitational/teleport/api/types"
target_pkg = "tfschema"
pwd = $(shell pwd)

.PHONY: terraform
terraform: build
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
		--terraform_out=types=${types},exclude_fields=${exclude_fields},\
pkg=types,custom_duration=Duration,custom_imports=${custom_imports},\
target_pkg=${target_pkg},required=${required},computed=${computed}:${out_dir} \
		types.proto

.PHONY: test
test: build
	@protoc \
		-I$(pwd)/test \
		-I$(pwd) \
		-I./vendor/github.com/gogo/protobuf \
		-I$(srcpath) \
		--plugin=./_build/protoc-gen-terraform \
		--terraform_out=target_pkg=test,types=Test,custom_duration=Duration:test \
		--gogo_out=test \
		test.proto
	@go test -v ./test