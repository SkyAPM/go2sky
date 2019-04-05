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
	"context"
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
	tracer, _ := NewTracer("segmentTest", WithReporter(&mr))
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
	tracer, _ := NewTracer("segmentTest", WithReporter(&mr))
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
	tracer, _ := NewTracer("segmentTest", WithReporter(&mr))
	ctx := context.Background()
	span, ctx, _ := tracer.CreateEntrySpan(ctx, MockExtractor)
	span.End()
	reportWg.Wait()
	reportWg.Add(2)
	go func() {
		oSpan, subCtx, _ := tracer.CreateLocalSpan(ctx)
		eSpan, _ := tracer.CreateExitSpan(subCtx, MockInjector)
		eSpan.End()
		oSpan.End()
	}()
	go func() {
		oSpan, subCtx, _ := tracer.CreateLocalSpan(ctx)
		eSpan, _ := tracer.CreateExitSpan(subCtx, MockInjector)
		eSpan.End()
		oSpan.End()
	}()
	reportWg.Wait()
	if err := mr.Verify(1, 2, 2); err != nil {
		t.Error(err)
	}
}

func MockExtractor() (c propagation.DownstreamContext, e error) {
	return
}

func MockInjector(carrier propagation.UpstreamContext) (e error) {
	return
}

type Segment []ReportedSpan

type MockReporter struct {
	Reporter
	Message []Segment
	*sync.WaitGroup
	sync.Mutex
}

func (r *MockReporter) Register(service string, instance string) (int32, int32, error) {
	return 0, 0, nil
}

func (r *MockReporter) Send(spans []ReportedSpan) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	r.Message = append(r.Message, spans)
	r.WaitGroup.Done()
}

func (r *MockReporter) Verify(mm ...int) error {
	if len(mm) != len(r.Message) {
		return fmt.Errorf("message size mismatch. expected %d actual %d", len(mm), len(r.Message))
	}
	for i, m := range mm {
		if m != len(r.Message[i]) {
			return fmt.Errorf("span size mismatch. expected %d actual %d", m, len(r.Message[i]))
		}
	}
	return nil
}
