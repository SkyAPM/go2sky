package go2sky

import (
	"runtime"
	"time"
)

type RunTimeMetric struct {
	time            int64
	service         string
	serviceInstance string
	heapAlloc       int64
	stackAlloc      int64
	gcNum           int64
	gcPauseTime     int64
	goroutineNum    int64
	threadNum       int64
	cpuUsedRate     int64
	memUsedRate     int64
}

type MetricCollector struct {
	service         string
	serviceInstance string
	sendCh          chan *RunTimeMetric
}

func initMetricCollector(service, serviceInstance string) {
	collector := &MetricCollector{
		service:         service,
		serviceInstance: serviceInstance,
		sendCh:          make(chan *RunTimeMetric, 1000),
	}

	go collector.collect()
	go collector.send()
}

func (c *MetricCollector) collect() {
	for {

		var rtm runtime.MemStats
		runtime.ReadMemStats(&rtm)
		r := &RunTimeMetric{
			time:            time.Now().Unix(),
			service:         c.service,
			serviceInstance: c.serviceInstance,
			heapAlloc:       int64(rtm.HeapAlloc),
			stackAlloc:      0,
			gcNum:           0,
			gcPauseTime:     0,
			goroutineNum:    0,
			threadNum:       0,
			cpuUsedRate:     0,
			memUsedRate:     0,
		}

		c.sendCh <- r
		time.Sleep(1 * time.Second)
	}
}

func (c *MetricCollector) send() {

}
