package go2sky

import (
	"context"
	"errors"
	"reflect"
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
	tracer, _ := NewTracer("", WithReporter(&mockRegisterReporter{
		success: true,
	}))
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	defer span.End()
	if err != nil {
		t.Error(err)
	}
	subSpan, _, err := tracer.CreateLocalSpan(ctx)
	defer subSpan.End()
	if err != nil {
		t.Error(err)
	}
	verifySpans(t, span, subSpan)
}

func TestTracer_CreateLocalSpanAsync(t *testing.T) {
	tracer, _ := NewTracer("", WithReporter(&mockRegisterReporter{
		success: true,
	}))
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	defer span.End()
	if err != nil {
		t.Error(err)
	}
	retCh := make(chan int32, 10)
	defer close(retCh)
	for i := 0; i < 10; i++ {
		go func() {
			subSpan, _, err := tracer.CreateLocalSpan(ctx)
			defer subSpan.End()
			if err != nil {
				t.Error(err)
			}
			verifySpans(t, span, subSpan)
			retCh <- subSpan.Context().SpanID
		}()
	}
	m := map[int32]interface{}{}
	for i := 0; i < 10; i++ {
		select {
		case a := <-retCh:
			m[a] = 0
		}
	}
	if len(m) != 10 {
		t.Error("duplicated span id")
	}
}

func verifySpans(t *testing.T, span Span, subSpan Span) {
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
}

func (r *mockRegisterReporter) Send(spans []ReportedSpan) {
}

func (r *mockRegisterReporter) Close() {
}

func (r *mockRegisterReporter) Register(service string, instance string) (int32, int32, error) {
	if r.success {
		return 1, 1, nil
	}
	return 0, 0, errRegister
}
