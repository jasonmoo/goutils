package goutils

import (
	"runtime"
	"strings"
	"testing"
)

const (
	onemb int = 1 << 20
)

type diffcase struct {
	n    int
	a, b []byte
}

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

func TestDiffBytes(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	cases := []diffcase{
		{0, []byte{0}, []byte{0}},
		{3, []byte{0, 0, 0}, []byte{1, 1, 1}},
		{1, []byte{0, 1, 1}, []byte{1, 1, 1}},
		{1, []byte{0}, []byte{1}},
	}

	for _, c := range cases {
		if d := DiffBytes(c.a, c.b); d != c.n {
			t.Errorf("Diff produced %d, expected %d, (%v,%v)", d, c.n, c.a, c.b)
		}
	}

}

func TestZeroBytes(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	aa, bb := make([]byte, onemb), make([]byte, onemb)
	ZeroBytes(bb)
	if ct := DiffBytes(aa, bb); ct > 0 {
		t.Errorf("[]byte arrays are different.  Found %d non-matching bytes.", ct)
	}

}

func TestFillBytes(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	a, b := make([]byte, onemb), make([]byte, onemb)
	FillBytes(a)
	FillBytes(b)
	if ct := DiffBytes(a, b); ct > 0 {
		t.Errorf("[]byte arrays are different.  Found %d non-matching bytes.", ct)
	}

}

func TestByteSizeToHumanReadable(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	cases := map[uint64]string{
		uint64(0):                 "0B",
		uint64(13):                "13B",
		uint64(1 << 10):           "1KB",
		uint64(onemb):             "1MB",
		uint64(1 << 30):           "1GB",
		uint64(1 << 40):           "1TB",
		uint64(1 << 50):           "1PB",
		uint64(1 << 60):           "1EB",
		uint64(42675243822):       "39.74GB",
		uint64(55555):             "54.25KB",
		uint64(11111111111111111): "9.87PB",
	}

	for b, name := range cases {
		precision := 0
		if b > 1024 && b%1024 > 0 {
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
		"0B":      uint64(0),
		"13B":     uint64(13),
		"1KB":     uint64(1 << 10),
		"1MB":     uint64(onemb),
		"1GB":     uint64(1 << 30),
		"1TB":     uint64(1 << 40),
		"1PB":     uint64(1 << 50),
		"1EB":     uint64(1 << 60),
		"39.74GB": uint64(42670500085),
		"54.25KB": uint64(55552),
		"9.87PB":  uint64(11112632080536698),
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

func BenchmarkVersionInfo(b *testing.B) {

	for i := 0; i < b.N; i++ {
		VersionInfo("bench")
	}

}
func BenchmarkZeroBytes(b *testing.B) {

	b.StopTimer()
	buf, size := make([]byte, onemb), onemb
	b.SetBytes(int64(size * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ZeroBytes(buf)
	}

}
func BenchmarkFillBytes(b *testing.B) {

	b.StopTimer()
	buf := make([]byte, onemb)
	b.SetBytes(int64(onemb * b.N))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		FillBytes(buf)
	}

}
func BenchmarkDiffBytes(b *testing.B) {

	b.StopTimer()
	aa, bb := make([]byte, onemb), make([]byte, onemb)
	b.SetBytes(int64(onemb * b.N))
	for i := 0; i < b.N; i++ {
		b.StartTimer()
		DiffBytes(aa, bb)
		b.StopTimer()
	}

}
func BenchmarkByteSizeToHumanReadable(b *testing.B) {

	b.StopTimer()

	ct := 0
	for size := uint64(1); size < 1<<10; size <<= 1 {
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
	for size := uint64(1); size < 1<<10; size <<= 1 {
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
