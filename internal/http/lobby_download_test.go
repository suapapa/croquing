package httpserver

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
	"github.com/suapapa/croquis-king/internal/ws"
)

func TestDownloadPhotosHandler(t *testing.T) {
	t.Parallel()

	imgServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "image/jpeg")
		_, _ = w.Write([]byte("jpeg-bytes"))
	}))
	t.Cleanup(imgServer.Close)

	store := lobby.NewMemoryStore()
	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	photos := []lobby.Photo{
		{PixabayID: 1, LargeImageURL: imgServer.URL + "/1"},
		{PixabayID: 2, LargeImageURL: imgServer.URL + "/2"},
	}
	if err := store.SetSelectedPhotos(context.Background(), created.ID, photos); err != nil {
		t.Fatalf("SetSelectedPhotos() error = %v", err)
	}
	if err := store.MarkReady(context.Background(), created.ID); err != nil {
		t.Fatalf("MarkReady() error = %v", err)
	}
	if err := store.StartSession(context.Background(), created.ID, time.Now()); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if err := store.FinishSession(context.Background(), created.ID); err != nil {
		t.Fatalf("FinishSession() error = %v", err)
	}

	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))
	req := httptest.NewRequest(http.MethodGet, "/api/lobbies/"+created.ID+"/photos/download", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Type"); got != "application/zip" {
		t.Fatalf("Content-Type = %q, want application/zip", got)
	}
	if got := rec.Header().Get("Content-Disposition"); got == "" {
		t.Fatal("Content-Disposition header is empty")
	}

	reader, err := zip.NewReader(bytes.NewReader(rec.Body.Bytes()), int64(rec.Body.Len()))
	if err != nil {
		t.Fatalf("zip.NewReader() error = %v", err)
	}
	if len(reader.File) != 2 {
		t.Fatalf("len(zip files) = %d, want 2", len(reader.File))
	}

	for _, file := range reader.File {
		rc, err := file.Open()
		if err != nil {
			t.Fatalf("Open(%q) error = %v", file.Name, err)
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			t.Fatalf("ReadAll(%q) error = %v", file.Name, err)
		}
		if string(data) != "jpeg-bytes" {
			t.Fatalf("file %q content = %q, want jpeg-bytes", file.Name, string(data))
		}
	}
}

func TestDownloadPhotosHandlerRequiresFinishedPhase(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	created, _ := createReadyLobby(t, store)

	router := newTestRouter(store, 5*time.Minute, pixabay.NewClient("test-key"), nil, ws.NewSnapshotSync(store))
	req := httptest.NewRequest(http.MethodGet, "/api/lobbies/"+created.ID+"/photos/download", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusConflict, rec.Body.String())
	}
}

func TestDownloadPhotosHandlerNotFound(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	handler := newLobbyHandler(store, 5*time.Minute, nil)

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/lobbies/missing/photos/download", nil)
	c.Params = gin.Params{{Key: "id", Value: "missing"}}
	handler.downloadPhotos(c)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}

func TestDownloadPhotosHandlerUsesMockFetcher(t *testing.T) {
	t.Parallel()

	store := lobby.NewMemoryStore()
	created, _ := createReadyLobby(t, store)

	if err := store.StartSession(context.Background(), created.ID, time.Now()); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}
	if err := store.FinishSession(context.Background(), created.ID); err != nil {
		t.Fatalf("FinishSession() error = %v", err)
	}

	handler := newLobbyHandler(store, 5*time.Minute, nil)
	handler.imageFetcher = func(_ context.Context, _ string) ([]byte, error) {
		return []byte("jpeg-bytes"), nil
	}

	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/lobbies/"+created.ID+"/photos/download", nil)
	c.Params = gin.Params{{Key: "id", Value: created.ID}}
	handler.downloadPhotos(c)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}
