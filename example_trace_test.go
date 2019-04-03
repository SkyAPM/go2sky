package go2sky_test

import (
	"context"
	"fmt"
	"github.com/tetratelabs/go2sky"
	"github.com/tetratelabs/go2sky/reporter"
	"log"
	"time"
)

func ExampleNewTracer() {
	r, err := reporter.NewGRPCReporter("localhost:11800")
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer r.Close()
	tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(r))
	if err != nil {
		log.Fatalf("create tracer error %v \n", err)
	}
	span, ctx, err := tracer.CreateLocalSpan(context.Background())
	if err != nil {
		log.Fatalf("create new local span error %v \n", err)
	}
	span.SetOperationName("invoke data")
	span.Tag("kind", "outer")
	time.Sleep(2 * time.Second)
	subSpan, _, err := tracer.CreateLocalSpan(ctx)
	if err != nil {
		log.Fatalf("create new sub local span error %v \n", err)
	}
	subSpan.SetOperationName("invoke inner")
	subSpan.Log(time.Now(), "inner", "this is right")
	time.Sleep(2 * time.Second)
	subSpan.End()
	time.Sleep(1 * time.Second)
	span.End()
	time.Sleep(time.Minute)
	fmt.Print("aa")
	// Output: aa
}