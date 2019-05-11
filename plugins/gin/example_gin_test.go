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

package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/tetratelabs/go2sky"
	"github.com/tetratelabs/go2sky/reporter"
	"log"
	"net/http"
	"net/http/httptest"
)


func ExampleNewGinServerMiddleware() {
	// Use gRPC reporter for production
	re, err := reporter.NewLogReporter()
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer re.Close()

	tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(re))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}
	tracer.WaitUntilRegister()
	gin.SetMode(gin.ReleaseMode)
	r :=gin.New()
	r.Use(Middleware(tracer))

	r.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	req := httptest.NewRequest("GET", "/user/skywalking", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

}