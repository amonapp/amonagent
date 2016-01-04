package sensu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/amonapp/amonagent/plugins"
	"github.com/gonuts/go-shellquote"
)

// Sensu - XXX
type Sensu struct {
	Config Config
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

// Config - XXX
type Config struct {
	Commands []string `mapstructure:"commands"`
}

// SetConfigDefaults - XXX
func (s *Sensu) SetConfigDefaults(configPath string) error {
	jsonFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		fmt.Printf("Can't read config file: %s %v\n", configPath, err)
	}
	var Commands []string
	if err := json.Unmarshal(jsonFile, &Commands); err != nil {
		fmt.Printf("Can't decode JSON file: %v\n", err)
	}

	s.Config.Commands = Commands

	return nil
}

// Run - XXX
func Run(command string) (string, error) {
	splitCmd, err := shellquote.Split(command)
	if err != nil || len(splitCmd) == 0 {
		return "", fmt.Errorf("exec: unable to parse command, %s", err)
	}

	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("exec: %s for command '%s'", err, command)
	}

	return out.String(), nil
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

// Command - XXX
type Command struct {
	Command string `json:"command"`
	Name    string `json:"name"`
}

// Collect - XXX
// This should return the following: sensu.plugin_name: {"gauges": {}}, sensu.another_plugin :{"gauges":{}}
func (s *Sensu) Collect(configPath string) (interface{}, error) {
	s.SetConfigDefaults(configPath)
	var wg sync.WaitGroup
	plugins := make(map[string]interface{})

	for _, command := range s.Config.Commands {
		wg.Add(1)

		go func(command string) {
			defer wg.Done()
			var result []Metric
			gauges := make(map[string]interface{})
			GaugesWrapper := make(map[string]interface{})
			plugin := ""
			Commandresult, err := Run(command)
			if err != nil {
				fmt.Println("Can't execute command: ", err)
			}
			lines := strings.Split(Commandresult, "\n")

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

		}(command)
	}

	wg.Wait()

	return plugins, nil
}

func init() {
	plugins.Add("sensu", func() plugins.Plugin {
		return &Sensu{}
	})
}
