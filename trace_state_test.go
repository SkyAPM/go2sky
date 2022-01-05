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

import (
	"testing"
)

func TestGlobalTracer(t *testing.T) {
	// globalTracer == nil
	verifyTracer(t, nil, GetGlobalTracer())

	// globalTracer ==tracer
	tracer2, _ := NewTracer("service")
	SetGlobalTracer(tracer2)
	verifyTracer(t, tracer2, GetGlobalTracer())
}

func verifyTracer(t *testing.T, expect *Tracer, actual *Tracer) {
	if expect != actual {
		t.Errorf("expect: %v, actual: %v", expect, actual)
	}
}
