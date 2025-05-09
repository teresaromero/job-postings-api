package storage

import (
	"fmt"
	"job-postings-api/internal/errors"
	"job-postings-api/internal/models"
	"time"

	"github.com/google/uuid"
)

type sorterInterface interface {
	Sort(input *models.SortInput)
}

type Storage struct {
	postMap      map[string]models.Post
	companyIndex map[string]int

	sorter sorterInterface
}

func NewStorage(sorter sorterInterface) *Storage {
	return &Storage{
		postMap:      make(map[string]models.Post),
		companyIndex: make(map[string]int),
		sorter:       sorter,
	}
}

func (s *Storage) GetPost(id string) (*models.Post, error) {
	item, ok := s.postMap[id]
	if !ok {
		return nil, fmt.Errorf("%w: post not found", errors.ErrNotFound)
	}

	return &item, nil
}

func (s *Storage) CreatePost(post *models.PostCreateRequest) (*models.Post, error) {
	p := models.Post{
		ID:          fmt.Sprintf("%v", uuid.New().ID()),
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
	s.companyIndex[p.Company]++
	return &p, nil

}

func (s *Storage) UpdatePost(id string, post *models.PostPutRequest) error {
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
	if item.Company != post.Company {
		s.companyIndex[item.Company]--
		if s.companyIndex[item.Company] == 0 {
			delete(s.companyIndex, item.Company)
		}
		s.companyIndex[post.Company]++
	}

	return nil
}

func (s *Storage) DeletePost(id string) error {
	item, ok := s.postMap[id]
	if !ok {
		return fmt.Errorf("%w: post not found", errors.ErrNotFound)
	}

	delete(s.postMap, item.ID)
	s.companyIndex[item.Company]--
	if s.companyIndex[item.Company] == 0 {
		delete(s.companyIndex, item.Company)
	}
	return nil
}

func (s *Storage) ListPosts(filter models.ListRequestQueryParams) ([]*models.Post, error) {
	list := make([]*models.Post, 0, len(s.postMap))
	for id := range s.postMap {
		if ok := isAllowedByFilter(s.postMap[id], filter); !ok {
			continue
		}

		job := s.postMap[id]
		list = append(list, &job)
	}
	// sort the list
	input := &models.SortInput{
		Posts:      list,
		CompanyMap: s.companyIndex,
	}
	s.sorter.Sort(input)
	return list, nil
}
