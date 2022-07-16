package go2sky

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"testing"
)

func TestCpu(t *testing.T) {
	cpuPercent, _ := cpu.Percent(0, false)
	fmt.Println(cpuPercent)
}

func TestMem(t *testing.T) {
	v, _ := mem.VirtualMemory()
	fmt.Println(v)
}
