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

	allMetrics := collectors.CollectSystem()

	fmt.Println("Collected Metrics: ")
	fmt.Println(allMetrics)
	fmt.Println("\n------------------")

	url := remote.SystemURL()
	fmt.Printf("\nSending data to %s ", url)

	err := remote.SendData(allMetrics)
	if err != nil {
		return fmt.Errorf("Can't connect to the Amon API on %s\n", err.Error())
	}

	return nil
}

// GatherAndSend - XXX
func (a *Agent) GatherAndSend() error {

	allMetrics := collectors.CollectSystem()
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

	for {

		if err := a.GatherAndSend(); err != nil {
			agentLogger.Info("Flusher routine failed, exiting: %s\n", err.Error())
		} else {
			agentLogger.Info("Collecting and sending data:%s\n", a.Interval)
		}
		select {
		case <-shutdown:
			return nil
		case <-ticker.C:
			continue
		}
	}
}
