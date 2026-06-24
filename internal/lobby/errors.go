package lobby

import "errors"

var (
	// ErrNotFound is returned when a lobby does not exist.
	ErrNotFound = errors.New("lobby: not found")

	// ErrInvalidTransition is returned when a phase transition is not allowed.
	ErrInvalidTransition = errors.New("lobby: invalid phase transition")

	// ErrEmptyPhotos is returned when a photo selection request has no photos.
	ErrEmptyPhotos = errors.New("lobby: photos are required")

	// ErrPhotosNotReady is returned when session photos are not available for download.
	ErrPhotosNotReady = errors.New("lobby: session photos are not ready")

	// ErrInvalidDrawDuration is returned when draw duration minutes are out of range.
	ErrInvalidDrawDuration = errors.New("lobby: invalid draw duration")
)
