package collectors

import (
	"encoding/json"
)

func (p SystemDataStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

func (p AllMetricsStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// AllMetricsStruct -XXX
type AllMetricsStruct struct {
	System    SystemDataStruct `json:"system"`
	Processes ProcessesList    `json:"processes"`
}

// SystemDataStruct - collect all system metrics
type SystemDataStruct struct {
	CPU     CPUUsageStruct   `json:"cpu"`
	Network NetworkUsageList `json:"network"`
	Disk    DiskUsageList    `json:"disk"`
	Load    LoadStruct       `json:"loadavg"`
	Uptime  UptimeStruct     `json:"uptime"`
	Memory  MemoryStruct     `json:"memory"`
}

// CollectSystem - XXX
func CollectSystem() AllMetricsStruct {
	networkUsage, _ := NetworkUsage()
	cpuUsage := CPUUsage()
	Load := LoadAverage()
	diskUsage, _ := DiskUsage()
	Uptime := Uptime()
	Memory := MemoryUsage()
	Processes, _ := Processes()

	System := SystemDataStruct{
		CPU:     cpuUsage,
		Network: networkUsage,
		Disk:    diskUsage,
		Load:    Load,
		Uptime:  Uptime,
		Memory:  Memory,
	}

	allMetrics := AllMetricsStruct{
		System:    System,
		Processes: Processes,
	}

	return allMetrics
}
