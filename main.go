package main

import (
	"fmt"

	"github.com/martinrusev/amonagent/collectors"
	"github.com/martinrusev/amonagent/logging"
	"github.com/martinrusev/amonagent/util"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
)

// AmonAgentLogger for the main file
var AmonAgentLogger = logging.GetLogger("amonagent")

func main() {
	v, _ := mem.VirtualMemory()
	fmt.Println(v)

	n, _ := net.NetIOCounters(true)
	fmt.Println(n)

	memoryTotalMB, _ := util.ToMegabytes(v.Total)
	fmt.Println(memoryTotalMB)

	l, _ := load.LoadAvg()
	fmt.Println(l)

	fmt.Println(v.Total)
	s, _ := host.HostInfo()
	fmt.Println(s)

	d, _ := collectors.DiskSpace()
	fmt.Println(d)
	// for _, volume := range d {
	// 	fmt.Println(volume)
	//
	// }

	collectors.Processes()

}
