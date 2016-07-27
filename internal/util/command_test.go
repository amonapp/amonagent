package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecWithExitCode(t *testing.T) {

	command := Command{Command: "check-disk-usage.rb -w 80 -c 90"}

	result := ExecWithExitCode(command)

	assert.Equal(t, result.Output, "CheckDisk OK: All disk usage under 80% and inode usage under 85%\n")

}

func TestExecWithExitCode_NonExistingCommand(t *testing.T) {

	command := Command{Command: "check-dummy-usage.rb -w 80 -c 90"}

	result := ExecWithExitCode(command)

	// Return empty string, don't panic
	assert.Equal(t, result.Output, "")

}
