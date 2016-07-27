package sensu

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/amonapp/amonagent/internal/util"
	"github.com/amonapp/amonagent/plugins"
)

// Start - XXX
func (s *Sensu) Start() error {
	return nil
}

// Stop - XXX
func (s *Sensu) Stop() {
}

// Description - XXX
func (s *Sensu) Description() string {
	return "Collects data from Sensu plugins"
}

var sampleConfig = `
#   Available config options:
#
#    [
#        "metrics-es-node-graphite.rb",
#        "metrics-net.rb",
#        "metrics-redis-graphite.rb",
#        "metrics-iostat-extended.rb"
#    ]
#
#    List of preinstalled sensu plugins + params
#
# Config location: /etc/opt/amonagent/plugins-enabled/sensu.conf
`

// SampleConfig - XXX
func (s *Sensu) SampleConfig() string {
	return sampleConfig
}

// Sensu - XXX
type Sensu struct {
	Config Config
}

// Config - XXX
type Config struct {
	Commands []util.Command `json:"commands"`
}

// SetConfigDefaults - XXX
func (s *Sensu) SetConfigDefaults() error {
	// Commands already set. For example - in the test suite
	if len(s.Config.Commands) > 0 {
		return nil
	}
	configFile, err := plugins.ReadPluginConfig("sensu")
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": "sensu",
			"error":  err,
		}).Error("Can't read config file")
	}

	var Commands []util.Command

	if e := json.Unmarshal(configFile, &Commands); e != nil {
		log.WithFields(log.Fields{"plugin": "sensu", "error": e.Error()}).Error("Can't decode JSON file")
	}

	s.Config.Commands = Commands

	return nil
}

func (m Metric) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

// Metric - XXX
type Metric struct {
	Plugin string `json:"plugin"`
	Gauge  string `json:"gauge"`
	Value  string `json:"value"`
}

// ParsedLine - XXX
type ParsedLine struct {
	Plugin   string
	Elements []string
	Value    string
}

// ParseLine - XXX
func ParseLine(s string) (Metric, error) {
	// split by space
	f := func(c rune) bool {
		return c == ' '
	}
	// split metric name by .
	dot := func(c rune) bool {
		return c == '.'
	}
	//split metric name by _
	underscore := func(c rune) bool {
		return c == '_'
	}

	fields := strings.FieldsFunc(s, f)
	line := ParsedLine{}
	m := Metric{}
	if len(fields) == 3 {
		toFloat, _ := strconv.ParseFloat(fields[1], 64)
		value := strconv.FormatFloat(toFloat, 'f', -1, 64)

		metricFields := strings.FieldsFunc(fields[0], dot)

		var cleanName string
		// Eliminate host and plugin name here
		// Example ubuntu.elasticsearch.thread_pool......
		if len(metricFields) > 2 {
			cleanName = strings.Join(metricFields[2:], ".")
		} else {
			cleanName = strings.Join(metricFields[:], ".")
		}
		CleanMetricFields := strings.FieldsFunc(cleanName, dot)
		splitOnUnderscore := strings.FieldsFunc(cleanName, underscore)

		// Standart use case here
		// Example: thread_pool.search.active
		if len(CleanMetricFields) > 1 {
			line.Elements = CleanMetricFields

		} else {
			line.Elements = splitOnUnderscore

		}
		elements := line.Elements

		if len(elements) > 2 {
			chart := strings.Join(elements[:2], "_")
			line := strings.Join(elements[2:], "_")
			m.Gauge = chart + "." + line

		} else {
			chart := elements[0]
			line := strings.Join(elements[1:], "_")
			m.Gauge = chart + "." + line

		}
		m.Value = value
		m.Plugin = "sensu." + metricFields[1]

	}

	return m, nil
}

// Collect - XXX
// This should return the following: sensu.plugin_name: {"gauges": {}}, sensu.another_plugin :{"gauges":{}}
func (s *Sensu) Collect() (interface{}, error) {
	s.SetConfigDefaults()
	var wg sync.WaitGroup
	plugins := make(map[string]interface{})

	resultChan := make(chan util.CommandResult, len(s.Config.Commands))
	for _, command := range s.Config.Commands {
		wg.Add(1)

		go func(command util.Command) {

			CheckResult := util.ExecWithExitCode(command)

			resultChan <- CheckResult
			defer wg.Done()

		}(command)
	}

	wg.Wait()
	close(resultChan)

	for command := range resultChan {
		var result []Metric
		gauges := make(map[string]string)
		GaugesWrapper := make(map[string]map[string]string)
		plugin := ""
		lines := strings.Split(command.Output, "\n")

		for _, line := range lines {
			metric, _ := ParseLine(line)
			if len(metric.Gauge) > 0 {
				result = append(result, metric)
			}
		}

		for _, r := range result {
			gauges[r.Gauge] = r.Value
			plugin = r.Plugin
		}

		GaugesWrapper["gauges"] = gauges

		if len(plugin) > 0 {
			plugins[plugin] = GaugesWrapper
		}

	}

	return plugins, nil
}

func init() {
	plugins.Add("sensu", func() plugins.Plugin {
		return &Sensu{}
	})
}
