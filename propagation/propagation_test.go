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

package propagation

import (
	"reflect"
	"testing"
)

type fields struct {
	TraceID               string
	ParentSegmentID       string
	ParentService         string
	ParentServiceInstance string
	ParentEndpoint        string
	AddressUsedAtClient   string
	ParentSpanID          int32
	Sample                int8
}

type args struct {
	header string
}

func TestSpanContext_DecodeSW8(t *testing.T) {
	tests := []struct {
		name    string
		fields  *fields
		args    args
		wantErr bool
	}{
		{
			name:   "Empty Header",
			fields: nil,
			args: args{
				header: "",
			},
			wantErr: true,
		},
		{
			name:   "Insufficient Header Entities",
			fields: nil,
			args: args{
				header: "1-MWYyZDRiZjQ3YmY3MTFlYWI3OTRhY2RlNDgwMDExMjI=-MWU3YzIwNGE3YmY3MTFlYWI4NThhY2RlNDgwMDExMjI=",
			},
			wantErr: true,
		},
		{
			name: "normal",
			fields: &fields{
				Sample:                1,
				TraceID:               "1f2d4bf47bf711eab794acde48001122",
				ParentSegmentID:       "1e7c204a7bf711eab858acde48001122",
				ParentSpanID:          0,
				ParentService:         "service",
				ParentServiceInstance: "instance",
				ParentEndpoint:        "propagation",
				AddressUsedAtClient:   "propagation:5566",
			},
			args: args{
				header: "1-MWYyZDRiZjQ3YmY3MTFlYWI3OTRhY2RlNDgwMDExMjI=-MWU3YzIwNGE3YmY3MTFlYWI4NThhY2RlNDgwMDExMjI=" +
					"-0-c2VydmljZQ==-aW5zdGFuY2U=-cHJvcGFnYXRpb24=-cHJvcGFnYXRpb246NTU2Ng==",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &SpanContext{}
			if err := tc.DecodeSW8(tt.args.header); (err != nil) != tt.wantErr {
				t.Errorf("DecodeSW8() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.fields != nil {
				if tc.Sample != tt.fields.Sample {
					t.Fail()
				}
				if tc.TraceID != tt.fields.TraceID {
					t.Fail()
				}
				if tc.ParentSegmentID != tt.fields.ParentSegmentID {
					t.Fail()
				}
				if tc.ParentService != tt.fields.ParentService {
					t.Fail()
				}
				if tc.ParentServiceInstance != tt.fields.ParentServiceInstance {
					t.Fail()
				}
				if tc.ParentEndpoint != tt.fields.ParentEndpoint {
					t.Fail()
				}
				if tc.AddressUsedAtClient != tt.fields.AddressUsedAtClient {
					t.Fail()
				}
				if tc.ParentSpanID != tt.fields.ParentSpanID {
					t.Fail()
				}
			}
		})
	}
}

func TestSpanContext_EncodeSW8(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "normal",
			fields: fields{
				Sample:                1,
				TraceID:               "1f2d4bf47bf711eab794acde48001122",
				ParentSegmentID:       "1e7c204a7bf711eab858acde48001122",
				ParentSpanID:          0,
				ParentService:         "service",
				ParentServiceInstance: "instance",
				ParentEndpoint:        "propagation",
				AddressUsedAtClient:   "propagation:5566",
			},
			want: "1-MWYyZDRiZjQ3YmY3MTFlYWI3OTRhY2RlNDgwMDExMjI=-MWU3YzIwNGE3YmY3MTFlYWI4NThhY2RlNDgwMDExMjI=" +
				"-0-c2VydmljZQ==-aW5zdGFuY2U=-cHJvcGFnYXRpb24=-cHJvcGFnYXRpb246NTU2Ng==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &SpanContext{
				TraceID:               tt.fields.TraceID,
				ParentSegmentID:       tt.fields.ParentSegmentID,
				ParentService:         tt.fields.ParentService,
				ParentServiceInstance: tt.fields.ParentServiceInstance,
				ParentEndpoint:        tt.fields.ParentEndpoint,
				AddressUsedAtClient:   tt.fields.AddressUsedAtClient,
				ParentSpanID:          tt.fields.ParentSpanID,
				Sample:                tt.fields.Sample,
			}
			if got := tc.EncodeSW8(); got != tt.want {
				t.Errorf("EncodeSW8() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpanContext_DecodeSw8Correlation(t *testing.T) {
	tests := []struct {
		name string
		args args
		data map[string]string
	}{
		{
			name: "Empty Header",
			args: args{header: ""},
			data: make(map[string]string),
		},
		{
			name: "Insufficient Header Entities",
			args: args{header: "dGVzdC1rZXk="},
			data: make(map[string]string),
		},
		{
			name: "normal",
			args: args{header: "dGVzdC1rZXk=:dGVzdC12YWx1ZQ=="},
			data: func() map[string]string {
				m := make(map[string]string)
				m["test-key"] = "test-value"
				return m
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &SpanContext{}
			if err := tc.DecodeSW8Correlation(tt.args.header); err != nil {
				t.Errorf("DecodeSW8() error = %v, wantErr %v", err, err)
			}
			if !reflect.DeepEqual(tc.CorrelationContext, tt.data) {
				t.Fail()
			}
		})
	}
}

func TestSpanContext_EncodeSW8Correlation(t *testing.T) {
	tests := []struct {
		name string
		data map[string]string
		want string
	}{
		{
			name: "empty",
			data: make(map[string]string),
			want: "",
		},
		{
			name: "empty value",
			data: func() map[string]string {
				m := make(map[string]string)
				m["test-key"] = ""
				return m
			}(),
			want: "dGVzdC1rZXk=:",
		},
		{
			name: "normal",
			data: func() map[string]string {
				m := make(map[string]string)
				m["test-key"] = "test-value"
				return m
			}(),
			want: "dGVzdC1rZXk=:dGVzdC12YWx1ZQ==",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := &SpanContext{
				CorrelationContext: tt.data,
			}
			if got := tc.EncodeSW8Correlation(); got != tt.want {
				t.Errorf("EncodeSW8Correlation() = %v, want %v", got, tt.want)
			}
		})
	}
}
