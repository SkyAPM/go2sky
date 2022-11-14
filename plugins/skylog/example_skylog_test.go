package skylog

import (
	"context"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"log"
	"testing"
)

var (
	skyapmOAPAddr = "oap-skywalking:11800"
)

func TestSkyapmLog(t *testing.T) {

	//if with gin.Context,ctx=ginContext.Request.Context(),then we can log with the trace
	ctx := context.Background()

	r, err := reporter.NewGRPCReporter(skyapmOAPAddr)
	if err != nil {
		log.Printf("new rpc reporter error %v \n", err)
	}

	r, err = reporter.NewLogReporter()
	if err != nil {
		log.Fatalf("new log reporter error %v \n", err)
	}

	defer func() {

		if r != nil {
			r.Close()
		}
	}()

	skyapmLogger, skyapmError := go2sky.NewSkyLogger(r)
	if skyapmError != nil {
		log.Fatalf("new SkyLogger error %v \n", skyapmError)
	}

	logData := "your application log need to send to backend here..."

	skyapmLogger.WriteLogWithContext(ctx, go2sky.LogLevelError, logData)
}
