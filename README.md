<!-- markdownlint-disable MD033 MD041 -->

<div align="center">

[![Release](https://img.shields.io/github/v/release/aileron-projects/go-backoff?sort=semver)](https://github.com/aileron-projects/go-backoff/releases)
[![Reference](https://pkg.go.dev/badge/github.com/aileron-projects/go-backoff.svg)](https://pkg.go.dev/github.com/aileron-projects/go-backoff)
[![DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/aileron-projects/go-backoff)
[![Test](https://github.com/aileron-projects/go-backoff/actions/workflows/test.yaml/badge.svg)](https://github.com/aileron-projects/go-backoff/actions/workflows/test.yaml)

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

See the [API doc](https://pkg.go.dev/github.com/aileron-projects/go-backoff) for details.

- **`Fixed backoff`**: *d = interval*
- **`Random backoff`**: *d = offset + random(limit-offset)*
- **`Linear backoff`**: *d = offset + coeff * n*
- **`Polynomial backoff`**: *d = offset + coeff * n^exponent*
- **`Exponential backoff`**: *d = offset + coeff * base^n*
- **`Fibonacci backoff`**: *d = offset + coeff * fibonacci(n)*

## Usages

### Basic backoff usage

Instanciate a backoff provider with paramerters.
The `Attempt(n int)` returns the n-th backoff duration instead of blocking the call.

```go
// Create a config.
// All params are optional. Set to override default values.
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
```

### Run retryable functions

The `Retryer` provides running and retrying functions with specified backoff strategy.
Retry count can be obtained by `RetryCount()`.

Basic usage:

```go
// Create a new backoff.
bo, _ := (&backoff.FixedConfig{}).New()

// Create a new retryer with backoff and other params.
r := backoff.NewRetryer(bo)
r.MaxRetry = 100
r.MaxElapsedTime = 0 // No timeout

// Run a retryable function.
err = r.Run(func(n *backoff.Notify) error {
  fmt.Println("Retry count:", n.RetryCount())
  return nil
})
```

### Stop retrying with success

Tell retryer to stop retrying using notifier's `Success()`.
The returned err will always be nil.

See also [ExampleRetryer_success](./example_test.go).

```go
bo, _ := (&backoff.FixedConfig{}).New()

r := backoff.NewRetryer(bo)
r.MaxRetry = 10

err := r.Run(func(n *backoff.Notify) error {
  count := n.RetryCount()
  if count == 3 {
    n.Success() // Tell retryer to stop retrying.
    return nil
  }
  return fmt.Errorf("Err-%d", count)
})
```

### Stop retrying with non-success

Tell retryer to stop retrying using notifier's `Stop()`.

See also [ExampleRetryer_stop](./example_test.go).

```go
bo, _ := (&backoff.FixedConfig{}).New()

r := backoff.NewRetryer(bo)
r.MaxRetry = 10

err := r.Run(func(n *backoff.Notify) error {
  count := n.RetryCount()
  if count == 3 {
    n.Stop() // Tell retryer to stop retrying.
    return errors.New("Err-stop")
  }
  return fmt.Errorf("Err-%d", count)
})
```

### Error handling

Errors from the `Retryer` can be identified like below.

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
BenchmarkFixed-8       1000000000    0.6935 ns/op   0 B/op   0 allocs/op
BenchmarkRandom-8       125753828    9.833  ns/op   0 B/op   0 allocs/op
BenchmarkLinear-8        84624442   12.98   ns/op   0 B/op   0 allocs/op
BenchmarkPolynomial-8    32345972   37.23   ns/op   0 B/op   0 allocs/op
BenchmarkExponential-8   48786039   27.31   ns/op   0 B/op   0 allocs/op
BenchmarkFibonacci-8    121107502   10.13   ns/op   0 B/op   0 allocs/op
```

## References

- <https://github.com/cenkalti/backoff>
- <https://github.com/avast/retry-go>
