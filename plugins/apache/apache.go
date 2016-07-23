package apache

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/amonapp/amonagent/plugins"
	"github.com/mitchellh/mapstructure"
)

// Start - XXX
func (a *Apache) Start() error { return nil }

// Stop - XXX
func (a *Apache) Stop() {}

// Description - XXX
func (a *Apache) Description() string {
	return "Read Apache status information (mod_status)"
}

// Config - XXX
type Config struct {
	StatusURL string `mapstructure:"status_url"`
}

var sampleConfig = `
#   Available config options:
#
#    {"status_url": "http://127.0.0.1/server-status?auto"}
#
# Config location: /etc/opt/amonagent/plugins-enabled/apache.conf
`

// SampleConfig - XXX
func (a *Apache) SampleConfig() string {
	return sampleConfig
}

// SetConfigDefaults - XXX
func (a *Apache) SetConfigDefaults() error {
	configFile, err := plugins.UmarshalPluginConfig("apache")
	if err != nil {
		log.WithFields(log.Fields{"plugin": "apache", "error": err.Error()}).Error("Can't read config file")
	}

	var config Config
	decodeError := mapstructure.Decode(configFile, &config)
	if decodeError != nil {
		log.WithFields(log.Fields{"plugin": "apache", "error": decodeError.Error()}).Error("Can't decode config file")
	}

	if len(config.StatusURL) == 0 {
		config.StatusURL = "http://127.0.0.1/server-status?auto"
	}

	a.Config = config

	return nil
}

// Apache - XXX
type Apache struct {
	Config Config
}

var tr = &http.Transport{
	ResponseHeaderTimeout: time.Duration(3 * time.Second),
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // make that optional
}

var client = &http.Client{Transport: tr}

// Tracking - XXX
var Tracking = map[string]string{
	"ReqPerSec":   "requests.request_per_second",
	"BytesPerSec": "bytes_per_second.bytes",
	"BytesPerReq": "bytes_per_request.bytes",
	"BusyWorkers": "workers.busy",
	"IdleWorkers": "workers.idle",
	"Scoreboard":  "scoreboard",
}

func (p PerformanceStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// PerformanceStruct - XXX
type PerformanceStruct struct {
	Gauges map[string]interface{} `json:"gauges"`
	// Counters map[string]interface{} `json:"counters"`
}

// Collect - XXX
func (a *Apache) Collect() (interface{}, error) {
	PerformanceStruct := PerformanceStruct{}
	a.SetConfigDefaults()

	addr, err := url.Parse(a.Config.StatusURL)
	resp, err := client.Get(addr.String())
	if err != nil {
		return PerformanceStruct, fmt.Errorf("error making HTTP request to %s: %s", addr.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return PerformanceStruct, fmt.Errorf("%s returned HTTP status %s", addr.String(), resp.Status)
	}

	sc := bufio.NewScanner(resp.Body)
	gauges := make(map[string]interface{})
	for sc.Scan() {
		line := sc.Text()

		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			rawKey := strings.Replace(parts[0], " ", "", -1)
			key, _ := Tracking[rawKey]

			part := strings.TrimSpace(parts[1])

			switch key {
			case "Scoreboard":
				fmt.Print(part)
				for field, value := range gatherScores(part) {
					gauges[field] = value
				}
			default:
				value, err := strconv.ParseFloat(part, 64)
				if err != nil {
					continue
				}
				if len(key) > 0 {
					gauges[key] = value
				}

			}
		}
	}
	PerformanceStruct.Gauges = gauges

	return PerformanceStruct, nil

}

func gatherScores(data string) map[string]interface{} {
	var waiting, open int = 0, 0
	var S, R, W, K, D, C, L, G, I int = 0, 0, 0, 0, 0, 0, 0, 0, 0

	for _, s := range strings.Split(data, "") {

		switch s {
		case "_":
			waiting++
		case "S":
			S++
		case "R":
			R++
		case "W":
			W++
		case "K":
			K++
		case "D":
			D++
		case "C":
			C++
		case "L":
			L++
		case "G":
			G++
		case "I":
			I++
		case ".":
			open++
		}
	}

	fields := map[string]interface{}{
		"scoreboard.waiting":      float64(waiting),
		"scoreboard.starting":     float64(S),
		"scoreboard.reading":      float64(R),
		"scoreboard.sending":      float64(W),
		"scoreboard.keepalive":    float64(K),
		"scoreboard.dnslookup":    float64(D),
		"scoreboard.closing":      float64(C),
		"scoreboard.logging":      float64(L),
		"scoreboard.finishing":    float64(G),
		"scoreboard.idle_cleanup": float64(I),
		"scoreboard.open":         float64(open),
	}

	return fields
}

func init() {
	plugins.Add("apache", func() plugins.Plugin {
		return &Apache{}
	})
}
