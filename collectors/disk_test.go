package collectors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiskStats(t *testing.T) {
	disk, _ := DiskUsage()
	var aString interface{} = "string"
	var aFloat interface{} = float64(1)
	assert.NotEmpty(t, disk)
	for _, volume := range disk {
		assert.IsType(t, aString, volume.Name)
		assert.IsType(t, aString, volume.Path)
		assert.IsType(t, aString, volume.Free)
		assert.IsType(t, aString, volume.Used)
		assert.IsType(t, aFloat, volume.UsedPercent)
	}

}
