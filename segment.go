package go2sky

import "sync/atomic"

func newSegmentSpan(defaultSpan *defaultSpan, parentSpan Span) Span {
	s := &segmentSpanImpl{
		defaultSpan: *defaultSpan,
		segmentContext: &segmentContext{},
	}
	if parentSpan == nil {
		return newSegmentRoot(s)
	}
	if rootSpan, ok := parentSpan.(segmentSpan); ok {
		if rootSpan.segmentRegister() {
			s.segmentContext = rootSpan.context()
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
	collect chan<- Span
	refNum  *int32
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
	go func() {
		s.collect <- s
	}()
}

type rootSegmentSpan struct {
	*segmentSpanImpl
	notify  <-chan Span
	segment []Span
	doneCh  chan int32
}

func (rs *rootSegmentSpan) End() {
	go func() {
		rs.doneCh <- atomic.SwapInt32(rs.refNum, -1)
	}()
}

func newSegmentRoot(segmentSpan *segmentSpanImpl) *rootSegmentSpan {
	s := &rootSegmentSpan{
		segmentSpanImpl: segmentSpan,
	}
	var init int32
	s.refNum = &init
	ch := make(chan Span)
	s.collect = ch
	s.notify = ch
	s.segment = make([]Span, 0, 10)
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
