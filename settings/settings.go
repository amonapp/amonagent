package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// Struct -XXX
type Struct struct {
	AmonInstance string `json:"amon_instance"`
	Interval     int    `json:"interval"`
	APIKey       string `json:"api_key"`
	ServerKey    string `json:"server_key"`
}

// Settings -XXX
func Settings() Struct {

	var settings Struct

	var SettingsPath = path.Join(ConfigPath, "amonagent.conf")

	configFile, err := os.Open(SettingsPath)
	if err != nil {
		fmt.Println("Error opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(&settings); err != nil {
		fmt.Println("Error while parsing /etc/opt/amonagent/amonagent.conf", err.Error())
	}

	// Set defaults
	if settings.Interval == 0 {
		settings.Interval = 60
	}

	// Remote trailing slash from the url
	if strings.HasSuffix(settings.AmonInstance, "/") {
		cutOffLastCharLen := len(settings.AmonInstance) - 1
		settings.AmonInstance = settings.AmonInstance[:cutOffLastCharLen]
	}

	return settings
}
