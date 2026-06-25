package httpserver

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

const defaultStaticDir = "frontend/dist"

// registerStaticRoutes serves the built React SPA when frontend/dist exists.
// Returns true when static routes were registered.
func registerStaticRoutes(r *gin.Engine, staticDir string, appName string, appLogo string) bool {
	if strings.TrimSpace(staticDir) == "" {
		staticDir = defaultStaticDir
	}

	info, err := os.Stat(staticDir)
	if err != nil || !info.IsDir() {
		return false
	}

	indexPath := filepath.Join(staticDir, "index.html")
	indexHTML, err := os.ReadFile(indexPath)
	if err != nil {
		return false
	}

	assetsDir := filepath.Join(staticDir, "assets")
	if info, err := os.Stat(assetsDir); err == nil && info.IsDir() {
		r.Static("/assets", assetsDir)
	}

	for _, name := range []string{"favicon.svg", "icons.svg"} {
		filePath := filepath.Join(staticDir, name)
		if _, err := os.Stat(filePath); err == nil {
			r.StaticFile("/"+name, filePath)
		}
	}

	if logoPath := resolveAppLogo(appLogo); strings.HasPrefix(logoPath, "/") {
		filePath := filepath.Join(staticDir, strings.TrimPrefix(logoPath, "/"))
		if _, err := os.Stat(filePath); err == nil {
			r.StaticFile(logoPath, filePath)
		}
	}

	serveSPA := func(c *gin.Context) {
		html := withAppTitle(string(indexHTML), appName)
		html = injectBeforeHeadClose(html, buildSocialMetaTags(c, appName, appLogo))
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}

	r.GET("/", serveSPA)

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/ws/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		serveSPA(c)
	})

	return true
}
