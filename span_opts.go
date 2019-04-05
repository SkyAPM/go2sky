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
		tc := &propagation.TraceContext{}
		err := tc.DecodeSW6(cc.Header())
		if err != nil {
			return
		}
		s.tc = tc
	}
}

// WithSpanType setup span type of a span
func WithSpanType(spanType SpanType) SpanOption {
	return func(s *defaultSpan) {
		s.spanType = spanType
	}
}
