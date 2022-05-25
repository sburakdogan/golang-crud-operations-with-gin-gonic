package routes

import (
	"github.com/burak/product-api/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(routes *gin.Engine) {
	routes.POST("/auth/register", controllers.Register())
	routes.POST("/auth/login", controllers.Login())
}
