package outbox

import "time"

type Config struct {

	PollInterval time.Duration

	BatchSize int

	MaxRetries int

	MaxBatchDuration time.Duration
}

func DefaultConfig() Config {

	return Config{

		PollInterval: 2 * time.Second,

		BatchSize: 100,

		MaxRetries: 5,

		MaxBatchDuration: 30 * time.Second,
	}
}