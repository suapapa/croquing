package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPublicConfigHandler(t *testing.T) {
	t.Parallel()

	router := newRouter(nil, 0, nil, nil, nil, nil, "My Studio")

	req := httptest.NewRequest(http.MethodGet, "/api/config", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body publicConfigResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.AppName != "My Studio" {
		t.Fatalf("app_name = %q, want %q", body.AppName, "My Studio")
	}
}
