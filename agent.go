package amonagent

import (
	"log"
	"time"

	"github.com/martinrusev/amonagent/collectors"
	"github.com/martinrusev/amonagent/core"
	"github.com/martinrusev/amonagent/remote"
)

// Agent - XXX
type Agent struct {
	// Interval at which to gather information
	Interval time.Duration
}

// GatherAndSend - XXX
func (a *Agent) GatherAndSend() error {

	allMetrics := collectors.CollectSystem()
	remote.SendData(allMetrics)
	return nil
}

// NewAgent - XXX
func NewAgent(config core.SettingsStruct) (*Agent, error) {
	agent := &Agent{
		Interval: 10 * time.Second,
	}

	return agent, nil
}

// Run runs the agent daemon, gathering every Interval
func (a *Agent) Run(shutdown chan struct{}) error {

	log.Printf("Agent Config: Interval:%s\n", a.Interval)

	ticker := time.NewTicker(a.Interval)

	for {

		if err := a.GatherAndSend(); err != nil {
			log.Printf("Flusher routine failed, exiting: %s\n", err.Error())
		}
		select {
		case <-shutdown:
			return nil
		case <-ticker.C:
			continue
		}
	}
}
