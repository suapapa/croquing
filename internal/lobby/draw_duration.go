package lobby

import "time"

const (
	// MinDrawDurationMinutes is the shortest allowed draw round length.
	MinDrawDurationMinutes = 1
	// MaxDrawDurationMinutes is the longest allowed draw round length.
	MaxDrawDurationMinutes = 60
)

// DrawDurationToMinutes converts a draw duration to whole minutes.
func DrawDurationToMinutes(d time.Duration) int {
	return int(d / time.Minute)
}

// MinutesToDrawDuration converts whole minutes to a draw duration.
func MinutesToDrawDuration(minutes int) time.Duration {
	return time.Duration(minutes) * time.Minute
}

// ValidateDrawDurationMinutes reports whether minutes is within allowed bounds.
func ValidateDrawDurationMinutes(minutes int) error {
	if minutes < MinDrawDurationMinutes || minutes > MaxDrawDurationMinutes {
		return ErrInvalidDrawDuration
	}
	return nil
}
