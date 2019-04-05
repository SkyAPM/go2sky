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

package propagation

import "errors"

var errEmptyHeader = errors.New("empty header")

// DownstreamContext define the trace context from downstream
type DownstreamContext interface {
	Header() string
}

// UpstreamContext define the trace context to upstream
type UpstreamContext interface {
	SetHeader(header string)
}

// Extractor is a tool specification which define how to
// extract trace parent context from propagation context
type Extractor func() (DownstreamContext, error)

// Injector is a tool specification which define how to
// inject trace context into propagation context
type Injector func(carrier UpstreamContext) error

// TraceContext defines propagation specification of SkyWalking
type TraceContext struct {
	sample                  int8
	traceID                 []int64
	parentSegmentID         []int64
	parentSpanID            int32
	parentServiceInstanceID int32
	entryServiceInstanceID  int32
}

// DecodeSW6 converts string header to TraceContext
func (tc *TraceContext) DecodeSW6(header string) error {
	if header == "" {
		return errEmptyHeader
	}
	return nil
}
