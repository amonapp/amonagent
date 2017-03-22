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

	ConfiguredPlugins []plugins.ConfiguredPlugin
}

// TestPlugin - XXX
func (a *Agent) TestPlugin(pluginName string) error {

	_, err := plugins.GetConfigPath(pluginName)

	if err != nil {
		fmt.Printf("Can't get config file for plugin: %s", err)
	}

	for _, p := range a.ConfiguredPlugins {
		if pluginName == p.Name {
			start := time.Now()

			PluginResult, err := p.Plugin.Collect()
			if err != nil {
				fmt.Printf("Can't get stats for plugin: %s", err)
			}
			fmt.Println(PluginResult)

			elapsed := time.Since(start)
			fmt.Printf("\nExecuted in \033[92m%s\033[0m", elapsed)
		}
	}

	return nil

}

// Test - XXX
func (a *Agent) Test(config settings.Struct) error {

	allMetrics := collectors.CollectAllData(a.ConfiguredPlugins)

	ProcessesData := collectors.CollectProcessData()
	SystemData := collectors.CollectSystemData()
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

	for _, p := range a.ConfiguredPlugins {
		start := time.Now()
		PluginResult, err := p.Plugin.Collect()
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

	err := remote.SendData(allMetrics, true)
	if err != nil {
		return fmt.Errorf("%s\n", err.Error())
	}

	return nil
}

// GatherAndSend - XXX
func (a *Agent) GatherAndSend(debug bool) error {

	allMetrics := collectors.CollectAllData(a.ConfiguredPlugins)

	log.Infof("Metrics collected (Interval:%s)\n", a.Interval)

	err := remote.SendData(allMetrics, debug)
	if err != nil {
		return fmt.Errorf("Can't connect to the Amon API on %s\n", err.Error())
	}

	return nil
}

// NewAgent - XXX
func NewAgent(config settings.Struct) (*Agent, error) {

	var configuredPlugins = []plugins.ConfiguredPlugin{}

	EnabledPlugins, _ := plugins.GetAllEnabledPlugins()
	for _, p := range EnabledPlugins {
		creator, _ := plugins.Plugins[p.Name]
		plugin := creator()

		t := plugins.ConfiguredPlugin{Name: p.Name, Plugin: plugin}
		configuredPlugins = append(configuredPlugins, t)

	}

	agent := &Agent{
		Interval:          time.Duration(config.Interval) * time.Second,
		ConfiguredPlugins: configuredPlugins,
	}

	return agent, nil
}

// Run runs the agent daemon, gathering every Interval
func (a *Agent) Run(shutdown chan struct{}, debug bool) error {

	log.Infof("Agent Config: Interval:%s\n", a.Interval)

	ticker := time.NewTicker(a.Interval)

	for _, p := range a.ConfiguredPlugins {

		if err := p.Plugin.Start(); err != nil {
			log.WithFields(log.Fields{
				"plugin": p.Name,
				"error":  err.Error(),
			}).Error("Service for plugin failed to start, exiting")

		}

		defer p.Plugin.Stop()

	}

	for {
		select {
		case <-shutdown:
			log.Info("Shutting down Amon Agent")
			ticker.Stop()

			return nil
		case <-ticker.C:
			if err := a.GatherAndSend(debug); err != nil {
				log.Infof("Can not collect and send metrics, exiting: %s\n", err.Error())
			}
		}
	}
}
