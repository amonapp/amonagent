package collectors

import (
	"encoding/json"
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

// CollectPlugins - XXX
func CollectPlugins() (interface{}, interface{}) {
	PluginResults := make(map[string]interface{})
	CheckResults := make(map[string]interface{})
	var wg sync.WaitGroup
	EnabledPlugins, _ := plugins.GetAllEnabledPlugins()
	for _, p := range EnabledPlugins {
		creator, ok := plugins.Plugins[p.Name]
		if ok {
			wg.Add(1)
			plugin := creator()

			go func(p plugins.PluginConfig) {
				defer wg.Done()
				PluginResult, err := plugin.Collect(p.Path)
				if err != nil {
					CollectorLogger.Errorf("Can't get stats for plugin: %s", err)

				}
				if p.Name == "checks" {
					CheckResults["checks"] = PluginResult
				} else {
					PluginResults[p.Name] = PluginResult
				}
			}(p)

		} else {
			CollectorLogger.Errorf("Non existing plugin: %s", p.Name)
		}
	}
	wg.Wait()

	return PluginResults, CheckResults
}

// CollectSystem - XXX
func CollectSystem() AllMetricsStruct {
	networkUsage, _ := NetworkUsage()
	cpuUsage := CPUUsage()
	Load := LoadAverage()
	diskUsage, _ := DiskUsage()
	Uptime := Uptime()
	Memory := MemoryUsage()
	Processes, _ := Processes()
	Plugins, Checks := CollectPlugins()

	System := SystemDataStruct{
		CPU:     cpuUsage,
		Network: networkUsage,
		Disk:    diskUsage,
		Load:    Load,
		Uptime:  Uptime,
		Memory:  Memory,
	}

	// Load settings
	settings := settings.Settings()

	host := Host()
	machineID := MachineID()
	distro := Distro()
	ip := IPAddress()
	InstanceID := CloudID()

	hoststruct := HostDataStruct{
		Host:       host,
		MachineID:  machineID,
		Distro:     distro,
		IPAddress:  ip,
		ServerKey:  settings.ServerKey,
		InstanceID: InstanceID,
	}

	allMetrics := AllMetricsStruct{
		System:    System,
		Processes: Processes,
		Host:      hoststruct,
		Plugins:   Plugins,
		Checks:    Checks,
	}

	return allMetrics
}
