package handlers

import "github.com/gin-gonic/gin"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}
func (h *Handler) GetPosts(c *gin.Context) {
	// Handler logic for getting posts
}
func (h *Handler) GetPost(c *gin.Context) {
	// Handler logic for getting a single post
}
func (h *Handler) CreatePost(c *gin.Context) {
	// Handler logic for creating a post
}
func (h *Handler) UpdatePost(c *gin.Context) {
	// Handler logic for updating a post
}
func (h *Handler) DeletePost(c *gin.Context) {
	// Handler logic for deleting a post
}
