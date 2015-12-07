package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/martinrusev/amonagent/logging"
)

// AmonAgentLogger for the main file
var AmonAgentLogger = logging.GetLogger("amonagent")

// the header values are %idle, %wait
var headerRE = regexp.MustCompile(`([%][a-zA-Z0-9]+)[\s+]?`)
var valueRE = regexp.MustCompile(`\d+[\.,]\d+`)

// Just for testing
func main() {
	c1, _ := exec.Command("sar", "1", "1").Output()

	sarOutput := string(c1)
	sarLines := strings.Split(sarOutput, "\n")
	header := []string{}
	values := []float64{}

	// var result map[string]float64

	result := make(map[string]float64)

	// Get the header
	for _, line := range sarLines {
		// Replace the regex with something faster
		matches := headerRE.FindAllStringSubmatch(line, -1)

		for _, m := range matches {
			if len(m) == 2 {
				result := strings.Replace(m[0], "%", "", -1) // replace % in idle (%idle)
				header = append(header, result)
			}
		}
	}

	// Get values
	for _, line := range sarLines {
		if strings.Contains(line, "Average") {
			matches := valueRE.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				if len(m) == 1 {
					valueFloat, _ := strconv.ParseFloat(m[0], 64)
					values = append(values, valueFloat)
				}

			}
		}

	}

	if len(header) == len(values) {
		for i := range header {
			result[header[i]] = values[i]
		}

	}

	jsonString, _ := json.Marshal(result)

	fmt.Println(string(jsonString))
	// if err != nil {
	// 	return fmt.Errorf("error getting CPU info: %s", err)
	// }
}
