#
# Licensed to the SkyAPM org under one or more
# contributor license agreements.  See the NOTICE file distributed with
# this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
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
	go test -v -race -cover `go list ./... | grep -v github.com/powerapm/go2sky/reporter/grpc`

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
	  mockgen github.com/powerapm/go2sky/reporter/grpc/register RegisterClient > mock_register/Register.mock.go && \
	  mockgen github.com/powerapm/go2sky/reporter/grpc/register ServiceInstancePingClient > mock_register/InstancePing.mock.go
	cd $(GRPC_PATH)/language-agent-v2 && \
    	  mkdir -p mock_trace && \
    	  mockgen github.com/powerapm/go2sky/reporter/grpc/language-agent-v2 TraceSegmentReportServiceClient > mock_trace/trace.mock.go

LINTER := bin/golangci-lint
$(LINTER):
	wget -q -O- https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.20.1

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