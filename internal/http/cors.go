package httpserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/suapapa/croquis-king/internal/lobby"
)

func corsMiddleware(origins []string) gin.HandlerFunc {
	allowed := normalizeCORSOrigins(origins)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" && isOriginAllowed(origin, allowed) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, "+lobby.AdminTokenHeader)
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			if origin != "" && isOriginAllowed(origin, allowed) {
				c.Status(http.StatusNoContent)
			} else {
				c.Status(http.StatusForbidden)
			}
			c.Abort()
			return
		}

		c.Next()
	}
}

func normalizeCORSOrigins(origins []string) []string {
	if len(origins) == 0 {
		return []string{"*"}
	}

	normalized := make([]string, 0, len(origins))
	for _, origin := range origins {
		origin = strings.TrimSpace(origin)
		if origin == "" {
			continue
		}
		normalized = append(normalized, origin)
	}
	if len(normalized) == 0 {
		return []string{"*"}
	}
	return normalized
}

func isOriginAllowed(origin string, allowed []string) bool {
	for _, candidate := range allowed {
		if candidate == "*" || candidate == origin {
			return true
		}
	}
	return false
}

func originAllowed(origins []string) func(*http.Request) bool {
	allowed := normalizeCORSOrigins(origins)
	return func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		return isOriginAllowed(origin, allowed)
	}
}
