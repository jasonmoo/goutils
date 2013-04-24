package goutils

import (
	"strconv"
	"strings"
)

var (
	SizeNames []string = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
)

// writes zeroes to every byte in []byte
func Zero(b []byte) {
	for i, l := 0, len(b); i < l; i++ {
		b[i] = byte(0)
	}
}
// writes incrementing short as byte to []byte 0-255,0-255
func Fill(b []byte) {
	for i, l := 0, len(b); i < l; i++ {
		b[i] = byte(i)
	}
}
// takes a number of bytes and outputs it in human readable format
// ie: 1024 -> 1KB
func ByteSizeToHumanReadable(b uint64, precision int) string {
	const divisor float64 = float64(1024)
	i, n := 0, float64(b)
	for ; n >= divisor; i, n = i+1, n/divisor {}
	return strconv.FormatFloat(n, 'f', precision, 64) + SizeNames[i]
}

func HumanReadableSizeToBytes(s string) (uint64, error) {

	if i := strings.IndexAny(s, "bBkKmMgGtTpPeE"); i > -1 {
		var m uint64 = 1
		switch s[i] {
			case 'k','K': m = uint64(1 << 10); break
			case 'm','M': m = uint64(1 << 20); break
			case 'g','G': m = uint64(1 << 30); break
			case 't','T': m = uint64(1 << 40); break
			case 'p','P': m = uint64(1 << 50); break
			case 'e','E': m = uint64(1 << 60)
		}
		num, err := strconv.ParseFloat(s[:i], 64)
		if err == nil {
			num = num * float64(m)
		}
		return uint64(num), err
	}
	num, err := strconv.ParseFloat(s, 64)
	return uint64(num), err
}

func IsPowerOf2(i uint) bool {
	return bool((i == 1 || (i-1)&i == 0) && i != 0)
}

