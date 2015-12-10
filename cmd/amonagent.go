package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/martinrusev/amonagent"
	"github.com/martinrusev/amonagent/collectors"
	"github.com/martinrusev/amonagent/settings"
)

var fTest = flag.Bool("test", false, "gather metrics, print them out, and exit")
var fVersion = flag.Bool("version", false, "display the version")
var fPidfile = flag.String("pidfile", "", "file to write our pid to")

// Amonagent version
//	-ldflags "-X main.Version=`git describe --always --tags`"
var Version string

func main() {
	flag.Parse()

	if *fVersion {
		v := fmt.Sprintf("Amon - Version %s", Version)
		fmt.Println(v)
		return
	}
	config := settings.Settings()

	// Detect Machine ID or ask for a valid Server Key in Settings
	machineID := collectors.MachineID()
	serverKey := config.ServerKey

	if len(machineID) == 0 && len(serverKey) == 0 {
		log.Fatal("Can't detect Machine ID. Please define `server_key` in /etc/amonagent/amonagent.conf ")
	}

	ag, err := amonagent.NewAgent(config)
	if err != nil {
		log.Fatal(err)
	}

	if *fTest {
		err = ag.Test()
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	shutdown := make(chan struct{})
	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		close(shutdown)
	}()

	log.Printf("Starting Amon Agent (version %s)\n", Version)

	if *fPidfile != "" {
		f, err := os.Create(*fPidfile)
		if err != nil {
			log.Fatalf("Unable to create pidfile: %s", err)
		}

		fmt.Fprintf(f, "%d\n", os.Getpid())

		f.Close()
	}

	// ag.Run(shutdown)
}
