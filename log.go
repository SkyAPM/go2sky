package go2sky


import (
	"context"
	"errors"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type ReportedLogData interface {
	Context() context.Context
	ErrorLevel() LogLevel
	Data() string
}

type DefaultLogData struct {
	LogCtx context.Context
	LogErrLevel LogLevel
	LogContent string
}

func (l *DefaultLogData) Context() context.Context {
	return l.LogCtx
}


func (l *DefaultLogData) ErrorLevel() LogLevel {
	return l.LogErrLevel
}

func (l *DefaultLogData) Data() string {
	return l.LogContent
}

type SkyLogger struct {
	mReporter Reporter
}

func NewSkyLogger(reporter Reporter) (*SkyLogger,error)  {

	if reporter==nil{
		return nil,errors.New("invalid reporter.")
	}

	l:=new(SkyLogger)
	l.mReporter=reporter

	return l,nil
}


func (l *SkyLogger)WriteLogWithContext(ctx context.Context,level LogLevel,data string)  {

	xLogData:=DefaultLogData{}
	xLogData.LogCtx=ctx
	xLogData.LogErrLevel=level
	xLogData.LogContent=data

	l.mReporter.SendLog(&xLogData)
}