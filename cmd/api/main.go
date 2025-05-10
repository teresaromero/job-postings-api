package main

import (
	"job-postings-api/internal/config"
	"job-postings-api/internal/handlers"
	"job-postings-api/internal/middleware"
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

	cfg := config.Load()

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
	sorter := sorting.NewSorter(cfg.SortRulesOrder)
	storage := storage.NewStorage(sorter)

	if err := storage.Seed(); err != nil {
		log.Fatalf("failed to seed storage: %v", err)
	}
	repo := repository.NewRepository(storage)
	handler := handlers.NewHandler(repo, cfg.JwtSecret)

	// router for v1
	apiv1 := engine.Group("/api/v1")
	apiv1.GET("/posts", handler.ListPosts)
	apiv1.GET("/posts/:id", handler.GetPost)

	// Only employer can create, update and delete posts
	apiv1.POST("/posts",
		middleware.Authenticated(cfg.JwtSecret),
		middleware.Authorize(models.EmployerUserType.String()),
		handler.CreatePost)

	// TODO: implement ONLY the user that created the post can update or delete it
	apiv1.PUT("/posts/:id",
		middleware.Authenticated(cfg.JwtSecret),
		middleware.Authorize(models.EmployerUserType.String()),
		handler.UpdatePost)
	apiv1.DELETE("/posts/:id",
		middleware.Authenticated(cfg.JwtSecret),
		middleware.Authorize(models.EmployerUserType.String()),
		handler.DeletePost)

	// router for auth
	engine.POST("/auth/token", handler.GetToken)

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
