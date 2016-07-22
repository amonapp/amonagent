package custom

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/amonapp/amonagent/internal/logging"
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

// Start - XXX
func (c *Custom) Start() error { return nil }

// Stop - XXX
func (c *Custom) Stop() {
}

var sampleConfig = `
#   Available config options:
#
#  A JSON list with commands
#   [
#    {
#        "command":"python /path/to/yourplugin.py",
#        "name":"requests"
#    },
#    {
#        "command":"python /path/to/yourpingplugin.py",
#        "name":"ping"
#    },
#   ]
#
#  You can create a custom plugin in any language. Amon parses all the lines from STDOUT:
#
#   Format: metric.line_on_chart_name:value|type
#   Example:
#   requests.per_second:12|gauge
#
#   Available types: counter and gauge
#
#   To group metrics on a single chart:
#
#     connections.active:12|gauge
#     connections.waiting:12|gauge
#
#    In Python you can emmit metrics with:
#    print "requests.per_second:12|gauge"
#
#    In Bash with printf "requests.per_second:12|gauge\n", Ruby "puts requests.per_second:12|gauge"
#
# Config location: /etc/opt/amonagent/plugins-enabled/custom.conf
`

// SampleConfig - XXX
func (c *Custom) SampleConfig() string {
	return sampleConfig
}

// Custom - XXX
type Custom struct {
	Config Config
}

//
// Command - XXX
type Command struct {
	Command string `json:"command"`
	Name    string `json:"name"`
}

// Config - XXX
type Config struct {
	Commands []Command `json:"commands"`
}

// SetConfigDefaults - XXX
func (c *Custom) SetConfigDefaults() error {
	configFile, err := plugins.ReadPluginConfig("custom")
	if err != nil {
		fmt.Printf("Can't read config file: %s\n", err)
	}
	var Commands []string
	if err := json.Unmarshal(configFile, &Commands); err != nil {
		fmt.Printf("Can't decode JSON file: %v\n", err)
	}

	fmt.Println(Commands)

	c.Config.Commands = Commands

	return nil
}

func (p PerformanceStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// PerformanceStruct - XXX
type PerformanceStruct struct {
	Gauges   map[string]interface{} `json:"gauges,omitempty"`
	Counters map[string]interface{} `json:"counters,omitempty"`
}

func (p PerformanceStructBlock) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// PerformanceStructBlock - XXX
type PerformanceStructBlock map[string]PerformanceStruct

// Description - XXX
func (c *Custom) Description() string {
	return "Collects metrics from custom collectors"
}

// Collect - XXX
func (c *Custom) Collect() (interface{}, error) {
	c.SetConfigDefaults()
	var wg sync.WaitGroup
	results := make(PerformanceStructBlock, 0)

	for _, command := range c.Config.Commands {
		wg.Add(1)
		go func(command Command) {
			defer wg.Done()

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

			pluginName := "custom." + command.Name
			PerformanceStruct.Gauges = gauges
			PerformanceStruct.Counters = counters
			results[pluginName] = PerformanceStruct

		}(command)

	}
	wg.Wait()

	return results, nil
}
func init() {
	plugins.Add("custom", func() plugins.Plugin {
		return &Custom{}
	})
}
