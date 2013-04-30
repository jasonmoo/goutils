package goutils

import (
	"runtime"
	"strings"
	"testing"
)

func TestMultiValue(t *testing.T) {

	cases := map[string][]string{
		"1,2,3":    []string{"1", "2", "3"},
		"1,2,,3":   []string{"1", "2", "3"},
		"1,2,":     []string{"1", "2"},
		"1,  2, 3": []string{"1", "2", "3"},
		"1,2,3,":   []string{"1", "2", "3"},
	}

	for input, output := range cases {
		v := MultiValue{}
		v.Set(input)

		if len(v) != len(output) {
			t.Errorf("Expected length %d, got: %d\n", len(output), len(v))
		}

		if v.String() != strings.Join(output, ",") {
			t.Errorf("Expected '%s', got: %v\n", input, v)
		}
	}

}

func TestVersionInfo(t *testing.T) {
	info := VersionInfo("test")
	if !strings.Contains(info, "test") || !strings.Contains(info, runtime.Version()) {
		t.Errorf("Version Info did not contain the expected info:\n%s", info)
	}
}

func TestIsPowerOf2(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	cases := map[uint]bool{
		uint(0):    false,
		uint(1):    true,
		uint(3):    false,
		uint(4):    true,
		uint(8):    true,
		uint(1024): true,
	}

	for i, tf := range cases {
		if p := IsPowerOf2(i); p != tf {
			t.Errorf("Incorrect determination:  Expected %t, got %t for %d", tf, p, i)
		}
	}

}

func BenchmarkStandardLoop(b *testing.B) {

	for i := 0; i < b.N; i++ {
	}

}
func BenchmarkStandardLoopFloatCasting(b *testing.B) {

	for i := float64(0); i < float64(b.N); i++ {
	}

}

func BenchmarkVersionInfo(b *testing.B) {

	for i := 0; i < b.N; i++ {
		VersionInfo("bench")
	}

}

func BenchmarkIsPowerOf2(b *testing.B) {

	for i := 0; i < b.N; i++ {
		IsPowerOf2(uint(i))
	}

}
