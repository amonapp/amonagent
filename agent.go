package amonagent

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/amonapp/amonagent/collectors"
	"github.com/amonapp/amonagent/internal/remote"
	"github.com/amonapp/amonagent/internal/settings"
	"github.com/amonapp/amonagent/plugins"
)

// Agent - XXX
type Agent struct {
	// Interval at which to gather information
	Interval time.Duration
}

// Test - XXX
func (a *Agent) Test(config settings.Struct) error {

	allMetrics := collectors.CollectAllData()

	ProcessesData := collectors.CollectProcessData()
	SystemData := collectors.CollectSystemData()
	EnabledPlugins, _ := plugins.GetAllEnabledPlugins()
	HostData := collectors.CollectHostData()

	fmt.Println("\n------------------")
	fmt.Println("\033[92mSystem Metrics: \033[0m")
	fmt.Println("")
	fmt.Println(SystemData)
	fmt.Println("\n------------------")

	fmt.Println("\n------------------")
	fmt.Println("\033[92mProcess Metrics: \033[0m")
	fmt.Println("")
	fmt.Println(ProcessesData)
	fmt.Println("\n------------------")

	fmt.Println("\n------------------")
	fmt.Println("\033[92mPlugins: \033[0m")
	fmt.Println("")

	for _, p := range EnabledPlugins {

		creator, _ := plugins.Plugins[p.Name]
		plugin := creator()
		start := time.Now()
		PluginResult, err := plugin.Collect()
		if err != nil {
			log.Errorf("Can't get stats for plugin: %s", err)
		}

		fmt.Println("\n------------------")
		fmt.Print("\033[92mPlugin: ")
		fmt.Print(p.Name)
		fmt.Print("\033[0m \n")
		fmt.Println(PluginResult)

		elapsed := time.Since(start)
		fmt.Printf("\n Executed in %s", elapsed)

	}

	fmt.Println("\n------------------")
	fmt.Println("\033[92mHost Data: \033[0m")
	fmt.Println("")
	fmt.Println(HostData)
	fmt.Println("\n------------------")

	fmt.Println("\033[92mTesting settings: \033[0m")
	fmt.Println("")
	machineID := collectors.GetOrCreateMachineID()

	if len(machineID) == 0 && len(config.ServerKey) == 0 {
		fmt.Println("Can't find Machine ID (looking in /etc/opt/amonagent/machine-id).")
		fmt.Println("To solve this problem, run the following command:")
		fmt.Println("---")
		fmt.Println("amonagent -machineid")
		fmt.Println("---")

	} else {
		fmt.Println("Settings OK")
	}

	fmt.Println("\n------------------")

	// url := remote.SystemURL()

	err := remote.SendData(allMetrics, true)
	if err != nil {
		return fmt.Errorf("%s\n", err.Error())
	}

	return nil
}

// GatherAndSend - XXX
func (a *Agent) GatherAndSend(debug bool) error {
	allMetrics := collectors.CollectAllData()
	log.Infof("Metrics collected (Interval:%s)\n", a.Interval)

	err := remote.SendData(allMetrics, debug)
	if err != nil {
		return fmt.Errorf("Can't connect to the Amon API on %s\n", err.Error())
	}

	return nil
}

// NewAgent - XXX
func NewAgent(config settings.Struct) (*Agent, error) {
	agent := &Agent{
		Interval: time.Duration(config.Interval) * time.Second,
	}

	return agent, nil
}

// Run runs the agent daemon, gathering every Interval
func (a *Agent) Run(shutdown chan struct{}, debug bool) error {

	log.Infof("Agent Config: Interval:%s\n", a.Interval)

	ticker := time.NewTicker(a.Interval)
	defer ticker.Stop()

	EnabledPlugins, _ := plugins.GetAllEnabledPlugins()
	for _, p := range EnabledPlugins {
		creator, _ := plugins.Plugins[p.Name]
		plugin := creator()
		if err := plugin.Start(); err != nil {
			log.Printf("Service for plugin %s failed to start, exiting\n%s\n",
				p.Name, err.Error())
			return err
		}
		defer plugin.Stop()
	}

	for {
		if err := a.GatherAndSend(debug); err != nil {
			log.Infof("Can not collect and send metrics, exiting: %s\n", err.Error())
		}
		select {
		case <-shutdown:
			return nil
		case <-ticker.C:
			continue
		}
	}
}
