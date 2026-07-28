package retry

import (
	"math/rand"
	"time"
)

func ApplyJitter(delay time.Duration) time.Duration {
	offset := time.Duration(rand.Int63n(int64(delay / 2)))
	return delay + offset
}