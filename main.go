package main

import (
	"encoding/json"
	"fmt"

	"github.com/martinrusev/amonagent/collectors"
	"github.com/martinrusev/amonagent/logging"
)

// AmonAgentLogger for the main file
var AmonAgentLogger = logging.GetLogger("amonagent")

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
	System    SystemDataStruct         `json:"system"`
	Processes collectors.ProcessesList `json:"processes"`
}

// SystemDataStruct - collect all system metrics
type SystemDataStruct struct {
	CPU     collectors.CPUUsageStruct   `json:"cpu"`
	Network collectors.NetworkUsageList `json:"network"`
	Disk    collectors.DiskUsageList    `json:"disk"`
	Load    collectors.LoadStruct       `json:"load"`
	Uptime  collectors.UptimeStruct     `json:"uptime"`
	Memory  collectors.MemoryStruct     `json:"memory"`
}

// Just for testing
func main() {

	networkUsage, _ := collectors.NetworkUsage()
	cpuUsage := collectors.CPUUsage()
	Load := collectors.LoadAverage()
	diskUsage, _ := collectors.DiskUsage()
	Uptime := collectors.Uptime()
	Memory := collectors.MemoryUsage()
	Processes, _ := collectors.Processes()

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

	fmt.Print(allMetrics)

}
