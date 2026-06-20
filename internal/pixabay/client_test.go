package pixabay

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClientSearchSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/" {
			t.Fatalf("path = %q, want /api/", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("key") != "test-key" {
			t.Fatalf("key = %q, want test-key", query.Get("key"))
		}
		if query.Get("q") != "flower" {
			t.Fatalf("q = %q, want flower", query.Get("q"))
		}
		if query.Get("order") != "popular" {
			t.Fatalf("order = %q, want popular", query.Get("order"))
		}
		if query.Get("page") != "1" {
			t.Fatalf("page = %q, want 1", query.Get("page"))
		}
		if query.Get("per_page") != "20" {
			t.Fatalf("per_page = %q, want 20", query.Get("per_page"))
		}
		if query.Get("image_type") != "photo" {
			t.Fatalf("image_type = %q, want photo", query.Get("image_type"))
		}
		if query.Get("safesearch") != "true" {
			t.Fatalf("safesearch = %q, want true", query.Get("safesearch"))
		}

		w.Header().Set("X-RateLimit-Limit", "100")
		w.Header().Set("X-RateLimit-Remaining", "99")
		w.Header().Set("X-RateLimit-Reset", "42")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"total": 100,
			"totalHits": 50,
			"hits": [{
				"id": 123,
				"pageURL": "https://pixabay.com/photos/flower-123/",
				"previewURL": "https://cdn.pixabay.com/preview.jpg",
				"webformatURL": "https://cdn.pixabay.com/web.jpg",
				"largeImageURL": "https://cdn.pixabay.com/large.jpg",
				"imageWidth": 1920,
				"imageHeight": 1080,
				"views": 10,
				"downloads": 5,
				"likes": 2
			}]
		}`))
	}))
	t.Cleanup(server.Close)

	client := NewClient("test-key", WithBaseURL(server.URL+"/api/"))
	result, err := client.Search(context.Background(), SearchParams{
		Query:      "flower",
		SafeSearch: true,
	})
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}

	if result.Total != 100 {
		t.Fatalf("Total = %d, want 100", result.Total)
	}
	if result.TotalHits != 50 {
		t.Fatalf("TotalHits = %d, want 50", result.TotalHits)
	}
	if len(result.Hits) != 1 {
		t.Fatalf("len(Hits) = %d, want 1", len(result.Hits))
	}

	hit := result.Hits[0]
	if hit.ID != 123 {
		t.Fatalf("ID = %d, want 123", hit.ID)
	}
	if hit.LargeImageURL != "https://cdn.pixabay.com/large.jpg" {
		t.Fatalf("LargeImageURL = %q", hit.LargeImageURL)
	}
	if hit.Width != 1920 || hit.Height != 1080 {
		t.Fatalf("size = %dx%d, want 1920x1080", hit.Width, hit.Height)
	}

	if result.RateLimit.Limit != 100 {
		t.Fatalf("RateLimit.Limit = %d, want 100", result.RateLimit.Limit)
	}
	if result.RateLimit.Remaining != 99 {
		t.Fatalf("RateLimit.Remaining = %d, want 99", result.RateLimit.Remaining)
	}
	if result.RateLimit.Reset != 42 {
		t.Fatalf("RateLimit.Reset = %d, want 42", result.RateLimit.Reset)
	}
}

func TestClientSearchEmptyQuery(t *testing.T) {
	t.Parallel()

	client := NewClient("test-key")
	_, err := client.Search(context.Background(), SearchParams{})
	if !errors.Is(err, ErrEmptyQuery) {
		t.Fatalf("error = %v, want ErrEmptyQuery", err)
	}
}

func TestClientSearchRateLimited(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-RateLimit-Remaining", "0")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("API rate limit exceeded"))
	}))
	t.Cleanup(server.Close)

	client := NewClient("test-key", WithBaseURL(server.URL+"/api/"))
	result, err := client.Search(context.Background(), SearchParams{Query: "cat"})
	if !errors.Is(err, ErrRateLimited) {
		t.Fatalf("error = %v, want ErrRateLimited", err)
	}
	if !strings.Contains(err.Error(), "API rate limit exceeded") {
		t.Fatalf("error = %q, want rate limit message", err.Error())
	}
	if result.RateLimit.Remaining != 0 {
		t.Fatalf("RateLimit.Remaining = %d, want 0", result.RateLimit.Remaining)
	}
}

func TestClientSearchAPIError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid API key"))
	}))
	t.Cleanup(server.Close)

	client := NewClient("bad-key", WithBaseURL(server.URL+"/api/"))
	_, err := client.Search(context.Background(), SearchParams{Query: "cat"})
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("error = %T(%v), want *APIError", err, err)
	}
	if apiErr.StatusCode != http.StatusBadRequest {
		t.Fatalf("StatusCode = %d, want 400", apiErr.StatusCode)
	}
	if apiErr.Message != "Invalid API key" {
		t.Fatalf("Message = %q, want Invalid API key", apiErr.Message)
	}
}

func TestNormalizeSearchParams(t *testing.T) {
	t.Parallel()

	params := normalizeSearchParams(SearchParams{
		Query:   "test",
		PerPage: 500,
	})

	if params.Order != defaultOrder {
		t.Fatalf("Order = %q, want %q", params.Order, defaultOrder)
	}
	if params.Page != defaultPage {
		t.Fatalf("Page = %d, want %d", params.Page, defaultPage)
	}
	if params.PerPage != maxPerPage {
		t.Fatalf("PerPage = %d, want %d", params.PerPage, maxPerPage)
	}
	if params.ImageType != defaultImageType {
		t.Fatalf("ImageType = %q, want %q", params.ImageType, defaultImageType)
	}
}
