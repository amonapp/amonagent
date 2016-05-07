package collectors

import (
	"os"
	"path"
	"testing"

	"github.com/amonapp/amonagent/settings"
	"github.com/stretchr/testify/assert"
)

func TestUptime(t *testing.T) {
	uptime := Uptime()
	var aString interface{} = "string"
	assert.IsType(t, aString, uptime)

}

func TestDistro(t *testing.T) {
	d := Distro()
	assert.NotNil(t, d.Name)
	assert.NotNil(t, d.Version)
}

func TestHost(t *testing.T) {
	d := Host()
	var aString interface{} = "string"
	assert.IsType(t, aString, d)
}

// http://www.devdungeon.com/content/working-files-go#check_if_exists
func TestGetCreateMachineID(t *testing.T) {
	// Delete the file and start from scratch. Run the test suite only on machines that are not running the agent in production.
	var machineidPath = path.Join(settings.ConfigPath, "machine-id") // Default machine id path, generated on first install

	if _, err := os.Stat(machineidPath); err == nil {
		err := os.Remove(machineidPath)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
	}

	// Check if the machine id path has been properly deleted
	_, err := os.Stat(machineidPath)
	assert.Error(t, err)

	first := GetOrCreateMachineID()
	var aString interface{} = "string"
	assert.IsType(t, aString, first)

	for i := 1; i <= 100; i++ {
		runAgain := GetOrCreateMachineID()
		assert.Equal(t, first, runAgain, "Creates the Machine ID only once. Should not overwrite on subsequent runs")

	}

}
