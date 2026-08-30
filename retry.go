package backoff

import (
	"context"
	"errors"
	"math"
	"time"
)

var (
	ErrMaxRetry = errors.New("go-backoff/backoff: max retry count reached")
	ErrTimeout  = errors.New("go-backoff/backoff: max elapsed time reached")
)

// Notify reports the outcome of a retryable function to the retryer.
// The retryer stops retrying when either Stop or Success is called.
type Notify struct {
	stop    bool
	success bool
}

func (s *Notify) shouldStop() bool {
	return s.stop || s.success
}

// Stop tells the retryer to stop retrying because the operation
// cannot succeed by retrying.
// Use [Success] when the operation has completed successfully.
func (s *Notify) Stop() {
	s.stop = true
}

// Success tells the retryer to stop retrying because the operation
// has completed successfully.
// Use [Stop] when the operation should stop retrying without success.
func (s *Notify) Success() {
	s.success = true
}

// NewRetryer creates a new Retryer with given backoff.
// Returned retryer has no retry and has no timeout.
// It uses NewRandom(0, time.Second) when nil was given.
func NewRetryer(backoff Backoff) *Retryer {
	if backoff == nil {
		backoff, _ = NewRandom(0, time.Second)
	}
	return &Retryer{
		backoff: backoff,
	}
}

// Retryer runs retryable functions with the backoff algorithm.
type Retryer struct {
	backoff Backoff
	// MaxRetry specifies the maximum number of retries.
	// The retryable function is executed at most MaxRetry + 1 times.
	// A value of zero or less disables retries.
	MaxRetry int
	// MaxElapsedTime specifies the maximum duration allowed for retries.
	// The duration starts when Run or RunContext is called.
	// It does not cancel a retryable function that is already running.
	MaxElapsedTime time.Duration
}

// Run executes the retryable function using a background context.
func (r *Retryer) Run(retryable func(*Notify) error) error {
	return r.RunContext(context.Background(), retryable)
}

// RunContext executes the retryable function with the configured backoff strategy.
// It returns nil if [Notify.Success] is called.
// It returns the errors returned by the retryable function joined with
// [errors.Join] if [Notify.Stop] is called.
// It returns an error wrapping [ErrTimeout] and all errors returned by the
// retryable function if MaxElapsedTime is reached.
// It returns an error wrapping [ErrMaxRetry] and all errors returned by the
// retryable function if MaxRetry is reached.
func (r *Retryer) RunContext(ctx context.Context, retryable func(*Notify) error) error {
	var errs []error
	notify := &Notify{}

	var timeout <-chan time.Time
	if met := r.MaxElapsedTime; met > 0 {
		timeout = time.After(met)
	}

	maxRun := max(0, min(r.MaxRetry, math.MaxInt-1)) + 1
	for i := range maxRun {
		duration := r.backoff.Attempt(i)
		select {
		case <-time.After(duration):
		case <-ctx.Done():
			errs = append(errs, ctx.Err())
			return errors.Join(errs...)
		case <-timeout:
			errs = append(errs, ErrTimeout)
			return errors.Join(errs...)
		}

		err := retryable(notify)
		if err != nil {
			errs = append(errs, err)
		}
		if notify.shouldStop() {
			break
		}
	}
	if notify.success {
		return nil
	}
	if notify.stop {
		return errors.Join(errs...)
	}
	errs = append(errs, ErrMaxRetry)
	return errors.Join(errs...)
}
