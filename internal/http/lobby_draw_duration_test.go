package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/suapapa/croquing/internal/lobby"
	"github.com/suapapa/croquing/internal/pixabay"
	"github.com/suapapa/croquing/internal/ws"
)

func TestSetDrawDurationHandler(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	lobbySync := ws.NewSnapshotSync(store)
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), ws.NewHandler(lobbySync, nil), lobbySync)

	created, _ := createLobbyWithPhotos(t, store)
	markLobbyReady(t, router, created)

	body := bytes.NewBufferString(`{"minutes":10}`)
	req := httptest.NewRequest(http.MethodPut, "/api/lobbies/"+created.ID+"/draw-duration", body)
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
	if snapshot.DrawDurationMinutes != 10 {
		t.Fatalf("DrawDurationMinutes = %d, want 10", snapshot.DrawDurationMinutes)
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.DrawDuration != 10*time.Minute {
		t.Fatalf("DrawDuration = %v, want 10m", got.DrawDuration)
	}
}

func TestSetDrawDurationHandlerInvalidPhase(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))

	created, _ := createLobbyWithPhotos(t, store)

	body := bytes.NewBufferString(`{"minutes":10}`)
	req := httptest.NewRequest(http.MethodPut, "/api/lobbies/"+created.ID+"/draw-duration", body)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusConflict)
	}
}

func TestSetDrawDurationHandlerOutOfRange(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))

	created, _ := createLobbyWithPhotos(t, store)
	markLobbyReady(t, router, created)

	body := bytes.NewBufferString(`{"minutes":90}`)
	req := httptest.NewRequest(http.MethodPut, "/api/lobbies/"+created.ID+"/draw-duration", body)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func markLobbyReady(t *testing.T, router http.Handler, created *lobby.Lobby) {
	t.Helper()

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/ready", nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("mark ready status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}
