package backoff

import (
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestNewFixed(t *testing.T) {
	t.Parallel()
	t.Run("interval<0", func(t *testing.T) {
		b, err := NewFixed(-1)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("interval=0", func(t *testing.T) {
		b, err := NewFixed(0)
		tester.AssertEqualErr(t, nil, err)
		want := &Fixed{
			interval: 0,
		}
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("interval>0", func(t *testing.T) {
		b, err := NewFixed(1)
		tester.AssertEqualErr(t, nil, err)
		want := &Fixed{
			interval: 1,
		}
		tester.AssertDeepEqual(t, want, b)
	})
}

func TestNewRandom(t *testing.T) {
	t.Parallel()
	t.Run("common param error", func(t *testing.T) {
		b, err := NewRandom(-1, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("valid params", func(t *testing.T) {
		b, err := NewRandom(1, 10)
		tester.AssertEqualErr(t, nil, err)
		want := &Random{
			offset:   1,
			deltaMax: 10 - 1,
		}
		tester.AssertDeepEqual(t, want, b)
	})
}

func TestNewLinear(t *testing.T) {
	t.Parallel()
	t.Run("common param error", func(t *testing.T) {
		b, err := NewLinear(-1, dummy, dummy, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("valid params", func(t *testing.T) {
		b, err := NewLinear(1, 10, 100, 0.5)
		tester.AssertEqualErr(t, nil, err)
		want := &Linear{
			offset:   1,
			deltaMax: 10 - 1,
			coeff:    100,
			jitter:   0.5,
		}
		tester.AssertDeepEqual(t, want, b)
	})
}

func TestNewPolynomial(t *testing.T) {
	t.Parallel()
	t.Run("common param error", func(t *testing.T) {
		b, err := NewPolynomial(-1, dummy, dummy, dummy, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("exponent<0", func(t *testing.T) {
		b, err := NewPolynomial(dummy, dummy, dummy, -1, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("exponent=0", func(t *testing.T) {
		b, err := NewPolynomial(dummy, dummy, dummy, 0, dummy)
		tester.AssertEqualErr(t, nil, err)
		want := &Polynomial{
			offset:   dummy,
			deltaMax: dummy - 1,
			coeff:    dummy,
			exponent: 0,
			jitter:   dummy,
		}
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("valid params", func(t *testing.T) {
		b, err := NewPolynomial(1, 10, 100, 1000, 0.5)
		tester.AssertEqualErr(t, nil, err)
		want := &Polynomial{
			offset:   1,
			deltaMax: 10 - 1,
			coeff:    100,
			exponent: 1000,
			jitter:   0.5,
		}
		tester.AssertDeepEqual(t, want, b)
	})
}

func TestNewExponential(t *testing.T) {
	t.Parallel()
	t.Run("common param error", func(t *testing.T) {
		b, err := NewExponential(-1, dummy, dummy, dummy, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("base<1", func(t *testing.T) {
		b, err := NewExponential(dummy, dummy, dummy, 0, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("base=1", func(t *testing.T) {
		b, err := NewExponential(dummy, dummy, dummy, 1, dummy)
		tester.AssertEqualErr(t, nil, err)
		want := &Exponential{
			offset:   dummy,
			deltaMax: dummy - dummy,
			coeff:    dummy,
			base:     1,
			jitter:   dummy,
		}
		tester.AssertDeepEqual(t, want, b)
	})
	t.Run("valid params", func(t *testing.T) {
		b, err := NewExponential(1, 10, 100, 1000, 0.5)
		tester.AssertEqualErr(t, nil, err)
		want := &Exponential{
			offset:   1,
			deltaMax: 10 - 1,
			coeff:    100,
			base:     1000,
			jitter:   0.5,
		}
		tester.AssertDeepEqual(t, want, b)
	})
}

func TestNewFibonacci(t *testing.T) {
	t.Parallel()
	t.Run("common param error", func(t *testing.T) {
		b, err := NewFibonacci(-1, dummy, dummy, dummy)
		tester.AssertEqualErr(t, &Error{Type: "range"}, err)
		tester.AssertEqual(t, nil, b)
	})
	t.Run("valid params", func(t *testing.T) {
		b, err := NewFibonacci(1, 10, 100, 0.5)
		tester.AssertEqualErr(t, nil, err)
		want := &Fibonacci{
			offset:   1,
			deltaMax: 10 - 1,
			coeff:    100,
			jitter:   0.5,
		}
		tester.AssertDeepEqual(t, want, b)
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
