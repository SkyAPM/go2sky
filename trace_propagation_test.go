// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package go2sky

import (
	"context"
	"strings"
	"sync"
	"testing"
)

const header string = "1-MTU1NTY0NDg4Mjk2Nzg2ODAwMC4wLjU5NDYzNzUyMDYzMzg3NDkwODc=" +
	"-NS4xNTU1NjQ0ODgyOTY3ODg5MDAwLjM3NzUyMjE1NzQ0Nzk0NjM3NTg=" +
	"-1-2-3-I2NvbS5oZWxsby5IZWxsb1dvcmxk-Iy9yZXN0L2Fh-Iy9nYXRld2F5L2Nj"

func TestTracer_EntryAndExit(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	tracer, err := NewTracer("service", WithReporter(&NoopReporter{wg: wg}))
	if err != nil {
		t.Error(err)
	}
	tracer.WaitUntilRegister()
	entrySpan, ctx, err := tracer.CreateEntrySpan(context.Background(), "/rest/api", func() (string, error) {
		return "", nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan, err := tracer.CreateExitSpan(ctx, "/foo/bar", "foo.svc:8787", func(head string) error {
		if head == "" {
			t.Fail()
		}
		if len(strings.Split(head, "-")) != 9 {
			t.Fail()
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan.End()
	entrySpan.End()
	wg.Wait()
}

func TestTracer_Entry(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	reporter := &NoopReporter{wg: wg}
	tracer, err := NewTracer("service", WithReporter(reporter))
	if err != nil {
		t.Error(err)
	}
	tracer.WaitUntilRegister()
	entrySpan, _, err := tracer.CreateEntrySpan(context.Background(), "/rest/api", func() (string, error) {
		return header, nil
	})
	if err != nil {
		t.Error(err)
	}
	entrySpan.End()
	wg.Wait()
	span := reporter.Spans[0]
	if span.Context().TraceID[0] != 1555644882967868000 || span.Context().TraceID[1] != 0 ||
		span.Context().TraceID[2] != 5946375206338749087 {
		t.Fail()
	}
	if len(span.Refs()) != 1 {
		t.Fail()
	}
}

func TestTracer_EntryAndExitInTrace(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(2)
	tracer, err := NewTracer("service", WithReporter(&NoopReporter{wg: wg}))
	if err != nil {
		t.Error(err)
	}
	tracer.WaitUntilRegister()
	entrySpan, ctx, err := tracer.CreateEntrySpan(context.Background(), "/rest/api", func() (string, error) {
		return header, nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan, err := tracer.CreateExitSpan(ctx, "/foo/bar", "foo.svc:8787", func(head string) error {
		ss := strings.Split(head, "-")
		if ss[0] != "1" {
			t.Fail()
		}
		if ss[1] != "MTU1NTY0NDg4Mjk2Nzg2ODAwMC4wLjU5NDYzNzUyMDYzMzg3NDkwODc=" {
			t.Fail()
		}
		if ss[3] != "1" {
			t.Fail()
		}
		if ss[4] != "5" {
			t.Fail()
		}
		if ss[5] != "3" {
			t.Fail()
		}
		if ss[6] != "I2Zvby5zdmM6ODc4Nw==" {
			t.Fail()
		}
		if ss[7] != "Iy9yZXN0L2Fh" {
			t.Fail()
		}
		if ss[8] != "Iy9yZXN0L2FwaQ==" {
			t.Fail()
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan.End()
	entrySpan.End()
	wg.Wait()
}

type NoopReporter struct {
	wg    *sync.WaitGroup
	Spans []ReportedSpan
}

func (*NoopReporter) Register(service string, instance string) (int32, int32, error) {
	return 2, 5, nil
}

func (r *NoopReporter) Send(spans []ReportedSpan) {
	r.Spans = spans
	for range spans {
		r.wg.Done()
	}
}

func (NoopReporter) Close() {
}
