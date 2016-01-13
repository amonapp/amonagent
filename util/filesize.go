package util

import (
	"math"
	"strconv"
)

// conversion units
const (
	BYTE     = 1.0
	KILOBYTE = float64(1024 * BYTE)
	MEGABYTE = float64(1024 * KILOBYTE)
	GIGABYTE = float64(1024 * MEGABYTE)
	TERABYTE = float64(1024 * GIGABYTE)
)

// ConvertBytesTo and converts bytes to ..
func ConvertBytesTo(s interface{}, convert string, precision int) (float64, error) {
	var bytes float64

	switch v := s.(type) {
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			bytes = f
		}
	case uint64:
		bytes = float64(v)
	case int64:
		bytes = float64(v)
	case float64:
		bytes = v

	}

	bytes, err := ConvertBytesFloatTo(bytes, convert)
	bytes, err = FloatDecimalPoint(bytes, precision)

	return bytes, err
}

// ConvertBytesFloatTo converts bytes to ...
func ConvertBytesFloatTo(s float64, convert string) (float64, error) {
	var bytes float64
	switch convert {
	case "kb":
		bytes = s / KILOBYTE
	case "mb":
		bytes = s / MEGABYTE
	case "gb":
		bytes = s / GIGABYTE
	case "tb":
		bytes = s / TERABYTE
	default:
		bytes = s / BYTE
	}

	return bytes, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// FloatDecimalPoint parses a float64
func FloatDecimalPoint(num float64, precision int) (float64, error) {
	var decimal float64
	output := math.Pow(10, float64(precision))

	decimal = float64(round(num*output)) / output

	return decimal, nil
}

// FloatToString - XXX
func FloatToString(num float64) (string, error) {

	f := strconv.FormatFloat(num, 'f', 2, 64)

	return f, nil
}
