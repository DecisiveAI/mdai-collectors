# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Update this version to match new release tag and run helm targets
VERSION = 0.1.6

.PHONY: all
all: build

.PHONY: build
build:
	builder --config=config/observer-collector/observer-collector-builder.yaml

.PHONY: docker-build
docker-build:
	docker build -t  public.ecr.aws/decisiveai/observer-collector:${VERSION} .

.PHONY: build-mdai-collector
build-mdai-collector:
	builder --config=config/mdai-collector/mdai-collector-builder.yaml

.PHONY: docker-build-mdai-collector
docker-build-mdai-collector:
	docker build -t public.ecr.aws/decisiveai/mdai-collector:${VERSION} -f mdai-collector.Dockerfile .