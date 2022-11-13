package go2sky

import "context"

type LogData struct {
	LogCtx context.Context
	LogLevel string
	LogContent string
}
