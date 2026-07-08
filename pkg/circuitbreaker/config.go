package circuitbreaker

import "time"

type Config struct {

	FailureThreshold int

	RecoveryTimeout time.Duration

	SuccessThreshold int
}

func DefaultConfig() Config {

	return Config{

		FailureThreshold: 3,

		RecoveryTimeout: 10 * time.Second,

		SuccessThreshold: 1,
	}

}