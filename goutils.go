package goutils

import (
	"runtime"
	"strings"
)

// multivalue for easy parsin comma delimited command line flags
type MultiValue []string

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

// simple version info helper
func VersionInfo(extra string) string {
	return strings.Join([]string{
		"Go Version: " + runtime.Version(),
		extra,
	}, "\n")
}

// moar bitwise pls
func IsPowerOf2(i uint) bool {
	return bool((i == 1 || (i-1)&i == 0) && i != 0)
}

func IntConcat(a ...[]int) []int {
	sum := 0
	for i := 0; i < len(a); i++ {
		sum += len(a[i])
	}
	newbuf := make([]int, sum)
	for i := 0; i < len(a); i++ {
		n := copy(newbuf, a[i])
		newbuf = newbuf[n:]
	}
	return newbuf[:cap(newbuf)]

}

// func ReflectConcat(a ...[]interface{}) []interface{} {


// 	for i := 0; i < len(a); i++ {

// 	}

// }