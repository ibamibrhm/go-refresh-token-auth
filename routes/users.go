package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ibamibrhm/donation-server/controllers"
	"github.com/ibamibrhm/donation-server/middlewares"
)

// UserRouter -> group of all user routes
type UserRouter struct{}

// Routes -> routes for endpoint /users
func (route UserRouter) Routes(router *gin.Engine) {
	usersController := new(controllers.UserController)

	userRoutes := router.Group("/users")
	{
		userRoutes.GET("", usersController.FindUsers)
		userRoutes.GET("/:id", usersController.FindUser)
		userRoutes.POST("/register", usersController.Register)
		userRoutes.POST("/login", usersController.Login)
	}
	userRoutes.Use(middlewares.Auth())
	{
		userRoutes.DELETE("/:id", usersController.DeleteUser)
		userRoutes.PATCH("/:id", usersController.UpdateUser)
	}
}
