package goutils

import (
	"strconv"
	"strings"
	"runtime"
	"sync/atomic"
)

var (
	SizeNames []string = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB"}
)

type (
	MultiValue []string
	AtomicUInt64CommaStringer uint64
)


func (mv *MultiValue) String() string {
	return strings.Join(*mv, ",")
}
func (mv *MultiValue) Set(value string) error {
	for _, s := range strings.Split(value, ",") {
		name := strings.TrimSpace(s)
		if len(name) > 0 {
			*mv = append(*mv, name)
		}
	}
	return nil
}

func (i *AtomicUInt64CommaStringer) Add(delta int) uint64 {
	return atomic.AddUint64((*uint64)(i), uint64(delta))
}
func (acs *AtomicUInt64CommaStringer) String() string {

	// make the number a string
	n := strconv.FormatUint(atomic.LoadUint64((*uint64)(acs)), 10)

	if len(n) < 4 {
		return n
	}

	// max uint64 len("18,446,744,073,709,551,615") == 27
	nbuf, l, start := make([]byte, 27), len(n), len(n)%3

	// write out the leading digits
	copy(nbuf[:start], n[:start])

	// write out the rest
	var i, ii int
	for i, ii = start, start; i < l; i, ii = i+3, ii+4 {
		nbuf[ii], nbuf[ii+1], nbuf[ii+2], nbuf[ii+3] = ',', n[i], n[i+1], n[i+2]
	}

	return string(nbuf[:ii])
}


func VersionInfo(extra string) string {
	return strings.Join([]string{
		"Go Version: " + runtime.Version(),
		extra,
	}, "\n")
}

// writes zeroes to every byte in []byte
func ZeroBytes(b []byte) {
	for i, l := 0, len(b); i < l; i++ {
		b[i] = byte(0)
	}
}

// writes incrementing short as byte to []byte 0-255,0-255
func FillBytes(b []byte) {
	for i, l := 0, len(b); i < l; i++ {
		b[i] = byte(i)
	}
}
func DiffBytes(a, b []byte) int {
	ct, l := 0, 0
	if len(a) < len(b) {
		l = len(a)
	} else {
		l = len(b)
	}
	for i := 0; i < l; i++ {
		if a[i] != b[i] {
			ct++
		}
	}
	return ct
}

// takes a number of bytes and outputs it in human readable format
// ie: 1024 -> 1KB
func ByteSizeToHumanReadable(b uint64, precision int) string {
	const divisor float64 = float64(1024)
	i, n := 0, float64(b)
	for ; n >= divisor; i, n = i+1, n/divisor {
	}
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
