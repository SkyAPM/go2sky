// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package http

import (
	"fmt"
	"log"
	"net/http/httptest"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"github.com/gorilla/mux"
)

func ExampleNewGoResty() {
	// Use log reporter for production
	r, err := reporter.NewLogReporter()
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer r.Close()

	tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(r))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}

	sm, err := NewServerMiddleware(tracer)
	if err != nil {
		log.Fatalf("create server middleware error %v \n", err)
	}

	router := mux.NewRouter()

	// create test server
	ts := httptest.NewServer(sm(router))
	defer ts.Close()

	// add handlers
	router.Methods("GET").Path("/end").HandlerFunc(endFunc())

	// create go-resty client
	client := NewGoResty(tracer)
	resp, err := client.R().Get(fmt.Sprintf("%s/end", ts.URL))
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}

	_ = resp.RawResponse.Body.Close()
	time.Sleep(time.Second)
	// Output:
}
