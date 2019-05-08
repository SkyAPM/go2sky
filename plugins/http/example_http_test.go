// Copyright 2019 Tetrate Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/tetratelabs/go2sky"

	"github.com/gorilla/mux"

	"github.com/tetratelabs/go2sky/reporter"
)

func ExampleNewServerMiddleware() {
	r, err := reporter.NewGRPCReporter("hello.com:11800")
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

	client, err := NewClient(tracer)
	if err != nil {
		log.Fatalf("create client error %v \n", err)
	}

	router := mux.NewRouter()

	// create test server
	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatalf("listen error %v \n", err)
	}
	ts := &httptest.Server{
		Listener: l,
		Config:   &http.Server{Handler: sm(router)},
	}
	ts.Start()
	defer ts.Close()

	// add handlers
	router.Methods("GET").Path("/middle").HandlerFunc(middleFunc(client, ts.URL))
	router.Methods("POST").Path("/end").HandlerFunc(endFunc())

	// call end service
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/middle", ts.URL), nil)
	if err != nil {
		log.Fatalf("unable to create http request: %+v\n", err)
	}
	res, err := client.Do(request)
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}
	_ = res.Body.Close()
	time.Sleep(time.Minute)

	// Output:
}

func middleFunc(client *http.Client, url string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("middle func called with method: %s\n", r.Method)

		// do some operation
		time.Sleep(100 * time.Millisecond)

		newRequest, err := http.NewRequest("POST", url+"/end", nil)
		if err != nil {
			log.Printf("unable to create client: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}

		//Link the context of entry and exit spans
		newRequest = newRequest.WithContext(r.Context())

		res, err := client.Do(newRequest)
		if err != nil {
			log.Printf("call to end fund returned error: %+v\n", err)
			http.Error(w, err.Error(), 500)
			return
		}
		_ = res.Body.Close()
	}
}

func endFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("end func called with method: %s\n", r.Method)
		time.Sleep(50 * time.Millisecond)
	}
}
