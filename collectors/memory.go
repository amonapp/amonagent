package collectors

import (
	"encoding/json"

	"github.com/amonapp/amonagent/internal/util"
	psmem "github.com/shirou/gopsutil/mem"
)

//{'used_percent': 62, 'swap_free_mb': 0, 'used_mb': 2452, 'swap_used_percent': 0, 'swap_used_mb': 0, 'total_mb': 3939, 'swap_total_mb': 0, 'free_mb': 1487}

func (p MemoryStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// MemoryStruct - XXX
type MemoryStruct struct {
	UsedPercent     float64 `json:"used_percent"`
	UsedMB          float64 `json:"used_mb"`
	TotalMB         float64 `json:"total_mb"`
	FreeMB          float64 `json:"free_mb"`
	SwapUsedMB      float64 `json:"swap_used_mb"`
	SwapFreeMB      float64 `json:"swap_free_mb"`
	SwapTotalMB     float64 `json:"swap_total_mb"`
	SwapUsedPercent float64 `json:"swap_used_percent"`
}

// MemoryUsage - XXX
func MemoryUsage() MemoryStruct {
	mem, _ := psmem.VirtualMemory()
	swap, _ := psmem.SwapMemory()

	TotalMB, _ := util.ConvertBytesTo(mem.Total, "mb", 0)
	FreeMB, _ := util.ConvertBytesTo(mem.Free, "mb", 0)
	UsedMB, _ := util.ConvertBytesTo(mem.Used, "mb", 0)
	UsedPercent, _ := util.FloatDecimalPoint(mem.UsedPercent, 0)
	SwapUsedMB, _ := util.ConvertBytesTo(swap.Used, "mb", 0)
	SwapTotalMB, _ := util.ConvertBytesTo(swap.Total, "mb", 0)
	SwapFreeMB, _ := util.ConvertBytesTo(swap.Free, "mb", 0)
	SwapUsedPercent, _ := util.FloatDecimalPoint(swap.UsedPercent, 0)

	m := MemoryStruct{
		UsedMB:          UsedMB,
		TotalMB:         TotalMB,
		FreeMB:          FreeMB,
		UsedPercent:     UsedPercent,
		SwapUsedMB:      SwapUsedMB,
		SwapTotalMB:     SwapTotalMB,
		SwapFreeMB:      SwapFreeMB,
		SwapUsedPercent: SwapUsedPercent,
	}

	return m
}
