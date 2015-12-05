package util

import "math"

// conversion units
const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
)

// ConvertBytesTo parses a string formatted by ByteSize as kb
func ConvertBytesTo(s float64, convert string) (float64, error) {
	var bytes float64
	switch convert {
	case "kb":
		bytes = s / KILOBYTE
	case "mb":
		bytes = s / MEGABYTE
	case "gb":
		bytes = s / GIGABYTE
	default:
		bytes = s / BYTE
	}

	return bytes, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// FloatDecimalPoint parses a float64
func FloatDecimalPoint(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
