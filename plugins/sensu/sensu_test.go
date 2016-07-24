package sensu

import (
	"path"
	"runtime"
	"testing"

	"github.com/amonapp/amonagent/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSensuCollect(t *testing.T) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("testdata directory not found")
	}

	var pythonScript = path.Join("python ", path.Dir(filename), "testdata", "connections.py")

	config := Config{}
	configLine := util.Command{Name: "connections", Command: pythonScript}

	config.Commands = append(config.Commands, configLine)

	c := Custom{}
	c.Config = config

	result, err := c.Collect()
	require.NoError(t, err)

	fields := map[string]interface{}{
		"connections.active": float64(100),
		"connections.error":  float64(500),
	}

	expectedResults := make(PerformanceStructBlock, 0)
	p := PerformanceStruct{Gauges: fields}
	expectedResults["connections"] = p

	assert.Equal(t, result, expectedResults)

}
