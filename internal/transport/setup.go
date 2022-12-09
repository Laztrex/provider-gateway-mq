package transport

import (
	"github.com/gin-gonic/gin"
)

// SetupRouters Function to setup routers and router groups
func SetupRouters(WebApp *gin.Engine) {
	v1 := WebApp.Group("/v1")
	{
		v1.POST("/predict", Provider)
	}
	WebApp.GET("/health", Health)
}
