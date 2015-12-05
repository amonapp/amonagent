package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/martinrusev/amonagent/logging"
	"github.com/martinrusev/amonagent/util"
	"github.com/shirou/gopsutil/mem"
)

// AmonAgentLogger for the main file
var AmonAgentLogger = logging.GetLogger("amonagent")

func (p ProcessStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

type ProcessStruct struct {
	CPU     float64 `json:"cpu"`
	Memory  float64 `json:"memory_mb"`
	KBRead  float64 `json:"kb_read"`
	KBWrite float64 `json:"kb_write"`
	Command string  `json:"command"`
}

type ProcessesList []*ProcessStruct

func Processes() {
	c1, _ := exec.Command("pidstat", "-ruhtd").Output()

	ps := ProcessesList{}

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

					append(*ps, c)
					fmt.Println(c)

				}

			}

		}

	}

	fmt.Print(ps)
}
func main() {

	// lps, _ := getLinuxProccesses()
	//
	// for _, w := range lps {
	// 	p, err := process.NewProcess(int32(w.Pid))
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	//
	// 	fmt.Println(p.IOCounters())
	// 	duration := time.Duration(1000) * time.Microsecond
	// 	fmt.Println(p.CPUPercent(duration))
	// 	fmt.Println(p.MemoryInfo())
	//
	// 	fmt.Println(w.Command)
	// }

	// v, _ := mem.VirtualMemory()
	// fmt.Println(v)
	//
	// n, _ := net.NetIOCounters(true)
	// fmt.Println(n)
	//
	// memoryTotalMB, _ := util.ToMegabytes(v.Total)
	// fmt.Println(memoryTotalMB)

	// l, _ := load.LoadAvg()
	// fmt.Println(l)
	//
	// fmt.Println(v.Total)
	// s, _ := host.HostInfo()
	// fmt.Println(s)
	//
	// d, _ := collectors.DiskSpace()
	// fmt.Println(d)
	// for _, volume := range d {
	// 	fmt.Println(volume)
	//
	// }

	Processes()
	// collectors.Processes()

}
