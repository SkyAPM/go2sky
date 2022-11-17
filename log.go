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
	LogCtx      context.Context
	LogErrLevel LogLevel
	LogContent  string
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

type Logger struct {
	mReporter Reporter
}

func NewLogger(reporter Reporter) (*Logger, error) {

	if reporter == nil {
		return nil, errors.New("invalid reporter.")
	}

	l := new(Logger)
	l.mReporter = reporter

	return l, nil
}

func (l *Logger) WriteLogWithContext(ctx context.Context, level LogLevel, data string) {

	xLogData := DefaultLogData{}
	xLogData.LogCtx = ctx
	xLogData.LogErrLevel = level
	xLogData.LogContent = data

	l.mReporter.SendLog(&xLogData)
}
