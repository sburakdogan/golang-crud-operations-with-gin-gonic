package middleware

import (
	"net/http"

	"github.com/burak/product-api/helpers"
	"github.com/gin-gonic/gin"
)

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "Authorization token is not found"})
			c.Abort()
			return
		}

		err := helpers.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": err.Error()})
			c.Abort()
			return
		}

		c.Next()
	}
}
