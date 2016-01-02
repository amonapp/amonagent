package nginx

// Tracks basic nginx metrics via the status module
// 	* number of connections
// 	* number of requets per second
// 	Requires nginx to have the status option compiled.
// 	See http://wiki.nginx.org/HttpStubStatusModule for more details

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/amonapp/amonagent/logging"
	"github.com/amonapp/amonagent/plugins"
	"github.com/mitchellh/mapstructure"
)

var pluginLogger = logging.GetLogger("amonagent.nginx")

var tr = &http.Transport{
	ResponseHeaderTimeout: time.Duration(3 * time.Second),
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // move this to a config option
}

var client = &http.Client{Transport: tr}

// PerformanceStruct - XXX
type PerformanceStruct struct {
	Gauges map[string]interface{} `json:"gauges"`
}

// Config - XXX
type Config struct {
	StatusURL string `mapstructure:"status_url"`
}

var sampleConfig = `
#   Available config options:
#
#    {"status_url": "http://127.0.0.1/nginx_status"}
#
# Config location: /etc/opt/amonagent/plugins-enabled/nginx.conf
`

// SampleConfig - XXX
func (n *Nginx) SampleConfig() string {
	return sampleConfig
}

// SetConfigDefaults - XXX
func (n *Nginx) SetConfigDefaults(configPath string) error {
	c, err := plugins.ReadConfigPath(configPath)
	if err != nil {
		fmt.Printf("Can't read config file: %s %v\n", configPath, err)
	}
	var config Config
	decodeError := mapstructure.Decode(c, &config)
	if decodeError != nil {
		fmt.Print("Can't decode config file", decodeError.Error())
	}

	n.Config = config

	return nil
}

// Nginx - XXX
type Nginx struct {
	Config Config
}

// Description - XXX
func (n *Nginx) Description() string {
	return "Read metrics from a nginx server"
}

// Collect - XXX
func (n *Nginx) Collect(configPath string) (interface{}, error) {
	n.SetConfigDefaults(configPath)
	PerformanceStruct := PerformanceStruct{}
	// 	Active connections: 8
	// 	server accepts handled requests
	// 	 1156958 1156958 4491319
	// 	Reading: 0 Writing: 2 Waiting: 6

	addr, err := url.Parse(n.Config.StatusURL)
	if err != nil {
		pluginLogger.Errorf("Unable to parse address '%s': %s", n.Config.StatusURL, err)
		return PerformanceStruct, err

	}
	resp, err := client.Get(addr.String())
	if err != nil {
		pluginLogger.Errorf("error making HTTP request to %s: %s", addr.String(), err)
		return PerformanceStruct, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		pluginLogger.Errorf("%s returned HTTP status %s", addr.String(), resp.Status)
		return PerformanceStruct, err
	}
	r := bufio.NewReader(resp.Body)

	// Active connections
	_, err = r.ReadString(':')
	if err != nil {
		pluginLogger.Errorf("Can't parse active connections stats%s", err)
		return PerformanceStruct, err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		pluginLogger.Errorf("Can't read stats page %s", err)
		return PerformanceStruct, err
	}
	active, err := strconv.ParseUint(strings.TrimSpace(line), 10, 64)
	if err != nil {
		pluginLogger.Errorf("Can't read active connections stats%s", err)
	}

	// Server accepts handled requests
	_, err = r.ReadString('\n')
	if err != nil {
		pluginLogger.Errorf("Can't read accepts, handled requests stats%s", err)
	}
	line, err = r.ReadString('\n')
	if err != nil {
		pluginLogger.Errorf("Can't read stats page %s", err)
	}
	data := strings.SplitN(strings.TrimSpace(line), " ", 3)
	accepts, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		pluginLogger.Errorf("Can't read accepts stats%s", err)
	}
	handled, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		pluginLogger.Errorf("Can't read handled stats%s", err)
	}
	requests, err := strconv.ParseUint(data[2], 10, 64)
	if err != nil {
		pluginLogger.Errorf("Can't read requests stats%s", err)
	}

	// Reading/Writing/Waiting
	line, err = r.ReadString('\n')
	if err != nil {
		pluginLogger.Errorf("Can't read Reading/Writing/Waiting stats%s", err)
	}
	data = strings.SplitN(strings.TrimSpace(line), " ", 6)
	reading, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		pluginLogger.Errorf("Can't read Reading stats%s", err)
	}
	writing, err := strconv.ParseUint(data[3], 10, 64)
	if err != nil {
		pluginLogger.Errorf("Can't read Writing stats%s", err)
	}
	waiting, err := strconv.ParseUint(data[5], 10, 64)
	if err != nil {
		pluginLogger.Errorf("Can't read Waiting stats%s", err)
	}

	requestPerSecond := requests / handled

	gauges := map[string]interface{}{
		"connections.active":        active,
		"connections_total.accepts": accepts,
		"connections_total.handled": handled,
		"requests_per_second":       requestPerSecond,
		"connections.reading":       reading,
		"connections.writing":       writing,
		"connections.waiting":       waiting,
	}
	PerformanceStruct.Gauges = gauges

	return PerformanceStruct, nil

}
func init() {
	plugins.Add("nginx", func() plugins.Plugin {
		return &Nginx{}
	})
}
