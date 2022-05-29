package routers

import (
	"github.com/gin-gonic/gin"
	"provider_gateway_mq/controllers"
)

func SetupRouters(app *gin.Engine) {
	v1 := app.Group("/v1")
	{
		v1.POST("/predict", controllers.Provider)
	}

	app.GET("/health", controllers.Health)
}
