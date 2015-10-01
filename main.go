package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// conversion units
const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)

// ToMegabytes parses a string formatted by ByteSize as megabytes
func ToMegabytes(s uint64) (uint64, error) {

	bytes := s / MEGABYTE

	return bytes, nil
}

// IsDigit returns true if s consists of decimal digits.
func IsDigit(s string) bool {
	r := strings.NewReader(s)
	for {
		ch, _, err := r.ReadRune()
		if ch == 0 || err != nil {
			break
		} else if ch == utf8.RuneError {
			return false
		} else if !unicode.IsDigit(ch) {
			return false
		}
	}
	return true
}

// type SystemMemory struct {
// 	Path         string `json:"path"`
// 	Rss          uint64 `json:"rss"`
// 	Size         uint64 `json:"size"`
// 	Pss          uint64 `json:"pss"`
// 	SharedClean  uint64 `json:"shared_clean"`
// 	SharedDirty  uint64 `json:"shared_dirty"`
// 	PrivateClean uint64 `json:"private_clean"`
// 	PrivateDirty uint64 `json:"private_dirty"`
// 	Referenced   uint64 `json:"referenced"`
// 	Anonymous    uint64 `json:"anonymous"`
// 	Swap         uint64 `json:"swap"`
// }
//
// func (m SystemMemory) String() string {
// 	s, _ := json.Marshal(m)
// 	return string(s)
// }

func main() {
	v, _ := mem.VirtualMemory()
	fmt.Println(v)

	memoryTotalMB, _ := ToMegabytes(v.Total)
	fmt.Println(memoryTotalMB)

	fmt.Println(v.Total)
	s, _ := host.HostInfo()
	fmt.Println(s)

	// d, _ := disk.DiskPartitions(false)
	//
	// for _, partition := range d {
	// 	// fmt.Println(partition)
	// 	diskUsage, _ := disk.DiskUsage(partition.Device)
	//
	// 	fmt.Println(diskUsage)
	//
	// }

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

			fmt.Println(fs, fsType, spaceUsed, spaceTotal, spaceFree, mount)
		}
	}

	// TODO: support mount points with spaces in them. They mess up the field order
	// currently due to df's columnar output.
	// if len(fields) != 7 || !IsDigit(fields[2]) {
	// 	return nil
	// }
	//
	// fs := fields[0]
	// fsType := fields[1]
	// spaceTotal := fields[2]
	// spaceUsed := fields[3]
	// spaceFree := fields[4]
	// mount := fields[6]

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
	// header
	//   Time   UID      TGID       TID    %usr %system  %guest    %CPU   CPU  minflt/s  majflt/s     VSZ    RSS   %MEM   kB_rd/s   kB_wr/s kB_ccwr/s  Command

	// pidstat := sh.Command("pidstat", "-ruhtd").Command("grep", "-v", "Linux").Command("grep", "-v", "Command").Command("awk", "NF").Run()
	// //
	// fmt.Println(pidstat)
	//
	// process := strings.Split("pidstat", "\n")
	// ip, port := process[0], process[1]
	// fmt.Println(ip, port)

	// scanner := bufio.NewScanner(pidstat)
	// for scanner.Scan() {
	// 	fmt.Println(scanner.Text())
	// }

	// processList, _ := process.Pids()

	// for _, pid := range processList {

	// p, err := process.NewProcess(pid)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// memory, e := p.MemoryInfo()
	// if e != nil {
	// 	fmt.Println(e)
	// }
	// fmt.Println(memory)
	// fmt.Println(p.NumThreads())
	// s := strconv.Itoa(int(pid))
	// fmt.Println(s)

	// pidstat -r -p 129856 | grep -v Linux | grep -v Command | awk NF
	// sh.Command("pidstat", "-r", "-p", s).Command("grep", "-v", "Linux").Command("grep", "-v", "Command").Command("awk", "NF").Run()

	// sh.Command("pidstat", "-ruhtd").Command("grep", "-v", "Linux").Command("grep", "-v", "Command").Command("awk", "NF").Run()

	// pidstat := exec.Command("pidstat", "-r", "-p", s)
	//
	// remove header
	// grepHeader := exec.Command("grep -v Linux")
	//
	// remove second header
	// grepHeader2 := exec.Command("grep -v Command")
	//
	// removeEmptyLine := exec.Command("awk NF")

	// Run the pipeline
	// output, _, err := Pipeline(pidstat)
	// if err != nil {
	// 	// fmt.Println(stderr)
	// 	fmt.Println(output)
	// }

	// pidstat -d -p 13212
	// fmt.Println(p.MemoryPercent())
	// fmt.Println(p.IOCounters())
	// fmt.Println(p.Name())
	// fmt.Println(p.CPUPercent(0))

	// file_status, e := os.Stat("/proc/" + pid)
	// if e != nil {
	// 	w.Remove(pid)
	// 	continue
	// }
	// processCount++
	// stats_file, e := ioutil.ReadFile("/proc/" + pid + "/stat")
	// if e != nil {
	// 	w.Remove(pid)
	// 	continue
	// }
	// io_file, e := ioutil.ReadFile("/proc/" + pid + "/io")
	// if e != nil {
	// 	w.Remove(pid)
	// 	continue
	// }
	// limits, e := ioutil.ReadFile("/proc/" + pid + "/limits")
	// if e != nil {
	// 	w.Remove(pid)
	// 	continue
	// }
	// fd_dir, e := os.Open("/proc/" + pid + "/fd")
	// if e != nil {
	// 	w.Remove(pid)
	// 	continue
	// }
	// fds, e := fd_dir.Readdirnames(0)
	// fd_dir.Close()
	// if e != nil {
	// 	w.Remove(pid)
	// 	continue
	// }
	// stats := strings.Fields(string(stats_file))
	// if len(stats) < 24 {
	// 	err = fmt.Errorf("stats too short")
	// 	continue
	// }
	// var io []string
	// for _, line := range strings.Split(string(io_file), "\n") {
	// 	f := strings.Fields(line)
	// 	if len(f) == 2 {
	// 		io = append(io, f[1])
	// 	}
	// }
	// if len(io) < 6 {
	// 	err = fmt.Errorf("io too short")
	// 	continue
	// }
	// }

	// for _, processPid := range processList {
	// 	p, err := process.NewProcess(processPid)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println(p.MemoryInfo())
	// 	fmt.Println(p.IOCounters())
	// 	fmt.Println(p.Name())
	// 	fmt.Println(p.CPUPercent(0))
	//
	// }

}
