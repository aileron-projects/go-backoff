package backoff_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/aileron-projects/go-backoff"
)

func ExampleFixed() {
	fixed := &backoff.FixedConfig{
		Interval: 100 * time.Millisecond,
	}

	bo, err := fixed.New()
	if err != nil {
		panic(err)
	}

	for i := range 11 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
	}
	// Output:
	// 000: 0s
	// 001: 100ms
	// 002: 100ms
	// 003: 100ms
	// 004: 100ms
	// 005: 100ms
	// 006: 100ms
	// 007: 100ms
	// 008: 100ms
	// 009: 100ms
	// 010: 100ms
}

func ExampleRandom() {
	random := &backoff.RandomConfig{
		Offset: 0 * time.Millisecond,
		Limit:  100 * time.Millisecond,
	}

	bo, err := random.New()
	if err != nil {
		panic(err)
	}

	for i := range 11 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
	}
}

func ExampleLinear() {
	linear := &backoff.LinearConfig{
		Offset: 0 * time.Millisecond,
		Limit:  200 * time.Millisecond,
		Coeff:  5 * time.Millisecond,
		Jitter: backoff.NoJitter,
	}

	bo, err := linear.New()
	if err != nil {
		panic(err)
	}

	for i := range 11 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
	}
	// Output:
	// 000: 0s
	// 001: 5ms
	// 002: 10ms
	// 003: 15ms
	// 004: 20ms
	// 005: 25ms
	// 006: 30ms
	// 007: 35ms
	// 008: 40ms
	// 009: 45ms
	// 010: 50ms
}

func ExamplePolynomial() {
	polynomial := &backoff.PolynomialConfig{
		Offset:   0 * time.Millisecond,
		Limit:    600 * time.Second,
		Coeff:    10 * time.Millisecond,
		Exponent: 2.0,
		Jitter:   backoff.NoJitter,
	}

	bo, err := polynomial.New()
	if err != nil {
		panic(err)
	}

	for i := range 11 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
	}
	// Output:
	// 000: 0s
	// 001: 10ms
	// 002: 40ms
	// 003: 90ms
	// 004: 160ms
	// 005: 250ms
	// 006: 360ms
	// 007: 490ms
	// 008: 640ms
	// 009: 810ms
	// 010: 1s
}

func ExampleExponential() {
	exponential := &backoff.ExponentialConfig{
		Offset: 0 * time.Millisecond,
		Limit:  600 * time.Second,
		Coeff:  1 * time.Millisecond,
		Base:   1.2,
		Jitter: backoff.NoJitter,
	}
	bo, err := exponential.New()
	if err != nil {
		panic(err)
	}

	for i := range 11 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
	}
	// Output:
	// 000: 0s
	// 001: 1.2ms
	// 002: 1.44ms
	// 003: 1.728ms
	// 004: 2.0736ms
	// 005: 2.48832ms
	// 006: 2.985983ms
	// 007: 3.58318ms
	// 008: 4.299816ms
	// 009: 5.15978ms
	// 010: 6.191736ms
}

func ExampleFibonacci() {
	fibonacci := &backoff.FibonacciConfig{
		Offset: 0 * time.Millisecond,
		Limit:  600 * time.Second,
		Coeff:  500 * time.Microsecond,
		Jitter: backoff.NoJitter,
	}

	bo, err := fibonacci.New()
	if err != nil {
		panic(err)
	}

	for i := range 11 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
	}
	// Output:
	// 000: 0s
	// 001: 500µs
	// 002: 500µs
	// 003: 1ms
	// 004: 1.5ms
	// 005: 2.5ms
	// 006: 4ms
	// 007: 6.5ms
	// 008: 10.5ms
	// 009: 17ms
	// 010: 27.5ms
}

func ExampleRetryer_stop() {
	bo, err := (&backoff.FixedConfig{10 * time.Millisecond}).New()
	if err != nil {
		panic(err)
	}

	r := backoff.NewRetryer(bo)
	r.MaxRetry = 100
	r.MaxElapsedTime = 0 // No timeout

	err = r.Run(func(n *backoff.Notify) error {
		count := n.RetryCount()
		fmt.Printf("%d-th call\n", count)
		if count == 3 {
			n.Stop()
			fmt.Println("stop !!")
			return errors.New("Err-stop")
		}
		return fmt.Errorf("Err-%d", count)
	})

	fmt.Println("----- Error -----")
	fmt.Println(err)
	// Output:
	// 0-th call
	// 1-th call
	// 2-th call
	// 3-th call
	// stop !!
	// ----- Error -----
	// Err-0
	// Err-1
	// Err-2
	// Err-stop
}

func ExampleRetryer_success() {
	bo, err := (&backoff.FixedConfig{10 * time.Millisecond}).New()
	if err != nil {
		panic(err)
	}

	r := backoff.NewRetryer(bo)
	r.MaxRetry = 100
	r.MaxElapsedTime = 0 // No timeout

	err = r.Run(func(n *backoff.Notify) error {
		count := n.RetryCount()
		fmt.Printf("%d-th call\n", count)
		if count == 3 {
			n.Success()
			fmt.Println("success !!")
			return nil
		}
		return fmt.Errorf("Err-%d", count)
	})

	fmt.Println("----- Error -----")
	fmt.Println(err)
	// Output:
	// 0-th call
	// 1-th call
	// 2-th call
	// 3-th call
	// success !!
	// ----- Error -----
	// <nil>
}
