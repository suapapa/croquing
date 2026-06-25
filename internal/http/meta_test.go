package httpserver

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRenderSocialMetaTags(t *testing.T) {
	t.Parallel()

	meta := socialMeta{
		Title:       `Croquing "Live"`,
		Description: "Draw & share",
		Image:       "https://cdn.example/logo.png",
		URL:         "https://croquing.example/",
	}

	got := renderSocialMetaTags(meta)
	if !strings.Contains(got, `content="Croquing &#34;Live&#34;"`) {
		t.Fatalf("title not escaped: %q", got)
	}
	if !strings.Contains(got, `property="og:image" content="https://cdn.example/logo.png"`) {
		t.Fatalf("og:image missing: %q", got)
	}
	if !strings.Contains(got, `name="twitter:card" content="summary_large_image"`) {
		t.Fatalf("twitter:card missing: %q", got)
	}
}

func TestAbsoluteAssetURL(t *testing.T) {
	t.Parallel()

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "https://croquing.example/lobby/abc", nil)

	if got := absoluteAssetURL(c, "https://cdn.example/logo.png"); got != "https://cdn.example/logo.png" {
		t.Fatalf("absolute logo = %q", got)
	}
	if got := absoluteAssetURL(c, "/logo.png"); got != "https://croquing.example/logo.png" {
		t.Fatalf("relative logo = %q", got)
	}
}

func TestServeSPAInjectsSocialMeta(t *testing.T) {
	t.Parallel()

	staticDir := t.TempDir()
	indexHTML := `<!doctype html><html><head><title>Croquing</title></head><body></body></html>`
	if err := os.WriteFile(filepath.Join(staticDir, "index.html"), []byte(indexHTML), 0o644); err != nil {
		t.Fatalf("write index.html: %v", err)
	}
	if err := os.WriteFile(filepath.Join(staticDir, "logo.png"), []byte("png"), 0o644); err != nil {
		t.Fatalf("write logo.png: %v", err)
	}

	r := gin.New()
	if !registerStaticRoutes(r, staticDir, "My Croquing", "/logo.png") {
		t.Fatal("registerStaticRoutes returned false")
	}

	req := httptest.NewRequest(http.MethodGet, "https://croquing.example/", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	body := rec.Body.String()
	if !strings.Contains(body, "<title>My Croquing</title>") {
		t.Fatalf("title not updated: %q", body)
	}
	if !strings.Contains(body, `property="og:image" content="https://croquing.example/logo.png"`) {
		t.Fatalf("og:image not injected: %q", body)
	}
	if !strings.Contains(body, `property="og:title" content="My Croquing"`) {
		t.Fatalf("og:title not injected: %q", body)
	}
}

func TestInjectBeforeHeadClose(t *testing.T) {
	t.Parallel()

	html := "<html><head><meta charset=\"utf-8\"></head><body></body></html>"
	got := injectBeforeHeadClose(html, "<meta property=\"og:image\" content=\"x\" />")
	want := "<html><head><meta charset=\"utf-8\"><meta property=\"og:image\" content=\"x\" /></head><body></body></html>"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
