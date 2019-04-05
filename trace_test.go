package go2sky

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
)

var (
	errRegister = errors.New("register error")
)

func TestTracerInit(t *testing.T) {
	_, err := NewTracer("", WithReporter(&mockRegisterReporter{
		success: true,
	}))
	if err != nil {
		t.Fail()
	}
	_, err = NewTracer("", WithReporter(&mockRegisterReporter{
		success: false,
	}))
	if err != errRegister {
		t.Fail()
	}
}

func TestTracer_CreateLocalSpan(t *testing.T) {
	reporter := &mockRegisterReporter{
		success: true,
	}
	tracer, _ := NewTracer("", WithReporter(reporter))
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		t.Error(err)
	}
	subSpan, _, err := tracer.CreateLocalSpan(ctx)
	if err != nil {
		t.Error(err)
	}
	subSpan.End()
	span.End()
	reporter.wait()
	verifySpans(t, reporter.Spans[1], reporter.Spans[0])
}

func TestTracer_CreateLocalSpanAsync(t *testing.T) {
	reporter := &mockRegisterReporter{
		success: true,
	}
	tracer, _ := NewTracer("", WithReporter(reporter))
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			subSpan, _, err := tracer.CreateLocalSpan(ctx)
			if err != nil {
				t.Error(err)
			}
			subSpan.End()
			wg.Done()
		}()
	}
	wg.Wait()
	span.End()
	reporter.wait()
	if len(reporter.Spans) != 11 {
		t.Errorf("less spans")
	}
	rootSpan := reporter.Spans[len(reporter.Spans) - 1]
	for _, span := range reporter.Spans[:len(reporter.Spans) - 2] {
		verifySpans(t, rootSpan, span)
	}
}

func verifySpans(t *testing.T, span ReportedSpan, subSpan ReportedSpan) {
	if !reflect.DeepEqual(subSpan.Context().TraceID, span.Context().TraceID) {
		t.Errorf("trace id is different %v %v", subSpan.Context().TraceID, span.Context().TraceID)
	}
	if subSpan.Context().ParentSpanID != span.Context().SpanID {
		t.Errorf("span linking is wrong %d %d", subSpan.Context().ParentSpanID, span.Context().SpanID)
	}
	if subSpan.Context().SpanID == span.Context().SpanID {
		t.Errorf("same span id %d", span.Context().SpanID)
	}
}

type mockRegisterReporter struct {
	success bool
	wg      sync.WaitGroup
	Spans   []ReportedSpan
}

func (r *mockRegisterReporter) Send(spans []ReportedSpan) {
	r.Spans = spans
	r.wg.Done()
}

func (r *mockRegisterReporter) Close() {
}

func (r *mockRegisterReporter) Register(service string, instance string) (int32, int32, error) {
	r.wg = sync.WaitGroup{}
	r.wg.Add(1)
	if r.success {
		return 1, 1, nil
	}
	return 0, 0, errRegister
}

func (r *mockRegisterReporter) wait() {
	r.wg.Wait()
}
