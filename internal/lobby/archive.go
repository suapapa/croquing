package lobby

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// ImageFetcher downloads image bytes for a URL.
type ImageFetcher func(ctx context.Context, url string) ([]byte, error)

// DefaultImageFetcher returns an ImageFetcher backed by client.
func DefaultImageFetcher(client *http.Client) ImageFetcher {
	if client == nil {
		client = http.DefaultClient
	}

	return func(ctx context.Context, url string) ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("lobby: image download status %d", resp.StatusCode)
		}

		const maxImageBytes = 20 << 20
		data, err := io.ReadAll(io.LimitReader(resp.Body, maxImageBytes+1))
		if err != nil {
			return nil, err
		}
		if len(data) > maxImageBytes {
			return nil, fmt.Errorf("lobby: image exceeds size limit")
		}

		return data, nil
	}
}

// ArchiveBaseName returns the dated prefix used for zip and entry names.
func ArchiveBaseName(now time.Time, sessionID string) string {
	return "croquing_" + now.UTC().Format("20060102") + "_" + sessionID
}

// ArchiveEntryName returns a numbered file name inside the archive.
func ArchiveEntryName(baseName string, round int, total int) string {
	width := len(strconv.Itoa(total))
	if width < 2 {
		width = 2
	}

	return fmt.Sprintf("%s_%0*d.jpeg", baseName, width, round)
}

// OrderedPhotos returns session photos in drawing order.
func (l *Lobby) OrderedPhotos() ([]Photo, error) {
	if l == nil || len(l.PhotoOrder) == 0 {
		return nil, ErrPhotosNotReady
	}

	photos := make([]Photo, 0, len(l.PhotoOrder))
	for _, idx := range l.PhotoOrder {
		if idx < 0 || idx >= len(l.SelectedPhotos) {
			return nil, ErrPhotosNotReady
		}
		photos = append(photos, l.SelectedPhotos[idx])
	}

	return photos, nil
}

// BuildPhotosZIP downloads photos and returns a zip archive.
func BuildPhotosZIP(ctx context.Context, photos []Photo, baseName string, fetch ImageFetcher) ([]byte, error) {
	if len(photos) == 0 {
		return nil, ErrEmptyPhotos
	}
	if fetch == nil {
		return nil, fmt.Errorf("lobby: image fetcher is required")
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for i, photo := range photos {
		data, err := fetch(ctx, photo.LargeImageURL)
		if err != nil {
			return nil, fmt.Errorf("lobby: fetch photo %d: %w", i+1, err)
		}

		entryName := ArchiveEntryName(baseName, i+1, len(photos))
		entry, err := zipWriter.Create(entryName)
		if err != nil {
			return nil, err
		}
		if _, err := entry.Write(data); err != nil {
			return nil, err
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
