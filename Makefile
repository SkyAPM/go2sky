export GO111MODULE=on
export GO2SKY_GO := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GRPC_PATH := $(GO2SKY_GO)/reporter/grpc

.DEFAULT_GOAL := test


.PHONY: test
test:
	go test -v -race -cover ./...

.PHONY: proto-gen
proto-gen:
	cd $(GRPC_PATH) && \
	  protoc common/*.proto --go_out=plugins=grpc:$(GOPATH)/src
	cd $(GRPC_PATH) && \
      protoc language-agent-v2/*.proto --go_out=plugins=grpc:$(GOPATH)/src
	cd $(GRPC_PATH) && \
      protoc register/*.proto --go_out=plugins=grpc:$(GOPATH)/src

.PHONY: mock-gen
mock-gen:
	cd $(GRPC_PATH)/register && \
	  mkdir -p mock_register && \
	  mockgen github.com/tetratelabs/go2sky/reporter/grpc/register RegisterClient > mock_register/Register.mock.go && \
	  mockgen github.com/tetratelabs/go2sky/reporter/grpc/register ServiceInstancePingClient > mock_register/InstancePing.mock.go
	cd $(GRPC_PATH)/language-agent-v2 && \
    	  mkdir -p mock_trace && \
    	  mockgen github.com/tetratelabs/go2sky/reporter/grpc/language-agent-v2 TraceSegmentReportServiceClient > mock_trace/trace.mock.go

LINTER := bin/golangci-lint
$(LINTER):
	wget -q -O- https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.13

.PHONY: lint
lint: $(LINTER) ./golangci.yml  ## Run the linters
	@echo "linting..."
	$(LINTER) run --config ./golangci.yml

.PHONY: all
all: test lint
