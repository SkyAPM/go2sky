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

/*
Package http contains several client/server http plugin which can be used for integration with net/http.
*/

package gin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tetratelabs/go2sky"
	"github.com/tetratelabs/go2sky/propagation"
	"github.com/tetratelabs/go2sky/reporter/grpc/common"
	"strconv"
	"time"
)

const (
	httpServerComponentID int32 = 49
)


func  Middleware(tracer *go2sky.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		if tracer != nil {
			operationName :=fmt.Sprintf("/%s%s", c.Request.Method, c.Request.URL.Path)
			span, ctx, err :=tracer.CreateEntrySpan(c.Request.Context(),operationName,func() (string, error) {
				return c.Request.Header.Get(propagation.Header), nil
			})
			if err != nil {
				c.Next()
				return
			}
			span.SetComponent(httpServerComponentID)
			span.Tag(go2sky.TagHTTPMethod,c.Request.Method)
			span.Tag(go2sky.TagURL, fmt.Sprintf("%s%s", c.Request.Host, c.Request.URL.Path))
			span.SetSpanLayer(common.SpanLayer_Http)

			c.Request = c.Request.WithContext(ctx)

			c.Next()

			if len(c.Errors) > 0 {
				span.Error(time.Now(), c.Errors.String())
			}

			span.Tag(go2sky.TagStatusCode, strconv.Itoa(c.Writer.Status()))
			span.End()
		}

	}
}