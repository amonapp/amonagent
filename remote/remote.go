package remote

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/martinrusev/amonagent/core"
)

// DefaultTimeOut - 10 seconds
var DefaultTimeOut = 10 * time.Second

// SendData - XXX
func SendData(data interface{}) {
	settings := core.Settings()
	url := settings.Host + "/api/system/" + settings.ServerKey

	fmt.Println("URL:>", url)

	JSONBytes, err := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(JSONBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: DefaultTimeOut}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
