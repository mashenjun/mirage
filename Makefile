SHELL := /bin/bash
PROJECT=mirage
GOPATH ?= $(shell go env GOPATH)

# Ensure GOPATH is set before running build process.
ifeq "$(GOPATH)" ""
  $(error Please set the environment variable GOPATH before running `make`)
endif

GO                  := GO111MODULE=on go
GOBUILD             := $(GO) build $(BUILD_FLAG) -tags codes
GOTEST              := $(GO) test -v --count=1 --parallel=1 -p=1
GORUN               := $(GO) run
TEST_LDFLAGS        := ""

PACKAGE_LIST        := go list ./...| grep -vE "cmd"
PACKAGES            := $$($(PACKAGE_LIST))

# Targets
.PHONY: dev build_linux

CURDIR := $(shell pwd)
export PATH := $(CURDIR)/bin/:$(PATH)

# run starts the server with dev config
dev:
	$(GORUN) $(CURIR)cmd/main.go -cfgPath=$(CURIR)dev/config.yaml

build_linux:
	CGO=false GOOS=linux GOARCH=amd64 go build -o bin/miraged cmd/main.go
