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

func Test_GetPost(t *testing.T) {
	resp, err := http.Get("http://localhost:8080/api/v1/posts/1")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	require.NoError(t, err)

	assert.Equal(t, result["id"], "1")
	assert.Equal(t, result["title"], "Senior Engineer")
	assert.Equal(t, result["location"], "Remote")
	assert.Equal(t, result["company"], "CompanyA")
	assert.Equal(t, result["description"], "Build scalable services")
	assert.Equal(t, result["type"], "full-time")
	assert.Equal(t, result["salary"], []interface{}{70000.0, 90000.0})
	assert.Equal(t, result["perks"], []interface{}{"Health", "401k"})
	assert.Equal(t, result["extras"], "")
	assert.Equal(t, result["created_at"], "2025-05-10T12:00:00Z")
	assert.Equal(t, result["updated_at"], "2025-05-10T12:00:00Z")

}
