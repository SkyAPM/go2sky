package go2sky

import (
	"github.com/SkyAPM/go2sky/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	maxSendQueueSize             int32 = 30000
	defaultGolangCollectInterval       = 5 * time.Second
	defaultLogPrefix                   = "go2sky-golang-metric"
)

type RunTimeMetric struct {
	// the Unix time when metrics were collected
	Time int64
	// the bytes of allocated heap objects
	HeapAlloc int64
	// the bytes in stack spans.
	StackInUse int64
	// the number of completed GC cycles since instance started
	GcNum int64
	// the latest gc pause time(NS)
	GcPauseTime float64
	// the number of goroutines that currently exist
	GoroutineNum int64
	// the number of records in the thread creation profile
	ThreadNum int64
	// the cpu Used float64
	CpuUsedRate float64
	// the Percentage of RAM used by programs
	MemUsedRate float64
}

type MetricCollector struct {
	sendCh   chan RunTimeMetric
	reporter MetricsReporter
	logger   logger.Log
}

func InitMetricCollector(reporter MetricsReporter) {
	collector := &MetricCollector{
		sendCh:   make(chan RunTimeMetric, maxSendQueueSize),
		logger:   logger.NewDefaultLogger(log.New(os.Stderr, defaultLogPrefix, log.LstdFlags)),
		reporter: reporter,
	}

	go collector.collect()
	go collector.send()
}

func (c *MetricCollector) collect() {
	defer func() {
		// recover the panic caused by close sendCh
		if err := recover(); err != nil {
			c.logger.Errorf("collect metric err %v", err)
		}
	}()

	var lastGCNum = uint32(0)

	for {

		var rtm runtime.MemStats
		runtime.ReadMemStats(&rtm)
		v, _ := mem.VirtualMemory()
		cpuPercent, _ := cpu.Percent(0, false)
		threadNum, _ := runtime.ThreadCreateProfile(nil)
		runTimeMetric := RunTimeMetric{
			Time:       time.Now().UnixMilli(),
			HeapAlloc:  int64(rtm.HeapAlloc),
			StackInUse: int64(rtm.StackInuse),
			GcNum:      int64(rtm.NumGC - lastGCNum),
			// transfer ns to ms
			GcPauseTime:  float64(rtm.PauseNs[(rtm.NumGC+255)%256]) / float64(1000000),
			GoroutineNum: int64(runtime.NumGoroutine()),
			ThreadNum:    int64(threadNum),
			CpuUsedRate:  cpuPercent[0],
			MemUsedRate:  v.UsedPercent,
		}

		lastGCNum = rtm.NumGC

		select {
		case c.sendCh <- runTimeMetric:
		default:
			c.logger.Errorf("reach max send buffer")
		}
		c.logger.Infof("%+v", runTimeMetric)

		time.Sleep(defaultGolangCollectInterval)
	}
}

func (c *MetricCollector) send() {

	for m := range c.sendCh {
		c.reporter.SendMetrics(m)
	}

}

type MetricsReporter interface {
	SendMetrics(metrics RunTimeMetric)
}
