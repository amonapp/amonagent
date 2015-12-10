package remote

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/martinrusev/amonagent/core"
	"github.com/prometheus/common/log"
)

// DefaultTimeOut - 10 seconds
var DefaultTimeOut = 10 * time.Second

// SendData - XXX
func SendData(data interface{}) {
	settings := core.Settings()
	url := settings.AmonInstance + "/api/system/v2/?api_key=" + settings.APIKey

	JSONBytes, err := json.Marshal(data)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(JSONBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: DefaultTimeOut}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Can't connect to the Amon API on %s" + url)
	}
	defer resp.Body.Close()

}
