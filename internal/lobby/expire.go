package lobby

import (
	"context"
	"time"
)

// ExpireDrawingTimers moves expired DRAWING lobbies to BETWEEN_ROUNDS.
// Returns the lobby IDs that were transitioned.
func (s *MemoryStore) ExpireDrawingTimers(ctx context.Context, now time.Time) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var expired []string
	for id, lobby := range s.lobbies {
		if lobby.Phase != PhaseDrawing {
			continue
		}
		if !IsDrawExpired(lobby, now) {
			continue
		}

		lobby.Phase = PhaseBetweenRounds
		ClearDrawTimer(lobby)
		expired = append(expired, id)
	}

	return expired, nil
}
