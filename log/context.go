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
