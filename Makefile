version = 1

.PHONY: clean

clean:
	@mkdir -p ./_build
	rm -rf ./_build/*
	go clean

.PHONY: build
build: clean
	go build -o _build/protoc-gen-terraform -X "main.Sha `git rev-parse HEAD` -X main.Version=$(version)" myapp.go

.PHONY: install
install: build
	go install .

gopath = $(shell go env GOPATH)
srcpath = $(gopath)/src
teleport_url = github.com/gravitational/teleport
teleport_repo = https://$(teleport_url)
teleport_dir = $(srcpath)/$(teleport_url)
out_dir = "./_out"

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
		-I$(teleport_dir) \
		-I$(srcpath) \
		--plugin=./_build/protoc-gen-terraform \
		--terraform_out=${out_dir} \
		types.proto
