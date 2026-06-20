package lobby

import "time"

// Photo is a selected Pixabay image stored in a lobby.
type Photo struct {
	PixabayID     int    `json:"pixabay_id"`
	PreviewURL    string `json:"preview_url"`
	LargeImageURL string `json:"large_image_url"`
	PageURL       string `json:"page_url"`
	Width         int    `json:"width"`
	Height        int    `json:"height"`
}

// Lobby holds the full server-side lobby state.
type Lobby struct {
	ID             string
	AdminToken     string
	Phase          LobbyPhase
	SelectedPhotos []Photo
	PhotoOrder     []int
	CurrentRound   int
	DrawDuration   time.Duration
	DrawEndsAt     *time.Time
	CreatedAt      time.Time
}

// LobbySnapshot is the public view sent to clients.
type LobbySnapshot struct {
	ID               string     `json:"id"`
	Phase            LobbyPhase `json:"phase"`
	ParticipantCount int        `json:"participant_count"`
	SelectedCount    int        `json:"selected_count"`
	CurrentRound     int        `json:"current_round"`
	TotalRounds      int        `json:"total_rounds"`
	DrawEndsAt       *time.Time `json:"draw_ends_at,omitempty"`
	CurrentPhoto     *Photo     `json:"current_photo,omitempty"`
	ServerTime       time.Time  `json:"server_time"`
}

// CurrentPhoto returns the photo shown during DRAWING, if any.
func (l *Lobby) CurrentPhoto() *Photo {
	if l == nil || l.Phase != PhaseDrawing {
		return nil
	}
	if l.CurrentRound < 0 || l.CurrentRound >= len(l.PhotoOrder) {
		return nil
	}

	idx := l.PhotoOrder[l.CurrentRound]
	if idx < 0 || idx >= len(l.SelectedPhotos) {
		return nil
	}

	photo := l.SelectedPhotos[idx]
	return &photo
}

// Snapshot builds the client-facing snapshot for the lobby.
func (l *Lobby) Snapshot(participantCount int, now time.Time) LobbySnapshot {
	snap := LobbySnapshot{
		ID:               l.ID,
		Phase:            l.Phase,
		ParticipantCount: participantCount,
		SelectedCount:    len(l.SelectedPhotos),
		TotalRounds:      len(l.PhotoOrder),
		ServerTime:       now,
	}

	if l.Phase == PhaseDrawing || l.Phase == PhaseBetweenRounds || l.Phase == PhaseFinished {
		snap.CurrentRound = l.CurrentRound + 1
	}

	if l.Phase == PhaseDrawing {
		snap.DrawEndsAt = l.DrawEndsAt
		snap.CurrentPhoto = l.CurrentPhoto()
	}

	return snap
}
