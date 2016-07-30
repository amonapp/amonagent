package telegraf

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
	"testing"

	"github.com/amonapp/amonagent/internal/testing"
	"github.com/amonapp/amonagent/plugins"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	fmt.Println(result)

	resultReflect := reflect.ValueOf(result)
	i := resultReflect.Interface()
	pluginMap := i.(map[string]interface{})

	require.NotZero(t, pluginMap["telegraf.mem"])

	gaugesMapReflect := reflect.ValueOf(pluginMap["telegraf.mem"])
	j := gaugesMapReflect.Interface()
	gaugesMap := j.(map[string]map[string]string)

	require.NotZero(t, gaugesMap["gauges"])

	// something := []string{"sda1.iused", "sda1.avail", "sda1.capacity", "sda1.used"}

	// for _, v := range something {
	// 	require.NotZero(t, gaugesMap["gauges"][v])
	// }

}
