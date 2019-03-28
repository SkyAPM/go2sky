export GO111MODULE=on
export GO2SKY_GO := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GRPC_PATH := $(GO2SKY_GO)/reporter/grpc

.DEFAULT_GOAL := test


.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: lint
lint:
	# Ignore grep's exit code since no match returns 1.
	echo 'linting...' ; golint ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: proto-gen
proto-gen:
	cd $(GRPC_PATH) && \
	  protoc common/*.proto --go_out=plugins=grpc:$(GOPATH)/src
	cd $(GRPC_PATH) && \
      protoc language-agent-v2/*.proto --go_out=plugins=grpc:$(GOPATH)/src
	cd $(GRPC_PATH) && \
      protoc register/*.proto --go_out=plugins=grpc:$(GOPATH)/src

.PHONY: all
all: vet lint test

.PHONY: example
