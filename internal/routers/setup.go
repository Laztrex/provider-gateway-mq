package routers

import (
	"github.com/gin-gonic/gin"

	"gateway_mq/internal/controllers"
)

// SetupRouters Function to setup routers and router groups
func SetupRouters(WebApp *gin.Engine) {
	v1 := WebApp.Group("/v1")
	{
		v1.POST("/predict", controllers.Provider)
	}
	WebApp.GET("/health", Health)
}
