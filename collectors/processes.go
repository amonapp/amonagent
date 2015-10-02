package collectors

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/martinrusev/amonagent/logging"
)

var processLogger = logging.GetLogger("amonagent.processes")

// Processes get process usage
func Processes() {
	c1, _ := exec.Command("pidstat", "-ruhtd").Output()

	// Find header and ignore
	headerRegex, _ := regexp.Compile("d+")

	pidstatOutput := string(c1)
	pidstatLines := strings.Split(pidstatOutput, "\n")

	for _, processLine := range pidstatLines {

		if len(headerRegex.FindString(processLine)) == 0 {
			processData := strings.Fields(processLine)

			// Helper
			// Time(0)   UID(1)      TGID(2)       TID(3)
			// %usr{4} %system{5}  %guest{6}    %CPU{7}   CPU{8}
			// minflt/s{9}  majflt/s{10}     VSZ{11}    RSS{12}
			// %MEM{13}   kB_rd/s{14}   kB_wr/s{15} kB_ccwr/s{16}  Command{17}
			if len(processData) == 18 {
				// pid, masterthreadID, cpuPercent, memoryPercent, processName := processData[2], processData[3], processData[7], processData[13], processData[17]
				//
				// masterthreadIDtoINT, _ := strconv.Atoi(masterthreadID)
				//
				// if masterthreadIDtoINT == 0 {
				// 	fmt.Println(pid, masterthreadID, cpuPercent, memoryPercent, processName)
				//
				// }

			}

		}

	}
}
