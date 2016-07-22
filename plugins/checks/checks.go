package checks

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/amonapp/amonagent/internal/util"
	"github.com/amonapp/amonagent/plugins"
)

// Checks - XXX
type Checks struct {
	Config Config
}

// Start - XXX
func (c *Checks) Start() error { return nil }

// Stop - XXX
func (c *Checks) Stop() {}

// Description - XXX
func (c *Checks) Description() string {
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
# Config location: /etc/opt/amonagent/plugins-enabled/checks.conf
`

// SampleConfig - XXX
func (c *Checks) SampleConfig() string {
	return sampleConfig
}

// Config - XXX
type Config struct {
	Commands []string `mapstructure:"commands"`
}

// SetConfigDefaults - XXX
func (c *Checks) SetConfigDefaults() error {
	configFile, err := plugins.ReadPluginConfig("checks")
	if err != nil {
		fmt.Printf("Can't read config file: %s\n", err)
	}

	var Commands []string

	if e := json.Unmarshal(configFile, &Commands); e != nil {
		fmt.Printf("Can't decode JSON file: %v\n", e)
	}

	c.Config.Commands = Commands

	return nil
}

// Collect - XXX
func (c *Checks) Collect() (interface{}, error) {
	c.SetConfigDefaults()
	var wg sync.WaitGroup
	var result []util.CommandResult

	resultChan := make(chan util.CommandResult, len(c.Config.Commands))

	for _, v := range c.Config.Commands {
		wg.Add(1)

		go func(command string) {

			CheckResult := util.ExecWithExitCode(command)

			resultChan <- CheckResult
			defer wg.Done()
		}(v)

	}
	wg.Wait()
	close(resultChan)

	for i := range resultChan {
		result = append(result, i)
	}

	return result, nil
}

func init() {
	plugins.Add("checks", func() plugins.Plugin {
		return &Checks{}
	})
}
