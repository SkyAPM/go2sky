package http

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/powerapm/go2sky"
	"github.com/powerapm/go2sky/reporter"
)

const (
	//定义上报的服务端地址,即grpc发送的地址
	oap = "172.18.40.193:11800"
	// 定义相关应用服务的名称
	service = "go-http-server"
)

func Example_NewPostJava() {
	// Use gRPC reporter for production
	r, err := reporter.NewGRPCReporter(oap)
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer r.Close()

	tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(r), go2sky.WithInstance("huangyao213"))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}
	tracer.WaitUntilRegister()

	// sm, err := httpPlugin.NewServerMiddleware(tracer)
	// if err != nil {
	// 	log.Fatalf("create server middleware error %v \n", err)
	// }

	client, err := NewClient(tracer)
	if err != nil {
		log.Fatalf("create client error %v \n", err)
	}

	// call end service
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/remoting.call", "http://172.18.40.193:20001/power-adm-server"), nil)
	if err != nil {
		log.Fatalf("unable to create http request: %+v\n", err)
	}
	res, err := client.Do(request)
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}
	_ = res.Body.Close()
	time.Sleep(time.Second)

	// Output:
}
