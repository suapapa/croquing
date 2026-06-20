package httpserver

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/ws"
)

type lobbyHandler struct {
	store        lobby.Store
	drawDuration time.Duration
	lobbySync    *ws.SnapshotSync
}

type createLobbyResponse struct {
	ID         string `json:"id"`
	AdminToken string `json:"admin_token"`
	JoinURL    string `json:"join_url"`
}

type setPhotosRequest struct {
	Photos []lobby.Photo `json:"photos"`
}

func newLobbyHandler(store lobby.Store, drawDuration time.Duration, lobbySync *ws.SnapshotSync) *lobbyHandler {
	return &lobbyHandler{
		store:        store,
		drawDuration: drawDuration,
		lobbySync:    lobbySync,
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

func (h *lobbyHandler) setPhotos(c *gin.Context) {
	id := c.Param("id")
	if _, ok := authenticateAdmin(c, h.store, id); !ok {
		return
	}

	var req setPhotosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.store.SetSelectedPhotos(c.Request.Context(), id, req.Photos); err != nil {
		switch {
		case errors.Is(err, lobby.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "lobby not found"})
		case errors.Is(err, lobby.ErrEmptyPhotos):
			c.JSON(http.StatusBadRequest, gin.H{"error": "photos are required"})
		case errors.Is(err, lobby.ErrInvalidTransition):
			c.JSON(http.StatusConflict, gin.H{"error": "invalid lobby phase for photo selection"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save photos"})
		}
		return
	}

	if h.lobbySync != nil {
		if err := h.lobbySync.Broadcast(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to broadcast snapshot"})
			return
		}
	}

	participantCount := 0
	if h.lobbySync != nil {
		participantCount = h.lobbySync.Hub().ClientCount(id)
	}

	snapshot, err := h.store.Snapshot(c.Request.Context(), id, participantCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get lobby snapshot"})
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

func registerLobbyRoutes(r gin.IRoutes, store lobby.Store, drawDuration time.Duration, lobbySync *ws.SnapshotSync) {
	handler := newLobbyHandler(store, drawDuration, lobbySync)
	r.POST("/lobbies", handler.createLobby)
	r.GET("/lobbies/:id", handler.getLobby)
	r.PUT("/lobbies/:id/photos", handler.setPhotos)
}
