# GO2Sky

[![Build](https://github.com/SkyAPM/go2sky/workflows/Build/badge.svg?branch=master)](https://github.com/SkyAPM/go2sky/actions?query=branch%3Amaster+event%3Apush+workflow%3ABuild)
[![Coverage](https://codecov.io/gh/SkyAPM/go2sky/branch/master/graph/badge.svg)](https://codecov.io/gh/SkyAPM/go2sky)
[![GoDoc](https://godoc.org/github.com/SkyAPM/go2sky?status.svg)](https://godoc.org/github.com/SkyAPM/go2sky)


**GO2Sky** is an instrument SDK library, written in Go, by following [Apache SkyWalking](https://github.com/apache/incubator-skywalking) tracing and metrics formats.

# Installation
```
$ go get -u github.com/SkyAPM/go2sky
```

The API of this project is still evolving. The use of vendoring tool is recommended.

# Quickstart

By completing this quickstart, you will learn how to trace local methods. For more details, please view
[the example](example_trace_test.go).

## Configuration

GO2Sky can export traces to Apache SkyWalking OAP server or local logger. In the following example, we configure GO2Sky to export to OAP server,
which is listening on `oap-skywalking` port `11800`, and all the spans from this program will be associated with a service name `example`.
`reporter.GRPCReporter` can also adjust the behavior through `reporter.GRPCReporterOption`, [view all](docs/GRPC-Reporter-Option.md).

```go
r, err := reporter.NewGRPCReporter("oap-skywalking:11800")
if err != nil {
    log.Fatalf("new reporter error %v \n", err)
}
defer r.Close()
tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(r))
```
In some scenarios, we may need a filter to filter segments that do not need to be submitted, for example, to reduce the load of gRPC reporting, or only track the request of error.

```go
r, err := reporter.NewGRPCReporter("oap-skywalking:11800", reporter.WithReportStrategy(func(s *v3.SegmentObject) bool {
	var isReport bool
	for _, span := s.GetSpans() {
		if span.GetIsError() {
			isReport = true
			break
		}
	}
	
	return isReport
}))

```

You can also create tracer with sampling rate. It supports decimals between **0-1** (two decimal places), representing the sampling percentage of trace.
```go
tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(r), go2sky.WithSampler(0.5))
```

Also could customize correlation context config.
```go
tracer, err := go2sky.NewTracer("example", go2sky.WithReporter(r), go2sky.WithSampler(0.5), go2sky.WithCorrelation(3, 128))
```

## Create span

To create a span in a trace, we used the `Tracer` to start a new span. We indicate this as the root span because of
passing `context.Background()`. We must also be sure to end this span, which will be show in [End span](#end-span).

```go
span, ctx, err := tracer.CreateLocalSpan(context.Background())
```

## Create a sub span

A sub span created as the children of root span links to its parent with `Context`.

```go
subSpan, newCtx, err := tracer.CreateLocalSpan(ctx)
```

## Get correlation

Get custom data from tracing context.

```go
value := go2sky.GetCorrelation(ctx, key)
```

## Put correlation

Put custom data to tracing context.

```go
success := go2sky.PutCorrelation(ctx, key, value)
```

## End span

We must end the spans so they becomes available for sending to the backend by a reporter.

```go
subSpan.End()
span.End()
```

## Global Tracer

Set and get global Tracer

```go
// new tracer
tr, err := go2sky.NewTracer("example")

// registers `tracer` as the global Tracer
go2sky.SetGlobalTracer(tr)

// returns the registered global Tracer
// if none is registered then an instance of `nil` is returned
tracer := go2sky.GetGlobalTracer()
```

## Get Active Span

Get the `activeSpan` in the `Context`.

```go
go2sky.ActiveSpan(ctx)
```

## With Span

Save the `activeSpan` to `Context`

```go
go2sky.WithSpan(ctx, activeSpan)
```

## Get Global Service Name

Get the `ServiceName` of the `activeSpan` in the `Context`.

```go
go2sky.ServiceName(ctx)
```

## Get Global Service Instance Name

Get the `ServiceInstanceName` of the `activeSpan` in the `Context`.

```go
go2sky.ServiceInstanceName(ctx)
```

## Get Global TraceID

Get the `TraceID` of the `activeSpan` in the `Context`.

```go
go2sky.TraceID(ctx)
```

## Get Global TraceSegmentID

Get the `TraceSegmentID` of the `activeSpan` in the `Context`.

```go
go2sky.TraceSegmentID(ctx)
```

## Get Global SpanID

Get the `SpanID` of the `activeSpan` in the `Context`.

```go
go2sky.SpanID(ctx)
```

## Application Log  Report

```go

import (
	"context"
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
	"log"
)

func SkyapmLogTest(ctx context.Context)  {

	r, err := reporter.NewGRPCReporter("oap-skywalking:11800")
	if err != nil {
		log.Fatalf("new reporter error %v \n", err)
	}
	defer r.Close()

	skyapmLogger,skyapmError:=go2sky.NewLogger(r)
	if skyapmError!=nil{
		log.Fatalf("new Logger error %v \n", skyapmError)
	}
	
	logData:="your application log need to send to backend here..."
	
	skyapmLogger.WriteLogWithContext(ctx,go2sky.LogLevelError,logData)
}


```

## Periodically Report
Go2sky agent reports the segments periodically.
It would not wait for all finished segments reported when the service exits.


# Advanced Concepts

We cover some advanced topics about GO2Sky.

## Context propagation

Trace links spans belong to it by using context propagation which varies between different scenario.

### In process

We use `context` package to link spans. The root span usually pick `context.Background()`, and sub spans
will inject the context generated by its parent.

```go
//Create a new context
entrySpan, entryCtx, err := tracer.CreateEntrySpan(context.Background(), ...)

// Some operation
...

// Link two spans by injecting entrySpan context into exitSpan
exitSpan, err := tracer.CreateExitSpan(entryCtx, ...)
```

### Crossing process

We use `Entry` span to extract context from downstream service, and use `Exit` span to inject context to
upstream service.

`Entry` and `Exit` spans make sense to OAP analysis which generates topology map and service metrics.

```go
//Extract context from HTTP request header `sw8`
span, ctx, err := tracer.CreateEntrySpan(r.Context(), "/api/login", func(key string) (string, error) {
		return r.Header.Get(key), nil
})

// Some operation
...

// Inject context into HTTP request header `sw8`
span, err := tracer.CreateExitSpan(req.Context(), "/service/validate", "tomcat-service:8080", func(key, value string) error {
		req.Header.Set(key, value)
		return nil
})
```

## Tag

We set tags into a span which is stored in the backend, but some tags have special purpose. OAP server
may use them to aggregate metrics, generate topology map and etc.

They are defined as constant in root package with prefix `Tag`.

## Log x Trace context

Inject trace context into the log text. SkyWalking LAL(log analysis language) engine could extract the context from the text and correlate trace and logs.

```go
// Get trace context data
import go2skylog "github.com/SkyAPM/go2sky/log"
logContext = go2skylog.FromContext(ctx)

// Build context data string
// Inject context string into log
// Context format string: [$serviceName,$instanceName,$traceId,$traceSegmentId,$spanId]
contextString := logContext.String()
```

## Plugins

Go to go2sky-plugins repo to see all the plugins, [click here](https://github.com/SkyAPM/go2sky-plugins).

## Supported Environment Variables

Below is the full list of supported environment variables you can set to customize the agent behavior, please read the descriptions for what they can achieve.

|                  Environment Variable                  |                                                                                                                                               Description                                                                                                                                               |      Default       |
| :----------------------------------------------------: | :-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------: | :----------------: |
|                    `SW_AGENT_NAME`                     |                                                                                                                                       The name of the Go service                                                                                                                                        |       unset        |
|                    `SW_AGENT_LAYER`                    |                                                                                                                         Instance belong layer name which define in the backend                                                                                                                          |       unset        |
|                `SW_AGENT_INSTANCE_NAME`                |                                                                                                                                   The name of the Go service instance                                                                                                                                   | Randomly generated |
|                   `SW_AGENT_SAMPLE`                    |                                                                                                                           Sample rate, 1 setting it to 1 means full sampling                                                                                                                            |         1          |
|         `SW_AGENT_COLLECTOR_BACKEND_SERVICES`          |                                                                                                                                     The backend OAP server address                                                                                                                                      |       unset        |
|               `SW_AGENT_AUTHENTICATION`                | The authentication token to verify that the agent is trusted by the backend OAP, as for how to configure the backend, refer to [the yaml](https://github.com/apache/skywalking/blob/4f0f39ffccdc9b41049903cc540b8904f7c9728e/oap-server/server-bootstrap/src/main/resources/application.yml#L155-L158). |       unset        |
|         `SW_AGENT_COLLECTOR_HEARTBEAT_PERIOD`          |                                                                                                                               Agent heartbeat report period. Unit, second                                                                                                                               |         20         |
| `SW_AGENT_COLLECTOR_GET_AGENT_DYNAMIC_CONFIG_INTERVAL` |                                                                                                                         Sniffer get agent dynamic config interval. Unit, second                                                                                                                         |         20         |
|        `SW_AGENT_COLLECTOR_MAX_SEND_QUEUE_SIZE`        |                                                                                                                                      Send span queue buffer length                                                                                                                                      |       30000        |
|         `SW_AGENT_PROCESS_STATUS_HOOK_ENABLE`          |                                                                                                                                 Enable the Process Status Hook feature                                                                                                                                  |       false        |
|               `SW_AGENT_PROCESS_LABELS`                |                                                                                                                       The labels of the process, multiple labels split by ","                                                                                                                           |       unset        |


## CDS - Configuration Discovery Service

Configuration Discovery Service provides the dynamic configuration for the agent, defined in gRPC and stored in the backend.

### Available key(s) and value(s) in Golang Agent.
Golang agent supports the following dynamic configurations.

|    Config Key     |                                    Value Description                                     | Value Format Example |
|:-----------------:|:----------------------------------------------------------------------------------------:|:--------------------:|
| agent.sample_rate | The percentage of trace when sampling. It's `[0, 1]`, Same with `WithSampler` parameter. |         0.1          |

## Process Status Hook

This feature is used in cooperation with the [skywalking-rover](https://github.com/apache/skywalking-rover) project.

When go2sky keeps alive with the backend, it would write a metadata file to the local (temporary directory) at the same time, which describes the information of the current process.
The rover side scans all processes, find out which process contains this metadata file. Finally, the rover could collect, profiling, with this process.

### Metadata File

The metadata file use to save metadata with current process, it save in: `{TMPDIR}/apache_skywalking/process/{pid}/metadata.properties`.

Also, when the go2sky keep alive with backend, modify and open time of the metadata file would be updated.

| Key | Type | Description |
|-----|------|------------|
|layer|string|this process layer.|
|service_name|string|this process service name.|
|instance_name|string|this process instance name.|
|process_name|string|this process process name, it's same with the instance name.|
|properties|json|the properties in instance, the process labels also in the properties value.|
|labels|string|the process labels, multiple labels split by ",".|
|language|string|current process language, which is `golang`.|

Please read the [official documentation of rover](https://skywalking.apache.org/docs/skywalking-rover/latest/en/setup/configuration/process_discovery/scanner/#agent-sensor-mode) to get more information.

# License
Apache License 2.0. See [LICENSE](LICENSE) file for details.
