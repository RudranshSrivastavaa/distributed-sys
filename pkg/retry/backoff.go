package retry

import "time"

func NextDelay(current time.Duration,config Config) time.Duration {

	next := time.Duration(
		float64(current) * config.Multiplier,
	)

	if next > config.MaxDelay {
		return config.MaxDelay
	}

	return next
}