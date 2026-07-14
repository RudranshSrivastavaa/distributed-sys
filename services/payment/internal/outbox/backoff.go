package outbox

import "time"

func CalculateNextRetry(

	retryCount int,

) time.Time {
	const MaxRetryDelay = 5 * time.Minute

	delay := time.Second << retryCount

	if delay > MaxRetryDelay {
	 delay = MaxRetryDelay
    }

	return time.Now().UTC().Add(delay)
}