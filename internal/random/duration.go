package random

import (
	"math/rand"
	"time"
)

// Dur returns a pseudo-random Duration in [0, max)
func Dur(max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max)))
}

// Uniformly jitters the provided duration by +/- 10%.
func Jitter(period time.Duration) time.Duration {
	return JitterFraction(period, .2)
}

// Uniformly jitters the provided duration by +/- the given fraction.  NOTE:
// fraction must be in (0, 1].
func JitterFraction(period time.Duration, fraction float64) time.Duration {
	fixed := time.Duration(float64(period) * (1 - fraction))
	return fixed + Dur(2*(period-fixed))
}
