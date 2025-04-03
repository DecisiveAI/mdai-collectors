# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

.PHONY: all
all: build

.PHONY: build
build:
	ocb --config=config/builder/watcher-builder.yaml

.PHONY: build-debug-collector
build-debug-collector:
	ocb --config=config/builder/debug-builder.yaml