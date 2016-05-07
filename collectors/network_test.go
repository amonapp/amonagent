package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkStats(t *testing.T) {
	network, _ := NetworkUsage()
	var aString interface{} = "string"
	var aFloat interface{} = float64(1)
	assert.NotEmpty(t, network)
	for _, iface := range network {
		assert.IsType(t, aString, iface.Name)
		assert.IsType(t, aFloat, iface.Inbound)
		assert.IsType(t, aFloat, iface.Outbound)
	}

}
