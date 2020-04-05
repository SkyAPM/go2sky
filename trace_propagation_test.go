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

package go2sky

import (
	"context"
	"sync"
	"testing"

	"github.com/SkyAPM/go2sky/propagation"
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

func TestTracer_EntryAndExit(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	tracer, err := NewTracer("service", WithReporter(&NoopReporter{wg: wg}))
	if err != nil {
		t.Error(err)
	}
	entrySpan, ctx, err := tracer.CreateEntrySpan(context.Background(), "/rest/api", func() (string, error) {
		return "", nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan, err := tracer.CreateExitSpan(ctx, "/foo/bar", "foo.svc:8787", func(head string) error {
		scx := propagation.SpanContext{}
		err = scx.DecodeSW8(head)
		if err != nil {
			t.Fail()
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan.End()
	entrySpan.End()
	wg.Wait()
}

func TestTracer_Entry(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	reporter := &NoopReporter{wg: wg}
	tracer, err := NewTracer("service", WithReporter(reporter))
	if err != nil {
		t.Error(err)
	}
	entrySpan, _, err := tracer.CreateEntrySpan(context.Background(), "/rest/api", func() (string, error) {
		return header, nil
	})
	if err != nil {
		t.Error(err)
	}
	entrySpan.End()
	wg.Wait()
	span := reporter.Spans[0]
	if span.Context().TraceID != traceID {
		t.Fail()
	}
	if len(span.Refs()) != 1 {
		t.Fail()
	}
}

func TestTracer_EntryAndExitInTrace(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	tracer, err := NewTracer("service", WithInstance("instance"), WithReporter(&NoopReporter{wg: wg}))
	if err != nil {
		t.Error(err)
	}
	entrySpan, ctx, err := tracer.CreateEntrySpan(context.Background(), "/rest/api", func() (string, error) {
		return header, nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan, err := tracer.CreateExitSpan(ctx, "/foo/bar", "foo.svc:8786", func(head string) error {
		sc := propagation.SpanContext{}
		err = sc.DecodeSW8(head)
		if err != nil {
			t.Fail()
		}

		if sc.Sample != sample {
			t.Fail()
		}

		if sc.TraceID != traceID {
			t.Fail()
		}

		if sc.ParentSpanID != 1 {
			t.Fail()
		}

		if sc.ParentService != "service" {
			t.Fail()
		}

		if sc.ParentServiceInstance != "instance" {
			t.Fail()
		}

		if sc.ParentEndpoint != "/rest/api" {
			t.Fail()
		}

		if sc.AddressUsedAtClient != "foo.svc:8786" {
			t.Fail()
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan.End()
	entrySpan.End()
	wg.Wait()
}

type NoopReporter struct {
	wg    *sync.WaitGroup
	Spans []ReportedSpan
}

func (*NoopReporter) Boot(service string, instance string) {
}

func (r *NoopReporter) Send(spans []ReportedSpan) {
	r.Spans = spans
	for range spans {
		r.wg.Done()
	}
}

func (NoopReporter) Close() {
}
