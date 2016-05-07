package plugins

import (
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetConfigPath(t *testing.T) {
	PluginConfigPath = path.Join("/tmp", "plugins-enabled")
	PathReturn, _ := GetConfigPath("testplugin")

	var pluginPath = path.Join(PluginConfigPath, strings.Join([]string{"testplugin", "conf"}, "."))

	assert.Equal(t, PathReturn.Name, "testplugin")
	assert.Equal(t, PathReturn.Path, pluginPath)

}

func TestGetAllEnabledPlugins(t *testing.T) {
	PluginConfigPath = path.Join("/tmp", "plugins-enabled")
	_, err := GetAllEnabledPlugins()

	// First run, plugin directory doesn't exist - don't panic
	assert.Error(t, err)

}
