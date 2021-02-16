# Package version
version = "0.0.1"

.PHONY: clean

clean:
	@mkdir -p ./_build
	rm -rf ./_build/*
	go clean

.PHONY: build
build: clean
	go build -o _build/protoc-gen-terraform -ldflags "-X main.Sha=`git rev-parse HEAD` -X main.Version=$(version)"

.PHONY: install
install: build
	go install .

gopath = $(shell go env GOPATH)
srcpath = $(gopath)/src
teleport_url = github.com/gravitational/teleport
teleport_repo = https://$(teleport_url)
teleport_dir = $(srcpath)/$(teleport_url)
out_dir = "./_out"
types = "types.UserV2+types.UserSpecV2+types.RoleV3+types.RoleSpecV3+types.ProvisionTokenV2+types.ProvisionTokenSpecV2+types.Metadata+types.ExternalIdentity+types.RoleOptions+types.RoleConditions+types.BoolValue+types.ExternalIdentity+wrappers.LabelValues+types.CreatedBy+wrappers.StringValues+wrappers.ValuesEntry"
# types = "types.RoleOptions+types.BoolValue"

.PHONY: example
example: build
ifeq ("$(wildcard $(teleport_dir))", "")
	@echo "Teleport source code is required to build this example!"
	@echo "git clone ${teleport_repo} ${teleport_dir} to proceed"
endif
	@mkdir -p ./_out
	@protoc \
		-I$(teleport_dir)/api/types \
		-I$(teleport_dir)/vendor/github.com/gogo/protobuf \
		-I$(srcpath) \
		--plugin=./_build/protoc-gen-terraform \
		--terraform_out=types=${types}:./${out_dir} \
		types.proto