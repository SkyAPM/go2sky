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

package skylog

import (
	"context"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"log"
	"testing"
)

var (
	skyapmOAPAddr = "oap-skywalking:11800"
)

func TestSkyapmLog(t *testing.T) {

	//if with gin.Context,ctx=ginContext.Request.Context(),then we can log with the trace
	ctx := context.Background()

	r, err := reporter.NewGRPCReporter(skyapmOAPAddr)
	if err != nil {
		log.Printf("new rpc reporter error %v \n", err)
	}

	r, err = reporter.NewLogReporter()
	if err != nil {
		log.Fatalf("new log reporter error %v \n", err)
	}

	defer func() {

		if r != nil {
			r.Close()
		}
	}()

	skyapmLogger, skyapmError := go2sky.NewSkyLogger(r)
	if skyapmError != nil {
		log.Fatalf("new SkyLogger error %v \n", skyapmError)
	}

	logData := "your application log need to send to backend here..."

	skyapmLogger.WriteLogWithContext(ctx, go2sky.LogLevelError, logData)
}
