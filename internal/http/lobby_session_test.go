package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
	"github.com/suapapa/croquis-king/internal/ws"
)

func TestStartSessionHandler(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	lobbySync := ws.NewSnapshotSync(store)
	router := newRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), ws.NewHandler(lobbySync), lobbySync)

	created, photos := createReadyLobby(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/start", nil)
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
	if snapshot.Phase != lobby.PhaseDrawing {
		t.Fatalf("Phase = %q, want DRAWING", snapshot.Phase)
	}
	if snapshot.CurrentRound != 1 {
		t.Fatalf("CurrentRound = %d, want 1", snapshot.CurrentRound)
	}
	if snapshot.TotalRounds != len(photos) {
		t.Fatalf("TotalRounds = %d, want %d", snapshot.TotalRounds, len(photos))
	}
	if snapshot.CurrentPhoto == nil {
		t.Fatal("CurrentPhoto = nil, want photo in DRAWING")
	}
	if snapshot.DrawEndsAt == nil {
		t.Fatal("DrawEndsAt = nil, want timestamp in DRAWING")
	}
}

func TestStartSessionHandlerInvalidPhase(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/start", nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestNextRoundHandler(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))

	created, _ := createReadyLobby(t, store)
	if err := store.StartSession(context.Background(), created.ID, time.Now()); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if err := store.AdvanceToBetweenRounds(context.Background(), created.ID); err != nil {
		t.Fatalf("AdvanceToBetweenRounds() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/next", nil)
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
	if snapshot.Phase != lobby.PhaseDrawing {
		t.Fatalf("Phase = %q, want DRAWING", snapshot.Phase)
	}
	if snapshot.CurrentRound != 2 {
		t.Fatalf("CurrentRound = %d, want 2", snapshot.CurrentRound)
	}
	if snapshot.CurrentPhoto == nil {
		t.Fatal("CurrentPhoto = nil, want photo in DRAWING")
	}
}

func TestFinishSessionHandler(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))

	created, _ := createReadyLobby(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/finish", nil)
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
	if snapshot.Phase != lobby.PhaseFinished {
		t.Fatalf("Phase = %q, want FINISHED", snapshot.Phase)
	}
	if snapshot.CurrentRound != 1 {
		t.Fatalf("CurrentRound = %d, want 1", snapshot.CurrentRound)
	}
	if snapshot.CurrentPhoto != nil {
		t.Fatal("CurrentPhoto should be hidden in FINISHED")
	}
	if snapshot.DrawEndsAt != nil {
		t.Fatal("DrawEndsAt should be hidden in FINISHED")
	}
}

func TestStartSessionHandlerBroadcastsSnapshot(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	lobbySync := ws.NewSnapshotSync(store)
	wsHandler := ws.NewHandler(lobbySync)
	router := newRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), wsHandler, lobbySync)

	created, _ := createReadyLobby(t, store)

	conn := dialLobbyWSForRouter(t, router, created.ID)
	defer conn.Close()
	readSnapshotFromConn(t, conn)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/start", nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST status = %d, want %d", rec.Code, http.StatusOK)
	}

	snapshot := readSnapshotFromConn(t, conn)
	if snapshot.Phase != lobby.PhaseDrawing {
		t.Fatalf("Phase = %q, want DRAWING", snapshot.Phase)
	}
}

func createReadyLobby(t *testing.T, store *lobby.MemoryStore) (*lobby.Lobby, []lobby.Photo) {
	t.Helper()

	created, photos := createLobbyWithPhotos(t, store)
	if err := store.MarkReady(context.Background(), created.ID); err != nil {
		t.Fatalf("MarkReady() error = %v", err)
	}
	return created, photos
}
