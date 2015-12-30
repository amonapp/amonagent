package sensu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gonuts/go-shellquote"
)

// Run - XXX
func Run(command *Command) (string, error) {
	splitCmd, err := shellquote.Split(command.Command)
	if err != nil || len(splitCmd) == 0 {
		return "", fmt.Errorf("exec: unable to parse command, %s", err)
	}

	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("exec: %s for command '%s'", err, command.Command)
	}

	return out.String(), nil
}

func (m Metric) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

// Metric - XXX
type Metric struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ParseLine - XXX
func ParseLine(s string) (Metric, error) {
	// split by space
	f := func(c rune) bool {
		return c == ' '
	}
	// split metric name by .
	dot := func(c rune) bool {
		return c == '.'
	}
	fields := strings.FieldsFunc(s, f)
	m := Metric{}
	if len(fields) == 3 {
		toFloat, _ := strconv.ParseFloat(fields[1], 64)
		s = strconv.FormatFloat(toFloat, 'f', -1, 64)

		metricFields := strings.FieldsFunc(fields[0], dot)

		var name string
		if len(metricFields) > 2 {
			name = strings.Join(metricFields[2:], ".")
		} else {
			name = strings.Join(metricFields[:], ".")
		}

		m = Metric{Name: name, Value: s}

	}

	return m, nil
}

// Command - XXX
type Command struct {
	Command string `json:"command"`
	Name    string `json:"name"`
}

// Collect - XXX
func Collect() error {

	file, err := ioutil.ReadFile("/etc/opt/amonagent/plugins-enabled/sensu.conf")
	if err != nil {
		fmt.Printf("Can't read config file: %v\n", err)
	}
	var Commands []Command

	if err := json.Unmarshal(file, &Commands); err != nil {
		return err
	}
	for _, command := range Commands {

		result, err := Run(&command)
		if err != nil {
			fmt.Println("Can't execute command: ", err)
		}
		lines := strings.Split(result, "\n")

		for _, line := range lines {
			metric, _ := ParseLine(line)
			fmt.Println(metric)
		}
	}

	return nil
}
