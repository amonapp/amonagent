package checks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"syscall"
)

func (c CommandResult) String() string {
	s, _ := json.Marshal(c)
	return string(s)
}

// CommandResult - XXX
type CommandResult struct {
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output"`
	Command  string `json:"command"`
}

// ExecWithExitCode - XXX
// Source: http://stackoverflow.com/questions/10385551/get-exit-code-go
func ExecWithExitCode(command string) (CommandResult, error) {
	parts := strings.Fields(command)
	head := parts[0]
	parts = parts[1:len(parts)]
	cmd := exec.Command(head, parts...)
	output := CommandResult{Command: command}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		return output, err
	}

	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				output.ExitCode = status.ExitStatus()
			}
		} else {
			return output, err
		}
	}
	output.Output = out.String()

	return output, nil

}

// Collect - XXX
func Collect() error {

	file, err := ioutil.ReadFile("/etc/opt/amonagent/checks.conf")
	if err != nil {
		fmt.Printf("Can't read config file: %v\n", err)
	}
	var arrayData []string

	if err := json.Unmarshal(file, &arrayData); err != nil {

		return err
	}
	for _, v := range arrayData {
		result, err := ExecWithExitCode(v)
		if err != nil {
			fmt.Println("Can't execute command: ", err)
		}
		fmt.Println(result)

	}

	return nil
}
