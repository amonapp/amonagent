package main

import (
	"fmt"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

func main() {
	v, _ := mem.VirtualMemory()
	fmt.Println(v)

	fmt.Println(v.Total)
	s, _ := host.HostInfo()
	fmt.Println(s)

	processList, _ := process.Pids()

	for pid := range processList {
		fmt.Println(pid)
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
	}

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
