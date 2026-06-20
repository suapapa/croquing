package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
)

func TestHealthHandler(t *testing.T) {
	t.Parallel()

	router := newRouter(lobby.NewMemoryStore(), 5*time.Minute, pixabay.NewClient("test-key"))

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body := rec.Body.String(); body != "OK" {
		t.Fatalf("body = %q, want %q", body, "OK")
	}
}

func TestNoRoute(t *testing.T) {
	t.Parallel()

	router := newRouter(lobby.NewMemoryStore(), 5*time.Minute, pixabay.NewClient("test-key"))

	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
