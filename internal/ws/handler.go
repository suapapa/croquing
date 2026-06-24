package ws

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/suapapa/croquing/internal/lobby"
)

// Handler upgrades HTTP requests to WebSocket connections for lobby subscriptions.
type Handler struct {
	sync        *SnapshotSync
	checkOrigin func(*http.Request) bool
}

// NewHandler creates a WebSocket handler.
func NewHandler(sync *SnapshotSync, checkOrigin func(*http.Request) bool) *Handler {
	if checkOrigin == nil {
		checkOrigin = func(*http.Request) bool { return true }
	}
	return &Handler{sync: sync, checkOrigin: checkOrigin}
}

// Handle upgrades GET /ws/lobby/:id and registers the connection with the hub.
func (h *Handler) Handle(c *gin.Context) {
	lobbyID := c.Param("id")

	if err := h.sync.LobbyExists(c.Request.Context(), lobbyID); err != nil {
		if errors.Is(err, lobby.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lobby not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load lobby"})
		return
	}

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     h.checkOrigin,
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := newClient(h.sync, lobbyID, conn)
	if err := h.sync.RegisterClient(c.Request.Context(), lobbyID, client); err != nil {
		_ = conn.Close()
		return
	}

	go client.writePump()
	go client.readPump()
}
