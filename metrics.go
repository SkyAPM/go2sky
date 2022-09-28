package go2sky

import (
	"context"
	"github.com/SkyAPM/go2sky/logger"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"log"
	"os"
	"runtime"
	"time"
)

const (
	defaultLogPrefix           = "go2sky-golang-metric"
	InstanceGolangHeap         = "instance_golang_heap"
	InstanceGolangStack        = "instance_golang_stack"
	InstanceGolangGCTime       = "instance_golang_gc_time"
	InstanceGolangGCCount      = "instance_golang_gc_count"
	InstanceGolangThreadNum    = "instance_golang_thread_num"
	InstanceGolangGoroutineNum = "instance_golang_goroutine_num"
	InstanceGolangCPUUsedRate  = "instance_golang_cpu_used_rate"
	InstanceGolangMemUsedRate  = "instance_golang_mem_used_rate"
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
	ctx      context.Context
	reporter MetricsReporter
	instance string
	service  string
	interval time.Duration
	logger   logger.Log
}

func InitMetricCollector(reporter MetricsReporter, interval time.Duration, cancelCtx context.Context) {
	collector := &MetricCollector{
		ctx:      cancelCtx,
		logger:   logger.NewDefaultLogger(log.New(os.Stderr, defaultLogPrefix, log.LstdFlags)),
		reporter: reporter,
		interval: interval,
	}

	go collector.collect()
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

		select {
		case <-c.ctx.Done():
			c.logger.Infof("stop the meter collection")
			return
		default:
			go c.collectMeter()
		}

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

	c.reporter.SendMetrics(runTimeMetric)
}

type MetricsReporter interface {
	SendMetrics(runTimeMeter RunTimeMetric)
}
