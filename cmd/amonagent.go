package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/amonapp/amonagent"
	"github.com/amonapp/amonagent/collectors"
	"github.com/amonapp/amonagent/plugins"

	_ "github.com/amonapp/amonagent/plugins/all"
	"github.com/amonapp/amonagent/settings"
)

var fTest = flag.Bool("test", false, "gather metrics, print them out, and exit")
var fVersion = flag.Bool("version", false, "display the version")
var fPidfile = flag.String("pidfile", "", "file to write our pid to")
var fMachineID = flag.Bool("machineid", false, "Returns machine id, this value is used in the Salt minion config")

// Amonagent version
//	-ldflags "-X main.Version=`git describe --always --tags`"
var Version string

func main() {
	// var wg sync.WaitGroup

	EnabledPlugins, _ := plugins.GetAllEnabledPlugins()
	for _, p := range EnabledPlugins {
		fmt.Print(p.Name)
		creator, ok := plugins.Plugins[p.Name]
		if !ok {

			fmt.Println("Non existing plugin:", p.Name)

		} else {
			plugin := creator()

			PluginResult, err := plugin.Collect(p.Path)
			if err != nil {

				fmt.Println("\n-----------")
				fmt.Printf("Can't get stats for plugin: %s", err)
				fmt.Println("\n-----------")
			}
			fmt.Print(plugin.SampleConfig())

			fmt.Println(PluginResult)
		}

	}
	// result := make(map[string]interface{})
	// for name := range plugins.Plugins {
	// 	fmt.Println(name)
	// 	creator, ok := plugins.Plugins[name]
	// 	if !ok {
	// 		fmt.Printf("Undefined but requested plugin: %s", name)
	// 	}
	// 	wg.Add(1)
	// 	plugin := creator()
	// 	go func(name string) {
	// 		defer wg.Done()
	// 		PluginResult, _ := plugin.Collect()
	// 		result[name] = PluginResult
	// 	}(name)
	//
	// }
	// wg.Wait()
	//
	// fmt.Println(result)

	return

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

	ag, err := amonagent.NewAgent(config)
	if err != nil {
		log.Fatal(err)
	}

	if *fTest {
		err = ag.Test(config)
		if err != nil {
			log.Fatal(err)
		}
		return
	}

	if *fMachineID {
		fmt.Print(machineID)
		return
	}

	if len(machineID) == 0 && len(serverKey) == 0 {
		log.Fatal("Can't detect Machine ID. Please define `server_key` in /etc/opt/amonagent/amonagent.conf ")
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
		// Ensure the required directory structure exists.
		err := os.MkdirAll(filepath.Dir(*fPidfile), 0700)
		if err != nil {
			log.Fatal(3, "Failed to verify pid directory", err)
		}

		f, err := os.Create(*fPidfile)
		if err != nil {
			log.Fatalf("Unable to create pidfile: %s", err)
		}

		fmt.Fprintf(f, "%d\n", os.Getpid())

		f.Close()
	}

	ag.Run(shutdown)
}
