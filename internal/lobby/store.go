package lobby

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Store persists and retrieves lobbies.
type Store interface {
	Create(ctx context.Context, drawDuration time.Duration) (*Lobby, error)
	Get(ctx context.Context, id string) (*Lobby, error)
	Snapshot(ctx context.Context, id string, participantCount int) (LobbySnapshot, error)
	SetSelectedPhotos(ctx context.Context, id string, photos []Photo) error
	ReopenPhotoSelection(ctx context.Context, id string) error
	MarkReady(ctx context.Context, id string) error
	SetDrawDuration(ctx context.Context, id string, minutes int) error
	StartSession(ctx context.Context, id string, now time.Time) error
	AdvanceToBetweenRounds(ctx context.Context, id string) error
	NextRound(ctx context.Context, id string, now time.Time) error
	FinishSession(ctx context.Context, id string) error
	ExpireDrawingTimers(ctx context.Context, now time.Time) ([]string, error)
}

// MemoryStore is an in-memory lobby store protected by a mutex.
type MemoryStore struct {
	mu      sync.RWMutex
	lobbies map[string]*Lobby
}

// NewMemoryStore creates an empty in-memory lobby store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		lobbies: make(map[string]*Lobby),
	}
}

// Create inserts a new lobby in WAITING phase.
func (s *MemoryStore) Create(ctx context.Context, drawDuration time.Duration) (*Lobby, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	lobby := &Lobby{
		ID:             uuid.NewString(),
		AdminToken:     uuid.NewString(),
		Phase:          PhaseWaiting,
		SelectedPhotos: make([]Photo, 0),
		PhotoOrder:     make([]int, 0),
		DrawDuration:   drawDuration,
		CreatedAt:      time.Now(),
	}

	s.mu.Lock()
	s.lobbies[lobby.ID] = lobby
	s.mu.Unlock()

	return cloneLobby(lobby), nil
}

// Get returns a lobby by ID.
func (s *MemoryStore) Get(ctx context.Context, id string) (*Lobby, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	lobby, ok := s.lobbies[id]
	s.mu.RUnlock()
	if !ok {
		return nil, ErrNotFound
	}

	return cloneLobby(lobby), nil
}

// Snapshot returns the public snapshot for a lobby.
func (s *MemoryStore) Snapshot(ctx context.Context, id string, participantCount int) (LobbySnapshot, error) {
	if err := ctx.Err(); err != nil {
		return LobbySnapshot{}, err
	}

	s.mu.RLock()
	lobby, ok := s.lobbies[id]
	if !ok {
		s.mu.RUnlock()
		return LobbySnapshot{}, ErrNotFound
	}

	snap := lobby.Snapshot(participantCount, time.Now())
	s.mu.RUnlock()
	return snap, nil
}

// SetSelectedPhotos saves the admin's photo selection and moves the lobby to SELECTING.
func (s *MemoryStore) SetSelectedPhotos(ctx context.Context, id string, photos []Photo) error {
	if len(photos) == 0 {
		return ErrEmptyPhotos
	}

	return s.withLobby(ctx, id, func(lobby *Lobby) error {
		if err := ValidateTransition(lobby.Phase, PhaseSelecting); err != nil {
			return err
		}

		lobby.SelectedPhotos = append([]Photo(nil), photos...)
		lobby.Phase = PhaseSelecting
		lobby.PhotoOrder = nil
		lobby.CurrentRound = 0
		ClearDrawTimer(lobby)
		return nil
	})
}

// ReopenPhotoSelection moves a lobby from SELECTING back to WAITING so the admin can edit photos.
func (s *MemoryStore) ReopenPhotoSelection(ctx context.Context, id string) error {
	return s.withLobby(ctx, id, func(lobby *Lobby) error {
		if err := ValidateTransition(lobby.Phase, PhaseWaiting); err != nil {
			return err
		}
		if len(lobby.SelectedPhotos) == 0 {
			return ErrEmptyPhotos
		}

		lobby.Phase = PhaseWaiting
		lobby.PhotoOrder = nil
		lobby.CurrentRound = 0
		ClearDrawTimer(lobby)
		return nil
	})
}

// SetDrawDuration updates the per-round draw length while waiting to start or between rounds.
func (s *MemoryStore) SetDrawDuration(ctx context.Context, id string, minutes int) error {
	if err := ValidateDrawDurationMinutes(minutes); err != nil {
		return err
	}

	return s.withLobby(ctx, id, func(lobby *Lobby) error {
		switch lobby.Phase {
		case PhaseReady, PhaseBetweenRounds:
			lobby.DrawDuration = MinutesToDrawDuration(minutes)
			return nil
		default:
			return ErrInvalidTransition
		}
	})
}

// MarkReady shuffles selected photo indices and moves the lobby to READY.
func (s *MemoryStore) MarkReady(ctx context.Context, id string) error {
	return s.withLobby(ctx, id, func(lobby *Lobby) error {
		if lobby.Phase != PhaseSelecting {
			return ErrInvalidTransition
		}
		if len(lobby.SelectedPhotos) == 0 {
			return ErrEmptyPhotos
		}

		order, err := shuffleIndices(len(lobby.SelectedPhotos))
		if err != nil {
			return err
		}

		lobby.PhotoOrder = order
		lobby.Phase = PhaseReady
		lobby.CurrentRound = 0
		ClearDrawTimer(lobby)
		return nil
	})
}

func cloneLobby(lobby *Lobby) *Lobby {
	if lobby == nil {
		return nil
	}

	cloned := *lobby
	cloned.SelectedPhotos = append([]Photo(nil), lobby.SelectedPhotos...)
	cloned.PhotoOrder = append([]int(nil), lobby.PhotoOrder...)
	if lobby.DrawEndsAt != nil {
		endsAt := *lobby.DrawEndsAt
		cloned.DrawEndsAt = &endsAt
	}

	return &cloned
}
