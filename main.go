package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/martinrusev/amonagent/collectors"
	"github.com/martinrusev/amonagent/logging"
)

// AmonAgentLogger for the main file
var AmonAgentLogger = logging.GetLogger("amonagent")
var settings struct {
	Host        string `json:"host"`
	CheckPeriod int    `json:"system_check_period"`
	ServerKey   string `json:"server_key"`
}

// DefaultTimeOut - 10 seconds
var DefaultTimeOut = 10 * time.Second

func (p SystemDataStruct) String() string {
	s, err := json.Marshal(p)

	if err != nil {
		fmt.Print("Json Error", err.Error())
	}
	return string(s)
}

// SystemDataStruct - collect all system metrics
type SystemDataStruct struct {
	CPU     collectors.CPUUsageStruct   `json:"cpu"`
	Network collectors.NetworkUsageList `json:"network"`
	Disk    collectors.DiskUsageList    `json:"disk"`
}

// Just for testing
func main() {

	// n := collectors.Distro()
	//
	// fmt.Println(n)
	// cpu, _ := collectors.CP
	networkUsage, _ := collectors.NetworkUsage()
	cpuUsage := collectors.CPUUsage()
	diskUsage, _ := collectors.DiskUsage()
	fmt.Println(diskUsage)

	s, err := json.Marshal(diskUsage)
	if err != nil {
		fmt.Print(err.Error())
	}
	fmt.Println(s)
	c := SystemDataStruct{
		CPU:     cpuUsage,
		Network: networkUsage,
		// Disk:    diskUsage,
	}

	fmt.Print(c)

	// system_data_dict = {
	// 		'memory': get_memory_info(),
	// 		'cpu': get_cpu_utilization(),
	// 		'disk': disk_check.check(),
	// 		'network': get_network_traffic(),
	// 		'loadavg': get_load_average(),
	// 		'uptime': get_uptime(),
	// 	}

	// configFile, err := os.Open("/etc/amon-agent.conf")
	// if err != nil {
	// 	fmt.Print("opening config file", err.Error())
	// }
	//
	// jsonParser := json.NewDecoder(configFile)
	//
	// if err = jsonParser.Decode(&settings); err != nil {
	// 	fmt.Print("parsing config file", err.Error())
	// }
	//
	// url := settings.Host + "/api/system/golang"
	//
	// fmt.Println("URL:>", url)
	//
	// metricsBytes, err := json.Marshal(c)
	//
	// req, err := http.NewRequest("POST", url, bytes.NewBuffer(metricsBytes))
	// req.Header.Set("Content-Type", "application/json")
	//
	// client := &http.Client{Timeout: DefaultTimeOut}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	//
	// fmt.Println("response Status:", resp.Status)
	// fmt.Println("response Headers:", resp.Header)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Println("response Body:", string(body))

}
