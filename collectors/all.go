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

func (p HostDataStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// AllMetricsStruct -XXX
type AllMetricsStruct struct {
	System    SystemDataStruct `json:"system"`
	Processes ProcessesList    `json:"processes"`
	Host      HostDataStruct   `json:"host"`
}

// HostDataStruct -XXX
type HostDataStruct struct {
	Host      string `json:"host"`
	MachineID string `json:"machineid"`
}

// SystemDataStruct - collect all system metrics
type SystemDataStruct struct {
	CPU     CPUUsageStruct   `json:"cpu"`
	Network NetworkUsageList `json:"network"`
	Disk    DiskUsageList    `json:"disk"`
	Load    LoadStruct       `json:"loadavg"`
	Uptime  string           `json:"uptime"`
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

	host := Host()
	machineID := MachineID()
	hoststruct := HostDataStruct{
		Host:      host,
		MachineID: machineID,
	}
	allMetrics := AllMetricsStruct{
		System:    System,
		Processes: Processes,
		Host:      hoststruct,
	}

	return allMetrics
}
