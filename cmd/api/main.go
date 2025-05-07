package main

import (
	"job-postings-api/internal/handlers"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()
	engine.SetTrustedProxies(nil)

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	handler := handlers.NewHandler()
	
	// router for v1
	apiv1 := engine.Group("/api/v1")
	apiv1.GET("/posts", handler.GetPosts)
	apiv1.GET("/posts/:id", handler.GetPost)
	apiv1.POST("/posts", handler.CreatePost)
	apiv1.PUT("/posts/:id", handler.UpdatePost)
	apiv1.DELETE("/posts/:id", handler.DeletePost)

	// healthcheck
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
