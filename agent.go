package amonagent

import (
	"fmt"
	"log"
	"time"

	"github.com/martinrusev/amonagent/collectors"
	"github.com/martinrusev/amonagent/remote"
	"github.com/martinrusev/amonagent/settings"
)

// Agent - XXX
type Agent struct {
	// Interval at which to gather information
	Interval time.Duration
}

// Test - XXX
func (a *Agent) Test() error {

	allMetrics := collectors.CollectSystem()

	fmt.Println(allMetrics)
	return nil
}

// GatherAndSend - XXX
func (a *Agent) GatherAndSend() error {

	allMetrics := collectors.CollectSystem()
	remote.SendData(allMetrics)
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

	log.Printf("Agent Config: Interval:%s\n", a.Interval)

	ticker := time.NewTicker(a.Interval)

	for {

		if err := a.GatherAndSend(); err != nil {
			log.Printf("Flusher routine failed, exiting: %s\n", err.Error())
		} else {
			log.Printf("Collecting and sending data:%s\n", a.Interval)
		}
		select {
		case <-shutdown:
			return nil
		case <-ticker.C:
			continue
		}
	}
}
