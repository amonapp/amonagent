package core

import (
	"encoding/json"
	"fmt"
	"os"
)

// SettingsStruct -XXX
type SettingsStruct struct {
	Host        string `json:"host"`
	CheckPeriod int    `json:"system_check_period"`
	ServerKey   string `json:"server_key"`
}

// Settings -XXX
func Settings() SettingsStruct {

	var settings SettingsStruct

	configFile, err := os.Open("/etc/amon-agent.conf")
	if err != nil {
		fmt.Print("opening config file", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)

	if err = jsonParser.Decode(&settings); err != nil {
		fmt.Print("parsing config file", err.Error())
	}

	return settings
}
