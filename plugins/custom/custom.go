package custom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/amonapp/amonagent/logging"
	"github.com/amonapp/amonagent/plugins"
	"github.com/gonuts/go-shellquote"
)

// Metric - XXX
type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
}

var pluginLogger = logging.GetLogger("amonagent.custom")

// Run - XXX
func Run(command *Command) (string, error) {
	splitCmd, err := shellquote.Split(command.Command)
	if err != nil || len(splitCmd) == 0 {
		return "", fmt.Errorf("exec: unable to parse command, %s", err)
	}

	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("exec: %s for command '%s'", err, command.Command)
	}

	return out.String(), nil
}

// ParseLine - XXX
func ParseLine(s string) (Metric, error) {
	// Split line by : and | -> name:value|type (requests.per_second:100|gauge)
	f := func(c rune) bool {
		return c == '|' || c == ':'
	}
	fields := strings.FieldsFunc(s, f)
	m := Metric{}
	if len(fields) == 3 {
		if s, err := strconv.ParseFloat(fields[1], 64); err == nil {
			m = Metric{Name: fields[0], Value: s, Type: fields[2]}

		}

	}

	return m, nil
}

// Command - XXX
type Command struct {
	Command string `json:"command"`
	Name    string `json:"name"`
}

// Custom - XXX
type Custom struct {
}

// PerformanceStruct - XXX
type PerformanceStruct struct {
	Gauges   map[string]interface{} `json:"gauges"`
	Counters map[string]interface{} `json:"counters"`
}

// Description - XXX
func (c *Custom) Description() string {
	return "Collects metrics from custom collectors"
}

// Collect - XXX
func (c *Custom) Collect() (interface{}, error) {
	file, e := ioutil.ReadFile("/home/martin/temp/amonagent/custom_config.json")
	if e != nil {
		fmt.Printf("Config error: %v\n", e)
	}

	var commands []Command
	var results []PerformanceStruct
	json.Unmarshal(file, &commands)

	for _, command := range commands {
		PerformanceStruct := PerformanceStruct{}
		result, err := Run(&command)

		lines := strings.Split(result, "\n")
		gauges := make(map[string]interface{})
		counters := make(map[string]interface{})
		for _, line := range lines {
			metric, _ := ParseLine(line)
			if metric.Type == "gauge" {
				gauges[metric.Name] = metric.Value
			}

			if metric.Type == "counter" {
				counters[metric.Name] = metric.Value
			}
		}

		if err != nil {
			fmt.Printf("Unable to execute command, %s", err)
		}

		PerformanceStruct.Gauges = gauges
		PerformanceStruct.Counters = counters
		results = append(results, PerformanceStruct)
	}

	return results, nil
}
func init() {
	plugins.Add("custom", func() plugins.Plugin {
		return &Custom{}
	})
}
