package handlers

import (
	"job-postings-api/internal/errors"
	"job-postings-api/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type repositoryInterface interface {
	GetPost(id string) (*models.Post, error)
	CreatePost(post *models.PostCreateRequest) (*models.Post, error)
	UpdatePost(id string, post *models.PostPutRequest) error
	DeletePost(id string) error
	ListPosts() (*models.PostList, error)
}

type Handler struct {
	// bypassing the repository interface here although it should the service interface
	repo repositoryInterface
}

func NewHandler(repo repositoryInterface) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) ListPosts(c *gin.Context) {
	posts, err := h.repo.ListPosts()
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
	post := &models.PostCreateRequest{}
	if err := c.ShouldBindJSON(post); err != nil {
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

	post := &models.PostPutRequest{}
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

	c.Status(http.StatusNoContent)
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
