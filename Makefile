package_version = "0.0.1"

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
types = "types.UserV2+types.RoleV3"
exclude_fields = "types.UserSpecV2.LocalAuth"
# exclude_fields = ""
custom_duration = "Duration"
custom_imports = "github.com/gravitational/teleport/api/types"
target_pkg = "tfschema"

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
pkg=types,custom_duration=Duration,custom_imports=${custom_imports},target_pkg=${target_pkg}:${out_dir} \
		types.proto

.PHONY: test
test: build
	@protoc \
		-I$(srcpath)/github.com/gravitational/protoc-gen-terraform/test \
		-I$(srcpath)/github.com/gravitational/protoc-gen-terraform \
		-I./vendor/github.com/gogo/protobuf \
		-I$(srcpath) \
		--plugin=./_build/protoc-gen-terraform \
		--terraform_out=target_pkg=test,types=Test,custom_duration=Duration:test \
		--gogo_out=test \
		test.proto
	@go test -v ./test