package go2sky

import (
	"context"

	"github.com/tetratelabs/go2sky/propagation"
)

// Tracer is go2sky tracer implementation.
type Tracer struct {
	service  string
	instance string
	reporter Reporter
	// 0 not init 1 init
	initFlag int32
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
	if t.reporter != nil {
		err := t.reporter.Register(t.service, t.instance)
		if err != nil {
			return nil, err
		}
		t.initFlag = 1
	}
	return t, nil
}

// CreateEntrySpan creates and starts an entry span for incoming request
func (t *Tracer) CreateEntrySpan(ctx context.Context, extractor propagation.Extractor) (Span, context.Context, error) {
	cc, err := extractor()
	if err != nil {
		return nil, nil, err
	}
	return t.CreateLocalSpan(ctx, WithParent(cc))
}

// CreateLocalSpan creates and starts a span for local usage
func (t *Tracer) CreateLocalSpan(ctx context.Context, opts ...SpanOption) (s Span, c context.Context, err error) {
	parentSpan, ok := ctx.Value(key).(Span)
	if ok && parentSpan != nil {
		opts = append(opts, WithParent(parentSpan.Context()))
	}
	ds := &defaultSpan{
		tracer: t,
	}
	for _, opt := range opts {
		opt(ds)
	}
	s = newSegmentSpan(ds, parentSpan)
	return s, context.WithValue(ctx, key, s), nil
}

// CreateExitSpan creates and starts an exit span for client
func (t *Tracer) CreateExitSpan(ctx context.Context, injector propagation.Injector) (Span, error) {
	s, _, err := t.CreateLocalSpan(ctx)
	if err != nil {
		return nil, err
	}
	cc := s.Context()
	err = injector(&cc)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Span interface as common span specification
type Span interface {
	Context() propagation.ContextCarrier
	End()
}

type defaultSpan struct {
	propagation.ContextCarrier
	tracer *Tracer
}

func (s *defaultSpan) Context() propagation.ContextCarrier {
	return s.ContextCarrier
}

func (s *defaultSpan) End() {

}

// SpanOption allows for functional options to adjust behaviour
// of a Span to be created by CreateLocalSpan
type SpanOption func(s *defaultSpan)

type ctxKey struct{}

var key = ctxKey{}

//Reporter is a data transit specification
type Reporter interface {
	Register(service string, instance string) error
	Send(spans []Span)
	Close()
}
