package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// DefaultTimeOut - 10 seconds
var DefaultTimeOut = 10 * time.Second

// SendData - XXX
func SendData() {
	url := settings.Host + "/api/system/golang"

	fmt.Println("URL:>", url)

	metricsBytes, err := json.Marshal(c)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(metricsBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: DefaultTimeOut}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
