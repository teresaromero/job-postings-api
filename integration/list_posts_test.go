//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetPostsSorted(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/posts")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Contains(t, result, "posts")
	assert.Contains(t, result, "count")
	assert.Equal(t, 6, int(result["count"].(float64)))

	posts := result["posts"].([]interface{})
	assert.Equal(t, "1", posts[0].(map[string]interface{})["id"])
	assert.Equal(t, "2", posts[1].(map[string]interface{})["id"])
	assert.Equal(t, "5", posts[2].(map[string]interface{})["id"])
	assert.Equal(t, "4", posts[3].(map[string]interface{})["id"])
	assert.Equal(t, "6", posts[4].(map[string]interface{})["id"])
	assert.Equal(t, "3", posts[5].(map[string]interface{})["id"])
}

func Test_GetPostsFiltered(t *testing.T) {

	t.Run("Filter by location", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/v1/posts?location=London")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Contains(t, result, "posts")
		assert.Contains(t, result, "count")
		assert.Equal(t, 1, int(result["count"].(float64)))

		posts := result["posts"].([]interface{})
		assert.Equal(t, "2", posts[0].(map[string]interface{})["id"])
	})

	t.Run("Filter by company name", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/v1/posts?company=CompanyA")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Contains(t, result, "posts")
		assert.Contains(t, result, "count")
		assert.Equal(t, 3, int(result["count"].(float64)))

		posts := result["posts"].([]interface{})
		assert.Equal(t, "1", posts[0].(map[string]interface{})["id"])
		assert.Equal(t, "5", posts[1].(map[string]interface{})["id"])
		assert.Equal(t, "3", posts[2].(map[string]interface{})["id"])
	})

	t.Run("Filter by min salary", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/v1/posts?min_salary=60000")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Contains(t, result, "posts")
		assert.Contains(t, result, "count")
		assert.Equal(t, 3, int(result["count"].(float64)))

		posts := result["posts"].([]interface{})
		assert.Equal(t, "1", posts[0].(map[string]interface{})["id"])
		assert.Equal(t, "2", posts[1].(map[string]interface{})["id"])
		assert.Equal(t, "5", posts[2].(map[string]interface{})["id"])
	})

	t.Run("Filter by max salary", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/v1/posts?max_salary=50000")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Contains(t, result, "posts")
		assert.Contains(t, result, "count")
		assert.Equal(t, 2, int(result["count"].(float64)))

		posts := result["posts"].([]interface{})
		assert.Equal(t, "6", posts[0].(map[string]interface{})["id"])
		assert.Equal(t, "3", posts[1].(map[string]interface{})["id"])
	})
	t.Run("Filter by range of salaries", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/v1/posts?min_salary=20000&max_salary=50000")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Contains(t, result, "posts")
		assert.Contains(t, result, "count")
		assert.Equal(t, 2, int(result["count"].(float64)))

		posts := result["posts"].([]interface{})
		assert.Equal(t, "6", posts[0].(map[string]interface{})["id"])
		assert.Equal(t, "3", posts[1].(map[string]interface{})["id"])
	})

	t.Run("Filter by title", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8080/api/v1/posts?title=Engineer")
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)
		require.NoError(t, err)

		assert.Contains(t, result, "posts")
		assert.Contains(t, result, "count")
		assert.Equal(t, 2, int(result["count"].(float64)))

		posts := result["posts"].([]interface{})
		assert.Equal(t, "1", posts[0].(map[string]interface{})["id"])
		assert.Equal(t, "4", posts[1].(map[string]interface{})["id"])
	})

}
