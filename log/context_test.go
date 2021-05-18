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
	"testing"
)

func TestFromContext(t *testing.T) {
	ctx := context.Background()
	skyWalkingContext := FromContext(ctx)
	verifyContext(t, skyWalkingContext)
}

func verifyContext(t *testing.T, ctx *SkyWalkingContext) {
	if ctx == nil {
		t.Error("nil context")
		return
	}

	contextString := ctx.String()
	if contextString == "" {
		t.Error("empty context string")
	}

	exceptString := fmt.Sprintf("[%s,%s,%s,%s,%d]", ctx.ServiceName, ctx.ServiceInstanceName,
		ctx.TraceID, ctx.TraceSegmentID, ctx.SpanID)
	if contextString != exceptString {
		t.Errorf("wrong context string, excepted:%s, actual:%s", exceptString, contextString)
	}
}
