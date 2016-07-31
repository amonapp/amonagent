package sensu

import (
	"path"
	"reflect"
	"testing"

	"github.com/amonapp/amonagent/internal/testing"
	"github.com/amonapp/amonagent/internal/util"
	"github.com/amonapp/amonagent/plugins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSensuParseLine(t *testing.T) {

	s := Sensu{}
	r, err := s.ParseLine("ubuntu.system.process_count 532 1469793929")
	require.NoError(t, err)

	assert.Equal(t, r.Value, "532")
	assert.Equal(t, r.Plugin, "sensu.system")
	assert.Equal(t, r.Gauge, "process.count")

	r, err = s.ParseLine("ubuntu.dns.A.google_com.response_time 0.02355488 1469794071")
	require.NoError(t, err)

	assert.Equal(t, r.Value, "0.02355488")
	assert.Equal(t, r.Plugin, "sensu.dns")
	assert.Equal(t, r.Gauge, "A_google_com.response_time")

	// Error here, index out of range
	r, err = s.ParseLine("response_time 2 1469794071")
	require.NoError(t, err)

	assert.Equal(t, r.Value, "2")
	assert.Equal(t, r.Plugin, "sensu.response_time")
	assert.Equal(t, r.Gauge, "response.time")
}

func TestSensuConfigDefaults(t *testing.T) {

	plugins.PluginConfigPath = path.Join("/tmp", "plugins-enabled")
	pluginhelper.WritePluginConfig("sensu", "bogusstring")

	s := Sensu{}
	configErr := s.SetConfigDefaults()
	require.Error(t, configErr)

	assert.Len(t, s.Config.Commands, 0, "0 commands in the config file")

	pluginhelper.WritePluginConfig("sensu", "[\"metrics-dns.rb -d google.com\", \"check-memory-usage.rb\"]")

	c2 := Sensu{}
	configErr2 := c2.SetConfigDefaults()
	require.NoError(t, configErr2)

	assert.Len(t, c2.Config.Commands, 2, "2 commands in the config file")

}

func TestSensuCollect(t *testing.T) {

	config := Config{}
	configLine := util.Command{Command: "metrics-disk-capacity.rb"}

	config.Commands = append(config.Commands, configLine)

	c := Sensu{}
	c.Config = config

	result, err := c.Collect()
	require.NoError(t, err)

	resultReflect := reflect.ValueOf(result)
	i := resultReflect.Interface()
	pluginMap := i.(map[string]interface{})

	require.NotZero(t, pluginMap["sensu.disk"])

	gaugesMapReflect := reflect.ValueOf(pluginMap["sensu.disk"])
	j := gaugesMapReflect.Interface()
	gaugesMap := j.(map[string]map[string]string)

	require.NotZero(t, gaugesMap["gauges"])

	something := []string{"sda1.iused", "sda1.avail", "sda1.capacity", "sda1.used"}

	for _, v := range something {
		require.NotZero(t, gaugesMap["gauges"][v])
	}

}
