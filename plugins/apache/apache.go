package apache

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var tr = &http.Transport{
	ResponseHeaderTimeout: time.Duration(3 * time.Second),
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // remove that from the final plugin
}

var client = &http.Client{Transport: tr}

// Tracking - XXX
var Tracking = map[string]string{
	"ReqPerSec":   "requests.request_per_second",
	"BytesPerSec": "bytes.per_second",
	"BytesPerReq": "bytes.per_request",
	"BusyWorkers": "workers.busy",
	"IdleWorkers": "workers.idle",
}

// Collect - XXX
func Collect() error {
	u := "http://127.0.0.1:81/server-status?auto"
	addr, err := url.Parse(u)
	resp, err := client.Get(addr.String())
	if err != nil {
		return fmt.Errorf("error making HTTP request to %s: %s", addr.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned HTTP status %s", addr.String(), resp.Status)
	}

	sc := bufio.NewScanner(resp.Body)
	fields := make(map[string]interface{})
	for sc.Scan() {
		line := sc.Text()

		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			rawKey := strings.Replace(parts[0], " ", "", -1)
			key, _ := Tracking[rawKey]

			part := strings.TrimSpace(parts[1])

			switch key {
			case "Scoreboard":
				for field, value := range gatherScores(part) {
					fields[field] = value
				}
			default:
				value, err := strconv.ParseFloat(part, 64)
				if err != nil {
					continue
				}
				fields[key] = value
			}
		}
	}

	return nil

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
