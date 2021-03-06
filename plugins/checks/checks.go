package checks

import (
	"encoding/json"
	"sync"

	log "github.com/Sirupsen/logrus"
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
#       "check-dns.rb -d twitter.com",
#		"check-netstat-tcp.rb",
# 		"check-banner.rb",
# 		"check-ports.rb",
# 		"check-postgres-alive.rb"
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
	Commands []util.Command `json:"commands"`
}

// SetConfigDefaults - XXX
func (c *Checks) SetConfigDefaults() error {

	// Commands already set. For example - in the test suite
	if len(c.Config.Commands) > 0 {
		return nil
	}

	configFile, err := plugins.ReadPluginConfig("checks")
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": "checks",
			"error":  err,
		}).Error("Can't read config file")

		return err
	}

	var Commands []util.Command
	var CommandStrings []string

	if e := json.Unmarshal(configFile, &CommandStrings); e != nil {
		log.WithFields(log.Fields{"plugin": "checks", "error": e.Error()}).Error("Can't decode JSON file")
		return e
	}

	for _, str := range CommandStrings {
		var command = util.Command{Command: str}
		Commands = append(Commands, command)
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

		go func(command util.Command) {

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
