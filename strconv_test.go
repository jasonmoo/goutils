package goutils

import (
	"strconv"
	"testing"
)


func TestFormatHumanReadableSize(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	cases := map[int64]string{
		int64(-3):                "-3B",
		int64(-1024):             "-1KB",
		int64(-1079110533):       "-1GB",
		int64(0):                 "0B",
		int64(13):                "13B",
		int64(1 << 10):           "1KB",
		int64(1 << 20):           "1MB",
		int64(1 << 30):           "1GB",
		int64(1 << 40):           "1TB",
		int64(1 << 50):           "1PB",
		int64(1 << 60):           "1EB",
		int64(42675243822):       "39.74GB",
		int64(55555):             "54.25KB",
		int64(11111111111111111): "9.87PB",
	}

	for b, name := range cases {
		precision := 0
		if b > 1024 && b%1024 > 0 {
			precision = 2
		}
		if h := FormatHumanReadableSize(b, precision); h != name {
			t.Errorf("Conversion incorrect.  Expecting %s, got %s, for %d with %d precision", name, h, b, precision)
		}
	}

}

func TestParseHumanReadableSize(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	cases := map[string]int64{
		"-10B":    int64(-10),
		"-110m":   int64(-115343360),
		"-1.005g": int64(-1079110533),
		"0B":      int64(0),
		"13B":     int64(13),
		"1KB":     int64(1 << 10),
		"1MB":     int64(1 << 20),
		"1GB":     int64(1 << 30),
		"1TB":     int64(1 << 40),
		"1PB":     int64(1 << 50),
		"1EB":     int64(1 << 60),
		"39.74GB": int64(42670500085),
		"54.25KB": int64(55552),
		"9.87PB":  int64(11112632080536698),
	}

	for name, b := range cases {
		if bs, err := ParseHumanReadableSize(name); bs != b {
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Conversion incorrect.  Expecting %d, got %d, for %s", b, bs, name)
		}
	}

}

func BenchmarkNativeFormatFloat64(b *testing.B) {

	for i := 0; i < b.N; i++ {
		strconv.FormatFloat(float64(i), 'f', i%10, 64)
	}

}

func BenchmarkFormatHumanReadableSize(b *testing.B) {

	for i := 0; i < b.N; i++ {
		// rotate through 0-10 precision while benching
		FormatHumanReadableSize(int64(i), i%10)
	}

}

func BenchmarkParseHumanReadableSize(b *testing.B) {
	b.StopTimer()
	names := []string{
		"128",
		"128b",
		"1k",
		"1mb",
		"10mb",
		"100mb",
		"1gb",
		"1.01gb",
		"2t",
		"2e",
	}
	ct := len(names)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ParseHumanReadableSize(names[i%ct])
	}

}

func BenchmarkNativeParseFloat(b *testing.B) {
	b.StopTimer()
	names := []string{
		"128",
		"128",
		"1",
		"1.0",
		"10",
		"100",
		"1111111111111111111",
		"1.00000009",
		"2",
		"0.22",
	}
	ct := len(names)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		strconv.ParseFloat(names[i%ct], 64)
	}

}
