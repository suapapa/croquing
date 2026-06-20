package timer

import "time"

// EndsAt returns the draw expiry instant for a timer started at now.
func EndsAt(now time.Time, duration time.Duration) time.Time {
	return now.Add(duration)
}

// NewDeadline returns a draw expiry pointer for a timer started at now.
func NewDeadline(now time.Time, duration time.Duration) *time.Time {
	endsAt := EndsAt(now, duration)
	return &endsAt
}

// IsExpired reports whether now is at or after drawEndsAt.
// A nil drawEndsAt is never expired.
func IsExpired(drawEndsAt *time.Time, now time.Time) bool {
	if drawEndsAt == nil {
		return false
	}

	return !now.Before(*drawEndsAt)
}

// Remaining returns time left until drawEndsAt, or zero if unset or expired.
func Remaining(drawEndsAt *time.Time, now time.Time) time.Duration {
	if drawEndsAt == nil {
		return 0
	}

	remaining := drawEndsAt.Sub(now)
	if remaining <= 0 {
		return 0
	}

	return remaining
}
