package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUptime(t *testing.T) {
	uptime := Uptime()
	var aString interface{} = "string"
	assert.IsType(t, aString, uptime)

}

func TestDistro(t *testing.T) {
	d := Distro()
	assert.NotNil(t, d.Name)
	assert.NotNil(t, d.Version)
}

func TestHost(t *testing.T) {
	d := Host()
	var aString interface{} = "string"
	assert.IsType(t, aString, d)

}
