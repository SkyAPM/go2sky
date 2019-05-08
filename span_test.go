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

package go2sky

import (
	"math"
	"testing"
	"time"

	"github.com/tetratelabs/go2sky/reporter/grpc/common"
)

func Test_defaultSpan_SetOperationName(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"set operation name",
			struct{ name string }{name: "invoke method"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &defaultSpan{}
			ds.SetOperationName(tt.args.name)
			if ds.operationName != tt.args.name {
				t.Error("operation name is different")
			}
		})
	}
}

func Test_defaultSpan_SetPeer(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"set peer",
			struct{ name string }{name: "localhost:9999"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &defaultSpan{}
			ds.SetPeer(tt.args.name)
			if ds.peer != tt.args.name {
				t.Error("peer is different")
			}
		})
	}
}

func Test_defaultSpan_SetSpanLayer(t *testing.T) {
	type args struct {
		layer common.SpanLayer
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"Set SpanLayer_Unknown",
			struct{ layer common.SpanLayer }{layer: common.SpanLayer_Unknown},
		},
		{
			"Set SpanLayer_Database",
			struct{ layer common.SpanLayer }{layer: common.SpanLayer_Database},
		},
		{
			"Set SpanLayer_RPCFramework",
			struct{ layer common.SpanLayer }{layer: common.SpanLayer_RPCFramework},
		},
		{
			"Set SpanLayer_Http",
			struct{ layer common.SpanLayer }{layer: common.SpanLayer_Http},
		},
		{
			"Set SpanLayer_MQ",
			struct{ layer common.SpanLayer }{layer: common.SpanLayer_MQ},
		},
		{
			"Set SpanLayer_Cache",
			struct{ layer common.SpanLayer }{layer: common.SpanLayer_Cache},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &defaultSpan{}
			ds.SetSpanLayer(tt.args.layer)
			if ds.layer != tt.args.layer {
				t.Error("span layer is different")
			}
		})
	}
}

func Test_defaultSpan_Tag(t *testing.T) {
	type args struct {
		key   Tag
		value string
	}
	tests := []struct {
		name string
		args []*args
	}{
		{
			"set null",
			[]*args{{}},
		},
		{
			"set tag",
			[]*args{{key: "name", value: "value"}, {key: "name1", value: "value1"}},
		},
		{
			"set duplicated tag",
			[]*args{{key: "name", value: "value"}, {key: "name", value: "value"}},
		},
		{
			"set invalid tag",
			[]*args{{key: "name"}, {value: "value"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &defaultSpan{}
			for _, a := range tt.args {
				ds.Tag(a.key, a.value)
			}
			if len(ds.tags) != len(tt.args) {
				t.Error("tags are not set property")
			}
		})
	}
}

func Test_defaultSpan_Log(t *testing.T) {
	tests := []struct {
		name string
		ll   []string
	}{
		{
			"set null logs",
			[]string{},
		},
		{
			"set logs",
			[]string{"name", "value", "name1", "value1"},
		},
		{
			"set duplicated logs",
			[]string{"name", "value", "name1", "value1"},
		},
		{
			"set invalid logs",
			[]string{"name", "value", "name1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &defaultSpan{}
			ds.Log(time.Now(), tt.ll...)
			if len(ds.logs) != 1 {
				t.Error("logs are not set property")
			}
			for _, l := range ds.logs {
				if len(l.Data) != int(math.Ceil(float64(len(tt.ll))/2)) {
					t.Error("logs are not set property")
				}
			}
		})
	}
}

func Test_defaultSpan_Error(t *testing.T) {
	tests := []struct {
		name string
		ll   []string
	}{
		{
			"set errors",
			[]string{"name", "value", "name1", "value1"},
		},
		{
			"set duplicated errors",
			[]string{"name", "value", "name1", "value1"},
		},
		{
			"set invalid errors",
			[]string{"name", "value", "name1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := &defaultSpan{}
			ds.Error(time.Now(), tt.ll...)
			if !ds.isError {
				t.Error("errors are not set property")
			}
			if len(ds.logs) != 1 {
				t.Error("errors are not set property")
			}
			for _, l := range ds.logs {
				if len(l.Data) != int(math.Ceil(float64(len(tt.ll))/2)) {
					t.Error("errors are not set property")
				}
			}
		})
	}
}
