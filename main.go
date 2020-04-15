package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ibamibrhm/donation-server/models"
	"github.com/ibamibrhm/donation-server/routes"
)

func main() {
	router := gin.Default()
	db := models.SetupModels()

	//routes
	userRouter := new(routes.UserRouter)

	defer db.Close()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.Use(cors.Default())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	// all routes
	userRouter.Routes(router)

	router.Run()
}
