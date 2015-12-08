package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/martinrusev/amonagent"
	"github.com/martinrusev/amonagent/core"
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
	config := core.Settings()
	ag, err := amonagent.NewAgent(config)
	if err != nil {
		log.Fatal(err)
	}

	// if *fTest {
	// 	err = ag.Test()
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	return
	// }
	//
	// err = ag.Connect()
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

	ag.Run(shutdown)
}
