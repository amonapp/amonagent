package collectors

import (
	"fmt"

	"github.com/martinrusev/amonagent/logging"
)

var diskLogger = logging.GetLogger("amonagent.disk")

func main() {
	fmt.Println(diskspace_windows())
	fmt.Println(physical_disk_windows())
	// collectors = append(collectors, &IntervalCollector{F: c_physical_disk_windows})
	// collectors = append(collectors, &IntervalCollector{F: c_diskspace_windows})
}
