package go2sky

import "github.com/tetratelabs/go2sky/propagation"

func WithParent(cc propagation.ContextCarrier) SpanOption {
	return func(s *defaultSpan) {
		s.ContextCarrier = cc
	}
}

func WithSegment(sc segmentContext) SpanOption {
	return func(s *defaultSpan) {
		s.segmentContext = sc
	}
}
