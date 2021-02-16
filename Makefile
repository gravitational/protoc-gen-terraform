# Package version
version = "0.0.1"

.PHONY: clean

clean:
	@mkdir -p ./_build
	rm -rf ./_build/*
	go clean

.PHONY: build
build: clean
	pkger && go build -o _build/protoc-gen-terraform -ldflags "-X main.Sha=`git rev-parse HEAD` -X main.Version=$(version)"

.PHONY: install
install: build
	go install .

gopath = $(shell go env GOPATH)
srcpath = $(gopath)/src
teleport_url = github.com/gravitational/teleport
teleport_repo = https://$(teleport_url)
teleport_dir = $(srcpath)/$(teleport_url)
out_dir = "./_out"
types = "types.UserV2+types.UserSpecV2+types.RoleV3"
excludeFields = "types.UserSpecV2.LocalAuth"

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
		--terraform_out=types=${types},excludeFields=${excludeFields}:./${out_dir} \
		types.proto