package lobby

import "errors"

var (
	// ErrNotFound is returned when a lobby does not exist.
	ErrNotFound = errors.New("lobby: not found")

	// ErrInvalidTransition is returned when a phase transition is not allowed.
	ErrInvalidTransition = errors.New("lobby: invalid phase transition")
)
