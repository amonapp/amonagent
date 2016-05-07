package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// relies on sysstat
func TestProcessesCollector(t *testing.T) {
	processes, _ := Processes()
	var aString interface{} = "string"
	var aFloat interface{} = float64(1)
	assert.NotEmpty(t, processes)
	for _, p := range processes {
		assert.IsType(t, aString, p.Name)
		assert.IsType(t, aFloat, p.CPU)
		assert.IsType(t, aString, p.Memory)
	}

}
