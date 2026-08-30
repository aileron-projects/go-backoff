package backoff_test

import (
	"testing"
	"time"

	"github.com/aileron-projects/go-backoff"
)

const benchAttempt = 10

func BenchmarkFixed(b *testing.B) {
	backoff, _ := backoff.NewFixed(time.Second)
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkRandom(b *testing.B) {
	backoff, _ := backoff.NewRandom(time.Second, 10*time.Second)
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkLinear(b *testing.B) {
	backoff, _ := backoff.NewLinear(time.Second, 10*time.Second, 100*time.Millisecond, 0.5)
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkPolynomial(b *testing.B) {
	backoff, _ := backoff.NewPolynomial(time.Second, 10*time.Second, 10*time.Millisecond, 3, 0.5)
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkExponential(b *testing.B) {
	backoff, _ := backoff.NewExponential(time.Second, 10*time.Second, 10*time.Millisecond, 1.5, 0.5)
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}

func BenchmarkFibonacci(b *testing.B) {
	backoff, _ := backoff.NewFibonacci(time.Second, 10*time.Second, 10*time.Millisecond, 0.5)
	b.ResetTimer()
	for b.Loop() {
		backoff.Attempt(benchAttempt)
	}
}
