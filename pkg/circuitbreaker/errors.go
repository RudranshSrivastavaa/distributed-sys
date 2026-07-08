package circuitbreaker

import "errors"

var ErrCircuitOpen = errors.New(
	"circuit breaker is open",
)