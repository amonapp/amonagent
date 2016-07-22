package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/amonapp/amonagent/internal/settings"
)

// PluginConfig - XXX
type PluginConfig struct {
	Path string
	Name string
}

// PluginConfigPath - XXX
var PluginConfigPath = path.Join(settings.ConfigPath, "plugins-enabled")

// GetConfigPath - XXX
func GetConfigPath(plugin string) (PluginConfig, error) {
	config := PluginConfig{}

	// On Linux /etc/opt/amonagent/plugins-enabled/plugin.conf
	var pluginPath = path.Join(PluginConfigPath, strings.Join([]string{plugin, "conf"}, "."))
	config.Path = pluginPath
	config.Name = plugin

	return config, nil
}

// GetPluginConfigKeyValue - {"status_url": url, "ip_address": 8000}
func GetPluginConfigKeyValue(plugin string) (interface{}, error) {
	c, err := GetConfigPath(plugin)
	if err != nil {
		fmt.Printf("Can't read config file: %s %v\n", c.Path, err)
	}
	var data map[string]interface{}
	file, e := ioutil.ReadFile(c.Path)
	if e != nil {
		return data, e
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return data, err
	}

	return data, nil

}

// GetPluginConfigList - ["command", "another_command"]
func GetPluginConfigList(plugin string) ([]string, error) {
	c, err := GetConfigPath(plugin)
	if err != nil {
		fmt.Printf("Can't read config file: %s %v\n", c.Path, err)
	}
	var data []string
	file, e := ioutil.ReadFile(c.Path)
	if e != nil {
		return data, e
	}

	if err := json.Unmarshal(file, &data); err != nil {
		return data, err
	}

	return data, nil

}

// GetAllEnabledPlugins - XXX
func GetAllEnabledPlugins() ([]PluginConfig, error) {
	fileList := []PluginConfig{}

	if _, err := os.Stat(PluginConfigPath); os.IsNotExist(err) {
		if err != nil {
			if os.IsNotExist(err) {
				// Plugin config directory doesn't exist for some reason. Create
				PluginDirErr := os.MkdirAll(PluginConfigPath, os.ModePerm)

				if PluginDirErr != nil {
					fmt.Printf("Plugin directory doesn't exist: %s\n", PluginConfigPath)
				}

			}
			return fileList, err
		}

	}

	filepath.Walk(PluginConfigPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			// Only files ending with .conf
			fileName := strings.Split(f.Name(), ".conf")
			if len(fileName) == 2 {

				// Go over the list of available plugins and add the config only if it is for an existing plugin
				for name := range Plugins {
					if name == fileName[0] {
						f := PluginConfig{Path: path, Name: fileName[0]}
						fileList = append(fileList, f)
					}

				}

			}

		}
		return nil
	})

	return fileList, nil
}

// Plugin - XXX
type Plugin interface {

	// Description returns a one-sentence description on the Plugin
	Description() string

	SampleConfig() string

	// Collects all the metrics and returns a struct with the results
	Collect() (interface{}, error)

	// Start starts the service - Optional
	Start() error

	// Stop stops the services and closes any necessary channels and connections - Optional
	Stop()
}

// PluginRegistry - XXX
type PluginRegistry func() Plugin

// Plugins - XXX
var Plugins = map[string]PluginRegistry{}

// Add - XXX
func Add(name string, registry PluginRegistry) {
	Plugins[name] = registry
}
