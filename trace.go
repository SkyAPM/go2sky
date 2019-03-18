package go2sky

import (
	"context"
	"sync/atomic"

	"github.com/tetratelabs/go2sky/propagation"
)

// Tracer is go2sky tracer implementation.
type Tracer struct {
	serviceCode string
}

// TracerOption allows for functional options to adjust behaviour
// of a Tracer to be created by NewTracer
type TracerOption func(t *Tracer) error

// NewTracer return a new go2sky Tracer
func NewTracer(opts ...TracerOption) (tracer *Tracer, err error) {
	return &Tracer{}, nil
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
func (t *Tracer) CreateLocalSpan(ctx context.Context, opts ...SpanOption) (Span, context.Context, error) {
	root := true
	if parentSpan, ok := ctx.Value(key).(Span); ok && parentSpan != nil {
		opts = append(opts, WithParent(parentSpan.Context()))
		if parentRootSpan, okk := parentSpan.(SegmentSpan); okk {
			root = parentRootSpan.SegmentRegister()
			opts = append(opts, WithSegment(parentRootSpan.SegmentContext()))
		}
	}
	s := &defaultSpan{}
	for _, opt := range opts {
		opt(s)
	}
	if root {
		s.createSegment()
	}
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

// SegmentSpan interface as segment span specification
type SegmentSpan interface {
	SegmentRegister() bool
	SegmentContext() segmentContext
}

type defaultSpan struct {
	propagation.ContextCarrier
	segmentContext
	notify  <-chan Span
	segment []Span
	doneCh  chan int32
}

type segmentContext struct {
	collect chan<- Span
	refNum  *int32
}

func (s *defaultSpan) Context() propagation.ContextCarrier {
	return s.ContextCarrier
}

func (s *defaultSpan) SegmentRegister() bool {
	for {
		o := atomic.LoadInt32(s.refNum)
		if o < 0 {
			return false
		}
		if atomic.CompareAndSwapInt32(s.refNum, o, o+1) {
			return true
		}
	}
}

func (s *defaultSpan) SegmentContext() segmentContext {
	return s.segmentContext
}

func (s *defaultSpan) End() {
	go func() {
		if s.segment == nil {
			s.collect <- s
			return
		}
		s.doneCh <- atomic.SwapInt32(s.refNum, -1)
	}()
}

func (s *defaultSpan) createSegment() {
	atomic.StoreInt32(s.refNum, 0)
	ch := make(chan Span)
	s.collect = ch
	s.notify = ch
	s.segment = make([]Span, 0, 10)
	s.doneCh = make(chan int32)
	go func() {
		total := -1
		defer close(ch)
		close(s.doneCh)
		for {
			select {
			case span, ok := <-s.notify:
				if !ok {
					return
				}
				s.segment = append(s.segment, span)
			case total = <-s.doneCh:
			}
			if total == len(s.segment) {
				break
			}
		}
		// Todo compose segment
	}()
}

// SpanOption allows for functional options to adjust behaviour
// of a Span to be created by CreateLocalSpan
type SpanOption func(s *defaultSpan)

type ctxKey struct{}

var key = ctxKey{}
