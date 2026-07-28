package retry

import (
	"errors"
	"log"
	"time"

	"github.com/rudransh/distributed-commerce/pkg/circuitbreaker"
)

type Executor struct {
	config Config
}

func (r *Executor) Config() *Config {
	return &Config{
		MaxAttempts: 3,
		InitialDelay: time.Second,
		MaxDelay: 10 * time.Second,
		Multiplier: 3,
		Jitter: true,
	}
}

func NewExecutor(config Config) *Executor {
	return &Executor{
		config: config,
	}
}

func (r *Executor) Do(fn func() error) error {

	var err error

	delay := r.config.InitialDelay

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {

		err = fn()

		if err == nil {
			log.Printf("Attempt %d succeeded", attempt)
			return nil
		}

		//----------------------------------------------------
		// Circuit is already open.
		// Don't keep retrying.
		//----------------------------------------------------

		if errors.Is(err, circuitbreaker.ErrCircuitOpen) {
			log.Println("Circuit is OPEN. Stopping retries.")
			return err
		}

		//----------------------------------------------------
		// Non-retryable error
		//----------------------------------------------------

		if !IsRetryable(err) {
			return err
		}

		if attempt == r.config.MaxAttempts {
			break
		}

		if r.config.Jitter {
			delay = ApplyJitter(delay)
		}

		log.Printf("retry attempt %d failed: %v sleeping for %v",attempt,err,delay)

		time.Sleep(delay)

		delay = NextDelay(delay,r.config)
	}

	return err
}
