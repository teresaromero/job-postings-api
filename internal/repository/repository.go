package repository

import (
	"job-postings-api/internal/models"
	"job-postings-api/internal/sorting"
)

type storageInterface interface {
	GetPost(id string) (*models.Post, error)
	CreatePost(post *models.PostCreateRequest) (*models.Post, error)
	UpdatePost(id string, post *models.PostPutRequest) error
	DeletePost(id string) error
	ListPosts(filter models.ListRequestQueryParams) ([]*models.Post, error)
	CompanyMapCount() map[string]int
}

type Repository struct {
	storage storageInterface
}

func NewRepository(storage storageInterface) *Repository {
	return &Repository{
		storage: storage,
	}
}

func (r *Repository) GetPost(id string) (*models.Post, error) {
	return r.storage.GetPost(id)
}

func (r *Repository) CreatePost(post *models.PostCreateRequest) (*models.Post, error) {
	return r.storage.CreatePost(post)
}

func (r *Repository) UpdatePost(id string, post *models.PostPutRequest) error {
	return r.storage.UpdatePost(id, post)
}

func (r *Repository) DeletePost(id string) error {
	return r.storage.DeletePost(id)
}

func (r *Repository) ListPosts(filter models.ListRequestQueryParams) (*models.PostList, error) {
	posts, err := r.storage.ListPosts(filter)
	if err != nil {
		return nil, err
	}
	companyMap := r.storage.CompanyMapCount()

	sorting.SortPosts(posts, companyMap)

	return &models.PostList{
		Posts: posts,
		Count: len(posts),
	}, nil
}
