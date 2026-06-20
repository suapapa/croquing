package ws

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/suapapa/croquis-king/internal/lobby"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Handler upgrades HTTP requests to WebSocket connections for lobby subscriptions.
type Handler struct {
	hub   *Hub
	store lobby.Store
}

// NewHandler creates a WebSocket handler.
func NewHandler(hub *Hub, store lobby.Store) *Handler {
	return &Handler{
		hub:   hub,
		store: store,
	}
}

// Handle upgrades GET /ws/lobby/:id and registers the connection with the hub.
func (h *Handler) Handle(c *gin.Context) {
	lobbyID := c.Param("id")

	if _, err := h.store.Get(c.Request.Context(), lobbyID); err != nil {
		if errors.Is(err, lobby.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lobby not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load lobby"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := newClient(h.hub, lobbyID, conn)
	h.hub.Register(lobbyID, client)

	go client.writePump()
	go client.readPump()
}
