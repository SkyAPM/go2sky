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
	"reflect"
	"sync"
	"testing"

	"github.com/SkyAPM/go2sky/propagation"
)

func TestTracerInit(t *testing.T) {
	_, err := NewTracer("service", WithReporter(&mockRegisterReporter{
		success: true,
	}))
	if err != nil {
		t.Fail()
	}
}

func TestTracer_CreateLocalSpan(t *testing.T) {
	reporter := &mockRegisterReporter{
		success: true,
	}
	tracer, _ := NewTracer("service", WithReporter(reporter))
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		t.Error(err)
	}
	subSpan, _, err := tracer.CreateLocalSpan(ctx)
	if err != nil {
		t.Error(err)
	}
	subSpan.End()
	span.End()
	reporter.wait()
	verifySpans(t, reporter.Spans[1], reporter.Spans[0])
}

func TestTracer_CreateLocalSpanAsync(t *testing.T) {
	reporter := &mockRegisterReporter{
		success: true,
	}
	tracer, _ := NewTracer("service", WithReporter(reporter))
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			subSpan, _, err := tracer.CreateLocalSpan(ctx)
			if err != nil {
				t.Error(err)
			}
			subSpan.End()
			wg.Done()
		}()
	}
	wg.Wait()
	span.End()
	reporter.wait()
	if len(reporter.Spans) != 11 {
		t.Errorf("less spans")
	}
	rootSpan := reporter.Spans[len(reporter.Spans)-1]
	for _, span := range reporter.Spans[:len(reporter.Spans)-2] {
		verifySpans(t, rootSpan, span)
	}
}

func verifySpans(t *testing.T, span ReportedSpan, subSpan ReportedSpan) {
	if !reflect.DeepEqual(subSpan.Context().TraceID, span.Context().TraceID) {
		t.Errorf("trace id is different %v %v", subSpan.Context().TraceID, span.Context().TraceID)
	}
	if subSpan.Context().ParentSpanID != span.Context().SpanID {
		t.Errorf("span linking is wrong %d %d", subSpan.Context().ParentSpanID, span.Context().SpanID)
	}
	if subSpan.Context().SpanID == span.Context().SpanID {
		t.Errorf("same span id %d", span.Context().SpanID)
	}
}

func TestTrace_TraceID(t *testing.T) {
	// activeSpan == nil
	traceID := TraceID(context.Background())
	verifyTraceID(t, EmptyTraceID, traceID)

	// activeSpan == NoopSpan
	tracer, _ := NewTracer("service")
	_, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		t.Error(err)
	}
	traceID = TraceID(ctx)
	verifyTraceID(t, NoopTraceID, traceID)

	// activeSpan == segmentSpan
	reporter := &mockRegisterReporter{
		success: true,
	}
	tracer, _ = NewTracer("service", WithReporter(reporter))
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		t.Error(err)
	}
	traceID = TraceID(ctx)
	verifyTraceID(t, span.(segmentSpan).context().TraceID, traceID)
}

func verifyTraceID(t *testing.T, expectTraceID string, actualTraceID string) {
	if expectTraceID != actualTraceID {
		t.Errorf("expectTraceID: %v, actualTraceID: %v", expectTraceID, actualTraceID)
	}
}

type mockRegisterReporter struct {
	success bool
	wg      sync.WaitGroup
	Spans   []ReportedSpan
}

func (r *mockRegisterReporter) Send(spans []ReportedSpan) {
	r.Spans = spans
	r.wg.Done()
}

func (r *mockRegisterReporter) Close() {
}

func (r *mockRegisterReporter) Boot(service string, serviceInstance string) {
	r.wg = sync.WaitGroup{}
	r.wg.Add(1)
}

func (r *mockRegisterReporter) wait() {
	r.wg.Wait()
}

func TestNewTracer(t *testing.T) {
	type args struct {
		service string
		opts    []TracerOption
	}
	tests := []struct {
		name       string
		args       args
		wantTracer *Tracer
		wantErr    bool
	}{
		{
			"null service name",
			struct {
				service string
				opts    []TracerOption
			}{service: "", opts: nil},
			nil,
			true,
		},
		{
			"without reporter",
			struct {
				service string
				opts    []TracerOption
			}{service: "test", opts: nil},
			&Tracer{service: "test", sampler: NewConstSampler(true)},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTracer, err := NewTracer(tt.args.service, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTracer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotTracer, tt.wantTracer) {
				t.Errorf("NewTracer() = %v, want %v", gotTracer, tt.wantTracer)
			}
		})
	}
}

func TestTracer_CreateEntrySpan_Parameter(t *testing.T) {
	type args struct {
		ctx           context.Context
		operationName string
		extractor     propagation.Extractor
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"context is nil",
			struct {
				ctx           context.Context
				operationName string
				extractor     propagation.Extractor
			}{ctx: nil, operationName: "query type", extractor: func() (s string, e error) {
				return "", nil
			}},
			true,
		},
		{
			"OperationName is nil",
			struct {
				ctx           context.Context
				operationName string
				extractor     propagation.Extractor
			}{ctx: context.Background(), operationName: "", extractor: func() (s string, e error) {
				return "", nil
			}},
			true,
		},
		{
			"extractor is nil",
			struct {
				ctx           context.Context
				operationName string
				extractor     propagation.Extractor
			}{ctx: context.Background(), operationName: "query type", extractor: nil},
			true,
		},
		{
			"normal",
			struct {
				ctx           context.Context
				operationName string
				extractor     propagation.Extractor
			}{ctx: context.Background(), operationName: "query type", extractor: func() (s string, e error) {
				return "", nil
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracer := &Tracer{}
			_, _, err := tracer.CreateEntrySpan(tt.args.ctx, tt.args.operationName, tt.args.extractor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tracer.CreateEntrySpan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTracer_CreateLocalSpan_Parameter(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"context is nil",
			struct {
				ctx context.Context
			}{ctx: nil},
			true,
		},
		{
			"normal",
			struct {
				ctx context.Context
			}{ctx: context.Background()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracer := &Tracer{}
			_, _, err := tracer.CreateLocalSpan(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tracer.CreateLocalSpan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTracer_CreateExitSpan_Parameter(t *testing.T) {
	type args struct {
		ctx           context.Context
		operationName string
		peer          string
		injector      propagation.Injector
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"context is nil",
			struct {
				ctx           context.Context
				operationName string
				peer          string
				injector      propagation.Injector
			}{ctx: nil, operationName: "query type", peer: "localhost:8080", injector: func(header string) error {
				return nil
			}},
			true,
		},
		{
			"OperationName is nil",
			struct {
				ctx           context.Context
				operationName string
				peer          string
				injector      propagation.Injector
			}{ctx: context.Background(), operationName: "", peer: "localhost:8080", injector: func(header string) error {
				return nil
			}},
			true,
		},
		{
			"Peer is nil",
			struct {
				ctx           context.Context
				operationName string
				peer          string
				injector      propagation.Injector
			}{ctx: context.Background(), operationName: "query type", peer: "", injector: func(header string) error {
				return nil
			}},
			true,
		},
		{
			"injector is nil",
			struct {
				ctx           context.Context
				operationName string
				peer          string
				injector      propagation.Injector
			}{ctx: context.Background(), operationName: "", peer: "localhost:8080", injector: nil},
			true,
		},
		{
			"normal",
			struct {
				ctx           context.Context
				operationName string
				peer          string
				injector      propagation.Injector
			}{ctx: context.Background(), operationName: "query type", peer: "localhost:8080", injector: func(header string) error {
				return nil
			}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tracer := &Tracer{}
			_, err := tracer.CreateExitSpan(tt.args.ctx, tt.args.operationName, tt.args.peer, tt.args.injector)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tracer.CreateExitSpan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
