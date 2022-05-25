package routes

import (
	"github.com/burak/product-api/controllers"
	"github.com/burak/product-api/middleware"
	"github.com/gin-gonic/gin"
)

func ProductRoutes(routes *gin.Engine) {
	routes.Use(middleware.Authorization())
	routes.GET("/products", controllers.GetProducts())
	routes.GET("/product/:id", controllers.GetProduct())
	routes.DELETE("/product/:id", controllers.DeleteProduct())
	routes.POST("/product", controllers.CreateProduct())
	routes.PUT("/product/:id", controllers.UpdateProduct())
}
