//
// Copyright 2021 SkyAPM org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package go2sky

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

func TestSyncSegment(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	mr := MockReporter{
		WaitGroup: wg,
	}
	tracer, _ := NewTracer("segmentTest", WithReporter(&mr))
	ctx := context.Background()
	span, ctx, _ := tracer.CreateEntrySpan(ctx, "entry", MockExtractor)
	eSpan, _ := tracer.CreateExitSpan(ctx, "exit", "localhost:8080", MockInjector)
	eSpan.End()
	span.End()
	wg.Wait()
	if err := mr.Verify(2); err != nil {
		t.Error(err)
	}

	if eSpan.IsValid() {
		t.Error("exit span is still valid")
	}
	eSpan.End() // not be panic
	if span.IsValid() {
		t.Error("entry span is still valid")
	}
	span.End() // not be panic
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
	span, ctx, _ := tracer.CreateEntrySpan(ctx, "entry", MockExtractor)
	go func() {
		eSpan, _ := tracer.CreateExitSpan(ctx, "exit", "localhost:8080", MockInjector)
		eSpan.End()
		exitWg.Done()
	}()
	go func() {
		eSpan, _ := tracer.CreateExitSpan(ctx, "exit", "localhost:8080", MockInjector)
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
	span, ctx, _ := tracer.CreateEntrySpan(ctx, "entry", MockExtractor)
	span.End()
	reportWg.Wait()
	reportWg.Add(2)
	go func() {
		oSpan, subCtx, _ := tracer.CreateLocalSpan(ctx)
		eSpan, _ := tracer.CreateExitSpan(subCtx, "exit", "localhost:8080", MockInjector)
		eSpan.End()
		oSpan.End()
	}()
	go func() {
		oSpan, subCtx, _ := tracer.CreateLocalSpan(ctx)
		eSpan, _ := tracer.CreateExitSpan(subCtx, "exit", "localhost:8080", MockInjector)
		eSpan.End()
		oSpan.End()
	}()
	reportWg.Wait()
	if err := mr.Verify(1, 2, 2); err != nil {
		t.Error(err)
	}
}

func TestReportedSpan(t *testing.T) {
	reportWg := &sync.WaitGroup{}
	reportWg.Add(1)
	mr := MockReporter{
		WaitGroup: reportWg,
	}
	tracer, _ := NewTracer("service", WithInstance("instance"), WithReporter(&mr))
	ctx := context.Background()
	span, err := tracer.CreateExitSpan(ctx, "exit", "localhost:9999", MockInjector)
	if err != nil {
		t.Error(err)
	}
	span.SetSpanLayer(agentv3.SpanLayer_Http)
	span.Error(time.Now(), "error")
	span.Log(time.Now(), "log")
	span.Tag(TagURL, "http://localhost:9999/exit")
	span.SetComponent(1)
	span.End()
	reportWg.Wait()
	if err := mr.Verify(1); err != nil {
		t.Error(err)
	}
	reportSpan := mr.Message[0][0]
	if reportSpan.StartTime() < 0 || reportSpan.StartTime() > reportSpan.EndTime() {
		t.Error("errors are not start/end time")
	}
	if reportSpan.OperationName() != "exit" {
		t.Error("error are not set operation name")
	}
	if reportSpan.Peer() != "localhost:9999" {
		t.Error("error are not set peer")
	}
	if reportSpan.SpanType() != agentv3.SpanType_Exit {
		t.Error("error are not set span type")
	}
	if reportSpan.SpanLayer() != agentv3.SpanLayer_Http {
		t.Error("error are not set span layer")
	}
	if !reportSpan.IsError() {
		t.Error("error are not set isError")
	}
	if len(reportSpan.Tags()) != 1 {
		t.Error("error are not set tag")
	}
	if len(reportSpan.Logs()) != 2 {
		t.Error("error are not set log")
	}
	if reportSpan.ComponentID() != 1 {
		t.Error("error are not set component id")
	}
}

func MockExtractor(key string) (c string, e error) {
	return
}

func MockInjector(key, value string) (e error) {
	return
}

type Segment []ReportedSpan

type MockReporter struct {
	Reporter
	Message []Segment
	*sync.WaitGroup
	sync.Mutex
}

func (r *MockReporter) Boot(service string, serviceInstance string, cdsWatchers []AgentConfigChangeWatcher) {

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
