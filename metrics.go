package go2sky

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"runtime"
	"time"
)

type RunTimeMetric struct {
	// the Unix time when metrics were collected
	time int64
	// the service name that user input
	service string
	// the service instance id that generated automatic
	serviceInstance string
	// the bytes of allocated heap objects
	heapAlloc int64
	// the bytes in stack spans.
	stackInUse int64
	// the number of completed GC cycles since
	gcNum int64
	// the latest gc pause time
	gcPauseTime int64
	// the number of goroutines that currently exist
	goroutineNum int64
	// the number of records in the thread creation profile
	threadNum int64
	// the cpu Used float64
	cpuUsedRate float64
	// the Percentage of RAM used by programs
	memUsedRate float64
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
		c.sendCh <- c.getCurrentMetrics()
		time.Sleep(5 * time.Second)
	}
}

func (c *MetricCollector) send() {

	for m := range c.sendCh {
		fmt.Println(fmt.Sprintf("%+v", m))
	}

}

func (c *MetricCollector) getCurrentMetrics() *RunTimeMetric {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	v, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(0, false)
	threadNum, _ := runtime.ThreadCreateProfile(nil)
	return &RunTimeMetric{
		time:            time.Now().Unix(),
		service:         c.service,
		serviceInstance: c.serviceInstance,
		heapAlloc:       int64(rtm.HeapAlloc),
		stackInUse:      int64(rtm.StackInuse),
		gcNum:           int64(rtm.NumGC),
		gcPauseTime:     int64(rtm.PauseNs[(rtm.NumGC+255)%256]),
		goroutineNum:    int64(runtime.NumGoroutine()),
		threadNum:       int64(threadNum),
		cpuUsedRate:     cpuPercent[0],
		memUsedRate:     v.UsedPercent,
	}
}
