package collectors

import (
	"encoding/json"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
)

//{'cores': 1, 'fifteen_minutes': '0.14', 'five_minutes': '0.11', 'minute': '0.01'}
func (p LoadStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// LoadStruct - returns load avg
type LoadStruct struct {
	Minute         float64 `json:"minute"`
	FiveMinutes    float64 `json:"five_minutes"`
	FifteenMinutes float64 `json:"fifteen_minutes"`
	Cores          int     `json:"cores"`
}

// LoadAverage - returns load avg
func LoadAverage() LoadStruct {

	cores, _ := cpu.Counts(true)
	load, _ := load.Avg()

	l := LoadStruct{
		Minute:         load.Load1,
		FiveMinutes:    load.Load5,
		FifteenMinutes: load.Load15,
		Cores:          cores,
	}

	return l
}
