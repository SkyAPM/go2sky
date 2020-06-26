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

package gear

import (
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/propagation"
	v3 "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
	"github.com/teambition/gear"
	"strconv"
	"time"
)

func Middleware(tracer *go2sky.Tracer) gear.Middleware {
	return func(ctx *gear.Context) error {
		if tracer == nil {
			return nil
		}

		span, _, err := tracer.CreateEntrySpan(ctx, ctx.Path, func() (string, error) {
			return ctx.GetHeader(propagation.Header), nil
		})
		if err != nil {
			return nil
		}

		span.SetComponent(go2sky.ComponentIDHttpServer)
		span.Tag(go2sky.TagHTTPMethod, ctx.Method)
		span.Tag(go2sky.TagURL, ctx.Host+ctx.Path)
		span.SetSpanLayer(v3.SpanLayer_Http)

		ctx.OnEnd(func() {
			code := ctx.Res.Status()
			span.Tag(go2sky.TagStatusCode, strconv.Itoa(code))
			if code >= 400 {
				span.Error(time.Now(), string(ctx.Res.Body()))
			}
			span.End()
		})
		return nil
	}
}
