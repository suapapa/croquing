package ws

import (
	"context"

	"github.com/suapapa/croquis-king/internal/lobby"
)

// SnapshotSync pushes lobby snapshots to WebSocket clients.
type SnapshotSync struct {
	hub   *Hub
	store lobby.Store
}

// NewSnapshotSync creates a snapshot sync service with a new hub.
func NewSnapshotSync(store lobby.Store) *SnapshotSync {
	return &SnapshotSync{
		hub:   NewHub(),
		store: store,
	}
}

// Hub returns the underlying WebSocket hub.
func (s *SnapshotSync) Hub() *Hub {
	return s.hub
}

// RegisterClient adds a client and broadcasts the latest snapshot to the lobby.
func (s *SnapshotSync) RegisterClient(ctx context.Context, lobbyID string, client *Client) error {
	s.hub.Register(lobbyID, client)
	return s.Broadcast(ctx, lobbyID)
}

// UnregisterClient removes a client and broadcasts an updated snapshot when others remain.
func (s *SnapshotSync) UnregisterClient(ctx context.Context, lobbyID string, client *Client) {
	s.hub.Unregister(lobbyID, client)
	if s.hub.ClientCount(lobbyID) > 0 {
		_ = s.Broadcast(ctx, lobbyID)
	}
}

// Broadcast sends the current lobby snapshot to all connected clients.
func (s *SnapshotSync) Broadcast(ctx context.Context, lobbyID string) error {
	snap, err := s.store.Snapshot(ctx, lobbyID, s.hub.ClientCount(lobbyID))
	if err != nil {
		return err
	}

	message, err := MarshalEnvelope(MessageTypeSnapshot, snap)
	if err != nil {
		return err
	}

	s.hub.Broadcast(lobbyID, message)
	return nil
}

// LobbyExists reports whether the lobby can be subscribed to.
func (s *SnapshotSync) LobbyExists(ctx context.Context, lobbyID string) error {
	_, err := s.store.Get(ctx, lobbyID)
	return err
}
