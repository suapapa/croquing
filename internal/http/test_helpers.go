package httpserver

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquing/internal/lobby"
	"github.com/suapapa/croquing/internal/pixabay"
	"github.com/suapapa/croquing/internal/ws"
)

func newTestRouter(store lobby.Store, drawDuration time.Duration, pixabayClient *pixabay.Client, wsHandler *ws.Handler, lobbySync *ws.SnapshotSync) *gin.Engine {
	return newRouter(store, drawDuration, pixabayClient, wsHandler, lobbySync, nil, "Croquing", "", "")
}
