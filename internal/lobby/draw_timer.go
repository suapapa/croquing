package lobby

import (
	"time"

	"github.com/suapapa/croquis-king/internal/timer"
)

// StartDrawTimer sets DrawEndsAt when entering DRAWING.
func StartDrawTimer(l *Lobby, now time.Time) {
	if l == nil {
		return
	}

	l.DrawEndsAt = timer.NewDeadline(now, l.DrawDuration)
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
