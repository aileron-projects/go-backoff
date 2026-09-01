package backoff

import (
	"math"
	"testing"
	"time"

	"github.com/aileron-projects/go-tester"
)

func TestFixedConfig(t *testing.T) {
	t.Parallel()
	t.Run("default", func(t *testing.T) {
		b, err := (&FixedConfig{}).New()
		tester.AssertEqualErr(t, nil, err)
		want := &Fixed{
			interval: time.Second,
		}
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("all", func(t *testing.T) {
		c := &FixedConfig{Interval: 1}
		want := &Fixed{interval: 1}
		b, err := c.New()
		tester.AssertEqualErr(t, nil, err)
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("interval<0", func(t *testing.T) {
		b, err := (&FixedConfig{Interval: -1}).New()
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("interval>0", func(t *testing.T) {
		b, err := (&FixedConfig{Interval: 1}).New()
		tester.AssertEqualErr(t, nil, err)
		want := &Fixed{
			interval: 1,
		}
		tester.AssertDeepEqual(t, want, b)
	})
}

func TestRandomConfig(t *testing.T) {
	t.Parallel()
	t.Run("default", func(t *testing.T) {
		b, err := (&RandomConfig{}).New()
		tester.AssertEqualErr(t, nil, err)
		want := &Random{
			offset:   0,
			deltaMax: float64(time.Second),
		}
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("all", func(t *testing.T) {
		c := &RandomConfig{Offset: 1, Limit: 10}
		want := &Random{offset: 1, deltaMax: 10 - 1}
		b, err := c.New()
		tester.AssertEqualErr(t, nil, err)
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("common param error", func(t *testing.T) {
		b, err := (&RandomConfig{Offset: -1}).New()
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
}

func TestLinearConfig(t *testing.T) {
	t.Parallel()
	t.Run("default", func(t *testing.T) {
		b, err := (&LinearConfig{}).New()
		tester.AssertEqualErr(t, nil, err)
		want := &Linear{
			offset:   0,
			deltaMax: math.MaxInt64,
			coeff:    float64(time.Second),
			jitter:   0,
		}
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("all", func(t *testing.T) {
		c := &LinearConfig{
			Offset: 1,
			Limit:  10,
			Coeff:  100,
			Jitter: 0.5,
		}
		want := &Linear{
			offset:   1,
			deltaMax: 10 - 1,
			coeff:    100,
			jitter:   0.5,
		}
		b, err := c.New()
		tester.AssertEqualErr(t, nil, err)
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("common param error", func(t *testing.T) {
		b, err := (&LinearConfig{Offset: -1}).New()
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})

}

func TestPolynomialConfig(t *testing.T) {
	t.Parallel()
	t.Run("default", func(t *testing.T) {
		b, err := (&PolynomialConfig{}).New()
		tester.AssertEqualErr(t, nil, err)
		want := &Polynomial{
			offset:   0,
			deltaMax: math.MaxInt64,
			coeff:    float64(time.Millisecond),
			exponent: 2,
			jitter:   0,
		}
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("all", func(t *testing.T) {
		c := &PolynomialConfig{
			Offset:   1,
			Limit:    10,
			Coeff:    100,
			Jitter:   0.5,
			Exponent: 100,
		}
		want := &Polynomial{
			offset:   1,
			deltaMax: 10 - 1,
			coeff:    100,
			jitter:   0.5,
			exponent: 100,
		}
		b, err := c.New()
		tester.AssertEqualErr(t, nil, err)
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("common param error", func(t *testing.T) {
		b, err := (&PolynomialConfig{Offset: -1}).New()
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("exponent<0", func(t *testing.T) {
		b, err := (&PolynomialConfig{Exponent: -1}).New()
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
}

func TestExponentialConfig(t *testing.T) {
	t.Parallel()
	t.Run("default", func(t *testing.T) {
		b, err := (&ExponentialConfig{}).New()
		tester.AssertEqualErr(t, nil, err)
		want := &Exponential{
			offset:   0,
			deltaMax: math.MaxInt64,
			coeff:    float64(time.Millisecond),
			base:     2,
			jitter:   0,
		}
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("all", func(t *testing.T) {
		c := &ExponentialConfig{
			Offset: 1,
			Limit:  10,
			Coeff:  100,
			Jitter: 0.5,
			Base:   1000,
		}
		want := &Exponential{
			offset:   1,
			deltaMax: 10 - 1,
			coeff:    100,
			jitter:   0.5,
			base:     1000,
		}
		b, err := c.New()
		tester.AssertEqualErr(t, nil, err)
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("common param error", func(t *testing.T) {
		b, err := (&ExponentialConfig{Offset: -1}).New()
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("base<1", func(t *testing.T) {
		b, err := (&ExponentialConfig{Base: -1}).New()
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
}

func TestFibonacciConfig(t *testing.T) {
	t.Parallel()
	t.Run("default", func(t *testing.T) {
		b, err := (&FibonacciConfig{}).New()
		tester.AssertEqualErr(t, nil, err)
		want := &Fibonacci{
			offset:   0,
			deltaMax: math.MaxInt64,
			coeff:    float64(time.Millisecond),
			jitter:   0,
		}
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("all", func(t *testing.T) {
		c := &FibonacciConfig{
			Offset: 1,
			Limit:  10,
			Coeff:  100,
			Jitter: 0.5,
		}
		want := &Fibonacci{
			offset:   1,
			deltaMax: 10 - 1,
			coeff:    100,
			jitter:   0.5,
		}
		b, err := c.New()
		tester.AssertEqualErr(t, nil, err)
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("common param error", func(t *testing.T) {
		b, err := (&FibonacciConfig{Offset: -1}).New()
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
}

func TestCheckCommonParams(t *testing.T) {
	t.Parallel()
	t.Run("offset<0", func(t *testing.T) {
		err := checkCommonParams(-1, dummy, dummy, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
	})
	t.Run("offset=0", func(t *testing.T) {
		err := checkCommonParams(0, dummy, dummy, dummy)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("offset>0", func(t *testing.T) {
		err := checkCommonParams(0, dummy, dummy, dummy)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("limit<offset", func(t *testing.T) {
		err := checkCommonParams(10, 5, dummy, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
	})
	t.Run("limit=offset", func(t *testing.T) {
		err := checkCommonParams(1, 1, dummy, dummy)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("limit>offset", func(t *testing.T) {
		err := checkCommonParams(1, 2, dummy, dummy)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("coeff<0", func(t *testing.T) {
		err := checkCommonParams(dummy, dummy, -1, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
	})
	t.Run("coeff=0", func(t *testing.T) {
		err := checkCommonParams(dummy, dummy, 0, dummy)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("coeff<0", func(t *testing.T) {
		err := checkCommonParams(dummy, dummy, 1, dummy)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("jitter<0", func(t *testing.T) {
		err := checkCommonParams(dummy, dummy, dummy, -0.1)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
	})
	t.Run("jitter=0", func(t *testing.T) {
		err := checkCommonParams(dummy, dummy, dummy, 0)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("jitter=1", func(t *testing.T) {
		err := checkCommonParams(dummy, dummy, dummy, 1)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("jitter>1", func(t *testing.T) {
		err := checkCommonParams(dummy, dummy, dummy, 1.1)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
	})
}
