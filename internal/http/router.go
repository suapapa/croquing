package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/lobby"
	"github.com/suapapa/croquis-king/internal/pixabay"
	"github.com/suapapa/croquis-king/internal/ws"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func newRouter(store lobby.Store, drawDuration time.Duration, pixabayClient *pixabay.Client, wsHandler *ws.Handler, lobbySync *ws.SnapshotSync, corsOrigins []string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogMiddleware())
	r.Use(corsMiddleware(corsOrigins))

	r.GET("/health", healthHandler)

	api := r.Group("/api")
	registerLobbyRoutes(api, store, drawDuration, lobbySync)
	if pixabayClient != nil {
		registerPixabayRoutes(api, store, pixabayClient)
	}

	if wsHandler != nil {
		wsGroup := r.Group("/ws")
		registerWSRoutes(wsGroup, wsHandler)
	}

	if !registerStaticRoutes(r, "") {
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		})
	}

	return r
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
