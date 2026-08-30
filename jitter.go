package backoff

import (
	"math"
	"math/rand/v2"
	"time"
)

const (
	NoJitter     float64 = 0.0
	FullJitter   float64 = 1.0
	EqualJitter  float64 = 0.5
	HalfJitter   float64 = 0.5
	QuaterJitter float64 = 0.25
)

// maxDurationFloat is the largest float64 value that is strictly less than 2^63.
// float64(math.MaxInt64) itself rounds up to exactly 2^63 (since MaxInt64 = 2^63-1
// is not exactly representable in float64's 53-bit mantissa), which would overflow
// int64 on conversion. Using the nearest representable value below 2^63 avoids that.
var maxDurationFloat = math.Nextafter(1<<63, 0)

// wilthJitter applies jitter to the delta and returns result in duration.
// Given jitter must be 0.0<=jitter<=1.0.
// 1.0 for full jitter. 0.0 for no jitter.
func withJitter(jitter, offset, delta float64) time.Duration {
	x := (1 - jitter) * delta
	y := rand.Float64() * jitter * delta
	return time.Duration(min(offset+x+y, maxDurationFloat))
}
