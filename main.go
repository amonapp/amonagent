package main

import (
	"fmt"

	"github.com/martinrusev/amonagent/collectors"
	"github.com/martinrusev/amonagent/logging"
	"github.com/martinrusev/amonagent/util"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// AmonAgentLogger for the main file
var AmonAgentLogger = logging.GetLogger("amonagent")

func main() {
	v, _ := mem.VirtualMemory()
	fmt.Println(v)

	memoryTotalMB, _ := util.ToMegabytes(v.Total)
	fmt.Println(memoryTotalMB)

	fmt.Println(v.Total)
	s, _ := host.HostInfo()
	fmt.Println(s)

	collectors.DiskSpace()
	collectors.Processes()

}
