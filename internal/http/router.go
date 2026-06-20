package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/lobby"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func newRouter(store lobby.Store, drawDuration time.Duration) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", healthHandler)

	api := r.Group("/api")
	registerLobbyRoutes(api, store, drawDuration)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
	})

	return r
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
