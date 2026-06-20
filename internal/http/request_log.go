package httpserver

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func requestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.Request.URL.Path
		if path == "/health" {
			return
		}

		level := slog.LevelInfo
		status := c.Writer.Status()
		if status >= 500 {
			level = slog.LevelError
		} else if status >= 400 {
			level = slog.LevelWarn
		}

		slog.Log(c.Request.Context(), level, "http request",
			"method", c.Request.Method,
			"path", path,
			"status", status,
			"duration", time.Since(start),
			"client_ip", c.ClientIP(),
		)
	}
}
