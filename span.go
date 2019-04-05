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

package go2sky

import (
	"math"
	"time"

	"github.com/tetratelabs/go2sky/pkg"
	"github.com/tetratelabs/go2sky/propagation"
	"github.com/tetratelabs/go2sky/reporter/grpc/common"
	v2 "github.com/tetratelabs/go2sky/reporter/grpc/language-agent-v2"
)

// SpanType is used to identify entry, exit and local
type SpanType int32

const (
	// SpanTypeEntry is a entry span, eg http server
	SpanTypeEntry SpanType = 0
	// SpanTypeExit is a exit span, eg http client
	SpanTypeExit SpanType = 1
	// SpanTypeLocal is a local span, eg local method invoke
	SpanTypeLocal SpanType = 2
)

// BaseSpan is a base interface defines the common method sharding among different spans
type BaseSpan interface {
	Context() SpanContext
}

// Span interface as common span specification
type Span interface {
	BaseSpan
	SetOperationName(string)
	SetPeer(string)
	SetSpanLayer(common.SpanLayer)
	Tag(string, string)
	Log(time.Time, ...string)
	Error(time.Time, ...string)
	End()
}

// ReportedSpan is accessed by Reporter to load reported data
type ReportedSpan interface {
	BaseSpan
	TraceContext() *propagation.TraceContext
	StartTime() int64
	EndTime() int64
	OperationName() string
	Peer() string
	SpanType() common.SpanType
	SpanLayer() common.SpanLayer
	IsError() bool
	Tags() []*common.KeyStringValuePair
	Logs() []*v2.Log
}

func newSpanContext(parentSpan Span) SpanContext {
	var sc SpanContext
	if parentSpan == nil {
		sc = SpanContext{}
		sc.TraceID = pkg.GenerateGlobalID()
	} else {
		sc = parentSpan.Context()
		sc.ParentSpanID = parentSpan.Context().SpanID
	}
	return sc
}

// SpanContext defines the relationship between spans in one trace
type SpanContext struct {
	TraceID         []int64
	SegmentID       []int64
	SpanID          int32
	ParentSegmentID []int64
	ParentSpanID    int32
}

func newLocalSpan(t *Tracer) *defaultSpan {
	return &defaultSpan{
		tracer:    t,
		startTime: time.Now(),
		spanType:  SpanTypeLocal,
	}
}

type defaultSpan struct {
	tc            *propagation.TraceContext
	sc            SpanContext
	tracer        *Tracer
	startTime     time.Time
	endTime       time.Time
	operationName string
	peer          string
	layer         common.SpanLayer
	tags          []*common.KeyStringValuePair
	logs          []*v2.Log
	isError       bool
	spanType      SpanType
}

// For ReportedSpan

func (ds *defaultSpan) TraceContext() *propagation.TraceContext {
	return ds.tc
}

func (ds *defaultSpan) StartTime() int64 {
	return pkg.Millisecond(ds.startTime)
}

func (ds *defaultSpan) EndTime() int64 {
	return pkg.Millisecond(ds.endTime)
}

func (ds *defaultSpan) OperationName() string {
	return ds.operationName
}

func (ds *defaultSpan) Peer() string {
	return ds.peer
}

func (ds *defaultSpan) SpanType() common.SpanType {
	return common.SpanType(ds.spanType)
}

func (ds *defaultSpan) SpanLayer() common.SpanLayer {
	return ds.layer
}

func (ds *defaultSpan) IsError() bool {
	return ds.isError
}

func (ds *defaultSpan) Tags() []*common.KeyStringValuePair {
	return ds.tags
}

func (ds *defaultSpan) Logs() []*v2.Log {
	return ds.logs
}

// For Span

func (ds *defaultSpan) SetOperationName(name string) {
	ds.operationName = name
}

func (ds *defaultSpan) SetPeer(peer string) {
	ds.peer = peer
}

func (ds *defaultSpan) SetSpanLayer(layer common.SpanLayer) {
	ds.layer = layer
}

func (ds *defaultSpan) Tag(key string, value string) {
	ds.tags = append(ds.tags, &common.KeyStringValuePair{Key: key, Value: value})
}

func (ds *defaultSpan) Log(time time.Time, ll ...string) {
	data := make([]*common.KeyStringValuePair, 0, int32(math.Ceil(float64(len(ll))/2.0)))
	var kvp *common.KeyStringValuePair
	for i, l := range ll {
		if i%2 == 0 {
			kvp = &common.KeyStringValuePair{}
			data = append(data, kvp)
			kvp.Key = l
		} else {
			kvp.Value = l
		}
	}
	ds.logs = append(ds.logs, &v2.Log{Time: pkg.Millisecond(time), Data: data})
}

func (ds *defaultSpan) Error(time time.Time, ll ...string) {
	ds.isError = true
	ds.Log(time, ll...)
}

func (ds *defaultSpan) End() {
	ds.endTime = time.Now()
}

func (ds *defaultSpan) Context() SpanContext {
	return ds.sc
}

// SpanOption allows for functional options to adjust behaviour
// of a Span to be created by CreateLocalSpan
type SpanOption func(s *defaultSpan)
