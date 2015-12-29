package plugins

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// PluginConfig - XXX
type PluginConfig struct {
	Path string
	Name string
}

// ConfigPath - XXX
const ConfigPath = "/etc/opt/amonagent/plugins-enabled"

// ReadConfigPath - Works only with flat config files, do something different for nested configs
func ReadConfigPath(path string) (interface{}, error) {
	var data map[string]interface{}
	file, e := ioutil.ReadFile(path)
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
	filepath.Walk(ConfigPath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			// Only files ending with .conf
			fileName := strings.Split(f.Name(), ".conf")
			if len(fileName) == 2 {
				f := PluginConfig{Path: path, Name: fileName[0]}
				fileList = append(fileList, f)
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
	Collect(string) (interface{}, error)
}

// Creator - XXX
type Creator func() Plugin

// Plugins - XXX
var Plugins = map[string]Creator{}

// Add - XXX
func Add(name string, creator Creator) {
	Plugins[name] = creator
}
