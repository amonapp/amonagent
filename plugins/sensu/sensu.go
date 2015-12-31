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
func Run(command string) (string, error) {
	splitCmd, err := shellquote.Split(command)
	if err != nil || len(splitCmd) == 0 {
		return "", fmt.Errorf("exec: unable to parse command, %s", err)
	}

	cmd := exec.Command(splitCmd[0], splitCmd[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("exec: %s for command '%s'", err, command)
	}

	return out.String(), nil
}

func (m Metric) String() string {
	s, _ := json.Marshal(m)
	return string(s)
}

// Metric - XXX
type Metric struct {
	Chart string `json:"chart"`
	Line  string `json:"line"`
	Value string `json:"value"`
}

// ParsedLine - XXX
type ParsedLine struct {
	Plugin   string
	Elements []string
	Value    string
}

// AnalyzeLines - XXX
func AnalyzeLines(completeResult []ParsedLine, line ParsedLine) (Metric, error) {
	m := Metric{Value: line.Value}
	elements := line.Elements

	if len(elements) > 2 {

		m.Chart = strings.Join(elements[:2], "_")
		m.Line = strings.Join(elements[2:], "_")

	} else {
		m.Chart = elements[0]
		m.Line = strings.Join(elements[1:], "_")
	}

	if len(m.Line) == 0 {
		m.Line = m.Chart
	}

	return m, nil
}

// ParseLine - XXX
func ParseLine(s string) (ParsedLine, error) {
	// split by space
	f := func(c rune) bool {
		return c == ' '
	}
	// split metric name by .
	dot := func(c rune) bool {
		return c == '.'
	}
	//split metric name by _
	underscore := func(c rune) bool {
		return c == '_'
	}

	fields := strings.FieldsFunc(s, f)
	line := ParsedLine{}
	if len(fields) == 3 {
		toFloat, _ := strconv.ParseFloat(fields[1], 64)
		value := strconv.FormatFloat(toFloat, 'f', -1, 64)

		metricFields := strings.FieldsFunc(fields[0], dot)

		var cleanName string
		// Eliminiate host and plugin name here
		// Example ubuntu.elasticsearch.thread_pool......
		if len(metricFields) > 2 {
			cleanName = strings.Join(metricFields[2:], ".")
		} else {
			cleanName = strings.Join(metricFields[:], ".")
		}
		CleanMetricFields := strings.FieldsFunc(cleanName, dot)
		splitOnUnderscore := strings.FieldsFunc(cleanName, underscore)

		// Standart use case here
		// Example: thread_pool.search.active
		if len(CleanMetricFields) > 1 {
			line.Elements = CleanMetricFields
			line.Plugin = metricFields[1]
			line.Value = value
		} else {
			line.Elements = splitOnUnderscore
			line.Plugin = metricFields[1]
			line.Value = value

		}

	}

	return line, nil
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
	var Commands []string

	if err := json.Unmarshal(file, &Commands); err != nil {
		return err
	}
	for _, command := range Commands {

		result, err := Run(command)
		if err != nil {
			fmt.Println("Can't execute command: ", err)
		}
		lines := strings.Split(result, "\n")
		ParsedLines := []ParsedLine{}

		for _, line := range lines {
			metric, _ := ParseLine(line)
			if len(metric.Elements) > 0 {
				ParsedLines = append(ParsedLines, metric)
			}

		}

		for _, line := range ParsedLines {
			metric, _ := AnalyzeLines(ParsedLines, line)
			fmt.Println(metric)
		}
		// fmt.Println(ParsedLines)
	}

	return nil
}
