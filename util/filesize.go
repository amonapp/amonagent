package util

// conversion units
const (
	BYTE     = 1.0
	KILOBYTE = 1024 * BYTE
	MEGABYTE = 1024 * KILOBYTE
	GIGABYTE = 1024 * MEGABYTE
	TERABYTE = 1024 * GIGABYTE
)


// ToMegabytes parses a string formatted by ByteSize as megabytes
func ToMegabytes(s uint64) (uint64, error) {

	bytes := s / MEGABYTE

	return bytes, nil
}