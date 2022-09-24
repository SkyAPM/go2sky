package go2sky

import (
	"github.com/SkyAPM/go2sky/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"os"
	"runtime"
	"time"

	"skywalking.apache.org/repo/goapi/collect/language/agent/v3"
)

const (
	maxSendQueueSize           int32 = 30000
	defaultLogPrefix                 = "go2sky-golang-metric"
	InstanceGolangHeap               = "instance_golang_heap"
	InstanceGolangStack              = "instance_golang_stack"
	InstanceGolangGCTime             = "instance_golang_gc_time"
	InstanceGolangGCCount            = "instance_golang_gc_count"
	InstanceGolangThreadNum          = "instance_golang_thread_num"
	InstanceGolangGoroutineNum       = "instance_golang_goroutine_num"
	InstanceGolangCPUUsedRate        = "instance_golang_cpu_used_rate"
	InstanceGolangMemUsedRate        = "instance_golang_mem_used_rate"
)

type RunTimeMetric struct {
	// the Unix time when metrics were collected
	Time int64
	// the bytes of allocated heap objects
	HeapAlloc int64
	// the bytes in stack spans.
	StackInUse int64
	// the number of completed GC cycles since instance started
	GCCount int64
	// the latest gc pause time(NS)
	GCPauseTime int64
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
	instance string
	service  string
	interval time.Duration
	logger   logger.Log
}

func InitMetricCollector(reporter MetricsReporter, instance, service string, interval time.Duration) {
	collector := &MetricCollector{
		sendCh:   make(chan RunTimeMetric, maxSendQueueSize),
		logger:   logger.NewDefaultLogger(log.New(os.Stderr, defaultLogPrefix, log.LstdFlags)),
		reporter: reporter,
		instance: instance,
		service:  service,
		interval: interval,
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

	timer := time.NewTicker(c.interval)

	for {
		go c.collectMeter()
		<-timer.C
	}
}

func (c *MetricCollector) collectMeter() {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	v, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(0, false)
	threadNum, _ := runtime.ThreadCreateProfile(nil)
	runTimeMetric := RunTimeMetric{
		Time:         time.Now().UnixMilli(),
		HeapAlloc:    int64(rtm.HeapAlloc),
		StackInUse:   int64(rtm.StackInuse),
		GCCount:      int64(rtm.NumGC),
		GCPauseTime:  int64(rtm.PauseNs[(rtm.NumGC+255)%256]),
		GoroutineNum: int64(runtime.NumGoroutine()),
		ThreadNum:    int64(threadNum),
		CpuUsedRate:  cpuPercent[0],
		MemUsedRate:  v.UsedPercent,
	}

	select {
	case c.sendCh <- runTimeMetric:
	default:
		c.logger.Errorf("reach max send buffer")
	}
}

func (c *MetricCollector) send() {

	for m := range c.sendCh {
		meterDataList := make([]*v3.MeterData, 0)
		meterDataList = append(meterDataList, c.generateMeter(InstanceGolangHeap, float64(m.HeapAlloc), m.Time))
		meterDataList = append(meterDataList, c.generateMeter(InstanceGolangStack, float64(m.StackInUse), m.Time))
		meterDataList = append(meterDataList, c.generateMeter(InstanceGolangGCTime, float64(m.GCPauseTime), m.Time))
		meterDataList = append(meterDataList, c.generateMeter(InstanceGolangGCCount, float64(m.GCCount), m.Time))
		meterDataList = append(meterDataList, c.generateMeter(InstanceGolangThreadNum, float64(m.ThreadNum), m.Time))
		meterDataList = append(meterDataList, c.generateMeter(InstanceGolangGoroutineNum, float64(m.GoroutineNum), m.Time))
		meterDataList = append(meterDataList, c.generateMeter(InstanceGolangCPUUsedRate, m.CpuUsedRate, m.Time))
		meterDataList = append(meterDataList, c.generateMeter(InstanceGolangMemUsedRate, m.MemUsedRate, m.Time))
		c.reporter.SendMetrics(meterDataList)
	}

}

func (c *MetricCollector) generateMeter(name string, value float64, time int64) *v3.MeterData {
	return &v3.MeterData{
		Metric: &v3.MeterData_SingleValue{
			SingleValue: &v3.MeterSingleValue{
				Name:  name,
				Value: value,
			},
		},
		Timestamp:       time,
		Service:         c.service,
		ServiceInstance: c.instance,
	}
}

type MetricsReporter interface {
	SendMetrics(metrics []*v3.MeterData)
}
