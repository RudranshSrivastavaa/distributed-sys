package retry

import "time"

type Config struct {
	MaxAttempts int

	InitialDelay time.Duration

	MaxDelay time.Duration

	Multiplier float64

	Jitter bool
}

func DefaultConfig() Config {

	return Config{
		MaxAttempts: 3,

		InitialDelay: time.Second,

		MaxDelay: 10 * time.Second,

		Multiplier: 3,

		Jitter: true,
	}

}