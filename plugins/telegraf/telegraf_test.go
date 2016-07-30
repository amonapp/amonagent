package telegraf

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"sort"
	"testing"

	"github.com/amonapp/amonagent/internal/testing"
	"github.com/amonapp/amonagent/plugins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTelegrafParseLine(t *testing.T) {

	s := Telegraf{}
	r, err := s.ParseLine("> mem,host=ubuntu available_percent=78.43483331332489,buffered=199602176i,used=1802661888i,used_percent=21.56516668667511 1469886743")
	require.NoError(t, err)

	assert.Len(t, r.Elements, 4)
	for _, line := range r.Elements {
		assert.Equal(t, line.Plugin, "telegraf.mem")

		validGauges := []string{"mem_available.percent", "mem_buffered", "mem_used", "mem_used.percent"}
		sort.Strings(validGauges)
		i := sort.SearchStrings(validGauges, line.Gauge)
		var gaugeFound = i < len(validGauges) && validGauges[i] == line.Gauge
		assert.True(t, gaugeFound, "Valid Gauge Name")

	}

	r, err = s.ParseLine("> system,host=ubuntu load1=0.11,load15=0.06,load5=0.05,n_cpus=4i,n_users=2i,uptime=7252i 1469891972000000000")
	require.NoError(t, err)
	assert.Len(t, r.Elements, 6)

	for _, line := range r.Elements {
		assert.Equal(t, line.Plugin, "telegraf.system")

		validGauges := []string{"system_load1", "system_load15", "system_load5", "system_n.cpus", "system_n.users", "system_uptime"}
		sort.Strings(validGauges)
		i := sort.SearchStrings(validGauges, line.Gauge)
		var gaugeFound = i < len(validGauges) && validGauges[i] == line.Gauge
		assert.True(t, gaugeFound, "Valid Gauge Name")

	}

}

func TestTelegrafonfigDefaults(t *testing.T) {

	plugins.PluginConfigPath = path.Join("/tmp", "plugins-enabled")
	pluginhelper.WritePluginConfig("telegraf", "bogusstring")

	s := Telegraf{}
	assert.Equal(t, s.Config.Config, "")

	pluginhelper.WritePluginConfig("telegraf", "{\"config\": \"/path/to/telegraf.conf\"}")

	t2 := Telegraf{}
	configErr2 := t2.SetConfigDefaults()
	require.NoError(t, configErr2)

	assert.Equal(t, t2.Config.Config, "/path/to/telegraf.conf")

}

func TestTelegraf(t *testing.T) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("testdata directory not found")
	}

	var telegrafConfig = path.Join(path.Dir(filename), "testdata", "telegraf.conf")

	ConfigString := fmt.Sprintf("{\"config\": \"%s\"}", telegrafConfig)
	pluginhelper.WritePluginConfig("telegraf", ConfigString)

	c := Telegraf{}

	result, err := c.Collect()
	require.NoError(t, err)

	resultReflect := reflect.ValueOf(result)
	i := resultReflect.Interface()
	pluginMap := i.(map[string]interface{})

	require.NotZero(t, pluginMap["telegraf.mem"])

	gaugesMapReflect := reflect.ValueOf(pluginMap["telegraf.mem"])
	j := gaugesMapReflect.Interface()
	gaugesMap := j.(map[string]map[string]string)

	require.NotZero(t, gaugesMap["gauges"])

}
