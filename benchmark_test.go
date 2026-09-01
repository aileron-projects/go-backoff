package backoff_test

import (
	"testing"
	"time"

	"github.com/aileron-projects/go-backoff"
)

const benchAttempt = 10

func BenchmarkFixed(b *testing.B) {
	fixed := &backoff.FixedConfig{
		Interval: time.Second,
	}
	backoff, _ := fixed.New()
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkRandom(b *testing.B) {
	random := &backoff.RandomConfig{
		Offset: time.Second,
		Limit:  10 * time.Second,
	}
	backoff, _ := random.New()
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkLinear(b *testing.B) {
	linear := &backoff.LinearConfig{
		Offset: time.Second,
		Limit:  10 * time.Second,
		Coeff:  10 * time.Millisecond,
		Jitter: backoff.EqualJitter,
	}
	backoff, _ := linear.New()
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkPolynomial(b *testing.B) {
	polynomial := &backoff.PolynomialConfig{
		Offset:   time.Second,
		Limit:    10 * time.Second,
		Coeff:    10 * time.Millisecond,
		Exponent: 3,
		Jitter:   backoff.EqualJitter,
	}
	backoff, _ := polynomial.New()
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkExponential(b *testing.B) {
	exponential := &backoff.ExponentialConfig{
		Offset: time.Second,
		Limit:  10 * time.Second,
		Coeff:  10 * time.Millisecond,
		Base:   1.5,
		Jitter: backoff.EqualJitter,
	}
	backoff, _ := exponential.New()
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkFibonacci(b *testing.B) {
	fibonacci := &backoff.FibonacciConfig{
		Offset: time.Second,
		Limit:  10 * time.Second,
		Coeff:  10 * time.Millisecond,
		Jitter: backoff.EqualJitter,
	}
	backoff, _ := fibonacci.New()
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}
