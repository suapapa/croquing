package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type publicConfigResponse struct {
	AppName     string `json:"app_name"`
	AppLogo     string `json:"app_logo"`
	AppLogoLink string `json:"app_logo_link"`
}

func registerConfigRoutes(r *gin.RouterGroup, appName string, appLogo string, appLogoLink string) {
	r.GET("/config", publicConfigHandler(appName, appLogo, appLogoLink))
}

func publicConfigHandler(appName string, appLogo string, appLogoLink string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, publicConfigResponse{
			AppName:     appName,
			AppLogo:     appLogo,
			AppLogoLink: appLogoLink,
		})
	}
}
