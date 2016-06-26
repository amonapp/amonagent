package collectors

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/amonapp/amonagent/logging"
	"github.com/amonapp/amonagent/plugins"
	"github.com/amonapp/amonagent/settings"
)

// CollectorLogger - XXX
var CollectorLogger = logging.GetLogger("amonagent.collector")

func (p SystemDataStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

func (p AllMetricsStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)

}

func (p HostDataStruct) String() string {
	s, _ := json.Marshal(p)
	return string(s)
}

// AllMetricsStruct -XXX
type AllMetricsStruct struct {
	System    SystemDataStruct `json:"system"`
	Processes ProcessesList    `json:"processes"`
	Host      HostDataStruct   `json:"host"`
	Plugins   interface{}      `json:"plugins"`
	Checks    interface{}      `json:"checks"`
}

// HostDataStruct -XXX
type HostDataStruct struct {
	Host       string       `json:"host"`
	MachineID  string       `json:"machineid"`
	ServerKey  string       `json:"server_key"`
	Distro     DistroStruct `json:"distro"`
	IPAddress  string       `json:"ip_address"`
	InstanceID string       `json:"instance_id"`
}

// SystemDataStruct - collect all system metrics
type SystemDataStruct struct {
	CPU     CPUUsageStruct   `json:"cpu"`
	Network NetworkUsageList `json:"network"`
	Disk    DiskUsageList    `json:"disk"`
	Load    LoadStruct       `json:"loadavg"`
	Uptime  string           `json:"uptime"`
	Memory  MemoryStruct     `json:"memory"`
}

// CollectPluginsData - XXX
func CollectPluginsData() (interface{}, interface{}) {
	PluginResults := make(map[string]interface{})
	var CheckResults interface{}
	var wg sync.WaitGroup
	EnabledPlugins, _ := plugins.GetAllEnabledPlugins()

	resultChan := make(chan interface{}, len(EnabledPlugins))

	for _, p := range EnabledPlugins {
		wg.Add(1)
		creator, _ := plugins.Plugins[p.Name]
		plugin := creator()

		go func(p plugins.PluginConfig) {
			PluginResult, err := plugin.Collect(p.Path)
			if err != nil {
				CollectorLogger.Errorf("Can't get stats for plugin: %s", err)

			}

			resultChan <- PluginResult
			defer wg.Done()
		}(p)

		// if p.Name == "checks" {
		// 	CheckResults = resultChan
		// } else {
		// 	PluginResults[p.Name] = resultChan
		// }

	}

	wg.Wait()
	close(resultChan)

	for i := range resultChan {
		fmt.Println(i)
		fmt.Println("---------------------------------------------------")
		// result = append(result, i)
	}

	return PluginResults, CheckResults
}

// CollectHostData - XXX
func CollectHostData() HostDataStruct {

	host := Host()
	// Load settings
	settings := settings.Settings()

	var machineID string
	var InstanceID string
	var ip string
	var distro DistroStruct

	machineID = GetOrCreateMachineID()
	InstanceID = CloudID()
	ip = IPAddress()
	distro = Distro()

	hoststruct := HostDataStruct{
		Host:       host,
		MachineID:  machineID,
		Distro:     distro,
		IPAddress:  ip,
		ServerKey:  settings.ServerKey,
		InstanceID: InstanceID,
	}

	return hoststruct
}

// CollectSystemData - XXX
func CollectSystemData() SystemDataStruct {
	var networkUsage NetworkUsageList
	var cpuUsage CPUUsageStruct
	var diskUsage DiskUsageList
	var memoryUsage MemoryStruct
	var UptimeString string
	var Load LoadStruct

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		networkUsage, _ = NetworkUsage()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		cpuUsage = CPUUsage()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		diskUsage, _ = DiskUsage()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		memoryUsage = MemoryUsage()
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		UptimeString = Uptime()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Load = LoadAverage()
	}()

	wg.Wait()

	SystemData := SystemDataStruct{
		CPU:     cpuUsage,
		Network: networkUsage,
		Disk:    diskUsage,
		Load:    Load,
		Uptime:  UptimeString,
		Memory:  memoryUsage,
	}

	return SystemData

}

// CollectProcessData - XXX
func CollectProcessData() ProcessesList {
	var ProcessesUsage ProcessesList
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		ProcessesUsage, _ = Processes()
	}()

	wg.Wait()

	return ProcessesUsage
}

// CollectAllData - XXX
func CollectAllData() AllMetricsStruct {

	ProcessesData := CollectProcessData()
	SystemData := CollectSystemData()
	Plugins, Checks := CollectPluginsData()
	HostData := CollectHostData()

	allMetrics := AllMetricsStruct{
		System:    SystemData,
		Processes: ProcessesData,
		Host:      HostData,
		Plugins:   Plugins,
		Checks:    Checks,
	}

	return allMetrics
}
