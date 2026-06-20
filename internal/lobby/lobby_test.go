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
