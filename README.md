<!-- markdownlint-disable MD033 MD041 -->

<div align="center">

[![Release](https://img.shields.io/github/v/release/aileron-projects/go-backoff?sort=semver)](https://github.com/aileron-projects/go-backoff/releases)
[![Reference](https://pkg.go.dev/badge/github.com/aileron-projects/go-backoff.svg)](https://pkg.go.dev/github.com/aileron-projects/go-backoff)
[![DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aileron-projects/go-backoff)
[![Test](https://github.com/aileron-projects/go-backoff/actions/workflows/test.yaml/badge.svg)](https://github.com/aileron-projects/go/actions/workflows/test.yaml)

[![Insights](https://badgen.net/badge/Insights/open%2Fsource%2Finsights/cyan)](https://deps.dev/go/github.com%2Faileron-projects%2Fgo-backoff)
[![Insights](https://badgen.net/badge/Insights/OSS%2FInsight/orange)](https://ossinsight.io/analyze/aileron-projects/go-backoff)

</div>

# go-backoff

**Backoff algorithm implementations for Go.**

## Features

- Support various backoff algorithms
- Offset support
- Jitter support
- Function retrying
- Zero dependency

**Supported backoff algorithms:**

- **`Fixed backoff`**
- **`Random backoff`**
- **`Linear backoff`**
- **`Polynomial backoff`**
- **`Exponential backoff`**
- **`Fibonacci backoff`**

## Usages

### Basic backoff usage

Instanciate a backoff provider with paramerters.
The `Attempt(n int)` returns the n-th backoff duration instead of blocking the call.

```go
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
```

### Run retryable functions

The `Retryer` provides running and retrying functions with specified backoff strategy.

Basic usage:

```go
// Create a new backoff.
bo, _ := backoff.NewFixed(10 * time.Millisecond)

// Create a new retryer with backoff and other params.
r := backoff.NewRetryer(bo)
r.MaxRetry = 100
r.MaxElapsedTime = 0 // No timeout

// Run a retryable function.
err = r.Run(func(n *backoff.Notify) error {
  fmt.Println("Function called")
  return nil
})
```

### Stop retrying with success

Tell retryer to stop retrying using notifier's `Success()`.
The returned err will always be nil.

See also [ExampleRetryer_success](./example_test.go).

```go
bo, _ := backoff.NewFixed(10 * time.Millisecond)

r := backoff.NewRetryer(bo)
r.MaxRetry = 10

counter := 0
err := r.Run(func(n *backoff.Notify) error {
  counter++
  if counter == 3 {
    n.Success() // Tell retryer to stop retrying.
    return nil
  }
  return fmt.Errorf("Err-%d", counter)
})
```

### Stop retrying with non-success

Tell retryer to stop retrying using notifier's `Stop()`.

See also [ExampleRetryer_stop](./example_test.go).

```go
bo, _ := backoff.NewFixed(10 * time.Millisecond)

r := backoff.NewRetryer(bo)
r.MaxRetry = 10

counter := 0
err := r.Run(func(n *backoff.Notify) error {
  counter++
  if counter == 3 {
    n.Stop() // Tell retryer to stop retrying.
    return errors.New("Err-stop")
  }
  return fmt.Errorf("Err-%d", counter)
})
```

### Error handling

```go
err = r.RunContext(ctx, func(n *backoff.Notify) error {
  // Your codes.
  return nil
})

if err != nil {
  switch {
  case errors.Is(err, context.Canceled):
    // Given ctx was canceled.
    // This won't happen when Run, not RunContext, is used.
  case errors.Is(err, context.DeadlineExceeded):
    // Given ctx's deadline was exceeded.
    // This won't happen when Run, not RunContext, is used.
  case errors.Is(err, backoff.ErrMaxRetry):
    // Maximum retry count was reached.
  case errors.Is(err, backoff.ErrTimeout):
    // Maximum elapsed time was reached.
  default:
    // The retryable function told retryer to stop retrying.
    // The err contains all errors returned by the retryable function. 
  }
}
```

## Docs & Examples

- GoDoc: <https://pkg.go.dev/github.com/aileron-projects/go-backoff>
- Examples:
  - [example_test.go](./example_test.go)
  - Fixed backoff: [examples/fixed/](./examples/fixed/)
  - Random backoff: [examples/random/](./examples/random/)
  - Linear backoff: [examples/linear/](./examples/linear/)
  - Polynomial backoff: [examples/polynomial/](./examples/polynomial/)
  - Exponential backoff: [examples/exponential/](./examples/exponential/)
  - Fibonacci backoff: [examples/fibonacci/](./examples/fibonacci/)
  - Retry: [examples/retry/](./examples/retry/)

## Benchmarks

Benchmarks of calling `Attempt(10)`.

See [benchmark_test.go](./benchmark_test.go) for test conditions.

```txt
BenchmarkFixed-8         1000000000    0.948 ns/op   0 B/op   0 allocs/op
BenchmarkRandom-8         100000000   10.760 ns/op   0 B/op   0 allocs/op
BenchmarkLinear-8         100000000   12.660 ns/op   0 B/op   0 allocs/op
BenchmarkPolynomial-8      45811842   26.200 ns/op   0 B/op   0 allocs/op
BenchmarkExponential-8     48073263   28.350 ns/op   0 B/op   0 allocs/op
BenchmarkFibonacci-8      104355823   11.280 ns/op   0 B/op   0 allocs/op
```

## References

- <https://github.com/cenkalti/backoff>
- <https://github.com/avast/retry-go>
