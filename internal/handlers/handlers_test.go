package handlers

import (
	"encoding/json"
	"job-postings-api/internal/errors"
	"job-postings-api/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListPosts() (*models.PostList, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.PostList), args.Error(1)
}

func (m *MockRepository) GetPost(id string) (*models.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *MockRepository) CreatePost(post *models.PostCreateRequest) (*models.Post, error) {
	return nil, nil
}

func (m *MockRepository) UpdatePost(id string, post *models.PostPutRequest) error {
	return nil
}

func (m *MockRepository) DeletePost(id string) error {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
}

func TestHandler_ListPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockResponse   *models.PostList
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success",
			mockResponse: &models.PostList{
				Posts: []*models.Post{
					{
						ID:        "1",
						Title:     "Test Post",
						Type:      "full-time",
						Perks:     []string{},
						Salary:    []int{},
						CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
						UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				Count: 1,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"posts": []interface{}{
					map[string]interface{}{
						"id":          "1",
						"title":       "Test Post",
						"description": "",
						"extras":      "",
						"created_at":  "2023-10-01T00:00:00Z",
						"updated_at":  "2023-10-01T00:00:00Z",
						"company":     "",
						"location":    "",
						"perks":       []interface{}{},
						"salary":      []interface{}{},
						"type":        "full-time",
					},
				},
				"count": float64(1),
			},
		},
		{
			name:           "repository error",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to list posts"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("ListPosts").Return(tt.mockResponse, tt.mockError)

			handler := NewHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			handler.ListPosts(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestHandler_GetPost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		id             string
		mockResponse   *models.Post
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success",
			id:   "1",
			mockResponse: &models.Post{
				ID:        "1",
				Title:     "Test Post",
				Perks:     []string{},
				Salary:    []int{},
				CreatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":          "1",
				"title":       "Test Post",
				"description": "",
				"extras":      "",
				"created_at":  "2023-10-01T00:00:00Z",
				"updated_at":  "2023-10-01T00:00:00Z",
				"company":     "",
				"location":    "",
				"perks":       []interface{}{},
				"salary":      []interface{}{},
				"type":        "",
			},
		},
		{
			name:           "empty id",
			id:             "",
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   gin.H{"error": "id is required"},
		},
		{
			name:           "not found",
			id:             "999",
			mockResponse:   nil,
			mockError:      errors.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   gin.H{"error": "post not found"},
		},
		{
			name:           "repository error",
			id:             "1",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   gin.H{"error": "failed to get post"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tt.id != "" {
				mockRepo.On("GetPost", tt.id).Return(tt.mockResponse, tt.mockError)
			}

			handler := NewHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.id}}

			handler.GetPost(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockRepo.AssertExpectations(t)
		})
	}
}
func TestHandler_DeletePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		id             string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "success",
			id:             "1",
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			name:           "empty id",
			id:             "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "id is required"},
		},
		{
			name:           "not found",
			id:             "999",
			mockError:      errors.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]interface{}{"error": "post not found"},
		},
		{
			name:           "repository error",
			id:             "1",
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to delete post"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tt.id != "" {
				mockRepo.On("DeletePost", tt.id).Return(tt.mockError)
			}

			handler := NewHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = []gin.Param{{Key: "id", Value: tt.id}}

			handler.DeletePost(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody == nil {
				assert.Empty(t, w.Body.String())
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
