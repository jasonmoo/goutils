package goutils

import (
	_ "expvar"
	"strconv"
	"sync/atomic"
	"sync"
	"expvar"
	"fmt"
)

type (

	ExpvarInt struct {
		value int64
	}
	ExpvarIntFormatted struct {
		ExpvarInt
	}

	ExpvarRollingFloat struct {
		mut *sync.Mutex
		reports int
		value_sum, delta_sum float64
		value_window []float64
	}
	ExpvarRollingFloatReport struct {
		Reports, WindowSize *ExpvarIntFormatted
		ValueAvg, DeltaAvg,
		LastValue, RollingValueAvg, RollingDeltaAvg float64
	}

)


func NewExpvarInt(name string, n int64) *ExpvarInt {
	v := &ExpvarInt{value: n}
	expvar.Publish(name, v)
	return v
}
func (v *ExpvarInt) Add(n int64) {
	atomic.AddInt64(&(v.value), n)
}
func (v *ExpvarInt) Set(n int64) {
	atomic.StoreInt64(&(v.value), n)
}
func (v *ExpvarInt) Load() int64 {
	return atomic.LoadInt64(&(v.value))
}
func (v *ExpvarInt) String() string {
	return strconv.FormatInt(v.value, 10)
}



func NewExpvarIntFormatted(name string, n int64) *ExpvarIntFormatted {
	v := &ExpvarIntFormatted{ExpvarInt{value: n}}
	expvar.Publish(name, v)
	return v
}
func (v *ExpvarIntFormatted) String() string {
	// make the number a string
	n := strconv.FormatInt(v.Load(), 10)

	if len(n) < 4 {
		return n
	}

	// max uint64 len("18,446,744,073,709,551,615") == 27
	nbuf, l, start := make([]byte, 32), len(n), len(n)%3

	// write out the leading digits
	copy(nbuf[:start], n[:start])

	// write out the rest
	var i, ii int
	for i, ii = start, start; i < l; i, ii = i+3, ii+4 {
		nbuf[ii], nbuf[ii+1], nbuf[ii+2], nbuf[ii+3] = ',', n[i], n[i+1], n[i+2]
	}

	if start == 0 {
		return string(nbuf[1:ii])
	}
	return string(nbuf[:ii])
}




func NewExpvarRollingFloat(name string, n float64, window_size int) *ExpvarRollingFloat {
	if window_size < 0 {
		window_size = 100
	}
	v := &ExpvarRollingFloat{
		mut: new(sync.Mutex),
		value_sum: n,
		value_window: make([]float64, window_size),
	}
	expvar.Publish(name, v)
	return v
}
func (a *ExpvarRollingFloat) Add(n float64) {
	a.mut.Lock()
	a.delta_sum += n-a.value_window[a.reports % len(a.value_window)]
	a.value_sum += n
	a.reports++
	a.value_window[a.reports % len(a.value_window)] = n
	a.mut.Unlock()
}
func (a *ExpvarRollingFloat) Report() *ExpvarRollingFloatReport {
	a.mut.Lock()
	defer a.mut.Unlock()

	if a.reports == 0 {
		return &ExpvarRollingFloatReport{}
	}

	// sum only the deltas between each value
	delta_rolling_avg, window_len := float64(0), len(a.value_window)
	for i, j := a.reports, 0; j < window_len; i, j = i+1, j+1 {
		delta_rolling_avg += (a.value_window[i % window_len] - a.value_window[(i+1) % window_len])
	}
	delta_rolling_avg /= float64(len(a.value_window))

	reports, window := &ExpvarIntFormatted{}, &ExpvarIntFormatted{}
	reports.Set(int64(a.reports))
	window.Set(int64(len(a.value_window)))

	return &ExpvarRollingFloatReport{
		Reports: reports,
		WindowSize: window,
		LastValue: a.value_window[a.reports % len(a.value_window)],
		ValueAvg: a.value_sum/float64(a.reports),
		DeltaAvg: a.delta_sum/float64(a.reports),
		RollingValueAvg: sum_floats(a.value_window)/float64(window_len),
		RollingDeltaAvg: delta_rolling_avg,
	}
}
func (a *ExpvarRollingFloat) String() string {
	return a.Report().String()
}
func (a *ExpvarRollingFloatReport) String() string {
	return fmt.Sprintf("Reports: %s, Window Size: %s, Last Value: %f, Value Avg: %f, Delta Avg: %f, Rolling Value Avg: %f, Rolling Delta Avg: %f",
						a.Reports, a.WindowSize, a.LastValue, a.ValueAvg, a.DeltaAvg, a.RollingValueAvg, a.RollingDeltaAvg)
}
func sum_floats(a []float64) (s float64) {
	for _, n := range a {
		s += n
	}
	return
}
