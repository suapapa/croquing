package httpserver

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquing/internal/lobby"
)

const adminLobbyContextKey = "adminLobby"

func requireAdmin(store lobby.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		lob, ok := authenticateAdminRequest(c, store, adminLobbyID(c))
		if !ok {
			return
		}
		c.Set(adminLobbyContextKey, lob)
		c.Next()
	}
}

func adminLobbyID(c *gin.Context) string {
	if id := c.Param("id"); id != "" {
		return id
	}
	return c.Query("lobby_id")
}

func authenticateAdminRequest(c *gin.Context, store lobby.Store, lobbyID string) (*lobby.Lobby, bool) {
	token := c.GetHeader(lobby.AdminTokenHeader)
	if lobbyID == "" || token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "lobby_id and admin token are required"})
		return nil, false
	}

	lob, err := store.Get(c.Request.Context(), lobbyID)
	if err != nil {
		if errors.Is(err, lobby.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lobby not found"})
			return nil, false
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to verify admin"})
		return nil, false
	}

	if !lobby.ValidateAdminToken(lob, token) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid admin token"})
		return nil, false
	}

	return lob, true
}
