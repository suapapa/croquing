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

	"github.com/suapapa/croquing/internal/lobby"
	"github.com/suapapa/croquing/internal/pixabay"
	"github.com/suapapa/croquing/internal/ws"
)

func TestMarkReadyHandler(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	lobbySync := ws.NewSnapshotSync(store)
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), ws.NewHandler(lobbySync, nil), lobbySync)

	created, photos := createLobbyWithPhotos(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/ready", nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "photo_order") {
		t.Fatalf("response must not expose photo_order, body = %s", rec.Body.String())
	}

	var snapshot lobby.LobbySnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snapshot); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if snapshot.Phase != lobby.PhaseReady {
		t.Fatalf("Phase = %q, want READY", snapshot.Phase)
	}
	if snapshot.TotalRounds != len(photos) {
		t.Fatalf("TotalRounds = %d, want %d", snapshot.TotalRounds, len(photos))
	}
	if snapshot.CurrentPhoto != nil {
		t.Fatal("CurrentPhoto should be hidden in READY")
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Phase != lobby.PhaseReady {
		t.Fatalf("Phase = %q, want READY", got.Phase)
	}
	if len(got.PhotoOrder) != len(photos) {
		t.Fatalf("len(PhotoOrder) = %d, want %d", len(got.PhotoOrder), len(photos))
	}
}

func TestMarkReadyHandlerUnauthorized(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))

	created, _ := createLobbyWithPhotos(t, store)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/ready", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestMarkReadyHandlerInvalidPhase(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/ready", nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestMarkReadyHandlerBroadcastsSnapshot(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	lobbySync := ws.NewSnapshotSync(store)
	wsHandler := ws.NewHandler(lobbySync, nil)
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), wsHandler, lobbySync)

	created, photos := createLobbyWithPhotos(t, store)

	conn := dialLobbyWSForRouter(t, router, created.ID)
	defer conn.Close()
	readSnapshotFromConn(t, conn)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/ready", nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("POST status = %d, want %d", rec.Code, http.StatusOK)
	}

	snapshot := readSnapshotFromConn(t, conn)
	if snapshot.Phase != lobby.PhaseReady {
		t.Fatalf("Phase = %q, want READY", snapshot.Phase)
	}
	if snapshot.TotalRounds != len(photos) {
		t.Fatalf("TotalRounds = %d, want %d", snapshot.TotalRounds, len(photos))
	}
}

func createLobbyWithPhotos(t *testing.T, store *lobby.MemoryStore) (*lobby.Lobby, []lobby.Photo) {
	t.Helper()

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	photos := []lobby.Photo{
		{PixabayID: 1, LargeImageURL: "https://cdn.example/1.jpg"},
		{PixabayID: 2, LargeImageURL: "https://cdn.example/2.jpg"},
	}
	body, err := json.Marshal(setPhotosRequest{Photos: photos})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, nil)
	req := httptest.NewRequest(http.MethodPut, "/api/lobbies/"+created.ID+"/photos", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("setPhotos status = %d, want %d", rec.Code, http.StatusOK)
	}

	return created, photos
}
