package goutils

import (
	"bytes"
	"testing"
)

func TestZero(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	// test 1B - 1M, in powers of 2
	for i, b := 1, 1 << 20; i < b; i<<=1 {
		blank, tozero := make([]byte, i), make([]byte, i)
		Zero(tozero)
		if !bytes.Equal(blank, tozero) {
			var ct int
			for ct, ii := 0, 0; i < ii; ii++ {
				if blank[ii] != tozero[ii] {
					ct++
				}
			}
			t.Errorf("[]byte arrays are different.  Found %d non-matching bytes.",ct)
		}
	}

}

func TestFill(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	// test 1B - 1M, in powers of 2
	for i, b, diff_ct := 1, 1 << 20, 0; i < b; i<<=1 {
		tofill := make([]byte, i)
		Fill(tofill)
		for i, b := range tofill {
			if byte(i) != b {
				diff_ct++
			}
		}
		if diff_ct > 0 {
			t.Errorf("[]byte arrays are different.  Found %d non-matching bytes.", diff_ct)
		}
		diff_ct = 0
	}

}

func TestByteSizeToHumanReadable(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	cases := map[uint64]string{
		uint64(0): "0B",
		uint64(13): "13B",
		uint64(1 << 10): "1KB",
		uint64(1 << 20): "1MB",
		uint64(1 << 30): "1GB",
		uint64(1 << 40): "1TB",
		uint64(1 << 50): "1PB",
		uint64(1 << 60): "1EB",
		uint64(42675243822): "39.74GB",
		uint64(55555): "54.25KB",
		uint64(11111111111111111): "9.87PB",
	}

	for b, name := range cases {
		precision := 0
		if b > 1024 && b % 1024 > 0 {
			precision = 2
		}
		if h := ByteSizeToHumanReadable(b, precision); h != name {
			t.Errorf("Conversion incorrect.  Expecting %s, got %s, for %d with %d precision", name, h, b, precision)
		}
	}

}

func TestHumanReadableSizeToBytes(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	cases := map[string]uint64{
		"0B": uint64(0),
		"13B": uint64(13),
		"1KB": uint64(1 << 10),
		"1MB": uint64(1 << 20),
		"1GB": uint64(1 << 30),
		"1TB": uint64(1 << 40),
		"1PB": uint64(1 << 50),
		"1EB": uint64(1 << 60),
		"39.74GB": uint64(42670500085),
		"54.25KB": uint64(55552),
		"9.87PB": uint64(11112632080536698),
	}

	for name, b := range cases {
		if bs, err := HumanReadableSizeToBytes(name); bs != b {
			if err != nil {
				t.Error(err)
			}
			t.Errorf("Conversion incorrect.  Expecting %d, got %d, for %s", b, bs, name)
		}
	}

}

func TestIsPowerOf2(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	cases := map[uint]bool{
		uint(0): false,
		uint(1): true,
		uint(3): false,
		uint(4): true,
		uint(8): true,
		uint(1024): true,
	}


	for i, tf := range cases {
		if p := IsPowerOf2(i); p != tf {
			t.Errorf("Incorrect determination:  Expected %t, got %t for %d", tf, p, i)
		}
	}

}




func BenchmarkZero(b *testing.B) {

	b.StopTimer()
	buf, size := make([]byte, 1 << 20), 1 << 20
	b.SetBytes(int64(size * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Zero(buf)
	}

}
func BenchmarkFill(b *testing.B) {

	b.StopTimer()
	buf, size := make([]byte, 1 << 20), 1 << 20
	b.SetBytes(int64(size * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Fill(buf)
	}

}
func BenchmarkByteSizeToHumanReadable(b *testing.B) {

	b.StopTimer()

	ct := 0
	for size := uint64(1); size < 1 << 10; size<<=1 {
		for i := 0; i < b.N/2; i++ {
			b.StartTimer()
			name := ByteSizeToHumanReadable(size, 0)
			b.StopTimer()
			ct += len(name)
		}
		for i := 0; i < b.N/2; i++ {
			b.StartTimer()
			name := ByteSizeToHumanReadable(size, 2)
			b.StopTimer()
			ct += len(name)
		}
	}
	b.SetBytes(int64(ct))

}

func BenchmarkHumanReadableSizeToBytes(b *testing.B) {

	b.StopTimer()

	ct := 0
	for size := uint64(1); size < 1 << 10; size<<=1 {
		name := ByteSizeToHumanReadable(size, 2)
		for i := 0; i < b.N; i++ {
			b.StartTimer()
			t, _ := HumanReadableSizeToBytes(name)
			b.StopTimer()
			ct += int(t)
		}
	}
	b.SetBytes(int64(ct))

}

func BenchmarkIsPowerOf2(b *testing.B) {

	for i := 0; i < b.N; i++ {
		IsPowerOf2(uint(i))
	}

}
