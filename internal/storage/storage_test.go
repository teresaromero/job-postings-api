package storage

import (
	"job-postings-api/internal/errors"
	"job-postings-api/internal/models"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStorage_GetPost(t *testing.T) {
	s := NewStorage()

	// Setup test data
	testPost := models.Post{
		ID:          1,
		Title:       "Software Engineer",
		Company:     "Test Company",
		Description: "Test Description",
		Type:        "Full-time",
		Location:    "Remote",
		Salary:      []int{50000, 80000},
		Perks:       []string{"Health Insurance"},
		Extras:      "Flexible Hours",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.postMap[testPost.ID] = testPost

	tests := []struct {
		name    string
		id      uint
		want    *models.Post
		wantErr error
	}{
		{
			name:    "existing post",
			id:      1,
			want:    &testPost,
			wantErr: nil,
		},
		{
			name:    "non-existing post",
			id:      999,
			want:    nil,
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.GetPost(tt.id)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestStorage_CreatePost(t *testing.T) {
	s := NewStorage()

	tests := []struct {
		name    string
		post    *models.PostCreateRequest
		want    *models.Post
		wantErr error
	}{
		{
			name: "create valid post",
			post: &models.PostCreateRequest{
				Title:       "Software Engineer",
				Company:     "Test Company",
				Description: "Test Description",
				Type:        "Full-time",
				Location:    "Remote",
				MinSalary:   50000,
				MaxSalary:   80000,
				Perks:       []string{"Health Insurance"},
				Extras:      "Flexible Hours",
			},
			want:    nil, // Not comparing exact post due to dynamic ID and timestamps
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.CreatePost(tt.post)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, got)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.NotZero(t, got.ID)
			assert.Equal(t, tt.post.Title, got.Title)
			assert.Equal(t, tt.post.Company, got.Company)
			assert.Equal(t, tt.post.Description, got.Description)
			assert.Equal(t, tt.post.Type, got.Type)
			assert.Equal(t, tt.post.Location, got.Location)
			assert.Equal(t, []int{tt.post.MinSalary, tt.post.MaxSalary}, got.Salary)
			assert.Equal(t, tt.post.Perks, got.Perks)
			assert.Equal(t, tt.post.Extras, got.Extras)
			assert.NotZero(t, got.CreatedAt)
			assert.NotZero(t, got.UpdatedAt)
		})
	}
}

func TestStorage_UpdatePost(t *testing.T) {
	s := NewStorage()

	// Setup test data
	testPost := models.Post{
		ID:          1,
		Title:       "Software Engineer",
		Company:     "Test Company",
		Description: "Test Description",
		Type:        "Full-time",
		Location:    "Remote",
		Salary:      []int{50000, 80000},
		Perks:       []string{"Health Insurance"},
		Extras:      "Flexible Hours",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.postMap[testPost.ID] = testPost

	tests := []struct {
		name    string
		id      uint
		post    *models.PostPutRequest
		wantErr error
	}{
		{
			name: "update existing post",
			id:   1,
			post: &models.PostPutRequest{
				Title:       "Updated Title",
				Company:     "Updated Company",
				Description: "Updated Description",
				Type:        "Part-time",
				Location:    "On-site",
				MinSalary:   60000,
				MaxSalary:   90000,
				Perks:       []string{"Health Insurance", "401k"},
				Extras:      "Updated Extras",
			},
			wantErr: nil,
		},
		{
			name: "update non-existing post",
			id:   999,
			post: &models.PostPutRequest{
				Title: "Test",
			},
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.UpdatePost(tt.id, tt.post)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			updatedPost, err := s.GetPost(tt.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.post.Title, updatedPost.Title)
			assert.Equal(t, tt.post.Company, updatedPost.Company)
			assert.Equal(t, tt.post.Description, updatedPost.Description)
			assert.Equal(t, tt.post.Type, updatedPost.Type)
			assert.Equal(t, tt.post.Location, updatedPost.Location)
			assert.Equal(t, []int{tt.post.MinSalary, tt.post.MaxSalary}, updatedPost.Salary)
			assert.Equal(t, tt.post.Perks, updatedPost.Perks)
			assert.Equal(t, tt.post.Extras, updatedPost.Extras)
			assert.NotEqual(t, testPost.UpdatedAt, updatedPost.UpdatedAt)
		})
	}
}
func TestStorage_DeletePost(t *testing.T) {
	s := NewStorage()

	// Setup test data
	testPost := models.Post{
		ID:          1,
		Title:       "Software Engineer",
		Company:     "Test Company",
		Description: "Test Description",
		Type:        "Full-time",
		Location:    "Remote",
		Salary:      []int{50000, 80000},
		Perks:       []string{"Health Insurance"},
		Extras:      "Flexible Hours",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	s.postMap[testPost.ID] = testPost

	tests := []struct {
		name    string
		id      uint
		wantErr error
	}{
		{
			name:    "delete existing post",
			id:      1,
			wantErr: nil,
		},
		{
			name:    "delete non-existing post",
			id:      999,
			wantErr: errors.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.DeletePost(tt.id)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}
			assert.NoError(t, err)

			// Verify post was deleted
			_, err = s.GetPost(tt.id)
			assert.ErrorIs(t, err, errors.ErrNotFound)
		})
	}
}
func TestStorage_ListPosts(t *testing.T) {
	s := NewStorage()

	// Test empty storage
	t.Run("empty storage", func(t *testing.T) {
		posts, err := s.ListPosts()
		assert.NoError(t, err)
		assert.Empty(t, posts)
	})

	// Setup test data
	testPosts := []models.Post{
		{
			ID:          1,
			Title:       "Software Engineer",
			Company:     "Company A",
			Description: "Description A",
			Type:        "Full-time",
			Location:    "Remote",
			Salary:      []int{50000, 80000},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Title:       "Product Manager",
			Company:     "Company B",
			Description: "Description B",
			Type:        "Part-time",
			Location:    "On-site",
			Salary:      []int{60000, 90000},
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Add posts to storage
	for _, post := range testPosts {
		s.postMap[post.ID] = post
	}

	t.Run("storage with posts", func(t *testing.T) {
		posts, err := s.ListPosts()
		assert.NoError(t, err)
		assert.Len(t, posts, len(testPosts))

		// Verify each post was returned
		for _, got := range posts {
			original := testPosts[got.ID-1] // Index 0 based, IDs 1 based
			assert.Equal(t, original.ID, got.ID)
			assert.Equal(t, original.Title, got.Title)
			assert.Equal(t, original.Company, got.Company)
			assert.Equal(t, original.Description, got.Description)
			assert.Equal(t, original.Type, got.Type)
			assert.Equal(t, original.Location, got.Location)
			assert.Equal(t, original.Salary, got.Salary)
		}
	})
}
