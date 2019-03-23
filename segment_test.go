package go2sky

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/tetratelabs/go2sky/propagation"
)

func TestSyncSegment(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	mr := MockReporter{
		WaitGroup: wg,
	}
	tracer, _ := NewTracer(WithReporter(&mr))
	ctx := context.Background()
	span, ctx, _ := tracer.CreateEntrySpan(ctx, MockExtractor)
	eSpan, _ := tracer.CreateExitSpan(ctx, MockInjector)
	eSpan.End()
	span.End()
	wg.Wait()
	if err := mr.Verify(2); err != nil {
		t.Error(err)
	}
}

func TestAsyncSingleSegment(t *testing.T) {
	reportWg := &sync.WaitGroup{}
	exitWg := &sync.WaitGroup{}
	reportWg.Add(1)
	exitWg.Add(2)
	mr := MockReporter{
		WaitGroup: reportWg,
	}
	tracer, _ := NewTracer(WithReporter(&mr))
	ctx := context.Background()
	span, ctx, _ := tracer.CreateEntrySpan(ctx, MockExtractor)
	go func() {
		eSpan, _ := tracer.CreateExitSpan(ctx, MockInjector)
		eSpan.End()
		exitWg.Done()
	}()
	go func() {
		eSpan, _ := tracer.CreateExitSpan(ctx, MockInjector)
		eSpan.End()
		exitWg.Done()
	}()
	exitWg.Wait()
	span.End()
	reportWg.Wait()
	if err := mr.Verify(3); err != nil {
		t.Error(err)
	}
}

func TestAsyncMultipleSegments(t *testing.T) {
	reportWg := &sync.WaitGroup{}
	reportWg.Add(1)
	mr := MockReporter{
		WaitGroup: reportWg,
	}
	tracer, _ := NewTracer(WithReporter(&mr))
	ctx := context.Background()
	span, ctx, _ := tracer.CreateEntrySpan(ctx, MockExtractor)
	span.End()
	reportWg.Wait()
	reportWg.Add(2)
	go func() {
		oSpan, ctx, _ := tracer.CreateLocalSpan(ctx)
		eSpan, _ := tracer.CreateExitSpan(ctx, MockInjector)
		eSpan.End()
		oSpan.End()
	}()
	go func() {
		oSpan, ctx, _ := tracer.CreateLocalSpan(ctx)
		eSpan, _ := tracer.CreateExitSpan(ctx, MockInjector)
		eSpan.End()
		oSpan.End()
	}()
	reportWg.Wait()
	if err := mr.Verify(1, 2, 2); err != nil {
		t.Error(err)
	}
}

func MockExtractor() (c propagation.ContextCarrier, e error) {
	return
}

func MockInjector(carrier *propagation.ContextCarrier) (e error) {
	carrier.GetAllItems()
	return
}

type Segment []Span

type MockReporter struct {
	Message []Segment
	*sync.WaitGroup
	sync.Mutex
}

func (r *MockReporter) Send(spans []Span) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Message = append(r.Message, spans)
	r.WaitGroup.Done()
}

func (r *MockReporter) Verify(mm ...int) error {
	if len(mm) != len(r.Message) {
		return errors.New(fmt.Sprintf("message size mismatch. expected %d actual %d", len(mm), len(r.Message)))
	}
	for i, m := range mm {
		if m != len(r.Message[i]) {
			return errors.New(fmt.Sprintf("span size mismatch. expected %d actual %d", len(mm), len(r.Message)))
		}
	}
	return nil
}

func (r *MockReporter) Close() {
}
