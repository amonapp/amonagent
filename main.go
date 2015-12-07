package main

import (
	"fmt"

	"github.com/martinrusev/amonagent/collectors"
	"github.com/martinrusev/amonagent/logging"
	"github.com/shirou/gopsutil/net"
)

// AmonAgentLogger for the main file
var AmonAgentLogger = logging.GetLogger("amonagent")

// Just for testing
func main() {

	l := collectors.LoadAverage()
	fmt.Println(l)

	n, _ := net.NetIOCounters(true)

	fmt.Println(n)

	ifaces, _ := net.NetInterfaces()
	fmt.Print(ifaces)

}
