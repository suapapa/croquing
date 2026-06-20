package ws

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/suapapa/croquis-king/internal/lobby"
)

func TestHubRegisterUnregisterAndCount(t *testing.T) {
	t.Parallel()

	sync := NewSnapshotSync(lobby.NewMemoryStore())
	hub := sync.Hub()
	conn := &websocket.Conn{}
	client := &Client{
		sync:    sync,
		lobbyID: "lobby-1",
		conn:    conn,
		send:    make(chan []byte, clientSendBuffer),
	}

	hub.Register("lobby-1", client)
	if got := hub.ClientCount("lobby-1"); got != 1 {
		t.Fatalf("ClientCount() = %d, want 1", got)
	}

	hub.Unregister("lobby-1", client)
	if got := hub.ClientCount("lobby-1"); got != 0 {
		t.Fatalf("ClientCount() after unregister = %d, want 0", got)
	}
}

func TestHubBroadcast(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	sync := NewSnapshotSync(store)
	handler := NewHandler(sync, nil)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	router := gin.New()
	router.GET("/ws/lobby/:id", handler.Handle)

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/lobby/" + created.ID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	defer conn.Close()

	readSnapshotMessage(t, conn)

	payload, err := MarshalEnvelope(MessageTypeSnapshot, map[string]string{"phase": "WAITING"})
	if err != nil {
		t.Fatalf("MarshalEnvelope() error = %v", err)
	}
	sync.Hub().Broadcast(created.ID, payload)

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage() error = %v", err)
	}

	var envelope Envelope
	if err := json.Unmarshal(message, &envelope); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if envelope.Type != MessageTypeSnapshot {
		t.Fatalf("type = %q, want snapshot", envelope.Type)
	}
}

func TestHandlerRejectsMissingLobby(t *testing.T) {
	t.Parallel()

	handler := NewHandler(NewSnapshotSync(lobby.NewMemoryStore()), nil)

	router := gin.New()
	router.GET("/ws/lobby/:id", handler.Handle)

	req := httptest.NewRequest(http.MethodGet, "/ws/lobby/missing", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestHandlerUpgradesExistingLobby(t *testing.T) {
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

	if sync.Hub().ClientCount(created.ID) != 1 {
		t.Fatalf("ClientCount() = %d, want 1", sync.Hub().ClientCount(created.ID))
	}

	snapshot := readSnapshotMessage(t, conn)
	if snapshot.ParticipantCount != 1 {
		t.Fatalf("ParticipantCount = %d, want 1", snapshot.ParticipantCount)
	}
}
