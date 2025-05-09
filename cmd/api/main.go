package main

import (
	"job-postings-api/internal/handlers"
	"job-postings-api/internal/models"
	"job-postings-api/internal/repository"
	"job-postings-api/internal/sorting"
	"job-postings-api/internal/storage"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	engine := gin.New()
	engine.SetTrustedProxies(nil)

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// register a custom validator for jobtype enum values
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		v.RegisterValidation("jobtype", models.JobTypeValidator)
	}

	// init storage, repository and handler
	storage := storage.NewStorage()
	sorter := &sorting.Sorter{}
	repo := repository.NewRepository(storage, sorter)
	handler := handlers.NewHandler(repo)

	// router for v1
	apiv1 := engine.Group("/api/v1")
	apiv1.GET("/posts", handler.ListPosts)
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
