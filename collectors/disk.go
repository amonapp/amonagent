package collectors

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	"github.com/martinrusev/amonagent/logging"
	"github.com/martinrusev/amonagent/metrics"
	"github.com/martinrusev/amonagent/util"
)

type DiskUsageStat struct {
	Path        string  `json:"path"`
	Fstype      string  `json:"fstype"`
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
}

func (d DiskUsageStat) String() string {
	s, _ := json.Marshal(d)
	return string(s)
}

var diskLogger = logging.GetLogger("amonagent.disk")

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
		diskLogger.Errorf("can not read '/proc/filesystems': %v", err)
	}
	return
}

// DiskSpace disk data
func DiskSpace() (metrics.Block, error) {

	disk, _ := exec.Command("df", "-lPT", "--block-size", "1").Output()
	diskString := string(disk)
	diskLines := strings.Split(diskString, "\n")

	var mb metrics.Block
	for _, diskLine := range diskLines {
		fields := strings.Fields(diskLine)

		if len(fields) == 7 {
			var md metrics.MultiDataPoint

			// type DiskUsageStat struct {
			// 	Path        string  `json:"path"`
			// 	Fstype      string  `json:"fstype"`
			// 	Total       uint64  `json:"total"`
			// 	Free        uint64  `json:"free"`
			// 	Used        uint64  `json:"used"`
			// 	UsedPercent float64 `json:"used_percent"`
			// }

			// /dev/mapper/vg0-usr ext4 13384816 9996920 2815784 79% /usr
			fs := fields[0]
			fsType := fields[1]
			spaceTotal := fields[2]
			spaceUsed := fields[3]
			spaceFree := fields[4]
			mount := fields[6]

			if !isPseudoFS(fsType) && !removableFs(fs) {

				// d := DiskUsageStat{
				// 	Path:        fields[0],
				// 	Fstype:      fields[1],
				// 	Total:       fields[2],
				// 	Free:        fields[4],
				// 	Used:        fields[3],
				// 	UsedPercent: fields[3],
				// }

				md.Add("free", spaceFree)
				md.Add("total", spaceTotal)
				md.Add("used", spaceUsed)
				md.Add("path", mount)

				mg := metrics.Group{Name: fs, Metrics: md}
				mb = append(mb, mg)
			}

		}
	}

	return mb, nil

}
