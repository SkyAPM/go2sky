module github.com/SkyAPM/go2sky

go 1.14

require (
	github.com/golang/mock v1.5.0
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/uuid v1.1.2
	github.com/pkg/errors v0.8.1
	github.com/shirou/gopsutil/v3 v3.22.6
	google.golang.org/grpc v1.48.0
	skywalking.apache.org/repo/goapi v0.0.0-20220714130828-0d56d1f4c592
)

replace skywalking.apache.org/repo/goapi v0.0.0-20220714130828-0d56d1f4c592 => github.com/easonyipj/skywalking-goapi v0.0.0-20220823085224-de73735204d1
