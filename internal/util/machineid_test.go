package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMachineID(t *testing.T) {
	generatedID := GenerateMachineID()

	assert.Len(t, generatedID, 32, "Machine ID should be 32 symbols")
}
