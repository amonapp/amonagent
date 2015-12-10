package collectors

import (
	"encoding/json"
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
	Name    string  `json:"name"`
}

// ProcessesList - list of individual process data
type ProcessesList []ProcessStruct

// Processes - get data from sysstat, format and return the result
func Processes() (ProcessesList, error) {
	c1, _ := exec.Command("pidstat", "-ruhtd").Output()

	var ps ProcessesList
	v, _ := mem.VirtualMemory()
	memoryTotalMB, _ := util.ConvertBytesTo(float64(v.Total), "mb", 0)

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
				masterthreadID, cpuPercent, memoryPercent, ReadPerSecond, WritePerSecond, processName := processData[3], processData[7], processData[13], processData[14], processData[15], processData[17]

				masterthreadIDtoINT, _ := strconv.Atoi(masterthreadID)

				if masterthreadIDtoINT == 0 {
					cpuPercenttoINT, _ := strconv.ParseFloat(cpuPercent, 64)
					memoryPercenttoINT, _ := strconv.ParseFloat(memoryPercent, 64)

					var processMemoryMB, _ = util.FloatDecimalPoint(memoryTotalMB/100*memoryPercenttoINT, 2)

					ReadKBytes, _ := strconv.ParseFloat(ReadPerSecond, 64)
					WriteKBytes, _ := strconv.ParseFloat(WritePerSecond, 64)

					if ReadKBytes == -1.0 && WriteKBytes == -1.0 {
						ReadKBytes = 0.0
						WriteKBytes = 0.0
					}

					c := ProcessStruct{
						CPU:     cpuPercenttoINT,
						Memory:  processMemoryMB,
						Name:    processName,
						KBRead:  ReadKBytes,
						KBWrite: WriteKBytes,
					}

					ps = append(ps, c)

				}

			}

		}

	}

	return ps, nil
}
