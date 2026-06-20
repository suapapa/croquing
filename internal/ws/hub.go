package ws

import (
	"sync"
)

// Hub tracks WebSocket clients grouped by lobby ID.
type Hub struct {
	mu      sync.RWMutex
	lobbies map[string]map[*Client]struct{}
}

// NewHub creates an empty WebSocket hub.
func NewHub() *Hub {
	return &Hub{
		lobbies: make(map[string]map[*Client]struct{}),
	}
}

// Register adds a client to a lobby group.
func (h *Hub) Register(lobbyID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.lobbies[lobbyID] == nil {
		h.lobbies[lobbyID] = make(map[*Client]struct{})
	}
	h.lobbies[lobbyID][client] = struct{}{}
}

// Unregister removes a client and cleans up empty lobby groups.
func (h *Hub) Unregister(lobbyID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	clients, ok := h.lobbies[lobbyID]
	if !ok {
		return
	}

	if _, exists := clients[client]; exists {
		delete(clients, client)
		close(client.send)
	}

	if len(clients) == 0 {
		delete(h.lobbies, lobbyID)
	}
}

// Broadcast sends a message to all clients in a lobby.
func (h *Hub) Broadcast(lobbyID string, message []byte) {
	h.mu.RLock()
	clients := h.lobbies[lobbyID]
	snapshot := make([]*Client, 0, len(clients))
	for client := range clients {
		snapshot = append(snapshot, client)
	}
	h.mu.RUnlock()

	for _, client := range snapshot {
		select {
		case client.send <- message:
		default:
			go func(c *Client) {
				h.Unregister(lobbyID, c)
				_ = c.conn.Close()
			}(client)
		}
	}
}

// ClientCount returns the number of connected clients in a lobby.
func (h *Hub) ClientCount(lobbyID string) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.lobbies[lobbyID])
}
