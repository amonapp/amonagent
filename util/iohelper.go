package util 

import (
	"bufio"
	"os"
	// "strings"
	// "unicode"
	// "unicode/utf8"
)

// ReadLine savely reads the line and cleans up afterwards
func ReadLine(fname string, line func(string) error) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if err := line(scanner.Text()); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// IsDigit returns true if s consists of decimal digits.
// func IsDigit(s string) bool {
// 	r := strings.NewReader(s)
// 	for {
// 		ch, _, err := r.ReadRune()
// 		if ch == 0 || err != nil {
// 			break
// 		} else if ch == utf8.RuneError {
// 			return false
// 		} else if !unicode.IsDigit(ch) {
// 			return false
// 		}
// 	}
// 	return true
// }
