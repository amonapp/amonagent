package util

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Command - XXX
type Command struct {
	Command string `json:"command"`
	Name    string `json:"name"`
}

func (c CommandResult) String() string {
	s, _ := json.Marshal(c)
	return string(s)
}

// CommandResult - XXX
type CommandResult struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	ExitCode int    `json:"exit_code"`
	Output   string `json:"output"`
	Error    string `json:"error"`
}

// ExecWithExitCode - XXX
// Source: http://stackoverflow.com/questions/10385551/get-exit-code-go
func ExecWithExitCode(command Command) CommandResult {
	parts := strings.Fields(command.Command)
	head := parts[0]
	parts = parts[1:]
	cmd := exec.Command(head, parts...)
	output := CommandResult{Command: command.Command, Name: command.Name}

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Start(); err != nil {
		output.Error = err.Error()

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
			output.Error = err.Error()
		}
	}

	timer := time.AfterFunc(10*time.Second, func() {
		cmd.Process.Kill()
	})
	timer.Stop()

	output.Output = out.String()

	return output

}
