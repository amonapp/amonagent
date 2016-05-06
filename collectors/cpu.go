package collectors

import (
	"encoding/json"
	"time"

	"github.com/amonapp/amonagent/logging"
	"github.com/amonapp/amonagent/util"
	"github.com/shirou/gopsutil/cpu"
)

var cpuLogger = logging.GetLogger("amonagent.cpu")

func (p CPUUsageStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// CPUUsageStruct - returns CPU usage stats
type CPUUsageStruct struct {
	User   float64 `json:"user"`
	Idle   float64 `json:"idle"`
	Nice   float64 `json:"nice"`
	Steal  float64 `json:"steal"`
	System float64 `json:"system"`
	IOWait float64 `json:"iowait"`
}

// totalCpuTime - XXX
func totalCPUTime(t cpu.TimesStat) float64 {
	total := t.User + t.System + t.Nice + t.Iowait + t.Irq + t.Softirq + t.Steal +
		t.Guest + t.GuestNice + t.Idle
	return total
}

// CPUUsage - return a map with CPU usage stats
func CPUUsage() CPUUsageStruct {
	cpuTimes1, err := cpu.Times(false)
	if err != nil {
		cpuLogger.Errorf("error getting CPU info: %s", err)
	}

	c := CPUUsageStruct{}
	for i, lastCts := range cpuTimes1 {
		lastTotal := totalCPUTime(lastCts)
		time.Sleep(1 * time.Second)

		cpuTimes2, _ := cpu.Times(false)
		cts := cpuTimes2[i]
		total := totalCPUTime(cts)

		totalDelta := total - lastTotal

		system := 100 * (cts.System - lastCts.System) / totalDelta
		nice := 100 * (cts.Nice - lastCts.Nice) / totalDelta
		user := 100 * (cts.User - lastCts.User) / totalDelta
		idle := 100 * (cts.Idle - lastCts.Idle) / totalDelta
		iowait := 100 * (cts.Iowait - lastCts.Iowait) / totalDelta
		steal := 100 * (cts.Steal - lastCts.Steal) / totalDelta

		systemPercent, _ := util.FloatDecimalPoint(system, 2)
		nicePercent, _ := util.FloatDecimalPoint(nice, 2)
		userPercent, _ := util.FloatDecimalPoint(user, 2)
		idlePercent, _ := util.FloatDecimalPoint(idle, 2)
		iowaitPercent, _ := util.FloatDecimalPoint(iowait, 2)
		stealPercent, _ := util.FloatDecimalPoint(steal, 2)

		c = CPUUsageStruct{
			User:   userPercent,
			Idle:   idlePercent,
			Nice:   nicePercent,
			Steal:  stealPercent,
			System: systemPercent,
			IOWait: iowaitPercent,
		}

	}

	return c
}
