package httpserver

import (
	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/ws"
)

func registerWSRoutes(r gin.IRoutes, handler *ws.Handler) {
	r.GET("/lobby/:id", handler.Handle)
}
