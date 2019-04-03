package go2sky

import (
	"math"
	"sync/atomic"
	"time"

	"github.com/tetratelabs/go2sky/pkg"
	"github.com/tetratelabs/go2sky/propagation"
	"github.com/tetratelabs/go2sky/reporter/grpc/common"
	v2 "github.com/tetratelabs/go2sky/reporter/grpc/language-agent-v2"
)

type SpanType int32

const (
	SpanTypeEntry SpanType = 0
	SpanTypeExit  SpanType = 1
	SpanTypeLocal SpanType = 2
)

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
	sc := SpanContext{}
	if parentSpan == nil {
		sc.TraceID = pkg.GenerateGlobalID()
		var g int32
		sc.SpanIDGenerator = &g
	} else {
		parentContext := parentSpan.Context()
		sc.TraceID = parentContext.TraceID
		sc.ParentSpanID = parentContext.SpanID
		sc.SpanIDGenerator = parentContext.SpanIDGenerator
	}
	sc.SpanID = atomic.AddInt32(sc.SpanIDGenerator, 1)
	return sc
}

// SpanContext defines the relationship between spans in one trace
type SpanContext struct {
	TraceID         []int64
	SegmentID       []int64
	SpanID          int32
	ParentSegmentID []int64
	ParentSpanID    int32
	SpanIDGenerator *int32
}

func newLocalSpan(t *Tracer) *defaultSpan {
	return &defaultSpan{
		tracer: t,
		startTime: time.Now(),
		spanType: SpanTypeLocal,
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
	data := make([]*common.KeyStringValuePair, int32(math.Ceil(float64(len(ll)) / 2.0)))
	for i, l := range ll {
		var kvp *common.KeyStringValuePair
		if i % 2 == 0 {
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
