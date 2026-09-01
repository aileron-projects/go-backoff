package backoff_test

import (
	"math"
	"testing"
	"time"

	"github.com/aileron-projects/go-backoff"
	"github.com/aileron-projects/go-tester"
)

var maxDurationFloat = time.Duration(math.Nextafter(1<<63, 0))

func TestFixed(t *testing.T) {
	t.Parallel()
	attempts := []int{0, 1, 2, 10, math.MaxInt}
	testCases := map[string]struct {
		interval time.Duration
		want     []time.Duration
	}{
		"interval=0":   {0, []time.Duration{0, time.Second, time.Second, time.Second, time.Second}},
		"interval=1":   {1, []time.Duration{0, 1, 1, 1, 1}},
		"interval=2":   {2, []time.Duration{0, 2, 2, 2, 2}},
		"interval=max": {math.MaxInt64, []time.Duration{0, math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxInt64}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			bo, err := (&backoff.FixedConfig{Interval: tc.interval}).New()
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				t.Log("attempt=", n)
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
			}
		})
	}
}

func TestRandom(t *testing.T) {
	t.Parallel()
	attempts := []int{0, 1, 2, 10, math.MaxInt}
	testCases := map[string]struct {
		offset time.Duration
		limit  time.Duration
		want   []time.Duration
	}{
		"offset=10,limit=10":   {10, 10, []time.Duration{0, 10, 10, 10, 10}},
		"offset=max,limit=max": {math.MaxInt64, math.MaxInt64, []time.Duration{0, math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxInt64}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &backoff.RandomConfig{
				Offset: tc.offset,
				Limit:  tc.limit,
			}
			bo, err := c.New()
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				t.Log("attempt=", n)
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		c := &backoff.RandomConfig{
			Offset: math.MaxInt64,
			Limit:  math.MaxInt64,
		}
		bo, err := c.New()
		tester.AssertEqual(t, nil, err)
		got := bo.Attempt(math.MaxInt)
		tester.AssertEqual(t, math.MaxInt64, got)
	})
}

func TestLinear(t *testing.T) {
	t.Parallel()
	maxDurationInt := time.Duration(math.Nextafter(math.MaxInt, 0))
	attempts := []int{0, 1, 2, 10, math.MaxInt}
	testCases := map[string]struct {
		offset time.Duration
		limit  time.Duration
		coeff  time.Duration
		want   []time.Duration
	}{
		"offset=0,limit=0,coeff=0": {0, 0, 0, []time.Duration{0, time.Second, 2 * time.Second, 10 * time.Second, maxDurationInt}},
		"offset=0,limit>0,coeff=0": {0, 100, 0, []time.Duration{0, 100, 100, 100, 100}},
		"offset=0,limit=0,coeff>0": {0, 0, 1, []time.Duration{0, 1, 2, 10, maxDurationInt}},
		"offset=0,limit>0,coeff>0": {0, 100, 2, []time.Duration{0, 2, 4, 20, 100}},
		"offset>0,limit>0,coeff=0": {10, 100, 0, []time.Duration{0, 100, 100, 100, 100}},
		"offset>0,limit>0,coeff>0": {10, 100, 2, []time.Duration{0, 12, 14, 30, 100}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &backoff.LinearConfig{
				Offset: tc.offset,
				Limit:  tc.limit,
				Coeff:  tc.coeff,
			}
			bo, err := c.New()
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				t.Log("attempt=", n)
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		c := &backoff.LinearConfig{
			Offset: math.MaxInt64,
			Limit:  math.MaxInt64,
			Coeff:  math.MaxInt64,
		}
		bo, err := c.New()
		tester.AssertEqual(t, nil, err)
		got := bo.Attempt(math.MaxInt)
		tester.AssertEqual(t, maxDurationFloat, got)
	})
}

func TestPolynomial(t *testing.T) {
	t.Parallel()
	attempts := []int{0, 1, 2, 10, math.MaxInt}
	testCases := map[string]struct {
		offset time.Duration
		limit  time.Duration
		coeff  time.Duration
		want   []time.Duration
	}{
		"offset=0,limit=0,coeff=0": {0, 0, 0, []time.Duration{0, 1 * time.Millisecond, 4 * time.Millisecond, 100 * time.Millisecond, maxDurationFloat}},
		"offset=0,limit>0,coeff=0": {0, 100, 0, []time.Duration{0, 100, 100, 100, 100}},
		"offset=0,limit=0,coeff>0": {0, 0, 10, []time.Duration{0, 10, 40, 1000, maxDurationFloat}},
		"offset=0,limit>0,coeff>0": {0, 300, 2, []time.Duration{0, 2, 8, 200, 300}},
		"offset>0,limit>0,coeff=0": {10, 100, 0, []time.Duration{0, 100, 100, 100, 100}},
		"offset>0,limit>0,coeff>0": {10, 300, 2, []time.Duration{0, 12, 18, 210, 300}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &backoff.PolynomialConfig{
				Offset:   tc.offset,
				Limit:    tc.limit,
				Coeff:    tc.coeff,
				Exponent: 2,
			}
			bo, err := c.New()
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				t.Log("attempt=", n)
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		c := &backoff.PolynomialConfig{
			Offset:   math.MaxInt64,
			Limit:    math.MaxInt64,
			Coeff:    math.MaxInt64,
			Exponent: math.MaxFloat64,
		}
		bo, err := c.New()
		tester.AssertEqual(t, nil, err)
		got := bo.Attempt(math.MaxInt)
		tester.AssertEqual(t, maxDurationFloat, got)
	})
}

func TestExponential(t *testing.T) {
	t.Parallel()
	attempts := []int{0, 1, 2, 10, math.MaxInt}
	testCases := map[string]struct {
		offset time.Duration
		limit  time.Duration
		coeff  time.Duration
		want   []time.Duration
	}{
		"offset=0,limit=0,coeff=0": {0, 0, 0, []time.Duration{0, 2 * time.Millisecond, 4 * time.Millisecond, 1024 * time.Millisecond, maxDurationFloat}},
		"offset=0,limit>0,coeff=0": {0, 100, 0, []time.Duration{0, 100, 100, 100, 100}},
		"offset=0,limit=0,coeff>0": {0, 0, 10, []time.Duration{0, 20, 40, 10240, maxDurationFloat}},
		"offset=0,limit>0,coeff>0": {0, 3000, 2, []time.Duration{0, 4, 8, 2048, 3000}},
		"offset>0,limit>0,coeff=0": {10, 100, 0, []time.Duration{0, 100, 100, 100, 100}},
		"offset>0,limit>0,coeff>0": {10, 3000, 2, []time.Duration{0, 14, 18, 2058, 3000}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &backoff.ExponentialConfig{
				Offset: tc.offset,
				Limit:  tc.limit,
				Coeff:  tc.coeff,
				Base:   2,
			}
			bo, err := c.New()
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				t.Log("attempt=", n)
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		c := &backoff.ExponentialConfig{
			Offset: math.MaxInt64,
			Limit:  math.MaxInt64,
			Coeff:  math.MaxInt64,
			Base:   math.MaxFloat64,
		}
		bo, err := c.New()
		tester.AssertEqual(t, nil, err)
		got := bo.Attempt(math.MaxInt)
		tester.AssertEqual(t, maxDurationFloat, got)
	})
}

func TestFibonacci(t *testing.T) {
	t.Parallel()
	attempts := []int{0, 1, 2, 10, math.MaxInt}
	testCases := map[string]struct {
		offset time.Duration
		limit  time.Duration
		coeff  time.Duration
		want   []time.Duration
	}{
		"offset=0,limit=0,coeff=0": {0, 0, 0, []time.Duration{0, 1 * time.Millisecond, 1 * time.Millisecond, 55 * time.Millisecond, maxDurationFloat}},
		"offset=0,limit>0,coeff=0": {0, 100, 0, []time.Duration{0, 100, 100, 100, 100}},
		"offset=0,limit=0,coeff>0": {0, 0, 10, []time.Duration{0, 10, 10, 550, maxDurationFloat}},
		"offset=0,limit>0,coeff>0": {0, 300, 2, []time.Duration{0, 2, 2, 110, 300}},
		"offset>0,limit>0,coeff=0": {10, 100, 0, []time.Duration{0, 100, 100, 100, 100}},
		"offset>0,limit>0,coeff>0": {10, 300, 2, []time.Duration{0, 12, 12, 120, 300}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &backoff.FibonacciConfig{
				Offset: tc.offset,
				Limit:  tc.limit,
				Coeff:  tc.coeff,
			}
			bo, err := c.New()
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				t.Log("attempt=", n)
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		c := &backoff.FibonacciConfig{
			Offset: math.MaxInt64,
			Limit:  math.MaxInt64,
			Coeff:  math.MaxInt64,
		}
		bo, err := c.New()
		tester.AssertEqual(t, nil, err)
		got := bo.Attempt(math.MaxInt)
		tester.AssertEqual(t, maxDurationFloat, got)
	})
}
