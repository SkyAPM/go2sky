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

// Span interface as common span specification
type Span interface {
	SetOperationName(string)
	SetPeer(string)
	SetSpanLayer(common.SpanLayer)
	Tag(string, string)
	Log(time.Time, ...string)
	Error(time.Time, ...string)
	End()
}

func newLocalSpan(t *Tracer) *defaultSpan {
	return &defaultSpan{
		tracer:    t,
		startTime: time.Now(),
		spanType:  SpanTypeLocal,
	}
}

type defaultSpan struct {
	Refs          []*propagation.SpanContext
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
			if kvp != nil {
				kvp.Value = l
			}
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

// SpanOption allows for functional options to adjust behaviour
// of a Span to be created by CreateLocalSpan
type SpanOption func(s *defaultSpan)
