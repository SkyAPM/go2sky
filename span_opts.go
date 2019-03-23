package go2sky

import "github.com/tetratelabs/go2sky/propagation"

// WithParent setup parent context from propagation
func WithParent(cc propagation.ContextCarrier) SpanOption {
	return func(s *defaultSpan) {
		s.ContextCarrier = cc
	}
}
