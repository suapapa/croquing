package httpserver

import (
	"bytes"
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
	"github.com/suapapa/croquis-king/internal/pixabay"
	"github.com/suapapa/croquis-king/internal/ws"
)

func TestSetPhotosHandler(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	lobbySync := ws.NewSnapshotSync(store)
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), ws.NewHandler(lobbySync, nil), lobbySync)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	body, err := json.Marshal(setPhotosRequest{Photos: []lobby.Photo{{
		PixabayID:     42,
		PreviewURL:    "https://cdn.example/preview.jpg",
		LargeImageURL: "https://cdn.example/large.jpg",
		PageURL:       "https://pixabay.com/photos/test-42/",
		Width:         800,
		Height:        600,
	}}})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/lobbies/"+created.ID+"/photos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var snapshot lobby.LobbySnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snapshot); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if snapshot.Phase != lobby.PhaseSelecting {
		t.Fatalf("Phase = %q, want SELECTING", snapshot.Phase)
	}
	if snapshot.SelectedCount != 1 {
		t.Fatalf("SelectedCount = %d, want 1", snapshot.SelectedCount)
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if len(got.SelectedPhotos) != 1 || got.SelectedPhotos[0].PixabayID != 42 {
		t.Fatalf("SelectedPhotos = %+v", got.SelectedPhotos)
	}
}

func TestReopenPhotoSelectionHandler(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	lobbySync := ws.NewSnapshotSync(store)
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), ws.NewHandler(lobbySync, nil), lobbySync)

	created, photos := createLobbyWithPhotos(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/photos/reopen", nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var snapshot lobby.LobbySnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snapshot); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if snapshot.Phase != lobby.PhaseWaiting {
		t.Fatalf("Phase = %q, want WAITING", snapshot.Phase)
	}
	if snapshot.SelectedCount != len(photos) {
		t.Fatalf("SelectedCount = %d, want %d", snapshot.SelectedCount, len(photos))
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Phase != lobby.PhaseWaiting {
		t.Fatalf("Phase = %q, want WAITING", got.Phase)
	}
}

func TestSetPhotosHandlerUnauthorized(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/lobbies/"+created.ID+"/photos", strings.NewReader(`{"photos":[]}`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestSetPhotosHandlerBroadcastsSnapshot(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	lobbySync := ws.NewSnapshotSync(store)
	wsHandler := ws.NewHandler(lobbySync, nil)
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), wsHandler, lobbySync)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	conn := dialLobbyWSForRouter(t, router, created.ID)
	defer conn.Close()
	readSnapshotFromConn(t, conn)

	body, err := json.Marshal(setPhotosRequest{Photos: []lobby.Photo{{
		PixabayID:     99,
		LargeImageURL: "https://cdn.example/large.jpg",
	}}})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/api/lobbies/"+created.ID+"/photos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("PUT status = %d, want %d", rec.Code, http.StatusOK)
	}

	snapshot := readSnapshotFromConn(t, conn)
	if snapshot.Phase != lobby.PhaseSelecting {
		t.Fatalf("Phase = %q, want SELECTING", snapshot.Phase)
	}
	if snapshot.SelectedCount != 1 {
		t.Fatalf("SelectedCount = %d, want 1", snapshot.SelectedCount)
	}
}

func dialLobbyWSForRouter(t *testing.T, router *gin.Engine, lobbyID string) *websocket.Conn {
	t.Helper()

	server := httptest.NewServer(router)
	t.Cleanup(server.Close)

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws/lobby/" + lobbyID
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Dial() error = %v", err)
	}
	return conn
}

func readSnapshotFromConn(t *testing.T, conn *websocket.Conn) lobby.LobbySnapshot {
	t.Helper()

	_ = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("ReadMessage() error = %v", err)
	}

	var envelope struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.Unmarshal(message, &envelope); err != nil {
		t.Fatalf("Unmarshal envelope: %v", err)
	}

	var snapshot lobby.LobbySnapshot
	if err := json.Unmarshal(envelope.Payload, &snapshot); err != nil {
		t.Fatalf("Unmarshal snapshot: %v", err)
	}
	return snapshot
}
