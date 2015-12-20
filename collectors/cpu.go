package collectors

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/amonapp/amonagent/logging"
	"github.com/amonapp/amonagent/util"
)

var cpuLogger = logging.GetLogger("amonagent.cpu")

//{'iowait': '0.00', 'system': '0.50', 'idle': '98.50', 'user': '1.00', 'steal': '0.00', 'nice': '0.00'}
// the header values are %idle, %wait
var headerRE = regexp.MustCompile(`([%][a-zA-Z0-9]+)[\s+]?`)
var valueRE = regexp.MustCompile(`\d+[\.,]\d+`)

func (p CPUUsageStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// CPUUsageStruct - returns CPU usage stats
type CPUUsageStruct struct {
	User   float64 `json:"user"`
	Idle   float64 `json:"idle"`
	Nice   float64 `json:"nice"`
	Steal  float64 `json:"steal"`
	System float64 `json:"system"`
	IOWait float64 `json:"iowait"`
}

// CPUUsage - return a map with CPU usage stats
func CPUUsage() CPUUsageStruct {
	c1, _ := exec.Command("sar", "1", "1").Output()

	sarOutput := string(c1)
	sarLines := strings.Split(sarOutput, "\n")
	header := []string{}
	values := []float64{}

	result := make(map[string]float64)

	// Get the header
	for _, line := range sarLines {
		// Replace the regex with something faster
		matches := headerRE.FindAllStringSubmatch(line, -1)

		for _, m := range matches {
			if len(m) == 2 {
				result := strings.Replace(m[0], "%", "", -1) // replace % in idle (%idle)
				result = strings.Trim(result, " ")           // remove white space
				header = append(header, result)
			}
		}
	}

	// Get values
	for _, line := range sarLines {
		if strings.Contains(line, "Average") {
			matches := valueRE.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				if len(m) == 1 {
					valueFloat, _ := strconv.ParseFloat(m[0], 64)
					valueDecimal, _ := util.FloatDecimalPoint(valueFloat, 2)
					values = append(values, valueDecimal)
				}

			}
		}

	}

	if len(header) == len(values) {
		for i := range header {
			result[header[i]] = values[i]
		}

	}

	c := CPUUsageStruct{
		User:   result["user"],
		Idle:   result["idle"],
		Nice:   result["nice"],
		Steal:  result["steal"],
		System: result["system"],
		IOWait: result["iowait"],
	}

	return c
}
