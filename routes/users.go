package routes

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/ibamibrhm/donation-server/controllers"
	"github.com/ibamibrhm/donation-server/helpers"
	"github.com/ibamibrhm/donation-server/middlewares"
	"github.com/ibamibrhm/donation-server/models"
	"github.com/jinzhu/gorm"
)

// UserRouter -> group of all user routes
type UserRouter struct{}

// Routes -> routes for endpoint /users
func (route UserRouter) Routes(router *gin.Engine) {
	usersController := new(controllers.UserController)

	userRoutes := router.Group("/users")
	{
		userRoutes.GET("", usersController.FindUsers)
		userRoutes.POST("/register", usersController.Register)
		userRoutes.POST("/login", usersController.Login)
		userRoutes.POST("/logout", usersController.Logout)
		userRoutes.POST("/refresh_token", func(c *gin.Context) {
			cookie, err := c.Cookie("jid")

			if err != nil {
				c.JSON(http.StatusOK, gin.H{"data": ""})
				return
			}

			payload, err := helpers.TokenValidation("Bearer "+cookie, os.Getenv("JWT_REFRESH_SECRET"))
			if err != nil {
				c.JSON(http.StatusOK, gin.H{"data": ""})
				return
			}

			db := c.MustGet("db").(*gorm.DB)

			// Get model if exist
			var user models.User
			if err := db.Where("id = ?", payload.UserID).First(&user).Error; err != nil {
				c.JSON(http.StatusOK, gin.H{"data": ""})
				return
			}

			if user.TokenVersion != payload.TokenVersion {
				c.JSON(http.StatusOK, gin.H{"data": ""})
				return
			}

			token, _ := helpers.CreateToken(user)
			refreshToken, _ := helpers.CreateRefreshToken(user)
			c.SetCookie("jid", refreshToken, 60*60*24*7 /* 7 days */, "/", "localhost", false, true)

			c.JSON(http.StatusOK, gin.H{"data": token})
		})
	}
	userRoutes.Use(middlewares.Auth())
	{
		userRoutes.GET("/:id", usersController.FindUser)
		userRoutes.DELETE("/:id", usersController.DeleteUser)
		userRoutes.PATCH("/:id", usersController.UpdateUser)
	}
}
