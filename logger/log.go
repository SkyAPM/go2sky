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

package logger

import "log"

type Log interface {
	// Info logs to the INFO log.
	Info(args ...interface{})
	// Infof logs to the INFO log.
	Infof(format string, args ...interface{})
	// Warn logs to the WARNING and INFO logs.
	Warn(args ...interface{})
	// Warnf logs to the WARNING and INFO logs.
	Warnf(format string, args ...interface{})
	// Error logs to the ERROR, WARNING, and INFO logs.
	Error(args ...interface{})
	// Errorf logs to the ERROR, WARNING, and INFO logs.
	Errorf(format string, args ...interface{})
}

type defaultLogger struct {
	logger *log.Logger
}

func (d defaultLogger) Info(args ...interface{}) {
	d.logger.Print(args...)
}

func (d defaultLogger) Infof(format string, args ...interface{}) {
	d.logger.Printf(format, args...)
}

func (d defaultLogger) Warn(args ...interface{}) {
	d.logger.Print(args...)
}

func (d defaultLogger) Warnf(format string, args ...interface{}) {
	d.logger.Printf(format, args...)
}

func (d defaultLogger) Error(args ...interface{}) {
	d.logger.Print(args...)
}

func (d defaultLogger) Errorf(format string, args ...interface{}) {
	d.logger.Printf(format, args...)
}

// NewDefaultLogger Creates a new Log
func NewDefaultLogger(logger *log.Logger) Log {
	return &defaultLogger{
		logger: logger,
	}
}
