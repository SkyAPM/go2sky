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
	"io/ioutil"
	"log"
	"net/http"

	"github.com/SkyAPM/go2sky"
	httpPlugin "github.com/SkyAPM/go2sky/plugins/http"
	"github.com/SkyAPM/go2sky/reporter"
)

const (
	oap         = "mockoap:19876"
	service     = "http-client"
	upstreamURL = "http://httpserver:8080/helloserver"
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

	client, err := httpPlugin.NewClient(tracer)
	if err != nil {
		log.Fatalf("create client error %v \n", err)
	}

	route := http.NewServeMux()
	route.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		clientReq, err1 := http.NewRequest(http.MethodGet, upstreamURL, nil)
		if err1 != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to create http request error: %v \n", err)
			return
		}
		clientReq = clientReq.WithContext(request.Context())
		res, err1 := client.Do(clientReq)
		if err1 != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("unable to do http request error: %v \n", err)
			return
		}
		defer res.Body.Close()
		body, err1 := ioutil.ReadAll(res.Body)
		if err1 != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			log.Printf("read http response error: %v \n", err)
			return
		}
		writer.WriteHeader(res.StatusCode)
		_, _ = writer.Write(body)
	})

	sm, err := httpPlugin.NewServerMiddleware(tracer)
	if err != nil {
		log.Fatalf("create client error %v \n", err)
	}
	err = http.ListenAndServe(":8080", sm(route))
	if err != nil {
		log.Fatal(err)
	}
}
