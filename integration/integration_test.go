//go:build integration
// +build integration

package integration

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.M) {

	os.Setenv("NOW", "2025-05-11T00:00:00Z")
	os.Setenv("SEED_FILE", "integration/seed.json")

	cmdRun := exec.Command("make", "run")
	cmdRun.Dir = "../"
	cmdRun.Stdout = os.Stdout
	cmdRun.Stderr = os.Stderr
	if err := cmdRun.Start(); err != nil {
		log.Fatalf("failed to run make build: %v", err)
	}

	// Wait a bit for the server to start listening
	time.Sleep(2 * time.Second)

	code := t.Run()

	if err := cmdRun.Process.Kill(); err != nil {
		log.Fatalf("failed to kill main.go process: %v", err)
	}

	log.Default().Printf("Test completed with code: %d", code)
	os.Exit(code)
}

func Test_Heath(t *testing.T) {

	resp, err := http.Get("http://localhost:8080/ping")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
