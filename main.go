package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ibamibrhm/donation-server/controllers"
	"github.com/ibamibrhm/donation-server/models"
)

func main() {
	r := gin.Default()
	db := models.SetupModels()
	defer db.Close()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	r.GET("/users", controllers.FindUsers)
	r.GET("/users/:id", controllers.FindUser)
	r.POST("/users", controllers.Register)
	r.PATCH("/users/:id", controllers.UpdateUser)
	r.DELETE("/users/:id", controllers.DeleteUser)

	r.Run()
}
