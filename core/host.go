package main

import (
	"bufio"
	"fmt"
	"os"
)

//MachineID - XXX
func MachineID() string {
	var machineidPath = "/var/lib/dbus/machine-id"
	var MachineID string
	if _, err := os.Stat(machineidPath); os.IsNotExist(err) {
		machineidPath = "/etc/machine-id"
	}
	file, err := os.Open(machineidPath)
	if err != nil {
		fmt.Printf(err.Error())
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) > 0 {
		MachineID = lines[0]
	}

	// Can't detect, return an empty string and fallback to server key
	if len(MachineID) != 32 {
		MachineID = ""
	}

	return MachineID
}

func main() {
	m := MachineID()
	host, _ := os.Hostname()
	fmt.Println(m)
	fmt.Print(host)
}
