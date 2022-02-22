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

	"github.com/powerapm/go2sky"
	"github.com/powerapm/go2sky/reporter"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/mux"
)

func Example_newGoResty() {
	// Use gRPC reporter for production
	r, err := reporter.NewLogReporter()
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer r.Close()

	tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(r))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}
	tracer.WaitUntilRegister()

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
	client := newGoResty(tracer)
	resp, err := client.R().Get(fmt.Sprintf("%s/end", ts.URL))
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}

	_ = resp.RawResponse.Body.Close()
	time.Sleep(time.Second)

	// Output:
}

func newGoResty(tracer *go2sky.Tracer) *resty.Client {
	hc, err := NewClient(tracer)
	if err != nil {
		log.Fatalf("create client error %v \n", err)
	}
	return resty.NewWithClient(hc)
}
