package sensu

import (
	"reflect"
	"testing"

	"github.com/amonapp/amonagent/internal/util"
	"github.com/stretchr/testify/require"
)

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
