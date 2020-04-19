package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ibamibrhm/donation-server/models"
	"github.com/ibamibrhm/donation-server/routes"
)

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           time.Hour * 24 * 7,
	}))
	db := models.SetupModels()

	//routes
	userRouter := new(routes.UserRouter)

	defer db.Close()

	router.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"data": "hello world"})
	})

	// all routes
	userRouter.Routes(router)

	router.Run()
}
