package lobby

import "context"

func (s *MemoryStore) withLobby(ctx context.Context, id string, fn func(*Lobby) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	lobby, ok := s.lobbies[id]
	if !ok {
		return ErrNotFound
	}

	return fn(lobby)
}
