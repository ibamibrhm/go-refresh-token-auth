package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibamibrhm/donation-server/helpers"
)

// Auth -> token authentication
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		userID, err := helpers.TokenValidation(bearerToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		c.Set("userId", userID)

		c.Next()
	}
}
