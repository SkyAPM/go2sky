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
	"sync/atomic"
)

var (
	globalTracer = &atomic.Value{}
)

// SetGlobalTracer registers `tracer` as the global Tracer.
func SetGlobalTracer(tracer *Tracer) {
	globalTracer.Store(tracer)
}

// GetGlobalTracer returns the registered global Tracer.
// If none is registered then an instance of `nil` is returned.
func GetGlobalTracer() *Tracer {
	value := globalTracer.Load()
	if value != nil {
		if tracer, ok := value.(*Tracer); ok {
			return tracer
		}
	}
	return nil
}
