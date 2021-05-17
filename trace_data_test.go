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
	"testing"
)

func TestTraceData(t *testing.T) {
	// activeSpan == nil
	verifyTraceData(context.Background(), t, "", "", EmptyTraceID, EmptyTraceSegmentID, EmptySpanID)

	// activeSpan == NoopSpan
	tracer, _ := NewTracer("service")
	_, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		t.Error(err)
	}
	verifyTraceData(ctx, t, "", "", EmptyTraceID, EmptyTraceSegmentID, EmptySpanID)

	// activeSpan == segmentSpan
	reporter := &mockRegisterReporter{
		success: true,
	}
	tracer, _ = NewTracer("service", WithInstance("instance"), WithReporter(reporter))
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		t.Error(err)
	}
	segmentContext := span.(segmentSpan).context()
	verifyTraceData(ctx, t, "service", "instance", segmentContext.TraceID, segmentContext.SegmentID, segmentContext.SpanID)
}

func verifyTraceData(ctx context.Context, t *testing.T, serviceName, serviceInstanceName, traceID, traceSegmentID string, spanID int32) {
	verifyEqual(t, "ServiceName", serviceName, ServiceName(ctx))
	verifyEqual(t, "ServiceInstanceName", serviceInstanceName, ServiceInstanceName(ctx))
	verifyEqual(t, "TraceID", traceID, TraceID(ctx))
	verifyEqual(t, "TraceSegmentID", traceSegmentID, TraceSegmentID(ctx))
	verifyEqual(t, "SpanID", spanID, SpanID(ctx))
}

func verifyEqual(t *testing.T, equalsKey string, expect interface{}, actual interface{}) {
	if expect != actual {
		t.Errorf("expect%s: %v, actual%s: %v", equalsKey, expect, equalsKey, actual)
	}
}
