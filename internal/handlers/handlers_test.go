package handlers

import (
	"bytes"
	"encoding/json"
	"job-postings-api/internal/errors"
	"job-postings-api/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) ListPosts(filter models.ListRequestQueryParams) (*models.PostList, error) {
	args := m.Called(filter)
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
	args := m.Called(post)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Post), args.Error(1)
}

func (m *MockRepository) UpdatePost(id string, post *models.PostPutRequest) error {
	args := m.Called(id, post)
	if args.Get(0) == nil {
		return nil
	}
	return args.Error(0)
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
		mockRawQuery   string
		expectedStatus int
		expectedBody   map[string]interface{}
		expectedFilter models.ListRequestQueryParams
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
			expectedFilter: models.ListRequestQueryParams{},
		},
		{
			name:           "repository error",
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to list posts"},
			expectedFilter: models.ListRequestQueryParams{},
		},
		{
			name:         "success with filter",
			mockRawQuery: `title="test post"&location=remote&company=example&max_salary=100000&min_salary=50000`,
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
			expectedFilter: models.ListRequestQueryParams{
				Title:     "\"test post\"",
				Location:  "remote",
				Company:   "example",
				MaxSalary: 100000,
				MinSalary: 50000,
			},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			mockRepo.On("ListPosts", tt.expectedFilter).Return(tt.mockResponse, tt.mockError)

			handler := NewHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/posts", nil)
			c.Request.URL.RawQuery = tt.mockRawQuery

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
func TestHandler_CreatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockResponse   *models.Post
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success",
			requestBody: map[string]interface{}{
				"title":       "Test Post",
				"perks":       []string{"perk1", "perk2"},
				"max_salary":  70000,
				"min_salary":  50000,
				"type":        "full-time",
				"company":     "Test Company",
				"location":    "Test Location",
				"extras":      "Test Extras",
				"description": "Test Description",
			},
			mockResponse: &models.Post{
				ID:          "1",
				Title:       "Test Post",
				Type:        "full-time",
				Perks:       []string{"perk1", "perk2"},
				Salary:      []int{50000, 70000},
				Company:     "Test Company",
				Location:    "Test Location",
				Extras:      "Test Extras",
				Description: "Test Description",
				CreatedAt:   time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt:   time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":          "1",
				"title":       "Test Post",
				"description": "Test Description",
				"extras":      "Test Extras",
				"created_at":  "2023-10-01T00:00:00Z",
				"updated_at":  "2023-10-01T00:00:00Z",
				"company":     "Test Company",
				"location":    "Test Location",
				"perks":       []interface{}{"perk1", "perk2"},
				"salary":      []interface{}{float64(50000), float64(70000)},
				"type":        "full-time",
			},
		},
		{
			name: "invalid request",
			requestBody: map[string]interface{}{
				"title": "",
			},
			mockResponse:   nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "invalid request"},
		},
		{
			name: "repository error",
			requestBody: map[string]interface{}{
				"title":       "Test Post",
				"perks":       []string{"perk1", "perk2"},
				"max_salary":  70000,
				"min_salary":  50000,
				"type":        "full-time",
				"company":     "Test Company",
				"location":    "Test Location",
				"extras":      "Test Extras",
				"description": "Test Description",
			},
			mockResponse:   nil,
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to create post"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tt.mockResponse != nil || tt.mockError != nil {
				mockRepo.On("CreatePost", mock.AnythingOfType("*models.PostCreateRequest")).Return(tt.mockResponse, tt.mockError)
			}

			handler := NewHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			// register a custom validator for jobtype enum values
			v, ok := binding.Validator.Engine().(*validator.Validate)
			if ok {
				v.RegisterValidation("jobtype", models.JobTypeValidator)
			}

			jsonBytes, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreatePost(c)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, response)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestHandler_UpdatePost(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		id             string
		requestBody    map[string]interface{}
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "success",
			id:   "1",
			requestBody: map[string]interface{}{
				"title":       "Updated Post",
				"description": "Updated Description",
				"type":        "full-time",
				"company":     "Updated Company",
				"location":    "Updated Location",
				"extras":      "Updated Extras",
				"perks":       []string{"perk1", "perk2"},
				"min_salary":  50000,
				"max_salary":  70000,
			},
			mockError:      nil,
			expectedStatus: http.StatusNoContent,
			expectedBody:   nil,
		},
		{
			name:           "empty id",
			id:             "",
			requestBody:    map[string]interface{}{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "id is required"},
		},
		{
			name: "invalid request",
			id:   "1",
			requestBody: map[string]interface{}{
				"type": "invalid-type",
			},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]interface{}{"error": "invalid request"},
		},
		{
			name: "not found",
			id:   "999",
			requestBody: map[string]interface{}{
				"title": "Updated Post",
			},
			mockError:      errors.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   map[string]interface{}{"error": "post not found"},
		},
		{
			name: "repository error",
			id:   "1",
			requestBody: map[string]interface{}{
				"title": "Updated Post",
			},
			mockError:      assert.AnError,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]interface{}{"error": "failed to update post"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			if tt.id != "" && tt.expectedStatus != http.StatusBadRequest {
				mockRepo.On("UpdatePost", tt.id, mock.AnythingOfType("*models.PostPutRequest")).Return(tt.mockError)
			}

			handler := NewHandler(mockRepo)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			v, ok := binding.Validator.Engine().(*validator.Validate)
			if ok {
				v.RegisterValidation("jobtype", models.JobTypeValidator)
			}

			c.Params = []gin.Param{{Key: "id", Value: tt.id}}
			jsonBytes, _ := json.Marshal(tt.requestBody)
			c.Request = httptest.NewRequest("PUT", "/", bytes.NewBuffer(jsonBytes))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.UpdatePost(c)

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
