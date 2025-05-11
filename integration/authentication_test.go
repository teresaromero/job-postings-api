//go:build integration
// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ManagePosts_Unauthenticated(t *testing.T) {
	// Attempt to create post without authentication
	postData := []byte(`{"title":"Test","company":"X","location":"Y","description":"D","type":"full-time","salary":[10000,20000]}`)
	resp, err := http.Post("http://localhost:8080/api/v1/posts", "application/json", bytes.NewReader(postData))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Attempt to update post without authentication
	req, err := http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/posts/1", bytes.NewReader(postData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	// Attempt to delete post without authentication
	req, err = http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/posts/1", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func Test_ManagePosts_Unauthorized(t *testing.T) {
	token := getToken(t, "employee")

	postData := []byte(`{"title":"Test","company":"X","location":"Y","description":"D","type":"full-time","salary":[10000,20000]}`)

	// Create post with unauthorized user
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/posts", bytes.NewReader(postData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Update post with unauthorized user
	req, err = http.NewRequest(http.MethodPut, "http://localhost:8080/api/v1/posts/1", bytes.NewReader(postData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	// Delete post with unauthorized user
	req, err = http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/posts/1", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func Test_ManagePosts_Authorized(t *testing.T) {
	token := getToken(t, "employer")

	postData := []byte(`{"title":"Test","company":"X","location":"Y","description":"D","type":"full-time","salary":[10000,20000]}`)

	client := &http.Client{}

	// Create post with authorized user
	req, err := http.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/posts", bytes.NewReader(postData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)
	assert.Contains(t, result, "id")

	// Update post with authorized user
	req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:8080/api/v1/posts/%s", result["id"]), bytes.NewReader(postData))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Delete post with authorized user
	req, err = http.NewRequest(http.MethodDelete, fmt.Sprintf("http://localhost:8080/api/v1/posts/%s", result["id"]), nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func getToken(t *testing.T, username string) string {
	loginData := map[string]string{
		"username": username,
		"password": "password",
	}
	bodyBytes, err := json.Marshal(loginData)
	require.NoError(t, err)

	resp, err := http.Post("http://localhost:8080/auth/token", "application/json", bytes.NewReader(bodyBytes))
	require.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	token, ok := result["token"].(string)
	require.True(t, ok, "token not found in response")
	return token
}
