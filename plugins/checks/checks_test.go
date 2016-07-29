package checks

import (
	"path"
	"testing"

	"github.com/amonapp/amonagent/internal/testing"
	"github.com/amonapp/amonagent/internal/util"
	"github.com/amonapp/amonagent/plugins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChecksConfigDefaults(t *testing.T) {

	plugins.PluginConfigPath = path.Join("/tmp", "plugins-enabled")
	pluginhelper.WritePluginConfig("checks", "bogusstring")

	c := Checks{}
	configErr := c.SetConfigDefaults()
	require.Error(t, configErr)

	assert.Len(t, c.Config.Commands, 0, "0 commands in the config file")

	pluginhelper.WritePluginConfig("checks", "[\"check-disk-usage.rb\", \"check-memory-usage.rb\"]")

	c2 := Checks{}
	configErr2 := c2.SetConfigDefaults()
	require.NoError(t, configErr2)

	assert.Len(t, c2.Config.Commands, 2, "2 commands in the config file")

}

func TestChecksCollect(t *testing.T) {

	config := Config{}
	configLine := util.Command{Command: "check-disk-usage.rb -w 80 -c 90"}

	config.Commands = append(config.Commands, configLine)

	c := Checks{}
	c.Config = config

	result, err := c.Collect()
	require.NoError(t, err)

	var expectedResult []util.CommandResult
	p := util.CommandResult{
		Command:  "check-disk-usage.rb -w 80 -c 90",
		Output:   "CheckDisk OK: All disk usage under 80% and inode usage under 85%\n",
		ExitCode: 0,
	}

	expectedResult = append(expectedResult, p)

	assert.Equal(t, result, expectedResult)

}
