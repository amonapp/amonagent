package telegraf

import (
	"path"
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
