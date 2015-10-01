package collectors

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"

	"github.com/martinrusev/amonagent/logging"
	"github.com/martinrusev/amonagent/util"
)

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
func DiskSpace() {

	disk, _ := exec.Command("df", "-lPT", "--block-size", "1").Output()
	diskString := string(disk)
	diskLines := strings.Split(diskString, "\n")

	for _, diskLine := range diskLines {
		fields := strings.Fields(diskLine)

		if len(fields) == 7 {
			// fmt.Println(fields)

			// /dev/mapper/vg0-usr ext4 13384816 9996920 2815784 79% /usr
			fs := fields[0]
			fsType := fields[1]
			spaceTotal := fields[2]
			spaceUsed := fields[3]
			spaceFree := fields[4]
			mount := fields[6]

			if !isPseudoFS(fsType) && !removableFs(fs) {
				fmt.Println(fs, fsType, spaceUsed, spaceTotal, spaceFree, mount)
			}

		}
	}

}
