package ws

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/google/uuid"

	"github.com/suapapa/croquis-king/internal/lobby"
)

func TestInitialSnapshotOnConnect(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	sync := NewSnapshotSync(store)
	handler := NewHandler(sync, nil)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := dialLobbyWS(t, handler, created.ID)
	defer conn.Close()

	snapshot := readSnapshotMessage(t, conn)
	if snapshot.ID != created.ID {
		t.Fatalf("ID = %q, want %q", snapshot.ID, created.ID)
	}
	if snapshot.Phase != lobby.PhaseWaiting {
		t.Fatalf("Phase = %q, want WAITING", snapshot.Phase)
	}
	if snapshot.ParticipantCount != 1 {
		t.Fatalf("ParticipantCount = %d, want 1", snapshot.ParticipantCount)
	}
}

func TestParticipantCountUpdatesOnConnectAndDisconnect(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	sync := NewSnapshotSync(store)
	handler := NewHandler(sync, nil)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn1 := dialLobbyWS(t, handler, created.ID)
	defer conn1.Close()
	readSnapshotMessage(t, conn1)

	conn2 := dialLobbyWS(t, handler, created.ID)
	defer conn2.Close()

	snapshot1 := readSnapshotMessage(t, conn1)
	snapshot2 := readSnapshotMessage(t, conn2)
	if snapshot1.ParticipantCount != 2 || snapshot2.ParticipantCount != 2 {
		t.Fatalf("counts = %d/%d, want 2/2", snapshot1.ParticipantCount, snapshot2.ParticipantCount)
	}

	_ = conn2.Close()
	waitForClientCount(t, sync.Hub(), created.ID, 1)

	snapshot1 = readSnapshotMessage(t, conn1)
	if snapshot1.ParticipantCount != 1 {
		t.Fatalf("ParticipantCount after disconnect = %d, want 1", snapshot1.ParticipantCount)
	}
}

func TestBroadcastPushesUpdatedSnapshot(t *testing.T) {
	t.Parallel()

	store := newFakeStore()
	sync := NewSnapshotSync(store)
	handler := NewHandler(sync, nil)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := dialLobbyWS(t, handler, created.ID)
	defer conn.Close()
	readSnapshotMessage(t, conn)

	store.setPhase(created.ID, lobby.PhaseSelecting)
	if err := sync.Broadcast(context.Background(), created.ID); err != nil {
		t.Fatalf("Broadcast() error = %v", err)
	}

	snapshot := readSnapshotMessage(t, conn)
	if snapshot.Phase != lobby.PhaseSelecting {
		t.Fatalf("Phase = %q, want SELECTING", snapshot.Phase)
	}
}

func dialLobbyWS(t *testing.T, handler *Handler, lobbyID string) *websocket.Conn {
	t.Helper()

	router := gin.New()
	router.GET("/ws/lobby/:id", handler.Handle)

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/lobby/" + lobbyID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	return conn
}

func readSnapshotMessage(t *testing.T, conn *websocket.Conn) lobby.LobbySnapshot {
	t.Helper()

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage() error = %v", err)
	}

	var envelope Envelope
	if err := json.Unmarshal(message, &envelope); err != nil {
		t.Fatalf("Unmarshal envelope: %v", err)
	}
	if envelope.Type != MessageTypeSnapshot {
		t.Fatalf("type = %q, want %q", envelope.Type, MessageTypeSnapshot)
	}

	var snapshot lobby.LobbySnapshot
	if err := json.Unmarshal(envelope.Payload, &snapshot); err != nil {
		t.Fatalf("Unmarshal snapshot: %v", err)
	}
	return snapshot
}

func waitForClientCount(t *testing.T, hub *Hub, lobbyID string, want int) {
	t.Helper()

	deadline := time.Now().Add(time.Second)
	for hub.ClientCount(lobbyID) != want && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if got := hub.ClientCount(lobbyID); got != want {
		t.Fatalf("ClientCount() = %d, want %d", got, want)
	}
}

type fakeStore struct {
	mu      sync.RWMutex
	lobbies map[string]*lobby.Lobby
}

func newFakeStore() *fakeStore {
	return &fakeStore{lobbies: make(map[string]*lobby.Lobby)}
}

func (s *fakeStore) Create(ctx context.Context, drawDuration time.Duration) (*lobby.Lobby, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	lob := &lobby.Lobby{
		ID:             uuid.NewString(),
		AdminToken:     uuid.NewString(),
		Phase:          lobby.PhaseWaiting,
		SelectedPhotos: make([]lobby.Photo, 0),
		PhotoOrder:     make([]int, 0),
		DrawDuration:   drawDuration,
		CreatedAt:      time.Now(),
	}

	s.mu.Lock()
	s.lobbies[lob.ID] = lob
	s.mu.Unlock()

	cloned := *lob
	return &cloned, nil
}

func (s *fakeStore) Get(ctx context.Context, id string) (*lobby.Lobby, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.RLock()
	lob, ok := s.lobbies[id]
	s.mu.RUnlock()
	if !ok {
		return nil, lobby.ErrNotFound
	}

	cloned := *lob
	return &cloned, nil
}

func (s *fakeStore) Snapshot(ctx context.Context, id string, participantCount int) (lobby.LobbySnapshot, error) {
	lob, err := s.Get(ctx, id)
	if err != nil {
		return lobby.LobbySnapshot{}, err
	}
	return lob.Snapshot(participantCount, time.Now()), nil
}

func (s *fakeStore) SetSelectedPhotos(ctx context.Context, id string, photos []lobby.Photo) error {
	if len(photos) == 0 {
		return lobby.ErrEmptyPhotos
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lob, ok := s.lobbies[id]
	if !ok {
		return lobby.ErrNotFound
	}
	if err := lobby.ValidateTransition(lob.Phase, lobby.PhaseSelecting); err != nil {
		return err
	}

	lob.SelectedPhotos = append([]lobby.Photo(nil), photos...)
	lob.Phase = lobby.PhaseSelecting
	lob.PhotoOrder = nil
	lob.CurrentRound = 0
	lob.DrawEndsAt = nil
	return nil
}

func (s *fakeStore) MarkReady(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lob, ok := s.lobbies[id]
	if !ok {
		return lobby.ErrNotFound
	}
	if lob.Phase != lobby.PhaseSelecting {
		return lobby.ErrInvalidTransition
	}
	if len(lob.SelectedPhotos) == 0 {
		return lobby.ErrEmptyPhotos
	}

	order := make([]int, len(lob.SelectedPhotos))
	for i := range order {
		order[i] = i
	}

	lob.PhotoOrder = order
	lob.Phase = lobby.PhaseReady
	lob.CurrentRound = 0
	lob.DrawEndsAt = nil
	return nil
}

func (s *fakeStore) StartSession(ctx context.Context, id string, now time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lob, ok := s.lobbies[id]
	if !ok {
		return lobby.ErrNotFound
	}
	if lob.Phase != lobby.PhaseReady || len(lob.PhotoOrder) == 0 {
		return lobby.ErrInvalidTransition
	}

	lob.CurrentRound = 0
	lob.Phase = lobby.PhaseDrawing
	endsAt := now.Add(lob.DrawDuration)
	lob.DrawEndsAt = &endsAt
	return nil
}

func (s *fakeStore) AdvanceToBetweenRounds(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lob, ok := s.lobbies[id]
	if !ok {
		return lobby.ErrNotFound
	}
	if lob.Phase != lobby.PhaseDrawing {
		return lobby.ErrInvalidTransition
	}

	lob.Phase = lobby.PhaseBetweenRounds
	lob.DrawEndsAt = nil
	return nil
}

func (s *fakeStore) NextRound(ctx context.Context, id string, now time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lob, ok := s.lobbies[id]
	if !ok {
		return lobby.ErrNotFound
	}
	if lob.Phase != lobby.PhaseBetweenRounds {
		return lobby.ErrInvalidTransition
	}

	if lob.CurrentRound+1 >= len(lob.PhotoOrder) {
		lob.Phase = lobby.PhaseFinished
		lob.DrawEndsAt = nil
		return nil
	}

	lob.CurrentRound++
	lob.Phase = lobby.PhaseDrawing
	endsAt := now.Add(lob.DrawDuration)
	lob.DrawEndsAt = &endsAt
	return nil
}

func (s *fakeStore) FinishSession(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lob, ok := s.lobbies[id]
	if !ok {
		return lobby.ErrNotFound
	}
	if lob.Phase == lobby.PhaseFinished {
		return lobby.ErrInvalidTransition
	}

	lob.Phase = lobby.PhaseFinished
	lob.DrawEndsAt = nil
	return nil
}

func (s *fakeStore) ExpireDrawingTimers(ctx context.Context, now time.Time) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	var expired []string
	for id, lob := range s.lobbies {
		if lob.Phase != lobby.PhaseDrawing || lob.DrawEndsAt == nil {
			continue
		}
		if now.Before(*lob.DrawEndsAt) {
			continue
		}

		lob.Phase = lobby.PhaseBetweenRounds
		lob.DrawEndsAt = nil
		expired = append(expired, id)
	}

	return expired, nil
}

func (s *fakeStore) setPhase(id string, phase lobby.LobbyPhase) {
	s.mu.Lock()
	s.lobbies[id].Phase = phase
	s.mu.Unlock()
}
