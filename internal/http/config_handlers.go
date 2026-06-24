package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type publicConfigResponse struct {
	AppName string `json:"app_name"`
}

func registerConfigRoutes(r *gin.RouterGroup, appName string) {
	r.GET("/config", publicConfigHandler(appName))
}

func publicConfigHandler(appName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, publicConfigResponse{AppName: appName})
	}
}
