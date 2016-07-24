package checks

import (
	"testing"

	"github.com/amonapp/amonagent/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
