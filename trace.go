// Copyright 2019 Tetrate Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package go2sky

import (
	"context"

	"github.com/google/uuid"

	"github.com/tetratelabs/go2sky/propagation"
)

// Tracer is go2sky tracer implementation.
type Tracer struct {
	service  string
	instance string
	reporter Reporter
	// 0 not init 1 init
	initFlag   int32
	serviceID  int32
	instanceID int32
}

// TracerOption allows for functional options to adjust behaviour
// of a Tracer to be created by NewTracer
type TracerOption func(t *Tracer)

// NewTracer return a new go2sky Tracer
func NewTracer(service string, opts ...TracerOption) (tracer *Tracer, err error) {
	t := &Tracer{
		service:  service,
		initFlag: 0,
	}
	for _, opt := range opts {
		opt(t)
	}
	if t.instance == "" {
		id, err := uuid.NewUUID()
		if err != nil {
			return nil, err
		}
		t.instance = id.String()
	}
	if t.reporter != nil {
		serviceID, instanceID, err := t.reporter.Register(t.service, t.instance)
		if err != nil {
			return nil, err
		}
		t.initFlag = 1
		t.serviceID = serviceID
		t.instanceID = instanceID
	}
	return t, nil
}

// CreateEntrySpan creates and starts an entry span for incoming request
func (t *Tracer) CreateEntrySpan(ctx context.Context, extractor propagation.Extractor) (Span, context.Context, error) {
	dc, err := extractor()
	if err != nil {
		return nil, nil, err
	}
	return t.CreateLocalSpan(ctx, WithDownstream(dc), WithSpanType(SpanTypeEntry))
}

// CreateLocalSpan creates and starts a span for local usage
func (t *Tracer) CreateLocalSpan(ctx context.Context, opts ...SpanOption) (s Span, c context.Context, err error) {
	ds := newLocalSpan(t)
	for _, opt := range opts {
		opt(ds)
	}
	parentSpan, ok := ctx.Value(key).(Span)
	if !ok {
		parentSpan = nil
	}
	ds.sc = newSpanContext(parentSpan)
	s = newSegmentSpan(ds, parentSpan)
	return s, context.WithValue(ctx, key, s), nil
}

// CreateExitSpan creates and starts an exit span for client
func (t *Tracer) CreateExitSpan(ctx context.Context, injector propagation.Injector) (Span, error) {
	s, _, err := t.CreateLocalSpan(ctx, WithSpanType(SpanTypeExit))
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return s, nil
}

type ctxKey struct{}

var key = ctxKey{}

//Reporter is a data transit specification
type Reporter interface {
	Register(service string, instance string) (int32, int32, error)
	Send(spans []ReportedSpan)
	Close()
}
