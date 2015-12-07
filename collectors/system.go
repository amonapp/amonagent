package collectors

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/martinrusev/amonagent/logging"
	pshost "github.com/shirou/gopsutil/host"
)

// conversion units
const (
	MINUTE = 60
	HOUR   = MINUTE * 60
	DAY    = HOUR * 24
)

// UptimeStruct - returns uptime struct
type UptimeStruct struct {
	Uptime string
}

var systemLogger = logging.GetLogger("amonagent.system")

// Uptime - returns uptime string
// uptime = "{0} days {1} hours {2} minutes".format(days, hours, minutes)
func Uptime() UptimeStruct {
	boot, _ := pshost.BootTime()
	secondsFromBoot := uint64(time.Now().Unix()) - boot

	days := secondsFromBoot / DAY
	hours := (secondsFromBoot % DAY) / HOUR
	minutes := (secondsFromBoot % HOUR) / MINUTE

	s := fmt.Sprintf("%v days %v hours %v minutes", days, hours, minutes)

	uptime := UptimeStruct{
		Uptime: s,
	}
	return uptime
}

// IPAddress - returns machine IP
func IPAddress() string {
	c1, _ := exec.Command("hostname", "-I").Output()
	ipOutput := string(c1)
	ipList := strings.Split(ipOutput, " ")
	if len(ipList) > 0 {
		return ipList[0]
	}
	return ""

}

func (p DistroStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// DistroStruct - returns information about the currently instaled distro
type DistroStruct struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

// Distro - gets distro info
// {'version': '14.04', 'name': 'Ubuntu'}
func Distro() DistroStruct {
	host, _ := pshost.HostInfo()

	d := DistroStruct{
		Version: host.PlatformVersion,
		Name:    host.Platform,
	}

	return d
}
