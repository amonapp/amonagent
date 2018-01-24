package jmx

import (
	"encoding/json"
	//"sync"

	log "github.com/Sirupsen/logrus"
	//"github.com/amonapp/amonagent/internal/util"
	"github.com/amonapp/amonagent/plugins"
	"strconv"
	"os"
	"os/exec"
	"bytes"
	"fmt"
)

// JMX - XXX
type JMX struct {
	Config Config
}

// PerformanceStruct - XXX
type PerformanceStruct struct {
	Gauges map[string]interface{} `json:"gauges"`
}

// Start - XXX
func (c *JMX) Start() error { return nil }

// Stop - XXX
func (c *JMX) Stop() {}

// Description - XXX
func (c *JMX) Description() string {
	return "Collects data from Java Applications"
}

var sampleConfig = `
#   Available config options:
#
#    [
#       {
#         "name": "Application",
#         "hostName": "localhost",
#         "port": 1234"
#       }
#    ]
#
# Config location: /etc/opt/amonagent/plugins-enabled/jmx.conf
`

// SampleConfig - XXX
func (c *JMX) SampleConfig() string {
	return sampleConfig
}

type Endpoint struct {
	Name string `json:"name"`
	HostName string `json:"hostName"`
	Port int `json:"port"`
}

type MJBJson struct {
	ThreadCount int64 `json:"java.lang:type=Threading ThreadCount"`
	DaemonThreadCount int64 `json:"java.lang:type=Threading DaemonThreadCount"`
}

// Config - XXX
type Config struct {
	Endpoints []Endpoint `json:"endpoints"`
}

// SetConfigDefaults - XXX
func (c *JMX) SetConfigDefaults() error {

	// Commands already set. For example - in the test suite
	if len(c.Config.Endpoints) > 0 {
		return nil
	}

	configFile, err := plugins.ReadPluginConfig("jmx")
	if err != nil {
		log.WithFields(log.Fields{
			"plugin": "jmx",
			"error":  err,
		}).Error("Can't read config file")

		return err
	}

	var Endpoints []Endpoint

	if e := json.Unmarshal(configFile, &Endpoints); e != nil {
		log.WithFields(log.Fields{"plugin": "jmx", "error": e.Error()}).Error("Can't decode JSON file")
		return e
	}

	c.Config.Endpoints = Endpoints

	return nil
}

// Collect - XXX
func (c *JMX) Collect() (interface{}, error) {
	c.SetConfigDefaults()
	PerformanceStruct := PerformanceStruct{}
	gauges := map[string]interface{}{

	}

	for _, v := range c.Config.Endpoints {
		var rawJson, err= runJar(v.HostName, v.Port)
		if err != nil {
			continue
		}
		var data MJBJson

		if e := json.Unmarshal([]byte(rawJson), &data); e != nil {
			log.WithFields(log.Fields{"plugin": "jmx", "error": e.Error()}).Error("Can't decode jmx response")
			continue
		}
		gauges[v.Name + "_jmx.threadCount"] = data.ThreadCount
		gauges[v.Name + "_jmx.daemonThreadCount"] = data.DaemonThreadCount

	}

	PerformanceStruct.Gauges = gauges

	return PerformanceStruct, nil
}

func init() {
	plugins.Add("jmx", func() plugins.Plugin {
		return &JMX{}
	})
}

// RunJar runs the embedded mjb.jar returning the output from STDOUT
func runJar(host string, port int) (string, error) {
	nport := strconv.Itoa(port)

	_, err := os.Stat("mjb.jar")

	if err != nil {
		return "", err
	}

	// Check that /usr/bin/java exists
	//_, err = os.Stat("java")

	//if err != nil {
	//	return "", err
	//}

	cmd := exec.Command("java", "-jar", "mjb.jar", host, nport)
	var out bytes.Buffer
	var erro bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &erro

	err = cmd.Run()

	if err != nil {
		return "", fmt.Errorf("%s %s", err.Error(), out.String())
	}

	return out.String(), nil
}
