package lobby

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestLobbyPhaseTransitions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		from LobbyPhase
		to   LobbyPhase
		want bool
	}{
		{PhaseWaiting, PhaseSelecting, true},
		{PhaseWaiting, PhaseDrawing, false},
		{PhaseSelecting, PhaseReady, true},
		{PhaseReady, PhaseDrawing, true},
		{PhaseDrawing, PhaseBetweenRounds, true},
		{PhaseBetweenRounds, PhaseDrawing, true},
		{PhaseBetweenRounds, PhaseFinished, true},
		{PhaseFinished, PhaseWaiting, false},
		{PhaseDrawing, PhaseFinished, true},
		{PhaseWaiting, PhaseWaiting, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"_to_"+string(tt.to), func(t *testing.T) {
			t.Parallel()
			if got := tt.from.CanTransitionTo(tt.to); got != tt.want {
				t.Fatalf("CanTransitionTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateTransition(t *testing.T) {
	t.Parallel()

	if err := ValidateTransition(PhaseWaiting, PhaseSelecting); err != nil {
		t.Fatalf("ValidateTransition() error = %v", err)
	}
	if err := ValidateTransition(PhaseWaiting, PhaseDrawing); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("ValidateTransition() error = %v, want ErrInvalidTransition", err)
	}
}

func TestLobbySnapshotMasking(t *testing.T) {
	t.Parallel()

	endsAt := time.Now().Add(5 * time.Minute)
	lobby := &Lobby{
		ID:    "lobby-1",
		Phase: PhaseDrawing,
		SelectedPhotos: []Photo{{
			PixabayID:     1,
			PreviewURL:    "preview",
			LargeImageURL: "large",
			PageURL:       "page",
			Width:         100,
			Height:        200,
		}},
		PhotoOrder:   []int{0},
		CurrentRound: 0,
		DrawEndsAt:   &endsAt,
	}

	snap := lobby.Snapshot(3, time.Now())
	if snap.ParticipantCount != 3 {
		t.Fatalf("ParticipantCount = %d, want 3", snap.ParticipantCount)
	}
	if snap.CurrentRound != 1 {
		t.Fatalf("CurrentRound = %d, want 1", snap.CurrentRound)
	}
	if snap.CurrentPhoto == nil {
		t.Fatal("CurrentPhoto = nil, want photo")
	}
	if snap.DrawEndsAt == nil {
		t.Fatal("DrawEndsAt = nil, want timestamp")
	}

	lobby.Phase = PhaseReady
	readySnap := lobby.Snapshot(0, time.Now())
	if readySnap.CurrentPhoto != nil {
		t.Fatal("CurrentPhoto should be hidden in READY")
	}
	if readySnap.DrawEndsAt != nil {
		t.Fatal("DrawEndsAt should be hidden in READY")
	}
	if readySnap.CurrentRound != 0 {
		t.Fatalf("CurrentRound = %d, want 0 in READY", readySnap.CurrentRound)
	}

	lobby.Phase = PhaseBetweenRounds
	betweenSnap := lobby.Snapshot(0, time.Now())
	if betweenSnap.CurrentPhoto != nil {
		t.Fatal("CurrentPhoto should be hidden in BETWEEN_ROUNDS")
	}
	if betweenSnap.CurrentRound != 1 {
		t.Fatalf("CurrentRound = %d, want 1 in BETWEEN_ROUNDS", betweenSnap.CurrentRound)
	}
}

func TestMemoryStoreCreateAndGet(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID == "" || created.AdminToken == "" {
		t.Fatal("Create() returned empty id or admin token")
	}
	if created.Phase != PhaseWaiting {
		t.Fatalf("Phase = %q, want WAITING", created.Phase)
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("ID = %q, want %q", got.ID, created.ID)
	}

	_, err = store.Get(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get() error = %v, want ErrNotFound", err)
	}
}

func TestValidateAdminToken(t *testing.T) {
	t.Parallel()

	lobby := &Lobby{AdminToken: "secret-token"}
	if !ValidateAdminToken(lobby, "secret-token") {
		t.Fatal("ValidateAdminToken() = false, want true")
	}
	if ValidateAdminToken(lobby, "wrong-token") {
		t.Fatal("ValidateAdminToken() = true, want false")
	}
}

func TestMemoryStoreSetSelectedPhotos(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	photos := []Photo{{
		PixabayID:     1,
		LargeImageURL: "https://cdn.example/large.jpg",
	}}

	if err := store.SetSelectedPhotos(context.Background(), created.ID, photos); err != nil {
		t.Fatalf("SetSelectedPhotos() error = %v", err)
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Phase != PhaseSelecting {
		t.Fatalf("Phase = %q, want SELECTING", got.Phase)
	}
	if len(got.SelectedPhotos) != 1 {
		t.Fatalf("len(SelectedPhotos) = %d, want 1", len(got.SelectedPhotos))
	}

	if err := store.SetSelectedPhotos(context.Background(), created.ID, nil); !errors.Is(err, ErrEmptyPhotos) {
		t.Fatalf("SetSelectedPhotos(empty) error = %v, want ErrEmptyPhotos", err)
	}
}

func TestMemoryStoreMarkReady(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	photos := []Photo{
		{PixabayID: 1, LargeImageURL: "https://cdn.example/1.jpg"},
		{PixabayID: 2, LargeImageURL: "https://cdn.example/2.jpg"},
		{PixabayID: 3, LargeImageURL: "https://cdn.example/3.jpg"},
	}
	if err := store.SetSelectedPhotos(context.Background(), created.ID, photos); err != nil {
		t.Fatalf("SetSelectedPhotos() error = %v", err)
	}

	if err := store.MarkReady(context.Background(), created.ID); err != nil {
		t.Fatalf("MarkReady() error = %v", err)
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Phase != PhaseReady {
		t.Fatalf("Phase = %q, want READY", got.Phase)
	}
	if len(got.PhotoOrder) != len(photos) {
		t.Fatalf("len(PhotoOrder) = %d, want %d", len(got.PhotoOrder), len(photos))
	}

	seen := make([]bool, len(photos))
	for _, idx := range got.PhotoOrder {
		if idx < 0 || idx >= len(photos) {
			t.Fatalf("PhotoOrder index out of range: %d", idx)
		}
		if seen[idx] {
			t.Fatalf("duplicate PhotoOrder index: %d", idx)
		}
		seen[idx] = true
	}

	if err := store.MarkReady(context.Background(), created.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("MarkReady() from READY error = %v, want ErrInvalidTransition", err)
	}
}

func TestDrawTimerHelpers(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)
	lob := &Lobby{DrawDuration: 5 * time.Minute}

	StartDrawTimer(lob, now)
	if lob.DrawEndsAt == nil {
		t.Fatal("DrawEndsAt = nil, want timestamp")
	}
	wantEndsAt := now.Add(RoundCountdown + 5*time.Minute)
	if !lob.DrawEndsAt.Equal(wantEndsAt) {
		t.Fatalf("DrawEndsAt = %v, want %v", lob.DrawEndsAt, wantEndsAt)
	}

	if IsDrawExpired(lob, wantEndsAt.Add(-time.Second)) {
		t.Fatal("IsDrawExpired() = true before deadline")
	}
	if !IsDrawExpired(lob, wantEndsAt) {
		t.Fatal("IsDrawExpired() = false at deadline")
	}

	ClearDrawTimer(lob)
	if lob.DrawEndsAt != nil {
		t.Fatalf("DrawEndsAt = %v, want nil", lob.DrawEndsAt)
	}
}

func TestMemoryStoreStartSession(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	created, photos := createLobbyReadyForTest(t, store)
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)

	if err := store.StartSession(context.Background(), created.ID, now); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Phase != PhaseDrawing {
		t.Fatalf("Phase = %q, want DRAWING", got.Phase)
	}
	if got.CurrentRound != 0 {
		t.Fatalf("CurrentRound = %d, want 0", got.CurrentRound)
	}
	if got.DrawEndsAt == nil {
		t.Fatal("DrawEndsAt = nil, want timestamp")
	}
	wantEndsAt := now.Add(RoundCountdown + got.DrawDuration)
	if !got.DrawEndsAt.Equal(wantEndsAt) {
		t.Fatalf("DrawEndsAt = %v, want %v", got.DrawEndsAt, wantEndsAt)
	}

	snap := got.Snapshot(0, now)
	if snap.CurrentRound != 1 {
		t.Fatalf("CurrentRound = %d, want 1", snap.CurrentRound)
	}
	if snap.CurrentPhoto == nil {
		t.Fatal("CurrentPhoto = nil, want photo in DRAWING")
	}
	if snap.TotalRounds != len(photos) {
		t.Fatalf("TotalRounds = %d, want %d", snap.TotalRounds, len(photos))
	}
}

func TestMemoryStoreNextRound(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	created, photos := createLobbyReadyForTest(t, store)
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)

	if err := store.StartSession(context.Background(), created.ID, now); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	if err := store.AdvanceToBetweenRounds(context.Background(), created.ID); err != nil {
		t.Fatalf("AdvanceToBetweenRounds() error = %v", err)
	}

	if err := store.NextRound(context.Background(), created.ID, now); err != nil {
		t.Fatalf("NextRound() error = %v", err)
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Phase != PhaseDrawing {
		t.Fatalf("Phase = %q, want DRAWING", got.Phase)
	}
	if got.CurrentRound != 1 {
		t.Fatalf("CurrentRound = %d, want 1", got.CurrentRound)
	}

	if err := store.AdvanceToBetweenRounds(context.Background(), created.ID); err != nil {
		t.Fatalf("AdvanceToBetweenRounds() error = %v", err)
	}

	if err := store.NextRound(context.Background(), created.ID, now); err != nil {
		t.Fatalf("NextRound(last) error = %v", err)
	}

	got, err = store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Phase != PhaseFinished {
		t.Fatalf("Phase = %q, want FINISHED", got.Phase)
	}
	if got.CurrentRound != len(photos)-1 {
		t.Fatalf("CurrentRound = %d, want %d", got.CurrentRound, len(photos)-1)
	}
}

func TestMemoryStoreFinishSession(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	created, _ := createLobbyReadyForTest(t, store)
	now := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)

	if err := store.StartSession(context.Background(), created.ID, now); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	if err := store.FinishSession(context.Background(), created.ID); err != nil {
		t.Fatalf("FinishSession() error = %v", err)
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Phase != PhaseFinished {
		t.Fatalf("Phase = %q, want FINISHED", got.Phase)
	}
	if got.DrawEndsAt != nil {
		t.Fatalf("DrawEndsAt = %v, want nil", got.DrawEndsAt)
	}

	if err := store.FinishSession(context.Background(), created.ID); !errors.Is(err, ErrInvalidTransition) {
		t.Fatalf("FinishSession() again error = %v, want ErrInvalidTransition", err)
	}
}

func createLobbyReadyForTest(t *testing.T, store *MemoryStore) (*Lobby, []Photo) {
	t.Helper()

	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	photos := []Photo{
		{PixabayID: 1, LargeImageURL: "https://cdn.example/1.jpg"},
		{PixabayID: 2, LargeImageURL: "https://cdn.example/2.jpg"},
	}
	if err := store.SetSelectedPhotos(context.Background(), created.ID, photos); err != nil {
		t.Fatalf("SetSelectedPhotos() error = %v", err)
	}
	if err := store.MarkReady(context.Background(), created.ID); err != nil {
		t.Fatalf("MarkReady() error = %v", err)
	}

	return created, photos
}

func TestMemoryStoreExpireDrawingTimers(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	created, _ := createLobbyReadyForTest(t, store)
	start := time.Date(2026, 6, 20, 12, 0, 0, 0, time.UTC)

	if err := store.StartSession(context.Background(), created.ID, start); err != nil {
		t.Fatalf("StartSession() error = %v", err)
	}

	ids, err := store.ExpireDrawingTimers(context.Background(), start.Add(4*time.Minute+59*time.Second))
	if err != nil {
		t.Fatalf("ExpireDrawingTimers() error = %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("ExpireDrawingTimers() = %v, want none", ids)
	}

	ids, err = store.ExpireDrawingTimers(context.Background(), start.Add(RoundCountdown+5*time.Minute-time.Second))
	if err != nil {
		t.Fatalf("ExpireDrawingTimers() error = %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("ExpireDrawingTimers() before deadline = %v, want none", ids)
	}

	ids, err = store.ExpireDrawingTimers(context.Background(), start.Add(RoundCountdown+5*time.Minute))
	if err != nil {
		t.Fatalf("ExpireDrawingTimers() error = %v", err)
	}
	if len(ids) != 1 || ids[0] != created.ID {
		t.Fatalf("ExpireDrawingTimers() = %v, want [%q]", ids, created.ID)
	}

	got, err := store.Get(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if got.Phase != PhaseBetweenRounds {
		t.Fatalf("Phase = %q, want BETWEEN_ROUNDS", got.Phase)
	}

	ids, err = store.ExpireDrawingTimers(context.Background(), start.Add(RoundCountdown+5*time.Minute))
	if err != nil {
		t.Fatalf("ExpireDrawingTimers() again error = %v", err)
	}
	if len(ids) != 0 {
		t.Fatalf("ExpireDrawingTimers() again = %v, want none", ids)
	}
}

func TestMemoryStoreMarkReadyRequiresPhotos(t *testing.T) {
	t.Parallel()

	store := NewMemoryStore()
	created, err := store.Create(context.Background(), 5*time.Minute)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	store.mu.Lock()
	store.lobbies[created.ID].Phase = PhaseSelecting
	store.mu.Unlock()

	if err := store.MarkReady(context.Background(), created.ID); !errors.Is(err, ErrEmptyPhotos) {
		t.Fatalf("MarkReady() error = %v, want ErrEmptyPhotos", err)
	}
}
