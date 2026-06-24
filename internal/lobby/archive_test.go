package lobby

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"testing"
	"time"
)

func TestArchiveBaseName(t *testing.T) {
	t.Parallel()

	got := ArchiveBaseName(time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC), "abc-123")
	want := "croquing_20260621_abc-123"
	if got != want {
		t.Fatalf("ArchiveBaseName() = %q, want %q", got, want)
	}
}

func TestArchiveEntryName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		baseName string
		round    int
		total    int
		want     string
	}{
		{
			name:     "two digit padding",
			baseName: "croquing_20260621_abc-123",
			round:    3,
			total:    12,
			want:     "croquing_20260621_abc-123_03.jpeg",
		},
		{
			name:     "minimum width",
			baseName: "croquing_20260621_abc-123",
			round:    1,
			total:    5,
			want:     "croquing_20260621_abc-123_01.jpeg",
		},
		{
			name:     "wider padding for large sets",
			baseName: "croquing_20260621_abc-123",
			round:    42,
			total:    120,
			want:     "croquing_20260621_abc-123_042.jpeg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := ArchiveEntryName(tt.baseName, tt.round, tt.total); got != tt.want {
				t.Fatalf("ArchiveEntryName() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestOrderedPhotos(t *testing.T) {
	t.Parallel()

	lobby := &Lobby{
		SelectedPhotos: []Photo{
			{PixabayID: 10, LargeImageURL: "https://example.com/a.jpg"},
			{PixabayID: 20, LargeImageURL: "https://example.com/b.jpg"},
		},
		PhotoOrder: []int{1, 0},
	}

	photos, err := lobby.OrderedPhotos()
	if err != nil {
		t.Fatalf("OrderedPhotos() error = %v", err)
	}
	if len(photos) != 2 {
		t.Fatalf("len(photos) = %d, want 2", len(photos))
	}
	if photos[0].PixabayID != 20 || photos[1].PixabayID != 10 {
		t.Fatalf("photos order = %+v, want pixabay ids 20 then 10", photos)
	}
}

func TestOrderedPhotosNotReady(t *testing.T) {
	t.Parallel()

	_, err := (&Lobby{}).OrderedPhotos()
	if !errors.Is(err, ErrPhotosNotReady) {
		t.Fatalf("OrderedPhotos() error = %v, want ErrPhotosNotReady", err)
	}
}

func TestBuildPhotosZIP(t *testing.T) {
	t.Parallel()

	photos := []Photo{
		{LargeImageURL: "https://example.com/1.jpg"},
		{LargeImageURL: "https://example.com/2.jpg"},
	}
	baseName := ArchiveBaseName(time.Date(2026, 6, 21, 12, 0, 0, 0, time.UTC), "session-42")

	fetch := func(_ context.Context, url string) ([]byte, error) {
		switch url {
		case "https://example.com/1.jpg":
			return []byte("image-one"), nil
		case "https://example.com/2.jpg":
			return []byte("image-two"), nil
		default:
			return nil, errors.New("unexpected url")
		}
	}

	data, err := BuildPhotosZIP(context.Background(), photos, baseName, fetch)
	if err != nil {
		t.Fatalf("BuildPhotosZIP() error = %v", err)
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		t.Fatalf("zip.NewReader() error = %v", err)
	}
	if len(reader.File) != 2 {
		t.Fatalf("len(zip files) = %d, want 2", len(reader.File))
	}

	wantNames := []string{
		ArchiveEntryName(baseName, 1, 2),
		ArchiveEntryName(baseName, 2, 2),
	}
	for i, file := range reader.File {
		if file.Name != wantNames[i] {
			t.Fatalf("file[%d].Name = %q, want %q", i, file.Name, wantNames[i])
		}

		rc, err := file.Open()
		if err != nil {
			t.Fatalf("Open() error = %v", err)
		}
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(rc); err != nil {
			rc.Close()
			t.Fatalf("ReadFrom() error = %v", err)
		}
		rc.Close()

		switch i {
		case 0:
			if buf.String() != "image-one" {
				t.Fatalf("first file content = %q, want image-one", buf.String())
			}
		case 1:
			if buf.String() != "image-two" {
				t.Fatalf("second file content = %q, want image-two", buf.String())
			}
		}
	}
}
