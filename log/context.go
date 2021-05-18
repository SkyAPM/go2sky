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

package log

import (
	"context"
	"fmt"

	"github.com/SkyAPM/go2sky"
)

type SkyWalkingContext struct {
	ServiceName         string
	ServiceInstanceName string
	TraceID             string
	TraceSegmentID      string
	SpanID              int32
}

// FromContext from context for logging
func FromContext(ctx context.Context) *SkyWalkingContext {
	return &SkyWalkingContext{
		ServiceName:         go2sky.ServiceName(ctx),
		ServiceInstanceName: go2sky.ServiceInstanceName(ctx),
		TraceID:             go2sky.TraceID(ctx),
		TraceSegmentID:      go2sky.TraceSegmentID(ctx),
		SpanID:              go2sky.SpanID(ctx),
	}
}

func (context *SkyWalkingContext) String() string {
	return fmt.Sprintf("[%s,%s,%s,%s,%d]", context.ServiceName, context.ServiceInstanceName,
		context.TraceID, context.TraceSegmentID, context.SpanID)
}
