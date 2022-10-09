//
// Copyright 2022 SkyAPM org
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

package reporter

import (
	"log"
	"strings"
	"time"

	"github.com/SkyAPM/go2sky/logger"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

// ReportStrategy allowed to set a custom filter
// to filter the reported segment
type ReportStrategy func(s *agentv3.SegmentObject) bool

// GRPCReporterOption allows for functional options to adjust behaviour
// of a gRPC reporter to be created by NewGRPCReporter
type GRPCReporterOption func(r *gRPCReporter)

// WithLogger setup logger for gRPC reporter
// Deprecated: WithLog is recommended
func WithLogger(log *log.Logger) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.logger = logger.NewDefaultLogger(log)
	}
}

// WithLog setup log for gRPC reporter
func WithLog(logger logger.Log) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.logger = logger
	}
}

// WithCheckInterval setup service and endpoint registry check interval
func WithCheckInterval(interval time.Duration) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.checkInterval = interval
	}
}

// WithMaxSendQueueSize setup send span queue buffer length
func WithMaxSendQueueSize(maxSendQueueSize int) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.sendCh = make(chan *agentv3.SegmentObject, maxSendQueueSize)
	}
}

// WithInstanceProps setup service instance properties eg: org=SkyAPM
func WithInstanceProps(props map[string]string) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.instanceProps = props
	}
}

// WithTransportCredentials setup transport layer security
func WithTransportCredentials(creds credentials.TransportCredentials) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.creds = creds
	}
}

// WithAuthentication used Authentication for gRPC
func WithAuthentication(auth string) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.md = metadata.New(map[string]string{authKey: auth})
	}
}

// WithCDS setup Configuration Discovery Service to dynamic config
func WithCDS(interval time.Duration) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.cdsInterval = interval
	}
}

// WithLayer setup layer
func WithLayer(layer string) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.layer = layer
	}
}

// WithFAASLayer set layer to FAAS
func WithFAASLayer() GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.layer = "FAAS"
	}
}

// WithProcessLabels setup labels bind to process
func WithProcessLabels(labels []string) GRPCReporterOption {
	return func(t *gRPCReporter) {
		t.processLabels = labels
		if t.instanceProps == nil {
			t.instanceProps = make(map[string]string)
		}
		t.instanceProps[ProcessLabelKey] = strings.Join(labels, ",")
	}
}

// WithProcessStatusHook setup is enabled the process status
func WithProcessStatusHook(enable bool) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.processStatusHookEnable = enable
	}
}

// WithReportStrategy set report strategy
func WithReportStrategy(rs ReportStrategy) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.rs = rs
	}
}

// WithMeterCollectPeriod setup is set the meter collect interval
func WithMeterCollectPeriod(interval time.Duration) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.meterInterval = &interval
	}
}
