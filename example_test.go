package backoff_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/aileron-projects/go-backoff"
)

func ExampleFixed() {
	// Parameters
	interval := 100 * time.Millisecond

	// Instanciate a backoff provider.
	bo, err := backoff.NewFixed(interval)
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
	// Parameters
	offset := 0 * time.Millisecond
	limit := 100 * time.Millisecond

	// Instanciate a backoff provider.
	bo, err := backoff.NewRandom(offset, limit)
	if err != nil {
		panic(err)
	}

	for i := range 11 {
		d := bo.Attempt(i)
		fmt.Printf("%03d: %s\n", i, d)
	}
}

func ExampleLinear() {
	// Parameters
	offset := 0 * time.Millisecond
	limit := 200 * time.Millisecond
	coeff := 5 * time.Millisecond
	jitter := backoff.NoJitter

	// Instanciate a backoff provider.
	bo, err := backoff.NewLinear(offset, limit, coeff, jitter)
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
	// Parameters
	offset := 0 * time.Millisecond
	limit := 600 * time.Second
	coeff := 10 * time.Millisecond
	exponent := 2.0
	jitter := backoff.NoJitter

	// Instanciate a backoff provider.
	bo, err := backoff.NewPolynomial(offset, limit, coeff, exponent, jitter)
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
	// Parameters
	offset := 0 * time.Millisecond
	limit := 600 * time.Second
	coeff := 1 * time.Millisecond
	base := 1.2
	jitter := backoff.NoJitter

	// Instanciate a backoff provider.
	bo, err := backoff.NewExponential(offset, limit, coeff, base, jitter)
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
	// Parameters
	offset := 0 * time.Millisecond
	limit := 600 * time.Second
	coeff := 500 * time.Microsecond
	jitter := backoff.NoJitter

	// Instanciate a backoff provider.
	bo, err := backoff.NewFibonacci(offset, limit, coeff, jitter)
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
	bo, err := backoff.NewFixed(10 * time.Millisecond)
	if err != nil {
		panic(err)
	}

	r := backoff.NewRetryer(bo)
	r.MaxRetry = 100
	r.MaxElapsedTime = 0 // No timeout

	counter := 0
	err = r.Run(func(n *backoff.Notify) error {
		counter++
		fmt.Printf("%d-th call\n", counter)
		if counter == 3 {
			n.Stop()
			fmt.Println("stop !!")
			return errors.New("Err-stop")
		}
		return fmt.Errorf("Err-%d", counter)
	})

	fmt.Println("----- Error -----")
	fmt.Println(err)
	// Output:
	// 1-th call
	// 2-th call
	// 3-th call
	// stop !!
	// ----- Error -----
	// Err-1
	// Err-2
	// Err-stop
}

func ExampleRetryer_success() {
	bo, err := backoff.NewFixed(10 * time.Millisecond)
	if err != nil {
		panic(err)
	}

	r := backoff.NewRetryer(bo)
	r.MaxRetry = 100
	r.MaxElapsedTime = 0 // No timeout

	counter := 0
	err = r.Run(func(n *backoff.Notify) error {
		counter++
		fmt.Printf("%d-th call\n", counter)
		if counter == 3 {
			n.Success()
			fmt.Println("success !!")
			return nil
		}
		return fmt.Errorf("Err-%d", counter)
	})

	fmt.Println("----- Error -----")
	fmt.Println(err)
	// Output:
	// 1-th call
	// 2-th call
	// 3-th call
	// success !!
	// ----- Error -----
	// <nil>
}
