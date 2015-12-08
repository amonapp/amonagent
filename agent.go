package amonagent

import (
	"log"
	"sync"
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
	var wg sync.WaitGroup

	log.Printf("Agent Config: Interval:%s\n", a.Interval)

	ticker := time.NewTicker(a.Interval)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.GatherAndSend(); err != nil {
			log.Printf("Flusher routine failed, exiting: %s\n", err.Error())
			close(shutdown)
		}
	}()

	defer wg.Wait()

	for {

		select {
		case <-shutdown:
			return nil
		case <-ticker.C:
			continue
		}
	}
}
