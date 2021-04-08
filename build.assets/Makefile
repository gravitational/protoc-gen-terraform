LOCALDIR := $(dir $(CURDIR)/$(word $(words $(MAKEFILE_LIST)),$(MAKEFILE_LIST)))
TOP ?= $(abspath $(LOCALDIR)/..)
SRCDIR=/go/src/github.com/gravitational/protoc-gen-terraform
TELEPORT_SRCDIR=/go/src/github.com/gravitational/teleport

package_version = 0.0.3

HOSTNAME=buildbox
BUILDBOX ?= proto-gen-buildbox:$(package_version)
GOCACHE ?= $(HOME)/.cache/go-build
DOCKERFLAGS := --rm \
	-v "$(TOP)":$(SRCDIR) \
	-v "$(TOP)/../teleport":$(TELEPORT_SRCDIR) \
	-v /tmp:/tmp \
	-v $(GOCACHE):/go/cache \
	-e GOCACHE=/go/cache \
	-w $(SRCDIR) \
	-h $(HOSTNAME)
DOCKERPULLFLAGS ?= --pull
GO_VER = 1.16

PROTOC_VER ?= 3.15.3
PROTOC_PLATFORM := linux-x86_64 # NOTE: use linux-aarch_64 for Mac M1
GOGO_PROTO_TAG ?= v1.3.2

NOROOT=-u $$(id -u):$$(id -g)

.PHONY: build
build: buildbox
	docker run $(DOCKERFLAGS) $(NOROOT) $(BUILDBOX) \
		dumb-init make -C $(SRCDIR) build

.PHONY: terraform
terraform: build
	docker run $(DOCKERFLAGS) $(NOROOT) $(BUILDBOX) \
		dumb-init make -C $(SRCDIR) terraform

.PHONY: test
test: build
	docker run $(DOCKERFLAGS) $(NOROOT) $(BUILDBOX) \
		dumb-init make -C $(SRCDIR) test

.PHONY: buildbox
buildbox:
	docker build \
		--build-arg GO_VER=$(GO_VER) \
		--build-arg PROTOC_VER=$(PROTOC_VER) \
		--build-arg PROTOC_PLATFORM=$(PROTOC_PLATFORM) \
		--build-arg GOGO_PROTO_TAG=$(GOGO_PROTO_TAG) \
		--build-arg UID=$(shell id -u) \
		--build-arg GID=$(shell id -g) \
		$(DOCKERPULLFLAGS) --tag $(BUILDBOX) .