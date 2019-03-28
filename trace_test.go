package go2sky

import (
	"errors"
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

type mockRegisterReporter struct {
	Reporter
	success bool
}

func (r *mockRegisterReporter) Register(service string, instance string) error {
	if r.success {
		return nil
	}
	return errRegister
}
