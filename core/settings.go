package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// SettingsStruct -XXX
type SettingsStruct struct {
	AmonInstance string `json:"amon_instance"`
	Interval     int    `json:"interval"`
	APIKey       string `json:"api_key"`
}

// Settings -XXX
func Settings() SettingsStruct {

	var settings SettingsStruct

	configFile, err := os.Open("/etc/opt/amonagent/amonagent.conf")
	if err != nil {
		fmt.Print("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(&settings); err != nil {
		fmt.Print("parsing config file", err.Error())
	}

	return settings
}
