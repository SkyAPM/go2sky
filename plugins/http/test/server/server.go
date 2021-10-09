//
// Copyright 2021 SkyAPM org
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
	"log"
	"net/http"

	"github.com/SkyAPM/go2sky"
	httpPlugin "github.com/SkyAPM/go2sky/plugins/http"
	"github.com/SkyAPM/go2sky/reporter"
)

const (
	oap     = "mockoap:19876"
	service = "http-server"
)

func main() {
	report, err := reporter.NewGRPCReporter(oap)
	if err != nil {
		log.Fatalf("crate grpc reporter error: %v \n", err)
	}

	tracer, err := go2sky.NewTracer(service, go2sky.WithReporter(report))
	if err != nil {
		log.Fatalf("crate tracer error: %v \n", err)
	}

	route := http.NewServeMux()
	route.HandleFunc("/helloserver", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Hello World!"))
	})

	sm, err := httpPlugin.NewServerMiddleware(tracer)
	if err != nil {
		log.Fatalf("create server middleware error %v \n", err)
	}
	err = http.ListenAndServe(":8080", sm(route))
	if err != nil {
		log.Fatal(err)
	}
}
