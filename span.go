package go2sky

import (
	"github.com/tetratelabs/go2sky/propagation"
	"sync/atomic"
)

// Span interface as common span specification
type Span interface {
	Context() SpanContext
	End()
}

// NewSpanContext create a new span context with its parent
func NewSpanContext(parentSpan Span) SpanContext {
	sc := SpanContext{}
	if parentSpan == nil {
		sc.TraceID = generateGlobalID()
		var g int32
		sc.SpanIDGenerator = &g
	} else {
		parentContext := parentSpan.Context()
		sc.TraceID = parentContext.TraceID
		sc.ParentSpanID = parentContext.SpanID
		sc.SpanIDGenerator = parentContext.SpanIDGenerator
	}
	sc.SpanID = atomic.AddInt32(sc.SpanIDGenerator, 1)
	return sc
}

// SpanContext defines the relationship between spans in one trace
type SpanContext struct {
	TraceID         []int64
	SpanID          int32
	ParentSpanID    int32
	SpanIDGenerator *int32
}

type defaultSpan struct {
	Span
	tc     *propagation.TraceContext
	tracer *Tracer
	sc     SpanContext
}

func (ds *defaultSpan) Context() SpanContext {
	return ds.sc
}

// SpanOption allows for functional options to adjust behaviour
// of a Span to be created by CreateLocalSpan
type SpanOption func(s *defaultSpan)
