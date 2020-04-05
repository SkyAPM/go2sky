// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package reporter

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	v3 "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
)

const (
	sample                = 1
	traceID               = "1f2d4bf47bf711eab794acde48001122"
	parentSegmentID       = "1e7c204a7bf711eab858acde48001122"
	parentSpanID          = 0
	parentService         = "service"
	parentServiceInstance = "instance"
	parentEndpoint        = "/foo/bar"
	addressUsedAtClient   = "foo.svc:8787"
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
	service, instance, reporter := createMockReporter()
	reporter.sendCh = make(chan *v3.SegmentObject, 10)
	tracer, err := go2sky.NewTracer(service, go2sky.WithReporter(reporter), go2sky.WithInstance(instance))
	if err != nil {
		t.Error(err)
	}
	entrySpan, ctx, err := tracer.CreateEntrySpan(context.Background(), "/rest/api", func() (string, error) {
		return header, nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan, err := tracer.CreateExitSpan(ctx, "/foo/bar", "foo.svc:8787", func(head string) error {
		scx := propagation.SpanContext{}
		err = scx.DecodeSW8(head)
		if err != nil {
			t.Fatal(err)
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
	}
}

func TestGRPCReporter_Close(t *testing.T) {
	service, instance, reporter := createMockReporter()
	reporter.sendCh = make(chan *v3.SegmentObject, 1)
	tracer, err := go2sky.NewTracer(service, go2sky.WithReporter(reporter), go2sky.WithInstance(instance))
	if err != nil {
		t.Error(err)
	}
	entry, _, err := tracer.CreateEntrySpan(context.Background(), "/close", func() (s string, err error) {
		return header, nil
	})
	if err != nil {
		t.Error(err)
	}
	reporter.Close()
	entry.End()
	time.Sleep(time.Second)
}

func createMockReporter() (string, string, *gRPCReporter) {
	reporter := &gRPCReporter{
		logger: log.New(os.Stderr, "go2sky", log.LstdFlags),
	}
	return "service", "instance", reporter
}
