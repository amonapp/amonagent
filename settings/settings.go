package settings

import (
	"encoding/json"
	"fmt"
	"os"
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
