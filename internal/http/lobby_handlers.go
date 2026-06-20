package httpserver

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/lobby"
)

type lobbyHandler struct {
	store        lobby.Store
	drawDuration time.Duration
}

type createLobbyResponse struct {
	ID         string `json:"id"`
	AdminToken string `json:"admin_token"`
	JoinURL    string `json:"join_url"`
}

func newLobbyHandler(store lobby.Store, drawDuration time.Duration) *lobbyHandler {
	return &lobbyHandler{
		store:        store,
		drawDuration: drawDuration,
	}
}

func (h *lobbyHandler) createLobby(c *gin.Context) {
	created, err := h.store.Create(c.Request.Context(), h.drawDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create lobby"})
		return
	}

	c.JSON(http.StatusCreated, createLobbyResponse{
		ID:         created.ID,
		AdminToken: created.AdminToken,
		JoinURL:    joinURL(c, created.ID),
	})
}

func (h *lobbyHandler) getLobby(c *gin.Context) {
	id := c.Param("id")
	snapshot, err := h.store.Snapshot(c.Request.Context(), id, 0)
	if err != nil {
		if errors.Is(err, lobby.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "lobby not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get lobby"})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

func joinURL(c *gin.Context, lobbyID string) string {
	scheme := "http"
	if c.Request.TLS != nil || c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host + "/lobby/" + lobbyID
}

func registerLobbyRoutes(r gin.IRoutes, store lobby.Store, drawDuration time.Duration) {
	handler := newLobbyHandler(store, drawDuration)
	r.POST("/lobbies", handler.createLobby)
	r.GET("/lobbies/:id", handler.getLobby)
}
