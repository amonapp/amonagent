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

	cpu := collectors.CPUUsage()
	fmt.Println(cpu)
	// if err != nil {
	// 	return fmt.Errorf("error getting CPU info: %s", err)
	// }
}
