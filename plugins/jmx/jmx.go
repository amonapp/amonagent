package jmx

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/amonapp/amonagent/plugins"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
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
#         "port": 1234
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
	Name     string `json:"name"`
	HostName string `json:"hostName"`
	Port     int    `json:"port"`
}

type MJBJson struct {
	ThreadCount                 int64 `json:"java.lang:type=Threading ThreadCount"`
	DaemonThreadCount           int64 `json:"java.lang:type=Threading DaemonThreadCount"`
	HeapMemoryUsageMax          int64 `json:"java.lang:type=Memory HeapMemoryUsage max"`
	HeapMemoryUsageInit         int64 `json:"java.lang:type=Memory HeapMemoryUsage init"`
	HeapMemoryUsageCommitted    int64 `json:"java.lang:type=Memory HeapMemoryUsage committed"`
	HeapMemoryUsageUsed         int64 `json:"java.lang:type=Memory HeapMemoryUsage used"`
	NonHeapMemoryUsageMax       int64 `json:"java.lang:type=Memory NonHeapMemoryUsage max"`
	NonHeapMemoryUsageInit      int64 `json:"java.lang:type=Memory NonHeapMemoryUsage init" `
	NonHeapMemoryUsageCommitted int64 `json:"java.lang:type=Memory NonHeapMemoryUsage committed" `
	NonHeapMemoryUsageUsed      int64 `json:"java.lang:type=Memory NonHeapMemoryUsage used" `
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
	e := ensureJarExists()
	if e != nil {
		return "", e
	}
	c.SetConfigDefaults()
	PerformanceStruct := PerformanceStruct{}
	gauges := map[string]interface{}{}
	resultChan := make(chan map[string]interface{}, len(c.Config.Endpoints))
	var wg sync.WaitGroup

	for _, v := range c.Config.Endpoints {
		wg.Add(1)

		go func(endpoint Endpoint) {
			var rawJson, err = runJar(endpoint.HostName, endpoint.Port)
			if err != nil {
				log.WithFields(log.Fields{"plugin": "jmx", "error": e.Error()}).Error("Could not run command")
				defer wg.Done()
				return
			}
			var data MJBJson

			if e := json.Unmarshal([]byte(rawJson), &data); e != nil {
				log.WithFields(log.Fields{"plugin": "jmx", "error": e.Error()}).Error("Can't decode jmx response")
				defer wg.Done()
				return
			}

			m := map[string]interface{}{}

			m[endpoint.Name+"_threads.count"] = data.ThreadCount
			m[endpoint.Name+"_threads.daemonCount"] = data.DaemonThreadCount
			m[endpoint.Name+"_heapMemory.committed"] = data.HeapMemoryUsageCommitted
			m[endpoint.Name+"_heapMemory.init"] = data.HeapMemoryUsageInit
			m[endpoint.Name+"_heapMemory.max"] = data.HeapMemoryUsageMax
			m[endpoint.Name+"_heapMemory.used"] = data.HeapMemoryUsageUsed
			m[endpoint.Name+"_heapMemoryPercentUsed"] = float64(data.HeapMemoryUsageUsed) / float64(data.HeapMemoryUsageMax)
			m[endpoint.Name+"_nonHeapMemory.committed"] = data.NonHeapMemoryUsageCommitted
			m[endpoint.Name+"_nonHeapMemory.init"] = data.NonHeapMemoryUsageInit
			m[endpoint.Name+"_nonHeapMemory.max"] = data.NonHeapMemoryUsageMax
			m[endpoint.Name+"_nonHeapMemory.used"] = data.NonHeapMemoryUsageUsed

			resultChan <- m

			defer wg.Done()
		}(v)
	}
	wg.Wait()
	close(resultChan)

	for m := range resultChan {
		for k, v := range m {
			gauges[k] = v
		}
	}

	PerformanceStruct.Gauges = gauges

	return PerformanceStruct, nil
}

func init() {
	plugins.Add("jmx", func() plugins.Plugin {
		return &JMX{}
	})
}

// runJar runs the embedded jar returning the output from STDOUT
func runJar(host string, port int) (string, error) {
	nport := strconv.Itoa(port)

	cmd := exec.Command("java", "-jar", TmpJarFile, host, nport)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("%s %s", err.Error(), out.String())
	}

	return out.String(), nil
}

// ensureJarExists creates the jar file on the file system
func ensureJarExists() error {
	_, err := os.Stat(TmpJarFile)
	if err != nil {
		data, err := Asset("data/mjb.jar")
		if err != nil {
			return err
		}
		os.Mkdir(filepath.Dir(TmpJarFile), os.FileMode(0755))
		err = ioutil.WriteFile(TmpJarFile, data, os.FileMode(0644))
		if err != nil {
			return err
		}
	}
	return nil
}
