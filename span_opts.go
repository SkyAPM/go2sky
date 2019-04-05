package go2sky

import "github.com/tetratelabs/go2sky/propagation"

// WithDownstream setup trace sc from propagation
func WithDownstream(cc propagation.DownstreamContext) SpanOption {
	return func(s *defaultSpan) {
		if cc == nil {
			return
		}
		header := cc.Header()
		if header == "" {
			return
		}
		tc := &propagation.SpanContext{}
		err := tc.DecodeSW6(cc.Header())
		if err != nil {
			return
		}
		s.Refs = append(s.Refs, tc)
	}
}

// WithSpanType setup span type of a span
func WithSpanType(spanType SpanType) SpanOption {
	return func(s *defaultSpan) {
		s.spanType = spanType
	}
}
