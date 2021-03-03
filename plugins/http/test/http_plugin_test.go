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

package test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/SkyAPM/go2sky"
	httpplugin "github.com/SkyAPM/go2sky/plugins/http"
	framework "github.com/SkyAPM/go2sky/test/framework/plugin"
)

func TestRun(t *testing.T) {
	framework.CreateTestPlugin().
		WithExpectedDataFile("expected.data.yml").
		AddService("provider", func(ctx context.Context, tracer *go2sky.Tracer) error {
			sm, err := httpplugin.NewServerMiddleware(tracer)
			if err != nil {
				return err
			}

			listen, err := net.Listen("tcp", ":38080")
			if err != nil {
				return err
			}

			mux := &http.ServeMux{}

			mux.HandleFunc("/sayhi", func(writer http.ResponseWriter, request *http.Request) {
				_, _ = writer.Write([]byte("hi"))
			})

			server := &http.Server{
				Handler: sm(mux),
			}
			defer server.Close()
			go func() {
				_ = server.Serve(listen)
			}()

			<-ctx.Done()
			return nil
		}).
		AddService("consumer", func(ctx context.Context, tracer *go2sky.Tracer) error {
			client, err := httpplugin.NewClient(tracer)
			if err != nil {
				return err
			}
			timer := time.NewTimer(time.Second * 3)
			for {
				select {
				case <-ctx.Done():
					return nil
				case <-timer.C:
					req, err := http.NewRequest("GET", "http://127.0.0.1:38080/sayhi", nil)
					if err != nil {
						timer = time.NewTimer(time.Second * 3)
						continue
					}
					resp, err := client.Do(req)
					if err != nil {
						timer = time.NewTimer(time.Second * 3)
						continue
					}
					_ = resp.Body.Close()
				}
			}
		}).
		Run(t)
}
