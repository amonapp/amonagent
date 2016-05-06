package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadStats(t *testing.T) {
	load := LoadAverage()
	assert.NotNil(t, load.Minute)
	assert.NotNil(t, load.FiveMinutes)
	assert.NotNil(t, load.FifteenMinutes)
	assert.NotNil(t, load.Cores)
}
