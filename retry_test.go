package backoff

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"testing"
	"time"

	"github.com/aileron-projects/go-tester"
)

func TestNotify(t *testing.T) {
	t.Parallel()
	t.Run("should not stop", func(t *testing.T) {
		n := &Notify{}
		tester.AssertEqual(t, false, n.shouldStop())
	})
	t.Run("success", func(t *testing.T) {
		n := &Notify{}
		n.Success()
		tester.AssertEqual(t, true, n.shouldStop())
	})
	t.Run("stop", func(t *testing.T) {
		n := &Notify{}
		n.Stop()
		tester.AssertEqual(t, true, n.shouldStop())
	})
}

func TestNewRetryer(t *testing.T) {
	t.Parallel()
	t.Run("nil backoff", func(t *testing.T) {
		r := NewRetryer(nil)
		want := &Random{
			offset:   0,
			deltaMax: float64(time.Second),
		}
		tester.AssertDeepEqual(t, want, r.backoff.(*Random))
	})
	t.Run("nil backoff", func(t *testing.T) {
		want := &Linear{}
		r := NewRetryer(want)
		tester.AssertEqual(t, want, r.backoff.(*Linear))
	})
}

func TestRetryer(t *testing.T) {
	t.Parallel()
	t.Run("success", func(t *testing.T) {
		r := NewRetryer(&Fixed{interval: 10 * time.Millisecond})
		r.MaxRetry = math.MaxInt
		counter := 0
		err := r.Run(func(n *Notify) error {
			if counter++; counter == 3 {
				n.Success()
				return nil
			}
			return fmt.Errorf("Err%d", counter)
		})
		tester.AssertEqual(t, 3, counter)
		tester.AssertEqual(t, nil, err)
	})
	t.Run("stop with nil error", func(t *testing.T) {
		r := NewRetryer(&Fixed{interval: 10 * time.Millisecond})
		r.MaxRetry = math.MaxInt
		counter := 0
		err := r.Run(func(n *Notify) error {
			if counter++; counter == 3 {
				n.Stop()
				return nil
			}
			return io.EOF
		})
		tester.AssertEqual(t, 3, counter)
		tester.AssertEqualErr(t, io.EOF, err)
		tester.AssertEqual(t, false, errors.Is(err, ErrMaxRetry))
		tester.AssertEqual(t, false, errors.Is(err, ErrTimeout))
	})
	t.Run("stop with non-nil error", func(t *testing.T) {
		r := NewRetryer(&Fixed{interval: 10 * time.Millisecond})
		r.MaxRetry = math.MaxInt
		counter := 0
		err := r.Run(func(n *Notify) error {
			if counter++; counter == 3 {
				n.Stop()
				return io.ErrUnexpectedEOF
			}
			return io.EOF
		})
		tester.AssertEqual(t, 3, counter)
		tester.AssertEqualErr(t, io.ErrUnexpectedEOF, err)
		tester.AssertEqualErr(t, io.EOF, err)
		tester.AssertEqual(t, false, errors.Is(err, ErrMaxRetry))
		tester.AssertEqual(t, false, errors.Is(err, ErrTimeout))
	})
	t.Run("max retry", func(t *testing.T) {
		r := NewRetryer(&Fixed{interval: 10 * time.Millisecond})
		r.MaxRetry = 2
		counter := 0
		err := r.Run(func(n *Notify) error {
			counter++
			return io.EOF
		})
		tester.AssertEqual(t, 3, counter)
		tester.AssertEqualErr(t, io.EOF, err)
		tester.AssertEqualErr(t, ErrMaxRetry, err)
		tester.AssertEqual(t, false, errors.Is(err, ErrTimeout))
	})
	t.Run("timeout", func(t *testing.T) {
		r := NewRetryer(&Fixed{interval: 10 * time.Millisecond})
		r.MaxRetry = math.MaxInt
		r.MaxElapsedTime = 50 * time.Millisecond
		counter := 0
		err := r.Run(func(n *Notify) error {
			counter++
			return io.EOF
		})
		tester.AssertEqual(t, true, counter >= 2)
		tester.AssertEqualErr(t, io.EOF, err)
		tester.AssertEqualErr(t, ErrTimeout, err)
		tester.AssertEqual(t, false, errors.Is(err, ErrMaxRetry))
	})
	t.Run("context canceled", func(t *testing.T) {
		r := NewRetryer(&Fixed{interval: 10 * time.Millisecond})
		r.MaxRetry = math.MaxInt
		ctx, cancel := context.WithCancel(context.Background())
		counter := 0
		err := r.RunContext(ctx, func(n *Notify) error {
			if counter++; counter == 3 {
				cancel()
				return io.ErrUnexpectedEOF
			}
			return io.EOF
		})
		tester.AssertEqual(t, 3, counter)
		tester.AssertEqualErr(t, io.EOF, err)
		tester.AssertEqualErr(t, io.ErrUnexpectedEOF, err)
		tester.AssertEqualErr(t, context.Canceled, err)
		tester.AssertEqual(t, false, errors.Is(err, ErrMaxRetry))
		tester.AssertEqual(t, false, errors.Is(err, ErrTimeout))
	})
	t.Run("context deadline", func(t *testing.T) {
		r := NewRetryer(&Fixed{interval: 10 * time.Millisecond})
		r.MaxRetry = math.MaxInt
		ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		counter := 0
		err := r.RunContext(ctx, func(n *Notify) error {
			counter++
			return io.EOF
		})
		tester.AssertEqual(t, true, counter >= 2)
		tester.AssertEqualErr(t, io.EOF, err)
		tester.AssertEqualErr(t, context.DeadlineExceeded, err)
		tester.AssertEqual(t, false, errors.Is(err, ErrMaxRetry))
		tester.AssertEqual(t, false, errors.Is(err, ErrTimeout))
	})
}
