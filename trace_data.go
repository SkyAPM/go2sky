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

import "context"

const (
	EmptyServiceName         = ""
	EmptyServiceInstanceName = ""
	EmptyTraceID             = "N/A"
	EmptyTraceSegmentID      = "N/A"
	EmptySpanID              = -1
)

func ServiceName(ctx context.Context) string {
	span, failed, ok := extractSpanString(ctx, EmptyServiceName)
	if !ok {
		return failed
	}
	return (*span).tracer().service
}

func ServiceInstanceName(ctx context.Context) string {
	span, failed, ok := extractSpanString(ctx, EmptyServiceInstanceName)
	if !ok {
		return failed
	}
	return (*span).tracer().instance
}

func TraceID(ctx context.Context) string {
	span, failed, ok := extractSpanString(ctx, EmptyTraceID)
	if !ok {
		return failed
	}
	return (*span).context().TraceID
}

func TraceSegmentID(ctx context.Context) string {
	span, failed, ok := extractSpanString(ctx, EmptyTraceSegmentID)
	if !ok {
		return failed
	}
	return (*span).context().SegmentID
}

func SpanID(ctx context.Context) int32 {
	span, failed, ok := extractSpanInt32(ctx, EmptySpanID)
	if !ok {
		return failed
	}
	return (*span).context().SpanID
}

func extractSpanString(ctx context.Context, noopResult string) (*segmentSpan, string, bool) {
	activeSpan := ctx.Value(ctxKeyInstance)
	if activeSpan != nil {
		span, ok := activeSpan.(segmentSpan)
		if ok {
			return &span, "", true
		}
	}
	return nil, noopResult, false
}

func extractSpanInt32(ctx context.Context, noopResult int32) (*segmentSpan, int32, bool) {
	activeSpan := ctx.Value(ctxKeyInstance)
	if activeSpan != nil {
		span, ok := activeSpan.(segmentSpan)
		if ok {
			return &span, 0, true
		}
	}
	return nil, noopResult, false
}
