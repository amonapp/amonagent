package main

import (
	"fmt"

	"github.com/martinrusev/amonagent/collectors"
	"github.com/martinrusev/amonagent/logging"
)

// AmonAgentLogger for the main file
var AmonAgentLogger = logging.GetLogger("amonagent")

// Just for testing
func main() {
	dl, _ := collectors.DiskUsage()
	fmt.Println(dl)

	p, _ := collectors.Processes()
	fmt.Println(p)

}
