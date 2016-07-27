package plugins

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestUmarshalPluginConfig(t *testing.T) {
// 	PluginConfigPath = path.Join("/tmp", "plugins-enabled")
// 	PathReturn, _ := GetConfigPath("testplugin")

// 	var pluginPath = path.Join(PluginConfigPath, strings.Join([]string{"testplugin", "conf"}, "."))

// 	assert.Equal(t, PathReturn.Name, "testplugin")
// 	assert.Equal(t, PathReturn.Path, pluginPath)

// }

// func TestReadPluginConfig(t *testing.T) {
// 	PluginConfigPath = path.Join("/tmp", "plugins-enabled")
// 	PathReturn, _ := GetConfigPath("testplugin")

// 	var pluginPath = path.Join(PluginConfigPath, strings.Join([]string{"testplugin", "conf"}, "."))

// 	assert.Equal(t, PathReturn.Name, "testplugin")
// 	assert.Equal(t, PathReturn.Path, pluginPath)

// }

func TestGetConfigPath(t *testing.T) {
	PluginConfigPath = path.Join("/tmp", "plugins-enabled")
	PathReturn, _ := GetConfigPath("testplugin")

	var pluginPath = path.Join(PluginConfigPath, strings.Join([]string{"testplugin", "conf"}, "."))

	assert.Equal(t, PathReturn.Name, "testplugin")
	assert.Equal(t, PathReturn.Path, pluginPath)

}

func TestGetAllEnabledPlugins(t *testing.T) {
	PluginConfigPath = path.Join("/tmp/amonagent", "plugins-enabled")
	PluginDirCleanup := os.RemoveAll(PluginConfigPath)

	if PluginDirCleanup != nil {
		log.Fatal(PluginDirCleanup)
	}

	_, err := GetAllEnabledPlugins()

	// First run, plugin directory doesn't exist - don't panic
	assert.Error(t, err)

	PluginDirErr := os.MkdirAll(PluginConfigPath, os.ModePerm)

	if PluginDirErr != nil {
		log.Fatal(PluginDirErr)
	}

	validPLugins := []string{"apache", "checks", "haproxy", "redis", "sensu"}

	for _, name := range validPLugins {
		var pluginPath = path.Join(PluginConfigPath, fmt.Sprint(name, ".conf"))
		fmt.Println(pluginPath)
		_, err := os.Create(pluginPath)

		if err != nil {
			log.Fatal(err)
		}
	}

	PluginList, PluginListErr := GetAllEnabledPlugins()

	assert.Nil(t, PluginListErr)

	assert.Len(t, PluginList, 5, "5 config files found")
	var aString interface{} = "string"
	for _, plugin := range PluginList {
		assert.IsType(t, aString, plugin.Path)
		assert.IsType(t, aString, plugin.Name)
	}

	// Create bogus config files
	for i := 1; i <= 5; i++ {
		var pluginPath = path.Join(PluginConfigPath, fmt.Sprint("plugin", i, ".bogus"))
		_, err := os.Create(pluginPath)

		if err != nil {
			log.Fatal(err)
		}

	}

	PluginListTestBogus, PluginListTestBogusErr := GetAllEnabledPlugins()

	assert.Nil(t, PluginListTestBogusErr)
	assert.Len(t, PluginListTestBogus, 5, "Ignore bogus configs.")

}
