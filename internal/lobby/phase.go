package lobby

// LobbyPhase represents the current session state of a lobby.
type LobbyPhase string

const (
	PhaseWaiting       LobbyPhase = "WAITING"
	PhaseSelecting     LobbyPhase = "SELECTING"
	PhaseReady         LobbyPhase = "READY"
	PhaseDrawing       LobbyPhase = "DRAWING"
	PhaseBetweenRounds LobbyPhase = "BETWEEN_ROUNDS"
	PhaseFinished      LobbyPhase = "FINISHED"
)

var validTransitions = map[LobbyPhase][]LobbyPhase{
	PhaseWaiting:       {PhaseSelecting, PhaseFinished},
	PhaseSelecting:     {PhaseReady, PhaseFinished},
	PhaseReady:         {PhaseDrawing, PhaseFinished},
	PhaseDrawing:       {PhaseBetweenRounds, PhaseFinished},
	PhaseBetweenRounds: {PhaseDrawing, PhaseFinished},
	PhaseFinished:      {},
}

// CanTransitionTo reports whether the lobby may move from the current phase to next.
func (p LobbyPhase) CanTransitionTo(next LobbyPhase) bool {
	if p == next {
		return true
	}

	allowed, ok := validTransitions[p]
	if !ok {
		return false
	}

	for _, candidate := range allowed {
		if candidate == next {
			return true
		}
	}
	return false
}

// ValidateTransition returns ErrInvalidTransition when next is not allowed from current.
func ValidateTransition(current, next LobbyPhase) error {
	if current.CanTransitionTo(next) {
		return nil
	}
	return ErrInvalidTransition
}
