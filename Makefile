#
# Copyright 2022 SkyAPM org
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

export GO111MODULE=on
export GO2SKY_GO := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
GRPC_PATH := $(GO2SKY_GO)/reporter/grpc

.DEFAULT_GOAL := test

.PHONY: deps
deps:
	go get -v -t -d ./...

.PHONY: test
test:
	go test -v -race -cover -coverprofile=coverage.txt -covermode=atomic `go list ./... | grep -v github.com/SkyAPM/go2sky/reporter/grpc | grep -v github.com/SkyAPM/go2sky/test`

.PHONY: mock-gen
mock-gen:
	cd $(GRPC_PATH)/language-agent && \
    	  mkdir -p mock_trace && \
    	  mockgen skywalking.apache.org/repo/goapi/collect/language/agent/v3 TraceSegmentReportServiceClient > mock_trace/Tracing.mock.go
	cd $(GRPC_PATH)/management && \
    	  mkdir -p mock_management && \
    	  mockgen skywalking.apache.org/repo/goapi/collect/management/v3 ManagementServiceClient > mock_management/Management.mock.go

LINTER := bin/golangci-lint
$(LINTER):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.20.1

.PHONY: lint
lint: $(LINTER) ./golangci.yml  ## Run the linters
	@echo "linting..."
	$(LINTER) run --config ./golangci.yml

.PHONY: fix
fix: $(LINTER)
	@echo "fix..."
	$(LINTER) run -v --fix ./...

.PHONY: all
all: test lint
