package httpserver

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
)

func TestPixabaySearchHandler(t *testing.T) {
	t.Parallel()

	pixabayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "99")
		w.Header().Set("X-RateLimit-Reset", "30")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"total": 10,
			"totalHits": 5,
			"hits": [{
				"id": 42,
				"pageURL": "https://pixabay.com/photos/test-42/",
				"previewURL": "https://cdn.pixabay.com/preview.jpg",
				"webformatURL": "https://cdn.pixabay.com/web.jpg",
				"largeImageURL": "https://cdn.pixabay.com/large.jpg",
				"imageWidth": 800,
				"imageHeight": 600,
				"views": 100,
				"downloads": 50,
				"likes": 10
			}]
		}`))
	}))
	t.Cleanup(pixabayServer.Close)

	store := lobby.NewMemoryStore()
	client := pixabay.NewClient("test-key", pixabay.WithBaseURL(pixabayServer.URL+"/api/"))
	router := newRouterWithPixabay(store, 5*time.Minute, client)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/pixabay/search?q=flower&lobby_id="+created.ID, nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	var resp searchImagesResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Total != 10 || resp.TotalHits != 5 {
		t.Fatalf("total = %d/%d, want 10/5", resp.Total, resp.TotalHits)
	}
	if len(resp.Hits) != 1 {
		t.Fatalf("len(Hits) = %d, want 1", len(resp.Hits))
	}
	if resp.Hits[0].PixabayID != 42 {
		t.Fatalf("PixabayID = %d, want 42", resp.Hits[0].PixabayID)
	}
	if resp.RateLimit.Limit != 100 || resp.RateLimit.Remaining != 99 {
		t.Fatalf("RateLimit = %+v", resp.RateLimit)
	}
}

func TestPixabaySearchHandlerUnauthorized(t *testing.T) {
	t.Parallel()

	router := newRouterWithPixabay(lobby.NewMemoryStore(), 5*time.Minute, pixabay.NewClient("key"))

	req := httptest.NewRequest(http.MethodGet, "/api/pixabay/search?q=flower", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestPixabaySearchHandlerForbidden(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newRouterWithPixabay(store, 5*time.Minute, pixabay.NewClient("key"))

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/pixabay/search?q=flower&lobby_id="+created.ID, nil)
	req.Header.Set(lobby.AdminTokenHeader, "wrong-token")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestPixabaySearchHandlerMissingQuery(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	router := newRouterWithPixabay(store, 5*time.Minute, pixabay.NewClient("key"))

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/pixabay/search?lobby_id="+created.ID, nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestPixabaySearchHandlerRateLimited(t *testing.T) {
	t.Parallel()

	pixabayServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("API rate limit exceeded"))
	}))
	t.Cleanup(pixabayServer.Close)

	store := lobby.NewMemoryStore()
	client := pixabay.NewClient("test-key", pixabay.WithBaseURL(pixabayServer.URL+"/api/"))
	router := newRouterWithPixabay(store, 5*time.Minute, client)

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/pixabay/search?q=cat&lobby_id="+created.ID, nil)
	req.Header.Set(lobby.AdminTokenHeader, created.AdminToken)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
}

func newRouterWithPixabay(store lobby.Store, drawDuration time.Duration, client *pixabay.Client) *gin.Engine {
	return newTestRouter(store, drawDuration, client, nil, nil)
}
