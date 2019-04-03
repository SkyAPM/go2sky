package reporter

import (
	"context"
	"errors"
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
	maxSendQueueSize    int32  = 30000
	defaultPingInterval        = 20 * time.Second
)

var (
	errRegister = errors.New("fail to register reporter")
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
	err := r.registerService(service)
	if err != nil {
		return 0, 0, err
	}
	err = r.registerInstance(instance)
	if err == nil {
		r.initSendPipeline()
		r.ping()
	}
	return r.serviceID, r.instanceID, err
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
		return errRegister
	}
	r.serviceID = mapping.Services[0].Value
	r.logger.Printf("the id of service %s is %d", name, r.serviceID)
	return nil
}

func (r *gRPCReporter) registerInstance(name string) error {
	in := &register.ServiceInstances{
		Instances: []*register.ServiceInstance{
			{
				ServiceId:    r.serviceID,
				InstanceUUID: name,
			},
		},
	}
	mapping, err := r.registerClient.DoServiceInstanceRegister(context.Background(), in)
	if err != nil {
		return err
	}
	if len(mapping.ServiceInstances) < 1 {
		return errRegister
	}
	r.instanceID = mapping.ServiceInstances[0].Value
	r.instanceName = name
	r.logger.Printf("the id of instance %s 's id is %d", name, r.serviceID)
	return nil
}

func (r *gRPCReporter) Send(spans []go2sky.ReportedSpan) {
	if len(spans) < 1 {
		return
	}
	rootSpan := spans[0]
	segment := &common.UpstreamSegment{
		GlobalTraceIds: []*common.UniqueId{
			{
				IdParts: rootSpan.Context().TraceID,
			},
		},
	}
	segmentObject := v2.SegmentObject{
		ServiceId:         r.serviceID,
		ServiceInstanceId: r.instanceID,
		TraceSegmentId: &common.UniqueId{
			IdParts: rootSpan.Context().SegmentID,
		},
		Spans: make([]*v2.SpanObjectV2, len(spans)),
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
			IsError:     s.IsError(),
			Tags:        s.Tags(),
			Logs:        s.Logs(),
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
	var b []byte
	b, err := segmentObject.XXX_Marshal(b, true)
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
	err := r.conn.Close()
	if err != nil {
		r.logger.Print(err)
	}
}

func (r *gRPCReporter) initSendPipeline() {
	go func() {
		var err error
	StreamLoop:
		for {
			stream, _ := r.traceClient.Collect(context.Background())
			for {
				select {
				case s := <-r.sendCh:
					err = stream.Send(s)
					if err != nil {
						log.Printf("send segment error %v", err)
						err = stream.CloseSend()
						if err != nil {
							log.Printf("send closing error %v", err)
						}
						continue StreamLoop
					}
				}
			}
		}
	}()
}

func (r *gRPCReporter) ping() {
	for {
		_, err := r.pingClient.DoPing(context.Background(), &register.ServiceInstancePingPkg{
			Time:                pkg.Millisecond(time.Now()),
			ServiceInstanceId:   r.instanceID,
			ServiceInstanceUUID: r.instanceName,
		})
		if err != nil {
			log.Printf("pinging error %v", err)
		}
		time.Sleep(r.pingInterval)
	}
}
