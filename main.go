package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/martinrusev/amonagent/logging"
	"github.com/martinrusev/amonagent/util"
	"github.com/shirou/gopsutil/disk"
)

// AmonAgentLogger for the main file
var AmonAgentLogger = logging.GetLogger("amonagent")

func (p DiskUsageStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// DiskUsageStruct - individual process data
type DiskUsageStruct struct {
	Name        string  `json:"name"`
	Path        string  `json:"path"`
	Fstype      string  `json:"fstype"`
	Total       float64 `json:"total_mb"`
	Free        float64 `json:"free_mb"`
	Used        float64 `json:"used_mb"`
	UsedPercent float64 `json:"used_percent"`
}

// DiskUsageList - list of individual process data
type DiskUsageList []DiskUsageStruct

var sdiskRE = regexp.MustCompile(`/dev/(sd[a-z])[0-9]?`)

// removableFs checks if the volume is removable
func removableFs(name string) bool {
	s := sdiskRE.FindStringSubmatch(name)
	if len(s) > 1 {
		b, err := ioutil.ReadFile("/sys/block/" + s[1] + "/removable")
		if err != nil {
			return false
		}
		return strings.Trim(string(b), "\n") == "1"
	}
	return false
}

// isPseudoFS checks if it is a valid volume
func isPseudoFS(name string) (res bool) {
	err := util.ReadLine("/proc/filesystems", func(s string) error {
		if strings.Contains(s, name) && strings.Contains(s, "nodev") {
			res = true
			return nil
		}
		return nil
	})
	if err != nil {
		// diskLogger.Errorf("can not read '/proc/filesystems': %v", err)
	}
	return
}

// DiskUsage - return a list with disk usage structs
func DiskUsage() (DiskUsageList, error) {

	parts, err := disk.DiskPartitions(false)
	if err != nil {
		return nil, err
	}

	var usage DiskUsageList

	for _, p := range parts {
		if _, err := os.Stat(p.Mountpoint); err == nil {
			du, err := disk.DiskUsage(p.Mountpoint)
			if err != nil {
				return nil, err
			}

			if !isPseudoFS(du.Fstype) && !removableFs(du.Path) {

				TotalMB, _ := util.ConvertBytesTo(du.Total, "mb", 0)
				FreeMB, _ := util.ConvertBytesTo(du.Free, "mb", 0)
				UsedMB, _ := util.ConvertBytesTo(du.Used, "mb", 0)

				d := DiskUsageStruct{
					Name:        p.Device,
					Path:        du.Path,
					Fstype:      du.Fstype,
					Total:       TotalMB,
					Free:        FreeMB,
					Used:        UsedMB,
					UsedPercent: du.UsedPercent,
				}
				usage = append(usage, d)
			}
		}
	}

	return usage, err
}

// Just for testing
func main() {
	dl, _ := DiskUsage()
	fmt.Println(dl)
}
