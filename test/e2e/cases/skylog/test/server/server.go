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

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

const (
	oap = "mockoap:19876"
)

func main() {

	oapAddr := os.Getenv("GO2SKY_OAP")
	if len(oapAddr) < 1 {
		oapAddr = oap
	}

	report, err := reporter.NewGRPCReporter(oapAddr)
	if err != nil {
		log.Fatalf("create grpc reporter error: %v \n", err)
	}

	log.Println("create grpc reporter success.")

	skylogWriter, err := go2sky.NewLogger(report)
	if err != nil {
		log.Fatalf("crate logger error: %v \n", err)
	}

	log.Println("create logger success.")

	route := http.NewServeMux()

	route.HandleFunc("/healthCheck", func(writer http.ResponseWriter, request *http.Request) {

		skylogWriter.WriteLogWithContext(request.Context(), go2sky.LogLevelInfo, fmt.Sprintf("log data from path=[%s]", request.URL.Path))

		_, _ = writer.Write([]byte("I am fine!"))
	})

	route.HandleFunc("/helloserver", func(writer http.ResponseWriter, request *http.Request) {

		skylogWriter.WriteLogWithContext(request.Context(), go2sky.LogLevelInfo, fmt.Sprintf("log data from path=[%s]", request.URL.Path))

		_, _ = writer.Write([]byte("Hello World!"))
	})

	log.Println("create server router success.")

	err = http.ListenAndServe(":58080", route)
	if err != nil {
		log.Fatal(err)
	}
}
