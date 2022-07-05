package go2sky

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"testing"
	"time"
)

func TestCpu(t *testing.T) {
	cpuPercent, _ := cpu.Percent(0, false)
	fmt.Println(cpuPercent)
}

func TestMem(t *testing.T) {
	v, _ := mem.VirtualMemory()
	fmt.Println(v)
}

func TestGetCurrentMetrics(t *testing.T) {
	collector := &MetricCollector{
		service:         "service",
		serviceInstance: "serviceInstance",
	}
	for i := 0; i < 5; i++ {
		fmt.Printf("%+v", collector.getCurrentMetrics())
		//time.Sleep(5 * time.Second)
	}
}

func TestInitMetricCollector(t *testing.T) {
	initMetricCollector("service", "serviceInstance")
	for {
		s := ""
		for i := 0; i < 10000; i++ {
			s += "afff"
		}
		time.Sleep(5 * time.Second)
	}
}
