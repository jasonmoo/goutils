package goutils

import (
	"expvar"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func TestExpvarInt(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	n := NewExpvarInt("derp", 0)

	if n.Load() != int64(0) {
		t.Errorf("Expected 0 and got %d\n", n.Load())
	}
	n.Add(1)
	if n.Load() != int64(1) {
		t.Errorf("Expected 1 and got %d\n", n.Load())
	}
	n.Set(5)
	if n.Load() != int64(5) {
		t.Errorf("Expected 5 and got %d\n", n.Load())
	}
	n.Add(1)
	if n.Load() != int64(6) {
		t.Errorf("Expected 6 and got %d\n", n.Load())
	}
	if n.String() != "6" {
		t.Errorf("Expected 6 and got %s\n", n.String())
	}

}

func TestExpvarIntFormatted(t *testing.T) {

	// safe to run in parallel
	t.Parallel()

	cases := map[int64]string{
		int64(1):                  "1",
		int64(999):                "999",
		int64(1001):               "1,001",
		int64(122001):             "122,001",
		int64(11001):              "11,001",
		int64(446744073709551615): "446,744,073,709,551,615",
	}

	for input, output := range cases {
		n := NewExpvarIntFormatted("TestExpvarIntFormatted "+output, input)
		if n.String() != output {
			t.Errorf("String Test: expected %#v  got: %#v", output, n.String())
		}
	}

}

func TestNewExpvarRollingFloat(t *testing.T) {

	// test simple increment
	v := NewExpvarRollingFloat("TestNewExpvarRollingFloat", 0, 10)

	// for i := 1; i < 20; i++ {
	// 	v.Add(float64(i))
	// 	t.Log(v.Report().String())
	// }
	// t.Error("done")


	expected := "Reports: 1, Window Size: 10, Last Value: 0.000000, Value Avg: 0.000000, Delta Avg: 0.000000, Rolling Value Avg: 0.000000, Rolling Delta Avg: 0.000000"
	if s := v.String(); s != expected {
		t.Errorf("Adding Test: \nexp: %s\ngot: %s", expected, s)
	}
	v.Add(1)
	expected = "Reports: 2, Window Size: 10, Last Value: 1.000000, Value Avg: 0.500000, Delta Avg: 0.500000, Rolling Value Avg: 0.100000, Rolling Delta Avg: 0.000000"
	if s := v.String(); s != expected {
		t.Errorf("Adding Test: \nexp: %s\ngot: %s", expected, s)
	}
	v.Add(1)
	expected = "Reports: 3, Window Size: 10, Last Value: 1.000000, Value Avg: 0.666667, Delta Avg: 0.333333, Rolling Value Avg: 0.200000, Rolling Delta Avg: 0.000000"
	if s := v.String(); s != expected {
		t.Errorf("Adding Test: \nexp: %s\ngot: %s", expected, s)
	}

	// test range of numbers
	v = NewExpvarRollingFloat("TestNewExpvarRollingFloat2", 2, 10)
	expected = "Reports: 1, Last Value: 0.000000, Value Avg: 2.000000, Delta Avg: 2.000000"
	if s := v.String(); s != expected {
		t.Errorf("Adding Test: \nexp: %s\ngot: %s", expected, s)
	}
	v.Add(4)
	expected = "Reports: 2, Last Value: 4.000000, Value Avg: 3.000000, Delta Avg: 3.000000"
	if s := v.String(); s != expected {
		t.Errorf("Adding Test: \nexp: %s\ngot: %s", expected, s)
	}
	v.Add(6)
	expected = "Reports: 3, Last Value: 6.000000, Value Avg: 4.000000, Delta Avg: 2.666667"
	if s := v.String(); s != expected {
		t.Errorf("Adding Test: \nexp: %s\ngot: %s", expected, s)
	}

}


var (
	expint    *expvar.Int = expvar.NewInt("my var")
	globalint int64
)

func BenchmarkNativeInt64Addition(b *testing.B) {

	for i := 0; i < b.N; i++ {
		globalint++
	}

}
func BenchmarkNativeAtomicInt64Adder(b *testing.B) {
	b.StopTimer()
	globalint = 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		atomic.AddInt64(&globalint, 1)
	}
}
func BenchmarkExpvarInt64Adder(b *testing.B) {
	b.StopTimer()
	xi := &ExpvarInt{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		xi.Add(int64(i))
	}
}
func BenchmarkExpvarInt64FormattedAdder(b *testing.B) {
	b.StopTimer()
	xi := &ExpvarIntFormatted{}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		xi.Add(int64(i))
	}
}
func BenchmarkMutexedInt64Adder(b *testing.B) {
	b.StopTimer()
	var m sync.Mutex
	globalint = 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Lock()
		globalint++
		m.Unlock()
	}
}
func BenchmarkNativeExpvarInt64Adder(b *testing.B) {

	for i := 0; i < b.N; i++ {
		expint.Add(int64(i))
	}

}
func BenchmarkRWMutexedInt64Adder(b *testing.B) {
	b.StopTimer()
	var m sync.RWMutex
	globalint = 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		m.Lock()
		globalint++
		m.Unlock()
	}
}
func BenchmarkDeferredMutexedInt64Adder(b *testing.B) {
	b.StopTimer()
	var m sync.Mutex
	globalint = 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		func() {
			defer m.Unlock()
			m.Lock()
			globalint++
		}()
	}
}
func BenchmarkDeferredRWMutexedUInt64Adder(b *testing.B) {
	b.StopTimer()
	var m sync.RWMutex
	globalint = 0
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		func() {
			defer m.Unlock()
			m.Lock()
			globalint++
		}()
	}
}
func BenchmarkChannelInt64Adder(b *testing.B) {

	globalint = 0
	c := make(chan int64)
	go func() {
		for {
			globalint += <-c
		}
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c <- int64(i)
	}

}

func BenchmarkFormatInt64(b *testing.B) {

	for i := 0; i < b.N; i++ {
		strconv.FormatInt(int64(i), 10)
	}

}
func BenchmarkExpvarInt64Formatted(b *testing.B) {

	n := &ExpvarIntFormatted{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n.Set(int64(i))
		n.String()
	}

}
