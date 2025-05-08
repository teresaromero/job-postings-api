package storage

import (
	"fmt"
	"job-postings-api/internal/errors"
	"job-postings-api/internal/models"
	"time"

	"github.com/google/uuid"
)

type Storage struct {
	postMap map[uint]models.Post
}

func NewStorage() *Storage {
	return &Storage{
		postMap: make(map[uint]models.Post),
	}
}

func (s *Storage) GetPost(id uint) (*models.Post, error) {
	item, ok := s.postMap[id]
	if !ok {
		return nil, fmt.Errorf("%w: post not found", errors.ErrNotFound)
	}

	return &item, nil
}

func (s *Storage) CreatePost(post *models.PostCreateRequest) (*models.Post, error) {
	p := models.Post{
		ID:          uint(uuid.New().ID()),
		Title:       post.Title,
		Company:     post.Company,
		Description: post.Description,
		Type:        post.Type,
		Location:    post.Location,
		Salary:      []int{post.MinSalary, post.MaxSalary},
		Perks:       post.Perks,
		Extras:      post.Extras,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.postMap[p.ID] = p
	return &p, nil

}

func (s *Storage) UpdatePost(id uint, post *models.PostPutRequest) error {
	item, ok := s.postMap[id]
	if !ok {
		return fmt.Errorf("%w: post not found", errors.ErrNotFound)
	}

	item.Title = post.Title
	item.Company = post.Company
	item.Description = post.Description
	item.Type = post.Type
	item.Location = post.Location
	item.Salary = []int{post.MinSalary, post.MaxSalary}
	item.Perks = post.Perks
	item.Extras = post.Extras
	item.UpdatedAt = time.Now()

	s.postMap[id] = item

	return nil
}

func (s *Storage) DeletePost(id uint) error {
	item, ok := s.postMap[id]
	if !ok {
		return fmt.Errorf("%w: post not found", errors.ErrNotFound)
	}

	delete(s.postMap, item.ID)
	return nil
}

func (s *Storage) ListPosts() ([]*models.Post, error) {
	list := make([]*models.Post, 0, len(s.postMap))
	for id := range s.postMap {
		job := s.postMap[id]
		list = append(list, &job)
	}
	return list, nil
}
