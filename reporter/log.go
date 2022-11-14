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

package reporter

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/SkyAPM/go2sky"
)

func NewLogReporter() (go2sky.Reporter, error) {
	return &logReporter{logger: log.New(os.Stderr, "go2sky-log", log.LstdFlags)}, nil
}

type logReporter struct {
	logger *log.Logger
}

func (lr *logReporter) SendLog(logData go2sky.ReportedLogData) {

	if logData == nil {
		return
	}

	if strings.EqualFold(string(logData.ErrorLevel()), string(go2sky.LogLevelWarn)) {
		os.Stderr.WriteString(logData.Data())
		os.Stderr.WriteString("\n")
		return
	}

	if strings.EqualFold(string(logData.ErrorLevel()), string(go2sky.LogLevelError)) {
		os.Stderr.WriteString(logData.Data())
		os.Stderr.WriteString("\n")
		return
	}

	os.Stdout.WriteString(logData.Data())
	os.Stdout.WriteString("\n")

}

func (lr *logReporter) Boot(service string, serviceInstance string, cdsWatchers []go2sky.AgentConfigChangeWatcher) {

}

func (lr *logReporter) Send(spans []go2sky.ReportedSpan) {
	if spans == nil {
		return
	}
	b, err := json.Marshal(spans)
	if err != nil {
		lr.logger.Printf("Error: %s", err)
		return
	}
	root := spans[len(spans)-1]
	lr.logger.Printf("Segment-%v: %s \n", root.Context().SegmentID, b)
}

func (lr *logReporter) Close() {
	lr.logger.Println("Close log reporter")
}
