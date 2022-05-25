package main

import (
	"os"

	"github.com/burak/product-api/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRoutes(router)
	routes.ProductRoutes(router)

	router.Run(":" + port)
}
