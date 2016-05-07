package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStats(t *testing.T) {
	m := MemoryUsage()

	var aFloat interface{} = float64(1)

	assert.IsType(t, aFloat, m.UsedMB)
	assert.IsType(t, aFloat, m.FreeMB)
	assert.IsType(t, aFloat, m.TotalMB)
	assert.IsType(t, aFloat, m.UsedPercent)
	assert.IsType(t, aFloat, m.SwapUsedMB)
	assert.IsType(t, aFloat, m.SwapTotalMB)
	assert.IsType(t, aFloat, m.SwapFreeMB)
	assert.IsType(t, aFloat, m.SwapUsedPercent)
}
