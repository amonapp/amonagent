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
)

var tr = &http.Transport{
	ResponseHeaderTimeout: time.Duration(3 * time.Second),
	TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // remove that from the final plugin
}

var client = &http.Client{Transport: tr}

// Collect - XXX
func Collect() error {

	// 	Active connections: 8
	// 	server accepts handled requests
	// 	 1156958 1156958 4491319
	// 	Reading: 0 Writing: 2 Waiting: 6
	u := "https://www.localhost/nginx_status"
	addr, err := url.Parse(u)
	if err != nil {
		return fmt.Errorf("Unable to parse address '%s': %s", u, err)
	}
	resp, err := client.Get(addr.String())
	if err != nil {
		return fmt.Errorf("error making HTTP request to %s: %s", addr.String(), err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s returned HTTP status %s", addr.String(), resp.Status)
	}
	r := bufio.NewReader(resp.Body)

	// Active connectionsmain
	_, err = r.ReadString(':')
	if err != nil {
		return err
	}
	line, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	active, err := strconv.ParseUint(strings.TrimSpace(line), 10, 64)
	if err != nil {
		return err
	}

	// Server accepts handled requests
	_, err = r.ReadString('\n')
	if err != nil {
		return err
	}
	line, err = r.ReadString('\n')
	if err != nil {
		return err
	}
	data := strings.SplitN(strings.TrimSpace(line), " ", 3)
	accepts, err := strconv.ParseUint(data[0], 10, 64)
	if err != nil {
		return err
	}
	handled, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return err
	}
	requests, err := strconv.ParseUint(data[2], 10, 64)
	if err != nil {
		return err
	}

	// Reading/Writing/Waiting
	line, err = r.ReadString('\n')
	if err != nil {
		return err
	}
	data = strings.SplitN(strings.TrimSpace(line), " ", 6)
	reading, err := strconv.ParseUint(data[1], 10, 64)
	if err != nil {
		return err
	}
	writing, err := strconv.ParseUint(data[3], 10, 64)
	if err != nil {
		return err
	}
	waiting, err := strconv.ParseUint(data[5], 10, 64)
	if err != nil {
		return err
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

	fmt.Println(fields)

	return nil

}
