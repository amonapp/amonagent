package amonagent

import (
	"fmt"
	"time"

	"github.com/amonapp/amonagent/collectors"
	"github.com/amonapp/amonagent/logging"
	"github.com/amonapp/amonagent/remote"
	"github.com/amonapp/amonagent/settings"
)

var agentLogger = logging.GetLogger("amonagent.log")

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
	Plugins, Checks := collectors.CollectPluginsData()
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
	fmt.Println(Plugins)
	fmt.Println("\n------------------")

	fmt.Println("\n------------------")
	fmt.Println("\033[92mChecks: \033[0m")
	fmt.Println("")
	fmt.Println(Checks)
	fmt.Println("\n------------------")

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

	err := remote.SendData(allMetrics)
	if err != nil {
		return fmt.Errorf("%s\n", err.Error())
	}

	return nil
}

// GatherAndSend - XXX
func (a *Agent) GatherAndSend() error {
	allMetrics := collectors.CollectAllData()
	agentLogger.Info("Metrics collected (Interval:%s)\n", a.Interval)

	err := remote.SendData(allMetrics)
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
func (a *Agent) Run(shutdown chan struct{}) error {

	agentLogger.Info("Agent Config: Interval:%s\n", a.Interval)

	ticker := time.NewTicker(a.Interval)
	defer ticker.Stop()

	for {
		if err := a.GatherAndSend(); err != nil {
			agentLogger.Info("Flusher routine failed, exiting: %s\n", err.Error())
		}
		select {
		case <-shutdown:
			return nil
		case <-ticker.C:
			continue
		}
	}
}
