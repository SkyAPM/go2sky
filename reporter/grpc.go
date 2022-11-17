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
	"context"
	"io"
	"log"
	"os"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/internal/tool"
	glog "github.com/SkyAPM/go2sky/log"
	"github.com/SkyAPM/go2sky/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	configuration "skywalking.apache.org/repo/goapi/collect/agent/configuration/v3"
	commonv3 "skywalking.apache.org/repo/goapi/collect/common/v3"
	agentv3 "skywalking.apache.org/repo/goapi/collect/language/agent/v3"
	logv3 "skywalking.apache.org/repo/goapi/collect/logging/v3"
	managementv3 "skywalking.apache.org/repo/goapi/collect/management/v3"
)

const (
	maxSendQueueSize     int32 = 30000
	defaultCheckInterval       = 20 * time.Second
	defaultCDSInterval         = 20 * time.Second
	defaultLogPrefix           = "go2sky-gRPC"
	authKey                    = "Authentication"
)

func applyGRPCReporterOption(r *gRPCReporter, opts ...GRPCReporterOption) error {
	// read the options in the environment variable
	envOps, err := gRPCReporterOptionsFormEnv()
	if err != nil {
		return err
	}
	opts = append(opts, envOps...)
	for _, o := range opts {
		o(r)
	}
	return nil
}

// NewGRPCReporter create a new reporter to send data to gRPC oap server. Only one backend address is allowed.
func NewGRPCReporter(serverAddr string, opts ...GRPCReporterOption) (go2sky.Reporter, error) {
	r := &gRPCReporter{
		logger:        logger.NewDefaultLogger(log.New(os.Stderr, defaultLogPrefix, log.LstdFlags)),
		sendCh:        make(chan *agentv3.SegmentObject, maxSendQueueSize),
		logCh:         make(chan *logv3.LogData, maxSendQueueSize),
		checkInterval: defaultCheckInterval,
		cdsInterval:   defaultCDSInterval, // cds default on
	}

	// init a cancel ctx and cancel function
	r.ctx, r.cancelFunc = context.WithCancel(context.Background())

	if err := applyGRPCReporterOption(r, opts...); err != nil {
		return nil, err
	}

	// if user choose not close meter collection, make the meter channel
	if r.meterInterval == nil || (r.meterInterval != nil && *r.meterInterval > 0) {
		r.meterCh = make(chan []*agentv3.MeterData, maxSendQueueSize)
	}

	var credsDialOption grpc.DialOption
	if r.creds != nil {
		// use tls
		credsDialOption = grpc.WithTransportCredentials(r.creds)
	} else {
		credsDialOption = grpc.WithInsecure()
	}

	// read the backend service address in the environment variable
	serverAddr = serverAddrFormEnv(serverAddr)
	conn, err := grpc.Dial(serverAddr, credsDialOption)
	if err != nil {
		return nil, err
	}
	r.conn = conn
	r.traceClient = agentv3.NewTraceSegmentReportServiceClient(r.conn)
	r.managementClient = managementv3.NewManagementServiceClient(r.conn)
	r.meterClient = agentv3.NewMeterReportServiceClient(r.conn)
	r.logClient = logv3.NewLogReportServiceClient(r.conn)
	if r.cdsInterval > 0 {
		r.cdsClient = configuration.NewConfigurationDiscoveryServiceClient(r.conn)
		r.cdsService = go2sky.NewConfigDiscoveryService()
	}
	return r, nil
}

type gRPCReporter struct {
	ctx              context.Context
	cancelFunc       context.CancelFunc
	service          string
	serviceInstance  string
	instanceProps    map[string]string
	logger           logger.Log
	sendCh           chan *agentv3.SegmentObject
	meterCh          chan []*agentv3.MeterData
	logCh            chan *logv3.LogData
	conn             *grpc.ClientConn
	traceClient      agentv3.TraceSegmentReportServiceClient
	managementClient managementv3.ManagementServiceClient
	meterClient      agentv3.MeterReportServiceClient
	logClient        logv3.LogReportServiceClient
	checkInterval    time.Duration
	cdsInterval      time.Duration
	meterInterval    *time.Duration
	cdsService       *go2sky.ConfigDiscoveryService
	cdsClient        configuration.ConfigurationDiscoveryServiceClient

	// set report strategy
	rs ReportStrategy

	md    metadata.MD
	creds credentials.TransportCredentials

	// bootFlag is set if Boot be executed
	bootFlag bool

	// Instance belong layer name which define in the backend
	layer string

	// The process metadata and is enabled the process status hook
	processLabels           []string
	processStatusHookEnable bool
}

func (r *gRPCReporter) Boot(service string, serviceInstance string, cdsWatchers []go2sky.AgentConfigChangeWatcher) {
	r.service = service
	r.serviceInstance = serviceInstance
	r.initSendPipeline()
	r.check()
	r.initCDS(cdsWatchers)
	r.initMetricsCollector()
	r.initSendLogPipeline()
	r.bootFlag = true
}

func (r *gRPCReporter) Send(spans []go2sky.ReportedSpan) {
	spanSize := len(spans)
	if spanSize < 1 {
		return
	}
	rootSpan := spans[spanSize-1]
	rootCtx := rootSpan.Context()
	segmentObject := &agentv3.SegmentObject{
		TraceId:         rootCtx.TraceID,
		TraceSegmentId:  rootCtx.SegmentID,
		Spans:           make([]*agentv3.SpanObject, spanSize),
		Service:         r.service,
		ServiceInstance: r.serviceInstance,
	}
	for i, s := range spans {
		spanCtx := s.Context()
		segmentObject.Spans[i] = &agentv3.SpanObject{
			SpanId:        spanCtx.SpanID,
			ParentSpanId:  spanCtx.ParentSpanID,
			StartTime:     s.StartTime(),
			EndTime:       s.EndTime(),
			OperationName: s.OperationName(),
			Peer:          s.Peer(),
			SpanType:      s.SpanType(),
			SpanLayer:     s.SpanLayer(),
			ComponentId:   s.ComponentID(),
			IsError:       s.IsError(),
			Tags:          s.Tags(),
			Logs:          s.Logs(),
		}
		srr := make([]*agentv3.SegmentReference, 0)
		if i == (spanSize-1) && spanCtx.ParentSpanID > -1 {
			srr = append(srr, &agentv3.SegmentReference{
				RefType:               agentv3.RefType_CrossThread,
				TraceId:               spanCtx.TraceID,
				ParentTraceSegmentId:  spanCtx.ParentSegmentID,
				ParentSpanId:          spanCtx.ParentSpanID,
				ParentService:         r.service,
				ParentServiceInstance: r.serviceInstance,
			})
		}
		if len(s.Refs()) > 0 {
			for _, tc := range s.Refs() {
				srr = append(srr, &agentv3.SegmentReference{
					RefType:                  agentv3.RefType_CrossProcess,
					TraceId:                  spanCtx.TraceID,
					ParentTraceSegmentId:     tc.ParentSegmentID,
					ParentSpanId:             tc.ParentSpanID,
					ParentService:            tc.ParentService,
					ParentServiceInstance:    tc.ParentServiceInstance,
					ParentEndpoint:           tc.ParentEndpoint,
					NetworkAddressUsedAtPeer: tc.AddressUsedAtClient,
				})
			}
		}
		segmentObject.Spans[i].Refs = srr
	}
	defer func() {
		// recover the panic caused by close sendCh
		if err := recover(); err != nil {
			r.logger.Errorf("reporter segment err %v", err)
		}
	}()
	select {
	case r.sendCh <- segmentObject:
	default:
		r.logger.Errorf("reach max send buffer")
	}
}

func (r *gRPCReporter) Close() {
	// close meter collection goroutine
	r.cancelFunc()
	if r.meterCh != nil {
		// close meter send channel
		close(r.meterCh)
	}

	if r.sendCh != nil && r.bootFlag {
		close(r.sendCh)
	} else {
		r.closeGRPCConn()
		cleanupProcessDirectory(r)
	}

	if r.logCh != nil {
		close(r.logCh)
	}
}

func (r *gRPCReporter) closeGRPCConn() {
	if r.conn != nil {
		if err := r.conn.Close(); err != nil {
			r.logger.Error(err)
		}
	}
}

func (r *gRPCReporter) initSendPipeline() {
	if r.traceClient == nil {
		return
	}
	go func() {
	StreamLoop:
		for {
			stream, err := r.traceClient.Collect(metadata.NewOutgoingContext(context.Background(), r.md))
			if err != nil {
				r.logger.Errorf("open stream error %v", err)
				time.Sleep(5 * time.Second)
				continue StreamLoop
			}
			for s := range r.sendCh {
				if r.rs != nil && !r.rs(s) {
					continue
				}
				err = stream.Send(s)
				if err != nil {
					r.logger.Errorf("send segment error %v", err)
					r.closeStream(stream)
					continue StreamLoop
				}
			}
			r.closeStream(stream)
			r.closeGRPCConn()
			break
		}
	}()
}

func (r *gRPCReporter) initCDS(cdsWatchers []go2sky.AgentConfigChangeWatcher) {
	if r.cdsClient == nil {
		return
	}

	// bind watchers
	r.cdsService.BindWatchers(cdsWatchers)

	// fetch config
	go func() {
		for {
			if r.conn.GetState() == connectivity.Shutdown {
				break
			}

			configurations, err := r.cdsClient.FetchConfigurations(context.Background(), &configuration.ConfigurationSyncRequest{
				Service: r.service,
				Uuid:    r.cdsService.UUID,
			})

			if err != nil {
				r.logger.Errorf("fetch dynamic configuration error %v", err)
				time.Sleep(r.cdsInterval)
				continue
			}

			if len(configurations.GetCommands()) > 0 && configurations.GetCommands()[0].Command == "ConfigurationDiscoveryCommand" {
				command := configurations.GetCommands()[0]
				r.cdsService.HandleCommand(command)
			}

			time.Sleep(r.cdsInterval)
		}
	}()
}

func (r *gRPCReporter) initMetricsCollector() {
	if r.meterInterval != nil && *r.meterInterval <= 0 {
		r.logger.Info("user choose to close the meter collection")
		return
	}
	go2sky.InitMetricCollector(r, r.meterInterval, r.ctx)
	r.initSendMeterPipeline()
}

func (r *gRPCReporter) SendMetrics(m go2sky.RunTimeMetric) {

	meterDataList := make([]*agentv3.MeterData, 0)
	meterDataList = append(meterDataList, r.generateMeter(go2sky.InstanceGolangHeap, float64(m.HeapAlloc), m.Time))
	meterDataList = append(meterDataList, r.generateMeter(go2sky.InstanceGolangStack, float64(m.StackInUse), m.Time))
	meterDataList = append(meterDataList, r.generateMeter(go2sky.InstanceGolangGCTime, float64(m.GCPauseTime), m.Time))
	meterDataList = append(meterDataList, r.generateMeter(go2sky.InstanceGolangGCCount, float64(m.GCCount), m.Time))
	meterDataList = append(meterDataList, r.generateMeter(go2sky.InstanceGolangThreadNum, float64(m.ThreadNum), m.Time))
	meterDataList = append(meterDataList, r.generateMeter(go2sky.InstanceGolangGoroutineNum, float64(m.GoroutineNum), m.Time))
	meterDataList = append(meterDataList, r.generateMeter(go2sky.InstanceCPUUsedRate, m.CpuUsedRate, m.Time))
	meterDataList = append(meterDataList, r.generateMeter(go2sky.InstanceMemUsedRate, m.MemUsedRate, m.Time))

	defer func() {
		// recover the panic caused by close sendCh
		if err := recover(); err != nil {
			r.logger.Errorf("reporter meter err %v", err)
		}
	}()

	select {
	case <-r.ctx.Done():
		r.logger.Infof("send channel is closed")
		return
	case r.meterCh <- meterDataList:
	default:
		r.logger.Errorf("reach max send buffer")
	}
}

func (r *gRPCReporter) generateMeter(name string, value float64, time int64) *agentv3.MeterData {
	return &agentv3.MeterData{
		Metric: &agentv3.MeterData_SingleValue{
			SingleValue: &agentv3.MeterSingleValue{
				Name:  name,
				Value: value,
			},
		},
		Timestamp:       time,
		Service:         r.service,
		ServiceInstance: r.serviceInstance,
	}
}

func (r *gRPCReporter) initSendMeterPipeline() {
	if r.meterClient == nil {
		return
	}
	go func() {
	StreamLoop:
		for {
			if r.conn.GetState() == connectivity.Shutdown {
				break
			}

			stream, err := r.meterClient.CollectBatch(metadata.NewOutgoingContext(context.Background(), r.md))
			if err != nil {
				r.logger.Errorf("open stream error %v", err)
				time.Sleep(5 * time.Second)
				continue StreamLoop
			}
			for meters := range r.meterCh {
				err = stream.Send(&agentv3.MeterDataCollection{MeterData: meters})
				if err != nil {
					r.logger.Errorf("send meter error %v", err)
					r.closeMeterStream(stream)
					continue StreamLoop
				}
			}
			r.closeMeterStream(stream)
			break
		}
	}()
}

func (r *gRPCReporter) SendLog(logData go2sky.ReportedLogData) {

	if r.logClient == nil || logData == nil {
		return
	}

	reportLogData := logv3.LogData{}
	reportLogData.Service = r.service
	reportLogData.ServiceInstance = r.serviceInstance
	reportLogData.Layer = r.layer
	reportLogData.Body = &logv3.LogDataBody{Type: "text", Content: &logv3.LogDataBody_Text{Text: &logv3.TextLog{Text: logData.Data()}}}

	logLevelTag := &commonv3.KeyStringValuePair{
		Key:   "LEVEL",
		Value: string(logData.ErrorLevel()),
	}

	logTags := []*commonv3.KeyStringValuePair{}
	logTags = append(logTags, logLevelTag)
	reportLogData.Tags = &logv3.LogTags{Data: logTags}

	if logData.Context() != nil {
		traceContext := logv3.TraceContext{}

		skyCtx := glog.FromContext(logData.Context())

		traceContext.TraceId = skyCtx.TraceID
		traceContext.TraceSegmentId = skyCtx.TraceSegmentID
		traceContext.SpanId = skyCtx.SpanID

		reportLogData.TraceContext = &traceContext
	}

	defer func() {
		// recover the panic caused by close sendCh
		if err := recover(); err != nil {
			r.logger.Errorf("reporter log err %v", err)
		}
	}()

	select {
	case <-r.ctx.Done():
		r.logger.Infof("send channel is closed")
		return
	case r.logCh <- &reportLogData:
	default:
		r.logger.Errorf("reach max send buffer")
	}
}

func (r *gRPCReporter) initSendLogPipeline() {
	if r.logClient == nil {
		return
	}
	go func() {
	StreamLoop:
		for {
			stream, err := r.logClient.Collect(metadata.NewOutgoingContext(context.Background(), r.md))
			if err != nil {
				r.logger.Errorf("open stream error %v", err)
				time.Sleep(5 * time.Second)
				continue StreamLoop
			}
			for s := range r.logCh {
				err = stream.Send(s)
				if err != nil {
					r.logger.Errorf("send segment error %v", err)
					r.closeLogStream(stream)
					continue StreamLoop
				}
			}
			r.closeLogStream(stream)
			r.closeGRPCConn()
			break
		}
	}()
}

func (r *gRPCReporter) closeLogStream(stream logv3.LogReportService_CollectClient) {

	if r.conn.GetState() == connectivity.Shutdown {
		return
	}

	_, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		r.logger.Errorf("send log stream closing error %v", err)
	}
}

func (r *gRPCReporter) closeMeterStream(stream agentv3.MeterReportService_CollectBatchClient) {
	if r.conn.GetState() == connectivity.Shutdown {
		return
	}

	_, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		r.logger.Errorf("send meter stream closing error %v", err)
	}
}

func (r *gRPCReporter) closeStream(stream agentv3.TraceSegmentReportService_CollectClient) {
	_, err := stream.CloseAndRecv()
	if err != nil && err != io.EOF {
		r.logger.Errorf("send closing error %v", err)
	}
}

func (r *gRPCReporter) reportInstanceProperties() (err error) {
	props := buildOSInfo()
	if r.instanceProps != nil {
		for k, v := range r.instanceProps {
			props = append(props, &commonv3.KeyStringValuePair{
				Key:   k,
				Value: v,
			})
		}
	}
	_, err = r.managementClient.ReportInstanceProperties(metadata.NewOutgoingContext(context.Background(), r.md), &managementv3.InstanceProperties{
		Service:         r.service,
		ServiceInstance: r.serviceInstance,
		Properties:      props,
		Layer:           r.layer,
	})
	return err
}

func (r *gRPCReporter) check() {
	if r.checkInterval < 0 || r.conn == nil || r.managementClient == nil {
		return
	}
	go func() {
		instancePropertiesSubmitted := false
		for {
			if r.conn.GetState() == connectivity.Shutdown {
				break
			}

			// report the process status
			if r.processStatusHookEnable {
				reportProcess(r)
			}

			if !instancePropertiesSubmitted {
				err := r.reportInstanceProperties()
				if err != nil {
					r.logger.Errorf("report serviceInstance properties error %v", err)
					time.Sleep(r.checkInterval)
					continue
				}
				instancePropertiesSubmitted = true
			}

			_, err := r.managementClient.KeepAlive(metadata.NewOutgoingContext(context.Background(), r.md), &managementv3.InstancePingPkg{
				Service:         r.service,
				ServiceInstance: r.serviceInstance,
				Layer:           r.layer,
			})

			if err != nil {
				r.logger.Errorf("send keep alive signal error %v", err)
			}
			time.Sleep(r.checkInterval)
		}
	}()
}

func buildOSInfo() (props []*commonv3.KeyStringValuePair) {
	processNo := tool.ProcessNo()
	if processNo != "" {
		kv := &commonv3.KeyStringValuePair{
			Key:   "Process No.",
			Value: processNo,
		}
		props = append(props, kv)
	}

	hostname := &commonv3.KeyStringValuePair{
		Key:   "hostname",
		Value: tool.HostName(),
	}
	props = append(props, hostname)

	language := &commonv3.KeyStringValuePair{
		Key:   "language",
		Value: "go",
	}
	props = append(props, language)

	osName := &commonv3.KeyStringValuePair{
		Key:   "OS Name",
		Value: tool.OSName(),
	}
	props = append(props, osName)

	ipv4s := tool.AllIPV4()
	if len(ipv4s) > 0 {
		for _, ipv4 := range ipv4s {
			kv := &commonv3.KeyStringValuePair{
				Key:   "ipv4",
				Value: ipv4,
			}
			props = append(props, kv)
		}
	}
	return
}
