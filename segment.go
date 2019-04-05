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
	"sync/atomic"

	"github.com/tetratelabs/go2sky/pkg"
)

func newSegmentSpan(defaultSpan *defaultSpan, parentSpan Span) Span {
	s := &segmentSpanImpl{
		defaultSpan:    *defaultSpan,
		segmentContext: &segmentContext{},
	}
	if parentSpan == nil {
		return newSegmentRoot(s)
	}
	if rootSpan, ok := parentSpan.(segmentSpan); ok {
		if rootSpan.segmentRegister() {
			s.segmentContext = rootSpan.context()
			s.sc.SpanID = atomic.AddInt32(s.SpanIDGenerator, 1)
			return s
		}
		return newSegmentRoot(s)
	}
	return newSegmentRoot(s)
}

type segmentSpan interface {
	context() *segmentContext
	segmentRegister() bool
}

type segmentSpanImpl struct {
	defaultSpan
	*segmentContext
}

func (s *segmentSpanImpl) context() *segmentContext {
	return s.segmentContext
}

type segmentContext struct {
	collect         chan<- ReportedSpan
	refNum          *int32
	SpanIDGenerator *int32
}

func (s *segmentSpanImpl) segmentRegister() bool {
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

func (s *segmentSpanImpl) End() {
	s.defaultSpan.End()
	go func() {
		s.collect <- s
	}()
}

type rootSegmentSpan struct {
	*segmentSpanImpl
	notify  <-chan ReportedSpan
	segment []ReportedSpan
	doneCh  chan int32
}

func (rs *rootSegmentSpan) End() {
	rs.defaultSpan.End()
	go func() {
		rs.doneCh <- atomic.SwapInt32(rs.refNum, -1)
	}()
}

func newSegmentRoot(segmentSpan *segmentSpanImpl) *rootSegmentSpan {
	s := &rootSegmentSpan{
		segmentSpanImpl: segmentSpan,
	}
	s.sc.SegmentID = pkg.GenerateScopedGlobalID(int64(s.tracer.instanceID))
	g := int32(0)
	s.SpanIDGenerator = &g
	s.sc.SpanID = g
	s.sc.ParentSpanID = -1
	var init int32
	s.refNum = &init
	ch := make(chan ReportedSpan)
	s.collect = ch
	s.notify = ch
	s.segment = make([]ReportedSpan, 0, 10)
	s.doneCh = make(chan int32)
	go func() {
		total := -1
		defer close(ch)
		defer close(s.doneCh)
		for {
			select {
			case span := <-s.notify:
				s.segment = append(s.segment, span)
			case n := <-s.doneCh:
				total = int(n)
			}
			if total == len(s.segment) {
				break
			}
		}
		s.tracer.reporter.Send(append(s.segment, s))
	}()
	return s
}
