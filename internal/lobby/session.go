package lobby

import (
	"context"
	"time"
)

// StartSession begins the first drawing round from READY.
func (s *MemoryStore) StartSession(ctx context.Context, id string, now time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lobby, ok := s.lobbies[id]
	if !ok {
		return ErrNotFound
	}
	if lobby.Phase != PhaseReady {
		return ErrInvalidTransition
	}
	if len(lobby.PhotoOrder) == 0 {
		return ErrEmptyPhotos
	}

	lobby.CurrentRound = 0
	beginDrawing(lobby, now)
	return nil
}

// NextRound advances from BETWEEN_ROUNDS to the next DRAWING or FINISHED.
func (s *MemoryStore) NextRound(ctx context.Context, id string, now time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lobby, ok := s.lobbies[id]
	if !ok {
		return ErrNotFound
	}
	if lobby.Phase != PhaseBetweenRounds {
		return ErrInvalidTransition
	}

	if lobby.CurrentRound+1 >= len(lobby.PhotoOrder) {
		lobby.Phase = PhaseFinished
		ClearDrawTimer(lobby)
		return nil
	}

	lobby.CurrentRound++
	beginDrawing(lobby, now)
	return nil
}

// AdvanceToBetweenRounds ends the current drawing round and waits for the next action.
func (s *MemoryStore) AdvanceToBetweenRounds(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lobby, ok := s.lobbies[id]
	if !ok {
		return ErrNotFound
	}
	if lobby.Phase != PhaseDrawing {
		return ErrInvalidTransition
	}

	lobby.Phase = PhaseBetweenRounds
	ClearDrawTimer(lobby)
	return nil
}

// FinishSession ends the session immediately.
func (s *MemoryStore) FinishSession(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lobby, ok := s.lobbies[id]
	if !ok {
		return ErrNotFound
	}
	if lobby.Phase == PhaseFinished {
		return ErrInvalidTransition
	}

	lobby.Phase = PhaseFinished
	ClearDrawTimer(lobby)
	return nil
}

func beginDrawing(l *Lobby, now time.Time) {
	l.Phase = PhaseDrawing
	StartDrawTimer(l, now)
}
