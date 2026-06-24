package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/suapapa/croquing/internal/lobby"
	"github.com/suapapa/croquing/internal/pixabay"
)

func TestAdminMiddlewareRejectsMissingToken(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/start", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
}

func TestAdminMiddlewareRejectsInvalidToken(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/lobbies/"+created.ID+"/start", nil)
	req.Header.Set(lobby.AdminTokenHeader, "wrong-token")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusForbidden, rec.Body.String())
	}
}

func TestAdminMiddlewareUsesLobbyIDQueryForPixabay(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/pixabay/search?lobby_id="+created.ID+"&q=cat", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
}
