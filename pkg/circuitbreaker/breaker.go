package circuitbreaker

import (
	"log"
	"sync"
	"time"
)

type Breaker struct {
	mu sync.Mutex

	state State

	config Config

	failures int

	successes int

	lastFailure time.Time
}

func New(config Config) *Breaker {

	return &Breaker{

		state: StateClosed,

		config: config,
	}

}

func (b *Breaker) shouldAllow() bool {

	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {

	case StateClosed:
		return true

	case StateOpen:

		if time.Since(b.lastFailure) >= b.config.RecoveryTimeout {
			log.Println("[CircuitBreaker] OPEN -> HALF_OPEN")

			b.transition(StateHalfOpen)
			b.successes = 0

			return true
		}

		return false

	case StateHalfOpen:
		return true
	}

	return false
}

func (b *Breaker) recordSuccess() {

	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {

	case StateClosed:

		b.failures = 0

	case StateHalfOpen:

		b.successes++

		if b.successes >= b.config.SuccessThreshold {
			log.Println("[CircuitBreaker] HALF_OPEN -> CLOSED")
			b.transition(StateClosed)
			b.failures = 0
			b.successes = 0
		}
	}
}

func (b *Breaker) recordFailure() {

	b.mu.Lock()
	defer b.mu.Unlock()

	b.lastFailure = time.Now()

	switch b.state {

	case StateClosed:

		b.failures++

		if b.failures >= b.config.FailureThreshold {
			log.Println("[CircuitBreaker] CLOSED -> OPEN")

			b.transition(StateOpen)
		}

	case StateHalfOpen:

		b.transition(StateOpen)
		b.successes = 0
	}
}

func (b *Breaker) Execute(
	fn func() error,
) error {

	if !b.shouldAllow() {
		return ErrCircuitOpen
	}

	err := fn()

	if err != nil {

		b.recordFailure()

		return err
	}

	b.recordSuccess()

	return nil
}

func (b *Breaker) transition(
	newState State,
) {

	if b.state == newState {
		return
	}

	log.Printf(
		"[CircuitBreaker] Transition %s -> %s",
		b.state,
		newState,
	)

	b.state = newState
}