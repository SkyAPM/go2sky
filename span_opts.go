package go2sky

import "github.com/tetratelabs/go2sky/propagation"

// WithContext setup trace sc from propagation
func WithContext(sc *propagation.SpanContext) SpanOption {
	return func(s *defaultSpan) {
		if sc == nil {
			return
		}
		s.Refs = append(s.Refs, sc)
	}
}

// WithSpanType setup span type of a span
func WithSpanType(spanType SpanType) SpanOption {
	return func(s *defaultSpan) {
		s.spanType = spanType
	}
}
