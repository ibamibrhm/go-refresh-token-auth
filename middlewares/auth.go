package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibamibrhm/donation-server/helpers"
)

// Auth -> token authorization
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		userID, err := helpers.TokenValidation(bearerToken)

		c.Set("userId", userID)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Anda tidak memiliki akses untuk mengubah"})
			c.Abort()
			return
		}

		c.Next()
	}
}
