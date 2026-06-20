package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
)

func TestCORSPreflightAllowsViteDevOrigin(t *testing.T) {
	t.Parallel()

	router := newRouter(
		lobby.NewMemoryStore(),
		5*time.Minute,
		pixabay.NewClient("test-key"),
		nil,
		nil,
		[]string{"http://localhost:5173"},
	)

	req := httptest.NewRequest(http.MethodOptions, "/api/lobbies", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
	}
	if got := rec.Header().Get("Access-Control-Allow-Headers"); got == "" {
		t.Fatal("Access-Control-Allow-Headers is empty")
	}
}

func TestCORSRejectsUnknownOriginOnPreflight(t *testing.T) {
	t.Parallel()

	router := newRouter(
		lobby.NewMemoryStore(),
		5*time.Minute,
		pixabay.NewClient("test-key"),
		nil,
		nil,
		[]string{"http://localhost:5173"},
	)

	req := httptest.NewRequest(http.MethodOptions, "/api/lobbies", nil)
	req.Header.Set("Origin", "http://evil.example")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want empty", got)
	}
}

func TestCORSAllowsWildcardOrigin(t *testing.T) {
	t.Parallel()

	router := newTestRouter(lobby.NewMemoryStore(), 5*time.Minute, pixabay.NewClient("test-key"), nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:5173" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, "http://localhost:5173")
	}
}
