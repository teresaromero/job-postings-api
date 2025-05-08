package repository

import (
	"job-postings-api/internal/models"
)

type storageInterface interface {
	GetPost(id uint) (*models.Post, error)
	CreatePost(post *models.PostCreateRequest) (*models.Post, error)
	UpdatePost(id uint, post *models.PostPutRequest) error
	DeletePost(id uint) error
	ListPosts() ([]*models.Post, error)
}

type Repository struct {
	storage storageInterface
}

func NewRepository(storage storageInterface) *Repository {
	return &Repository{
		storage: storage,
	}
}

func (r *Repository) GetPost(id uint) (*models.Post, error) {
	return r.storage.GetPost(id)
}

func (r *Repository) CreatePost(post *models.PostCreateRequest) (*models.Post, error) {
	return r.storage.CreatePost(post)
}

func (r *Repository) UpdatePost(id uint, post *models.PostPutRequest) error {
	return r.storage.UpdatePost(id, post)
}

func (r *Repository) DeletePost(id uint) error {
	return r.storage.DeletePost(id)
}

func (r *Repository) ListPosts(queryParams models.ListRequestQueryParams) (*models.PostList, error) {
	posts, err := r.storage.ListPosts()
	if err != nil {
		return nil, err
	}
	return &models.PostList{
		Posts: posts,
		Count: len(posts),
	}, nil
}
