package go2sky


import (
	"context"
	"errors"
)

type ReportedLogData interface {
	Context() context.Context
	ErrorLevel() string
	Data() string
}

type DefaultLogData struct {
	LogCtx context.Context
	LogErrLevel string
	LogContent string
}

func (l *DefaultLogData) Context() context.Context {
	return l.LogCtx
}


func (l *DefaultLogData) ErrorLevel() string {
	return l.LogErrLevel
}

func (l *DefaultLogData) Data() string {
	return l.LogContent
}

type SkyLogger struct {
	mReporter Reporter
}

func NewSkyLogger(reporter Reporter) (error,*SkyLogger)  {

	if reporter==nil{
		return errors.New("invalid reporter."),nil
	}

	l:=new(SkyLogger)
	l.mReporter=reporter

	return nil,l
}

func (l *SkyLogger)WriteLogWithContext(ctx context.Context,level string,data string)  {

	xLogData:=DefaultLogData{}
	xLogData.LogCtx=ctx
	xLogData.LogErrLevel=level
	xLogData.LogContent=data

	l.mReporter.SendLog(&xLogData)
}