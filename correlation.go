//
// Copyright 2022 SkyAPM org
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

type CorrelationConfig struct {
	MaxKeyCount  int
	MaxValueSize int
}

func WithCorrelation(keyCount, valueSize int) TracerOption {
	return func(t *Tracer) {
		t.correlation = &CorrelationConfig{
			MaxKeyCount:  keyCount,
			MaxValueSize: valueSize,
		}
	}
}

func PutCorrelation(ctx context.Context, key, value string) bool {
	if key == "" {
		return false
	}

	activeSpan := ctx.Value(ctxKeyInstance)
	if activeSpan == nil {
		return false
	}

	span, ok := activeSpan.(segmentSpan)
	if !ok {
		return false
	}
	correlationContext := span.context().CorrelationContext
	// remove key
	if value == "" {
		delete(correlationContext, key)
		return true
	}
	// out of max value size
	if len(value) > span.tracer().correlation.MaxValueSize {
		return false
	}
	// already exists key
	if _, ok := correlationContext[key]; ok {
		correlationContext[key] = value
		return true
	}
	// out of max key count
	if len(correlationContext) >= span.tracer().correlation.MaxKeyCount {
		return false
	}
	span.context().CorrelationContext[key] = value
	return true
}

func GetCorrelation(ctx context.Context, key string) string {
	activeSpan := ctx.Value(ctxKeyInstance)
	if activeSpan == nil {
		return ""
	}

	span, ok := activeSpan.(segmentSpan)
	if !ok {
		return ""
	}
	return span.context().CorrelationContext[key]
}
