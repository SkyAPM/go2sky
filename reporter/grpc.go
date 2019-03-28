package reporter

import (
	"context"
	"errors"
	"log"
	"os"

	"google.golang.org/grpc"

	"github.com/tetratelabs/go2sky"
	"github.com/tetratelabs/go2sky/reporter/grpc/register"
)

var (
	errRegister = errors.New("fail to register reporter")
)

// NewGRPCReporter create a new reporter to send data to gRPC oap server
func NewGRPCReporter(serverAddr string, opts ...GRPCReporterOption) (go2sky.Reporter, error) {
	r := &gRPCReporter{
		conn: &grpc.ClientConn{},
		logger: log.New(os.Stderr, "go2sky", log.LstdFlags),
	}
	for _, o := range opts {
		o(r)
	}
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure()) //TODO add TLS
	if err != nil {
		return nil, err
	}
	r.conn = conn
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

type gRPCReporter struct {
	conn       *grpc.ClientConn
	serviceID  int32
	instanceID int32
	logger     *log.Logger
}

func (r *gRPCReporter) Register(service string, instance string) error {
	err := r.registerService(service)
	if err != nil {
		return err
	}
	return r.registerInstance(instance)
}

func (r *gRPCReporter) registerService(name string) error {
	client := register.NewRegisterClient(r.conn)
	in := &register.Services{
		Services: []*register.Service{
			{
				ServiceName: name,
			},
		},
	}
	mapping, err := client.DoServiceRegister(context.Background(), in)
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
	client := register.NewRegisterClient(r.conn)
	in := &register.ServiceInstances{
		Instances: []*register.ServiceInstance{
			{
				ServiceId:    r.serviceID,
				InstanceUUID: name,
			},
		},
	}
	mapping, err := client.DoServiceInstanceRegister(context.Background(), in)
	if err != nil {
		return err
	}
	if len(mapping.ServiceInstances) < 1 {
		return errRegister
	}
	r.instanceID = mapping.ServiceInstances[0].Value
	r.logger.Printf("the id of instance %s 's id is %d", name, r.serviceID)
	return nil
}

func (r *gRPCReporter) Send(spans []go2sky.Span) {

}

func (r *gRPCReporter) Close() {
	err := r.conn.Close()
	if err != nil {
		r.logger.Print(err)
	}
}
