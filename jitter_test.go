package backoff

import (
	"math"
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestWithJitter(t *testing.T) {
	t.Parallel()
	t.Run("NoJitter", func(t *testing.T) {
		t.Run("max int64 offset wo delta", func(t *testing.T) {
			d := withJitter(NoJitter, math.MaxInt64, 0)
			tester.AssertEqual(t, true, d >= 0)
		})
		t.Run("max int64 offset w delta", func(t *testing.T) {
			d := withJitter(NoJitter, math.MaxInt64, math.MaxInt64)
			tester.AssertEqual(t, true, d >= 0)
		})
		t.Run("max float64 offset wo delta", func(t *testing.T) {
			d := withJitter(NoJitter, math.MaxFloat64, 0)
			tester.AssertEqual(t, true, d >= 0)
		})
		t.Run("max float64 offset w delta", func(t *testing.T) {
			d := withJitter(NoJitter, math.MaxFloat64, 100)
			tester.AssertEqual(t, true, d >= 0)
		})
		t.Run("max float64 offset w delta", func(t *testing.T) {
			d := withJitter(NoJitter, math.MaxFloat64, math.MaxFloat64)
			tester.AssertEqual(t, true, d >= 0)
		})
	})
	t.Run("FullJitter", func(t *testing.T) {
		t.Run("max int64 offset wo delta", func(t *testing.T) {
			d := withJitter(FullJitter, math.MaxInt64, 0)
			tester.AssertEqual(t, true, d >= 0)
		})
		t.Run("max int64 offset w delta", func(t *testing.T) {
			d := withJitter(FullJitter, math.MaxInt64, math.MaxInt64)
			tester.AssertEqual(t, true, d >= 0)
		})
		t.Run("max float64 offset wo delta", func(t *testing.T) {
			d := withJitter(FullJitter, math.MaxFloat64, 0)
			tester.AssertEqual(t, true, d >= 0)
		})
		t.Run("max float64 offset w delta", func(t *testing.T) {
			d := withJitter(FullJitter, math.MaxFloat64, 100)
			tester.AssertEqual(t, true, d >= 0)
		})
		t.Run("max float64 offset w delta", func(t *testing.T) {
			d := withJitter(FullJitter, math.MaxFloat64, math.MaxFloat64)
			tester.AssertEqual(t, true, d >= 0)
		})
	})
}
