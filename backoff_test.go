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
		"interval=0":   {0, []time.Duration{0, 0, 0, 0, 0}},
		"interval=1":   {1, []time.Duration{0, 1, 1, 1, 1}},
		"interval=2":   {2, []time.Duration{0, 2, 2, 2, 2}},
		"interval=max": {math.MaxInt, []time.Duration{0, math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			bo, err := backoff.NewFixed(tc.interval)
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
				t.Log("attempt=", n)
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
		"offset=0,limit=0":     {0, 0, []time.Duration{0, 0, 0, 0, 0}},
		"offset=10,limit=10":   {10, 10, []time.Duration{0, 10, 10, 10, 10}},
		"offset=max,limit=max": {math.MaxInt64, math.MaxInt64, []time.Duration{0, math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			bo, err := backoff.NewRandom(tc.offset, tc.limit)
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
				t.Log("attempt=", n)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		bo, err := backoff.NewRandom(math.MaxInt64, math.MaxInt64)
		tester.AssertEqual(t, nil, err)
		got := bo.Attempt(math.MaxInt)
		tester.AssertEqual(t, math.MaxInt64, got)
	})
}

func TestLinear(t *testing.T) {
	t.Parallel()
	attempts := []int{0, 1, 2, 10, math.MaxInt}
	testCases := map[string]struct {
		offset time.Duration
		limit  time.Duration
		coeff  time.Duration
		want   []time.Duration
	}{
		"offset=0,limit=0,coeff=0": {0, 0, 0, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit>0,coeff=0": {0, 100, 0, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit=0,coeff>0": {0, 0, 10, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit>0,coeff>0": {0, 100, 2, []time.Duration{0, 2, 4, 20, 100}},
		"offset>0,limit>0,coeff=0": {10, 100, 0, []time.Duration{0, 10, 10, 10, 10}},
		"offset>0,limit>0,coeff>0": {10, 100, 2, []time.Duration{0, 12, 14, 30, 100}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			bo, err := backoff.NewLinear(tc.offset, tc.limit, tc.coeff, backoff.NoJitter)
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
				t.Log("attempt=", n)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		bo, err := backoff.NewLinear(math.MaxInt64, math.MaxInt64, math.MaxInt64, 0)
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
		"offset=0,limit=0,coeff=0": {0, 0, 0, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit>0,coeff=0": {0, 100, 0, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit=0,coeff>0": {0, 0, 10, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit>0,coeff>0": {0, 300, 2, []time.Duration{0, 2, 8, 200, 300}},
		"offset>0,limit>0,coeff=0": {10, 100, 0, []time.Duration{0, 10, 10, 10, 10}},
		"offset>0,limit>0,coeff>0": {10, 300, 2, []time.Duration{0, 12, 18, 210, 300}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			bo, err := backoff.NewPolynomial(tc.offset, tc.limit, tc.coeff, 2, backoff.NoJitter)
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
				t.Log("attempt=", n)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		bo, err := backoff.NewPolynomial(math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxFloat64, 0)
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
		"offset=0,limit=0,coeff=0": {0, 0, 0, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit>0,coeff=0": {0, 100, 0, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit=0,coeff>0": {0, 0, 10, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit>0,coeff>0": {0, 3000, 2, []time.Duration{0, 4, 8, 2048, 3000}},
		"offset>0,limit>0,coeff=0": {10, 100, 0, []time.Duration{0, 10, 10, 10, 10}},
		"offset>0,limit>0,coeff>0": {10, 3000, 2, []time.Duration{0, 14, 18, 2058, 3000}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			bo, err := backoff.NewExponential(tc.offset, tc.limit, tc.coeff, 2, backoff.NoJitter)
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
				t.Log("attempt=", n)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		bo, err := backoff.NewExponential(math.MaxInt64, math.MaxInt64, math.MaxInt64, math.MaxFloat64, 0)
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
		"offset=0,limit=0,coeff=0": {0, 0, 0, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit>0,coeff=0": {0, 100, 0, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit=0,coeff>0": {0, 0, 10, []time.Duration{0, 0, 0, 0, 0}},
		"offset=0,limit>0,coeff>0": {0, 300, 2, []time.Duration{0, 2, 2, 110, 300}},
		"offset>0,limit>0,coeff=0": {10, 100, 0, []time.Duration{0, 10, 10, 10, 10}},
		"offset>0,limit>0,coeff>0": {10, 300, 2, []time.Duration{0, 12, 12, 120, 300}},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			bo, err := backoff.NewFibonacci(tc.offset, tc.limit, tc.coeff, backoff.NoJitter)
			tester.AssertEqual(t, nil, err)
			for i, n := range attempts {
				got := bo.Attempt(n)
				tester.AssertEqual(t, tc.want[i], got)
				t.Log("attempt=", n)
			}
		})
	}

	t.Run("overflow", func(t *testing.T) {
		bo, err := backoff.NewFibonacci(math.MaxInt64, math.MaxInt64, math.MaxInt64, 0)
		tester.AssertEqual(t, nil, err)
		got := bo.Attempt(math.MaxInt)
		tester.AssertEqual(t, maxDurationFloat, got)
	})
}
