package httpserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
)

func TestCreateLobbyHandler(t *testing.T) {
	t.Parallel()

	router := newRouter(lobby.NewMemoryStore(), 5*time.Minute, pixabay.NewClient("test-key"), nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusCreated)
	}

	var resp createLobbyResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.ID == "" || resp.AdminToken == "" {
		t.Fatal("response missing id or admin_token")
	}
	if resp.JoinURL != "http://"+req.Host+"/lobby/"+resp.ID {
		t.Fatalf("JoinURL = %q", resp.JoinURL)
	}
}

func TestGetLobbyHandler(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, nil)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/lobbies/"+created.ID, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var snapshot lobby.LobbySnapshot
	if err := json.NewDecoder(rec.Body).Decode(&snapshot); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if snapshot.ID != created.ID {
		t.Fatalf("ID = %q, want %q", snapshot.ID, created.ID)
	}
	if snapshot.Phase != lobby.PhaseWaiting {
		t.Fatalf("Phase = %q, want WAITING", snapshot.Phase)
	}
}

func TestGetLobbyHandlerNotFound(t *testing.T) {
	t.Parallel()

	router := newRouter(lobby.NewMemoryStore(), 5*time.Minute, pixabay.NewClient("test-key"), nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/lobbies/missing", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestValidateAdminTokenHelper(t *testing.T) {
	t.Parallel()

	lob := &lobby.Lobby{AdminToken: "abc"}
	if !lobby.ValidateAdminToken(lob, "abc") {
		t.Fatal("ValidateAdminToken() = false, want true")
	}
	if lobby.ValidateAdminToken(lob, "def") {
		t.Fatal("ValidateAdminToken() = true, want false")
	}
	if lobby.ValidateAdminToken(nil, "abc") {
		t.Fatal("ValidateAdminToken(nil) = true, want false")
	}
}

func TestJoinURLUsesHTTPSWhenForwarded(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	c.Request.Host = "example.com"
	c.Request.Header.Set("X-Forwarded-Proto", "https")

	got := joinURL(c, "lobby-id")
	want := "https://example.com/lobby/lobby-id"
	if got != want {
		t.Fatalf("joinURL() = %q, want %q", got, want)
	}
}

func TestGetLobbyHandlerStoreError(t *testing.T) {
	t.Parallel()

	router := newRouter(errorStore{}, 5*time.Minute, pixabay.NewClient("test-key"), nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/lobbies/any", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}

type errorStore struct{}

func (errorStore) Create(_ context.Context, _ time.Duration) (*lobby.Lobby, error) {
	return nil, errors.New("boom")
}

func (errorStore) Get(_ context.Context, _ string) (*lobby.Lobby, error) {
	return nil, errors.New("boom")
}

func (errorStore) Snapshot(_ context.Context, _ string, _ int) (lobby.LobbySnapshot, error) {
	return lobby.LobbySnapshot{}, errors.New("boom")
}

func (errorStore) SetSelectedPhotos(_ context.Context, _ string, _ []lobby.Photo) error {
	return errors.New("boom")
}
