module github.com/SkyAPM/go2sky

go 1.14

require (
	github.com/carvalhorr/protoc-gen-mock v1.5.1 // indirect
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.8.1
	github.com/shirou/gopsutil/v3 v3.22.6
	google.golang.org/grpc v1.49.0
	skywalking.apache.org/repo/goapi v0.0.0-20220824100816-9c0fee7e3581
)

replace skywalking.apache.org/repo/goapi v0.0.0-20220824100816-9c0fee7e3581 => github.com/easonyipj/skywalking-goapi v0.0.0-20220901154300-150416212cda
