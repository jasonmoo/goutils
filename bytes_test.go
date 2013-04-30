package goutils

import (
	"testing"
)

type diffcase struct {
	n    int
	a, b []byte
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

	aa, bb := make([]byte, 1<<20), make([]byte, 1<<20)
	ZeroBytes(bb)
	if ct := DiffBytes(aa, bb); ct > 0 {
		t.Errorf("[]byte arrays are different.  Found %d non-matching bytes.", ct)
	}

}

func TestFillBytes(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	a, b := make([]byte, 1<<20), make([]byte, 1<<20)
	FillBytes(a)
	FillBytes(b)
	if ct := DiffBytes(a, b); ct > 0 {
		t.Errorf("[]byte arrays are different.  Found %d non-matching bytes.", ct)
	}

}

func BenchmarkZeroBytes(b *testing.B) {

	buf := make([]byte, 1<<20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ZeroBytes(buf)
	}

}
func BenchmarkFillBytes(b *testing.B) {

	buf := make([]byte, 1<<20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FillBytes(buf)
	}

}
func BenchmarkDiffBytes(b *testing.B) {

	aa, bb := make([]byte, 1<<20), make([]byte, 1<<20)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DiffBytes(aa, bb)
	}

}
