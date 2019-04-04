package reporter

import (
	"context"
	"errors"
	"github.com/golang/protobuf/proto"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"

	"github.com/tetratelabs/go2sky"
	"github.com/tetratelabs/go2sky/pkg"
	"github.com/tetratelabs/go2sky/reporter/grpc/common"
	v2 "github.com/tetratelabs/go2sky/reporter/grpc/language-agent-v2"
	"github.com/tetratelabs/go2sky/reporter/grpc/register"
)

const (
	maxSendQueueSize    int32 = 30000
	defaultPingInterval       = 20 * time.Second
)

var (
	errServiceRegister  = errors.New("fail to register service")
	errInstanceRegister = errors.New("fail to instance service")
)

// NewGRPCReporter create a new reporter to send data to gRPC oap server
func NewGRPCReporter(serverAddr string, opts ...GRPCReporterOption) (go2sky.Reporter, error) {
	r := &gRPCReporter{
		logger:       log.New(os.Stderr, "go2sky", log.LstdFlags),
		sendCh:       make(chan *common.UpstreamSegment, maxSendQueueSize),
		pingInterval: defaultPingInterval,
	}
	for _, o := range opts {
		o(r)
	}
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure()) //TODO add TLS
	if err != nil {
		return nil, err
	}
	r.conn = conn
	r.registerClient = register.NewRegisterClient(conn)
	r.pingClient = register.NewServiceInstancePingClient(conn)
	r.traceClient = v2.NewTraceSegmentReportServiceClient(r.conn)
	return r, nil
}

// GRPCReporterOption allows for functional options to adjust behaviour
// of a gRPC reporter to be created by NewGRPCReporter
type GRPCReporterOption func(r *gRPCReporter)

// WithLogger setup logger for gRPC reporter
func WithLogger(logger *log.Logger) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.logger = logger
	}
}

// WithPingInterval setup ping interval
func WithPingInterval(interval time.Duration) GRPCReporterOption {
	return func(r *gRPCReporter) {
		r.pingInterval = interval
	}
}

type gRPCReporter struct {
	serviceID      int32
	instanceID     int32
	instanceName   string
	logger         *log.Logger
	sendCh         chan *common.UpstreamSegment
	registerClient register.RegisterClient
	conn           *grpc.ClientConn
	traceClient    v2.TraceSegmentReportServiceClient
	pingClient     register.ServiceInstancePingClient
	pingInterval   time.Duration
}

func (r *gRPCReporter) Register(service string, instance string) (int32, int32, error) {
	r.retryRegister(func() error {
		return r.registerService(service)
	})
	r.retryRegister(func() error {
		return r.registerInstance(instance)
	})
	r.initSendPipeline()
	r.ping()
	return r.serviceID, r.instanceID, nil
}

type retryFunction func() error

func (r *gRPCReporter) retryRegister(f retryFunction) {
	for {
		err := f()
		if err == nil {
			break
		}
		r.logger.Printf("register error %v \n", err)
		time.Sleep(time.Second)
	}
}

func (r *gRPCReporter) registerService(name string) error {
	in := &register.Services{
		Services: []*register.Service{
			{
				ServiceName: name,
			},
		},
	}
	mapping, err := r.registerClient.DoServiceRegister(context.Background(), in)
	if err != nil {
		return err
	}
	if len(mapping.Services) < 1 {
		return errServiceRegister
	}
	r.serviceID = mapping.Services[0].Value
	r.logger.Printf("the id of service '%s' is %d", name, r.serviceID)
	return nil
}

func (r *gRPCReporter) registerInstance(name string) error {
	in := &register.ServiceInstances{
		Instances: []*register.ServiceInstance{
			{
				ServiceId:    r.serviceID,
				InstanceUUID: name,
				Time:         pkg.Millisecond(time.Now()),
			},
		},
	}
	mapping, err := r.registerClient.DoServiceInstanceRegister(context.Background(), in)
	if err != nil {
		return err
	}
	if len(mapping.ServiceInstances) < 1 {
		return errInstanceRegister
	}
	r.instanceID = mapping.ServiceInstances[0].Value
	r.instanceName = name
	r.logger.Printf("the id of instance '%s' id is %d", name, r.instanceID)
	return nil
}

func (r *gRPCReporter) Send(spans []go2sky.ReportedSpan) {
	spanSize := len(spans)
	if spanSize < 1 {
		return
	}
	rootSpan := spans[spanSize-1]
	segment := &common.UpstreamSegment{
		GlobalTraceIds: []*common.UniqueId{
			{
				IdParts: rootSpan.Context().TraceID,
			},
		},
	}
	segmentObject := &v2.SegmentObject{
		ServiceId:         r.serviceID,
		ServiceInstanceId: r.instanceID,
		TraceSegmentId: &common.UniqueId{
			IdParts: rootSpan.Context().SegmentID,
		},
		Spans: make([]*v2.SpanObjectV2, spanSize),
	}
	for i, s := range spans {
		segmentObject.Spans[i] = &v2.SpanObjectV2{
			SpanId:        s.Context().SpanID,
			ParentSpanId:  s.Context().ParentSpanID,
			StartTime:     s.StartTime(),
			EndTime:       s.EndTime(),
			OperationName: s.OperationName(),
			Peer:          s.Peer(),
			SpanType:      s.SpanType(),
			SpanLayer:     s.SpanLayer(),
			IsError:       s.IsError(),
			Tags:          s.Tags(),
			Logs:          s.Logs(),
		}
		srr := make([]*v2.SegmentReference, 0)
		if i == 0 && s.Context().ParentSpanID > -1 {
			srr = append(srr, &v2.SegmentReference{
				ParentSpanId: s.Context().ParentSpanID,
				ParentTraceSegmentId: &common.UniqueId{
					IdParts: s.Context().ParentSegmentID,
				},
				ParentServiceInstanceId: r.instanceID,
			})
		}
	}
	b, err := proto.Marshal(segmentObject)
	if err != nil {
		log.Printf("marshal segment object err %v", err)
		return
	}
	segment.Segment = b
	select {
	case r.sendCh <- segment:
	default:
		log.Printf("reach max send buffer")
	}
}

func (r *gRPCReporter) Close() {
	close(r.sendCh)
	err := r.conn.Close()
	if err != nil {
		r.logger.Print(err)
	}
}

func (r *gRPCReporter) initSendPipeline() {
	if r.traceClient == nil {
		return
	}
	go func() {
	StreamLoop:
		for {
			stream, err := r.traceClient.Collect(context.Background())
			for {
				select {
				case s, ok := <-r.sendCh:
					if !ok {
						r.closeStream(stream)
						return
					}
					err = stream.Send(s)
					if err != nil {
						r.logger.Printf("send segment error %v", err)
						r.closeStream(stream)
						continue StreamLoop
					}
				}
			}
		}
	}()
}

func (r *gRPCReporter) closeStream(stream v2.TraceSegmentReportService_CollectClient) {
	err := stream.CloseSend()
	if err != nil {
		r.logger.Printf("send closing error %v", err)
	}
}

func (r *gRPCReporter) ping() {
	if r.pingInterval < 0 || r.pingClient == nil {
		return
	}
	go func() {
		for {
			_, err := r.pingClient.DoPing(context.Background(), &register.ServiceInstancePingPkg{
				Time:                pkg.Millisecond(time.Now()),
				ServiceInstanceId:   r.instanceID,
				ServiceInstanceUUID: r.instanceName,
			})
			if err != nil {
				r.logger.Printf("pinging error %v", err)
			}
			time.Sleep(r.pingInterval)
		}
	}()

}
