package nginx

// Tracks basic nginx metrics via the status module
// 	* number of connections
// 	* number of requets per second
// 	Requires nginx to have the status option compiled.
// 	See http://wiki.nginx.org/HttpStubStatusModule for more details

import (
	"bufio"
	"crypto/tls"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/amonapp/amonagent/logging"
	"github.com/amonapp/amonagent/plugins"
)

var pluginLogger = logging.GetLogger("amonagent.nginx")

var tr = &http.Transport{
	ResponseHeaderTimeout: time.Duration(3 * time.Second),
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // remove that from the final plugin
}

var client = &http.Client{Transport: tr}

// PerformanceStruct - XXX
type PerformanceStruct struct {
	Gauges map[string]interface{} `json:"gauges"`
}

// Nginx - XXX
type Nginx struct {
}

// Description - XXX
func (n *Nginx) Description() string {
	return "Read metrics from a nginx server"
}

// Collect - XXX
func (n *Nginx) Collect() (interface{}, error) {
	PerformanceStruct := PerformanceStruct{}
	// 	Active connections: 8
	// 	server accepts handled requests
	// 	 1156958 1156958 4491319
	// 	Reading: 0 Writing: 2 Waiting: 6
	u := "https://www.localhost/nginx_status"
	addr, err := url.Parse(u)
	if err != nil {
		pluginLogger.Errorf("Unable to parse address '%s': %s", u, err)
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

	fields := map[string]interface{}{
		"connections.active":        active,
		"connections_total.accepts": accepts,
		"connections_total.handled": handled,
		"requests_per_second":       requestPerSecond,
		"connections.reading":       reading,
		"connections.writing":       writing,
		"connections.waiting":       waiting,
	}
	PerformanceStruct.Gauges = fields

	return PerformanceStruct, nil

}
func init() {
	plugins.Add("nginx", func() plugins.Plugin {
		return &Nginx{}
	})
}
