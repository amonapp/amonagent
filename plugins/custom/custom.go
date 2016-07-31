package custom

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/amonapp/amonagent/internal/util"
	"github.com/amonapp/amonagent/plugins"
)

// Metric - XXX
type Metric struct {
	Name  string  `json:"name"`
	Value float64 `json:"value"`
	Type  string  `json:"type"`
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
func (c *Custom) Stop() {}

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

// Config - XXX
type Config struct {
	Commands []util.Command `json:"commands"`
}

// SetConfigDefaults - XXX
func (c *Custom) SetConfigDefaults() error {
	// Commands already set. For example - in the test suite
	if len(c.Config.Commands) > 0 {
		return nil
	}
	configFile, err := plugins.ReadPluginConfig("custom")
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": "custom",
			"error":  err,
		}).Error("Can't read config file")
	}
	var Commands []util.Command
	if e := json.Unmarshal(configFile, &Commands); e != nil {
		log.WithFields(log.Fields{"plugin": "custom", "error": e.Error()}).Error("Can't decode JSON file")
	}

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

	resultChan := make(chan util.CommandResult, len(c.Config.Commands))

	for _, v := range c.Config.Commands {
		wg.Add(1)

		go func(v util.Command) {

			CheckResult := util.ExecWithExitCode(v)

			resultChan <- CheckResult
			defer wg.Done()
		}(v)

	}
	wg.Wait()
	close(resultChan)

	for i := range resultChan {

		PerformanceStruct := PerformanceStruct{}

		lines := strings.Split(i.Output, "\n")
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

			pluginName := i.Name
			if len(gauges) > 0 {

				PerformanceStruct.Gauges = gauges
			}
			if len(counters) > 0 {
				PerformanceStruct.Counters = counters
			}

			results[pluginName] = PerformanceStruct

		}

	}

	return results, nil
}
func init() {
	plugins.Add("custom", func() plugins.Plugin {
		return &Custom{}
	})
}
