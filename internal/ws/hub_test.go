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

	hub := NewHub()
	conn := &websocket.Conn{}
	client := newClient(hub, "lobby-1", conn)

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
	hub := NewHub()
	handler := NewHandler(hub, store)

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

	deadline := time.Now().Add(time.Second)
	for hub.ClientCount(created.ID) == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}

	payload, err := MarshalEnvelope("snapshot", map[string]string{"phase": "WAITING"})
	if err != nil {
		t.Fatalf("MarshalEnvelope() error = %v", err)
	}
	hub.Broadcast(created.ID, payload)

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage() error = %v", err)
	}

	var envelope Envelope
	if err := json.Unmarshal(message, &envelope); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if envelope.Type != "snapshot" {
		t.Fatalf("type = %q, want snapshot", envelope.Type)
	}
}

func TestHandlerRejectsMissingLobby(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	handler := NewHandler(NewHub(), store)

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
	hub := NewHub()
	handler := NewHandler(hub, store)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	router := gin.New()
	router.GET("/ws/lobby/:id", handler.Handle)

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/lobby/" + created.ID
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	defer conn.Close()

	if resp.StatusCode != http.StatusSwitchingProtocols {
		t.Fatalf("status = %d, want %d", resp.StatusCode, http.StatusSwitchingProtocols)
	}

	deadline := time.Now().Add(500 * time.Millisecond)
	for hub.ClientCount(created.ID) == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	if hub.ClientCount(created.ID) != 1 {
		t.Fatalf("ClientCount() = %d, want 1", hub.ClientCount(created.ID))
	}
}
