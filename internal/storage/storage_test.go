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
		ID:          "1",
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
		id      string
		want    *models.Post
		wantErr error
	}{
		{
			name:    "existing post",
			id:      "1",
			want:    &testPost,
			wantErr: nil,
		},
		{
			name:    "non-existing post",
			id:      "999",
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
		ID:          "1",
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
		id      string
		post    *models.PostPutRequest
		wantErr error
	}{
		{
			name: "update existing post",
			id:   "1",
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
			id:   "999",
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
		ID:          "1",
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
		id      string
		wantErr error
	}{
		{
			name:    "delete existing post",
			id:      "1",
			wantErr: nil,
		},
		{
			name:    "delete non-existing post",
			id:      "999",
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

	// Setup test data
	testPosts := []models.Post{
		{
			ID:          "1",
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
		},
		{
			ID:          "2",
			Title:       "Product Manager",
			Company:     "Another Company",
			Description: "Another Description",
			Type:        "Part-time",
			Location:    "On-site",
			Salary:      []int{60000, 90000},
			Perks:       []string{"401k"},
			Extras:      "Gym Membership",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Add test posts to storage
	for _, post := range testPosts {
		s.postMap[post.ID] = post
	}

	// Test listing posts
	t.Run("list all posts", func(t *testing.T) {
		got, err := s.ListPosts(models.ListRequestQueryParams{})
		assert.NoError(t, err)
		assert.Len(t, got, len(testPosts))

		// Verify each post is in the list
		for _, want := range testPosts {
			found := false
			for _, post := range got {
				if post.ID == want.ID {
					assert.Equal(t, want.Title, post.Title)
					assert.Equal(t, want.Company, post.Company)
					assert.Equal(t, want.Description, post.Description)
					assert.Equal(t, want.Type, post.Type)
					assert.Equal(t, want.Location, post.Location)
					assert.Equal(t, want.Salary, post.Salary)
					assert.Equal(t, want.Perks, post.Perks)
					assert.Equal(t, want.Extras, post.Extras)
					found = true
					break
				}
			}
			assert.True(t, found, "Post with ID %s not found in list", want.ID)
		}
	})

	// Test empty storage
	t.Run("list empty storage", func(t *testing.T) {
		s := NewStorage()
		got, err := s.ListPosts(models.ListRequestQueryParams{})
		assert.NoError(t, err)
		assert.Empty(t, got)
	})
}
func TestStorage_CompanyCount(t *testing.T) {
	s := NewStorage()

	// Setup test data
	post1 := models.Post{
		ID:      "1",
		Company: "Company A",
	}
	post2 := models.Post{
		ID:      "2",
		Company: "Company A",
	}
	post3 := models.Post{
		ID:      "3",
		Company: "Company B",
	}

	// Add test posts to storage and increment company counts
	s.postMap[post1.ID] = post1
	s.companyIndex[post1.Company]++
	s.postMap[post2.ID] = post2
	s.companyIndex[post2.Company]++
	s.postMap[post3.ID] = post3
	s.companyIndex[post3.Company]++

	tests := []struct {
		name    string
		company string
		want    int
	}{
		{
			name:    "company with multiple posts",
			company: "Company A",
			want:    2,
		},
		{
			name:    "company with single post",
			company: "Company B",
			want:    1,
		},
		{
			name:    "non-existing company",
			company: "Company C",
			want:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.CompanyCount(tt.company)
			assert.Equal(t, tt.want, got)
		})
	}
}
func TestNewStorage(t *testing.T) {
	s := NewStorage()
	assert.NotNil(t, s)
	assert.NotNil(t, s.postMap)
	assert.NotNil(t, s.companyIndex)
	assert.Empty(t, s.postMap)
	assert.Empty(t, s.companyIndex)
}

func Test_StorageIntegration(t *testing.T) {
	s := NewStorage()

	// Create a new post
	post := &models.PostCreateRequest{
		Title:       "Software Engineer",
		Company:     "Tech Company",
		Description: "Develop software solutions.",
		Type:        "Full-time",
		Location:    "Remote",
		MinSalary:   70000,
		MaxSalary:   120000,
		Perks:       []string{"Health Insurance", "401k"},
		Extras:      "Flexible Hours",
	}

	newPost, err := s.CreatePost(post)
	assert.NoError(t, err)
	assert.NotNil(t, newPost)

	count := s.CompanyCount(newPost.Company)
	assert.Equal(t, 1, count)

	// Get the created post
	gotPost, err := s.GetPost(newPost.ID)
	assert.NoError(t, err)
	assert.Equal(t, newPost.ID, gotPost.ID)

	// Update the post
	updateRequest := &models.PostPutRequest{
		Company:     "Tech Company",
		Perks:       []string{"Health Insurance", "401k"},
		Extras:      "Flexible Hours",
		Title:       "Senior Software Engineer",
		Description: "Lead software development projects.",
		Type:        "Full-time",
		Location:    "Remote",
		MinSalary:   80000,
		MaxSalary:   130000,
	}
	err = s.UpdatePost(newPost.ID, updateRequest)
	assert.NoError(t, err)

	gotUpdatedPost, err := s.GetPost(newPost.ID)
	assert.NoError(t, err)
	assert.Equal(t, updateRequest.Title, gotUpdatedPost.Title)

	// List posts
	listedPosts, err := s.ListPosts(models.ListRequestQueryParams{})
	assert.NoError(t, err)
	assert.Len(t, listedPosts, 1)

	for _, post := range listedPosts {
		assert.Equal(t, newPost.Company, post.Company)
	}

	// Delete the post
	err = s.DeletePost(newPost.ID)
	assert.NoError(t, err)

	count = s.CompanyCount(newPost.Company)
	assert.Equal(t, 0, count)

	gotDeletedPost, err := s.GetPost(newPost.ID)
	assert.ErrorIs(t, err, errors.ErrNotFound)
	assert.Nil(t, gotDeletedPost)
}
