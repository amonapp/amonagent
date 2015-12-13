package collectors

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	pshost "github.com/shirou/gopsutil/host"
)

// conversion units
const (
	MINUTE = 60
	HOUR   = MINUTE * 60
	DAY    = HOUR * 24
)

// Uptime - returns uptime string
// uptime = "{0} days {1} hours {2} minutes".format(days, hours, minutes)
func Uptime() string {
	boot, _ := pshost.BootTime()
	secondsFromBoot := uint64(time.Now().Unix()) - boot

	days := secondsFromBoot / DAY
	hours := (secondsFromBoot % DAY) / HOUR
	minutes := (secondsFromBoot % HOUR) / MINUTE

	s := fmt.Sprintf("%v days %v hours %v minutes", days, hours, minutes)

	return s
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

//MachineID - XXX
func MachineID() string {
	var machineidPath = "/var/lib/dbus/machine-id" // Default machine id path
	var MachineID string
	if _, err := os.Stat(machineidPath); os.IsNotExist(err) {
		machineidPath = "/etc/machine-id"
		// Does not exists, probably an older distro or docker container
		if _, err := os.Stat(machineidPath); os.IsNotExist(err) {
			machineidPath = ""
		}
	}

	if len(machineidPath) > 0 {
		file, err := os.Open(machineidPath)
		if err != nil {
			fmt.Printf(err.Error())
		}
		defer file.Close()
		scanner := bufio.NewScanner(file)
		var lines []string
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		if len(lines) > 0 {
			MachineID = lines[0]
		}
	}

	// Can't detect, return an empty string and ask for a server key
	if len(MachineID) != 32 {
		MachineID = ""
	}

	return MachineID
}

// Host - XXX
func Host() string {
	host, _ := os.Hostname()

	return host
}
