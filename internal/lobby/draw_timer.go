package lobby

import (
	"time"

	"github.com/suapapa/croquis-king/internal/timer"
)

// RoundCountdown is the pre-draw countdown shown before each round timer starts.
var RoundCountdown = 5 * time.Second

// StartDrawTimer sets DrawEndsAt when entering DRAWING.
// The deadline includes RoundCountdown plus the draw duration.
func StartDrawTimer(l *Lobby, now time.Time) {
	if l == nil {
		return
	}

	l.DrawEndsAt = timer.NewDeadline(now, RoundCountdown+l.DrawDuration)
}

// ClearDrawTimer removes the draw deadline from the lobby.
func ClearDrawTimer(l *Lobby) {
	if l == nil {
		return
	}

	l.DrawEndsAt = nil
}

// IsDrawExpired reports whether the lobby draw timer has expired at now.
func IsDrawExpired(l *Lobby, now time.Time) bool {
	if l == nil {
		return false
	}

	return timer.IsExpired(l.DrawEndsAt, now)
}
