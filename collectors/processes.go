package collectors

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/amonapp/amonagent/logging"
	"github.com/amonapp/amonagent/util"
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
	Memory  string  `json:"memory_mb"`
	KBRead  float64 `json:"kb_read"`
	KBWrite float64 `json:"kb_write"`
	Name    string  `json:"name"`
}

// ProcessesList - list of individual process data
type ProcessesList []ProcessStruct

// SliceFindStringIndex - XXX
func SliceFindStringIndex(s []string, e string) int {
	for v, a := range s {
		// match, ignore case
		if strings.EqualFold(a, e) {
			return v
		}
	}
	return -1
}

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

	// Helper
	// Time(0)   UID(1)      TGID(2)       TID(3)
	// %usr{4} %system{5}  %guest{6}    %CPU{7}   CPU{8}
	// minflt/s{9}  majflt/s{10}     VSZ{11}    RSS{12}
	// %MEM{13}   kB_rd/s{14}   kB_wr/s{15} kB_ccwr/s{16}  Command{17}
	var headerData []string
	for _, processLine := range pidstatLines {
		// Get the header
		if strings.Contains(processLine, "%CPU") || strings.Contains(processLine, "%Command") {
			headerData = strings.Fields(processLine)
			if len(headerData) > 0 {
				// remove the first column, if it has the # sign
				if headerData[0] == "#" {
					headerData = append(headerData[:0], headerData[0+1:]...)
				}
			}
			break
		}

	}
	for _, processLine := range pidstatLines {
		if len(headerRegex.FindString(processLine)) == 0 {
			processData := strings.Fields(processLine)

			if len(processData) == len(headerData) {
				masterthreadIDIndex := SliceFindStringIndex(headerData, "TID")
				if masterthreadIDIndex != -1 {
					masterthreadID := processData[masterthreadIDIndex]
					masterthreadIDtoINT, _ := strconv.Atoi(masterthreadID)

					// It is a master thread, proceed with the actual data
					if masterthreadIDtoINT == 0 {

						CPUIndex := SliceFindStringIndex(headerData, "%CPU")
						MEMIndex := SliceFindStringIndex(headerData, "%MEM")

						if CPUIndex != -1 && MEMIndex != -1 {
							cpuPercent := processData[CPUIndex]
							cpuPercenttoINT, _ := strconv.ParseFloat(cpuPercent, 64)

							memPercent := processData[MEMIndex]
							memPercenttoINT, _ := strconv.ParseFloat(memPercent, 64)

							var processMemoryMB, _ = util.FloatDecimalPoint(memoryTotalMB/100*memPercenttoINT, 2)

							ReadKBIndex := SliceFindStringIndex(headerData, "kB_rd/s")
							WriteKBIndex := SliceFindStringIndex(headerData, "kB_wr/s")
							var ReadKBytes = 0.0
							var WriteKBytes = 0.0

							if ReadKBIndex != -1 && WriteKBIndex != -1 {
								ReadPerSecond := processData[ReadKBIndex]
								ReadKBytes, _ = strconv.ParseFloat(ReadPerSecond, 64)

								WritePerSecond := processData[WriteKBIndex]
								WriteKBytes, _ = strconv.ParseFloat(WritePerSecond, 64)

								if ReadKBytes == -1.0 && WriteKBytes == -1.0 {
									ReadKBytes = 0.0
									WriteKBytes = 0.0
								}

							}

							processNameIndex := SliceFindStringIndex(headerData, "command")

							// Everything is find up to this point, append
							if processNameIndex != -1 {
								processName := processData[processNameIndex]

								formattedprocessMemoryMB, _ := util.FloatToString(processMemoryMB)

								c := ProcessStruct{
									CPU:     cpuPercenttoINT,
									Memory:  formattedprocessMemoryMB,
									Name:    processName,
									KBRead:  ReadKBytes,
									KBWrite: WriteKBytes,
								}

								ps = append(ps, c)

							}

						} else {
							processLogger.Info("Can't find mem/cpu data")
						}

					}
				}

			}

			//fmt.Print(headerData)

		}

	}

	return ps, nil
}
