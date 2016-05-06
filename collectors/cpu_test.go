package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCPUStats(t *testing.T) {
	cpuUsage := CPUUsage()
	assert.NotNil(t, cpuUsage.IOWait)
	assert.NotNil(t, cpuUsage.System)
	assert.NotNil(t, cpuUsage.Idle)
	assert.NotNil(t, cpuUsage.User)
	assert.NotNil(t, cpuUsage.Steal)
	assert.NotNil(t, cpuUsage.Nice)
}
