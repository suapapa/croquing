package httpserver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquing/internal/lobby"
	"github.com/suapapa/croquing/internal/pixabay"
	"github.com/suapapa/croquing/internal/ws"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func newRouter(store lobby.Store, drawDuration time.Duration, pixabayClient *pixabay.Client, wsHandler *ws.Handler, lobbySync *ws.SnapshotSync, corsOrigins []string, appName string, appLogo string, appLogoLink string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogMiddleware())
	r.Use(corsMiddleware(corsOrigins))

	r.GET("/health", healthHandler)

	api := r.Group("/api")
	registerConfigRoutes(api, appName, appLogo, appLogoLink)
	registerLobbyRoutes(api, store, drawDuration, lobbySync)
	if pixabayClient != nil {
		registerPixabayRoutes(api, store, pixabayClient)
	}

	if wsHandler != nil {
		wsGroup := r.Group("/ws")
		registerWSRoutes(wsGroup, wsHandler)
	}

	if !registerStaticRoutes(r, "", appName, appLogo) {
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		})
	}

	return r
}

func healthHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
