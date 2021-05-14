//
// Copyright 2021 SkyAPM org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package reporter

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	mock "github.com/SkyAPM/go2sky/reporter/grpc/management/mock_management"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/credentials"
	commonv3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	managementv3 "skywalking.apache.org/repo/goapi/collect/management/v3"
)

const (
	sample                = 1
	traceID               = "1f2d4bf47bf711eab794acde48001122"
	parentSegmentID       = "1e7c204a7bf711eab858acde48001122"
	parentSpanID          = 0
	parentService         = "service"
	parentServiceInstance = "serviceInstance"
	parentEndpoint        = "/foo/bar"
	addressUsedAtClient   = "foo.svc:8787"

	mockService         = "service"
	mockServiceInstance = "serviceInstance"
)

var header string

func init() {
	scx := propagation.SpanContext{
		Sample:                sample,
		TraceID:               traceID,
		ParentSegmentID:       parentSegmentID,
		ParentSpanID:          parentSpanID,
		ParentService:         parentService,
		ParentServiceInstance: parentServiceInstance,
		ParentEndpoint:        parentEndpoint,
		AddressUsedAtClient:   addressUsedAtClient,
	}
	header = scx.EncodeSW8()
}

func Test_e2e(t *testing.T) {
	reporter := createGRPCReporter()
	reporter.sendCh = make(chan *agentv3.SegmentObject, 10)
	tracer, err := go2sky.NewTracer(mockService, go2sky.WithReporter(reporter), go2sky.WithInstance(mockServiceInstance))
	if err != nil {
		t.Error(err)
	}
	entrySpan, ctx, err := tracer.CreateEntrySpan(context.Background(), "/rest/api", func(key string) (string, error) {
		return header, nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan, err := tracer.CreateExitSpan(ctx, "/foo/bar", "foo.svc:8787", func(key, value string) error {
		scx := propagation.SpanContext{}
		if key == propagation.Header {
			err = scx.DecodeSW8(value)
			if err != nil {
				t.Fatal(err)
			}
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan.End()
	entrySpan.End()
	for s := range reporter.sendCh {
		reporter.Close()
		if s.TraceId != traceID {
			t.Errorf("trace id parse error")
		}
		if len(s.Spans) == 0 {
			t.Error("empty spans")
		}
		if s.Service != mockService {
			t.Error("error are not set service")
		}
		if s.ServiceInstance != mockServiceInstance {
			t.Error("error are not set service instance")
		}
	}
}

func TestGRPCReporter_Close(t *testing.T) {
	reporter := createGRPCReporter()
	reporter.sendCh = make(chan *agentv3.SegmentObject, 1)
	tracer, err := go2sky.NewTracer(mockService, go2sky.WithReporter(reporter), go2sky.WithInstance(mockServiceInstance))
	if err != nil {
		t.Error(err)
	}
	entry, _, err := tracer.CreateEntrySpan(context.Background(), "/close", func(key string) (s string, err error) {
		return header, nil
	})
	if err != nil {
		t.Error(err)
	}
	reporter.Close()
	entry.End()
	time.Sleep(time.Second)
}

func TestGRPCReporterOption(t *testing.T) {
	// props
	instanceProps := make(map[string]string)
	instanceProps["org"] = "SkyAPM"

	// log
	logger := log.New(os.Stderr, "WithLogger", log.LstdFlags)

	// tls
	creds, err := credentials.NewClientTLSFromFile("../test/test-data/certs/cert.crt", "SkyAPM.org")
	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name       string
		option     GRPCReporterOption
		verifyFunc func(t *testing.T, reporter *gRPCReporter)
	}{
		{
			name:   "with check interval",
			option: WithCheckInterval(time.Second),
			verifyFunc: func(t *testing.T, reporter *gRPCReporter) {
				if reporter.checkInterval != time.Second {
					t.Error("error are not set checkInterval")
				}
			},
		},
		{
			name:   "with max send queue size",
			option: WithMaxSendQueueSize(50000),
			verifyFunc: func(t *testing.T, reporter *gRPCReporter) {
				if cap(reporter.sendCh) != 50000 {
					t.Error("error are not set WithMaxSendQueueSize")
				}
			},
		},
		{
			name:   "with serviceInstance props",
			option: WithInstanceProps(instanceProps),
			verifyFunc: func(t *testing.T, reporter *gRPCReporter) {
				var value string
				var ok bool
				if value, ok = reporter.instanceProps["org"]; !ok {
					t.Error("error are not set instanceProps")
				}
				if value != "SkyAPM" {
					t.Error("error are not set instanceProps")
				}
			},
		},
		{
			name:   "with logger",
			option: WithLogger(logger),
			verifyFunc: func(t *testing.T, reporter *gRPCReporter) {
				if reporter.logger != logger {
					t.Error("error are not set logger")
				}
			},
		},
		{
			name:   "with auth",
			option: WithAuthentication("test"),
			verifyFunc: func(t *testing.T, reporter *gRPCReporter) {
				if reporter.md.Get(authKey)[0] != "test" {
					t.Error("error are not set Authentication")
				}
			},
		},
		{
			name:   "with tls",
			option: WithTransportCredentials(creds),
			verifyFunc: func(t *testing.T, reporter *gRPCReporter) {
				if reporter.creds != creds {
					t.Error("error are not set TransportCredentials")
				}
			},
		},
		{
			name:   "with cds",
			option: WithCDS(10),
			verifyFunc: func(t *testing.T, reporter *gRPCReporter) {
				if reporter.cdsInterval != 10 {
					t.Error("error cds interval")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reporter := createGRPCReporter()
			tt.option(reporter)
			tt.verifyFunc(t, reporter)
		})
	}
}

func TestGRPCReporter_reportInstanceProperties(t *testing.T) {
	customProps := make(map[string]string)
	customProps["org"] = "SkyAPM"
	osProps := buildOSInfo()
	for k, v := range customProps {
		osProps = append(osProps, &commonv3.KeyStringValuePair{
			Key:   k,
			Value: v,
		})
	}
	instanceProperties := &managementv3.InstanceProperties{
		Service:         mockService,
		ServiceInstance: mockServiceInstance,
		Properties:      osProps,
	}

	ctrl := gomock.NewController(t)
	mockManagementServiceClient := mock.NewMockManagementServiceClient(ctrl)
	mockManagementServiceClient.EXPECT().ReportInstanceProperties(gomock.Any(), instancePropertiesMatcher{instanceProperties}).Return(nil, nil)

	reporter := createGRPCReporter()
	reporter.service = mockService
	reporter.serviceInstance = mockServiceInstance
	reporter.instanceProps = customProps
	reporter.managementClient = mockManagementServiceClient
	err := reporter.reportInstanceProperties()
	if err != nil {
		t.Error()
	}
}

func createGRPCReporter() *gRPCReporter {
	reporter := &gRPCReporter{
		logger: log.New(os.Stderr, "go2sky", log.LstdFlags),
	}
	return reporter
}

type instancePropertiesMatcher struct {
	x *managementv3.InstanceProperties
}

func (e instancePropertiesMatcher) Matches(x interface{}) bool {
	var props *managementv3.InstanceProperties
	var ok bool
	if props, ok = x.(*managementv3.InstanceProperties); !ok {
		return ok
	}
	if props.Service != e.x.Service {
		return false
	}
	if props.ServiceInstance != e.x.ServiceInstance {
		return false
	}
	if len(props.Properties) != len(e.x.Properties) {
		return false
	}
	return true
}

func (e instancePropertiesMatcher) String() string {
	return fmt.Sprintf("is equal to %v", e.x)
}
