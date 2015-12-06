package collectors

import (
	"encoding/json"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/martinrusev/amonagent/logging"
	"github.com/martinrusev/amonagent/util"
	"github.com/shirou/gopsutil/mem"
)

var processLogger = logging.GetLogger("amonagent.processes")

func (p ProcessStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// ProcessStruct - individual process data
type ProcessStruct struct {
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory_mb"`
	KBRead  float64 `json:"kb_read"`
	KBWrite float64 `json:"kb_write"`
	Command string  `json:"command"`
}

// ProcessesList - list of individual process data
type ProcessesList []ProcessStruct

// Processes - get data from sysstat, format and return the result
func Processes() ProcessesList {
	c1, _ := exec.Command("pidstat", "-ruhtd").Output()

	var ps ProcessesList
	v, _ := mem.VirtualMemory()
	memoryTotalMB, _ := util.ConvertBytesTo(float64(v.Total), "mb")

	// Find header and ignore
	headerRegex, _ := regexp.Compile("d+")

	pidstatOutput := string(c1)
	pidstatLines := strings.Split(pidstatOutput, "\n")
	// var err error
	for _, processLine := range pidstatLines {

		if len(headerRegex.FindString(processLine)) == 0 {
			processData := strings.Fields(processLine)

			// Helper
			// Time(0)   UID(1)      TGID(2)       TID(3)
			// %usr{4} %system{5}  %guest{6}    %CPU{7}   CPU{8}
			// minflt/s{9}  majflt/s{10}     VSZ{11}    RSS{12}
			// %MEM{13}   kB_rd/s{14}   kB_wr/s{15} kB_ccwr/s{16}  Command{17}
			if len(processData) == 18 {
				pid, masterthreadID, cpuPercent, memoryPercent, processName := processData[2], processData[3], processData[7], processData[13], processData[17]

				masterthreadIDtoINT, _ := strconv.Atoi(masterthreadID)
				cpuPercenttoINT, _ := strconv.ParseFloat(cpuPercent, 64)
				memoryPercenttoINT, _ := strconv.ParseFloat(memoryPercent, 64)

				if masterthreadIDtoINT == 0 {

					ioFile, e := ioutil.ReadFile("/proc/" + pid + "/io")
					if e != nil {
						continue
					}
					var io []string
					for _, line := range strings.Split(string(ioFile), "\n") {
						f := strings.Fields(line)
						if len(f) == 2 {
							io = append(io, f[1])
						}
					}

					var processMemoryMB = util.FloatDecimalPoint(memoryTotalMB/100*memoryPercenttoINT, 2)

					ReadBytesInt, _ := strconv.Atoi(io[4])
					ReadBytesFloat := float64(ReadBytesInt)
					processReadKB, _ := util.ConvertBytesTo(ReadBytesFloat, "kb")

					WriteBytesInt, _ := strconv.Atoi(io[5])
					WriteBytesFloat := float64(WriteBytesInt)
					processWriteKB, _ := util.ConvertBytesTo(WriteBytesFloat, "kb")

					c := ProcessStruct{
						CPU:     cpuPercenttoINT,
						Memory:  processMemoryMB,
						Command: processName,
						KBRead:  processReadKB,
						KBWrite: processWriteKB,
					}

					ps = append(ps, c)

				}

			}

		}

	}

	return ps
}
