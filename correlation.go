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

import "context"

func PutCorrelation(ctx context.Context, key, value string) {
	activeSpan := ctx.Value(ctxKeyInstance)
	if activeSpan == nil {
		return
	}

	span, ok := activeSpan.(segmentSpan)
	if ok {
		span.context().CorrelationContext[key] = value
	}
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
