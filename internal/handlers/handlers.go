package handlers

import (
	"job-postings-api/internal/config"
	"job-postings-api/internal/errors"
	"job-postings-api/internal/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type repositoryInterface interface {
	GetPost(id string) (*models.Post, error)
	CreatePost(post *models.PostRequestPayload) (*models.Post, error)
	UpdatePost(id string, post *models.PostRequestPayload) error
	DeletePost(id string) error
	ListPosts(filter models.ListRequestQueryParams) (*models.PostList, error)
}

type Handler struct {
	// bypassing the repository interface here although it should the service interface
	repo repositoryInterface

	jwtSecret []byte
}

func NewHandler(repo repositoryInterface, jwtSecret []byte) *Handler {
	return &Handler{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (h *Handler) ListPosts(c *gin.Context) {
	var filter models.ListRequestQueryParams
	if err := c.ShouldBindQuery(&filter); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	posts, err := h.repo.ListPosts(filter)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to list posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *Handler) GetPost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// if err != nil {
	post, err := h.repo.GetPost(id)
	if err != nil {
		if errors.IsNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to get post"})
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *Handler) CreatePost(c *gin.Context) {
	post := &models.PostRequestPayload{}
	if err := c.ShouldBindJSON(post); err != nil {
		// TODO: improve feedback for invalid request
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	newPost, err := h.repo.CreatePost(post)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, newPost)

}

func (h *Handler) UpdatePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	post := &models.PostRequestPayload{}
	if err := c.ShouldBindJSON(post); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.repo.UpdatePost(id, post)
	if err != nil {
		if errors.IsNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) DeletePost(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	err := h.repo.DeletePost(id)
	if err != nil {
		if errors.IsNotFound(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "post not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to delete post"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *Handler) GetToken(c *gin.Context) {
	var creds *models.LoginRequestPayload
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Only allow login for specific usernames
	// TODO: implement proper authentication
	if creds.Username != models.EmployeeUserType.String() &&
		creds.Username != models.EmployerUserType.String() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      config.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(h.jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}
