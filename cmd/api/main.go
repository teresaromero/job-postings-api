package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.New()

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// router for v1
	apiv1 := engine.Group("/api/v1")
	apiv1.GET("/posts")
	apiv1.GET("/posts/:id")
	apiv1.POST("/posts")
	apiv1.PUT("/posts/:id")
	apiv1.DELETE("/posts/:id")

	// healthcheck
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
}
