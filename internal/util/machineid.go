package util

import (
	"crypto/rand"
	"fmt"
)

// GenerateMachineID - generates unique machine id
func GenerateMachineID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
