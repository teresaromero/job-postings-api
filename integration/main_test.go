//go:build integration
// +build integration

package integration

import (
	"log"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestMain(t *testing.M) {

	buildBinary()

	// Integration tests are taking into account the default rule ordering
	cmdRun := exec.Command("./bin/api")
	cmdRun.Env = []string{
		"PORT=8080",
		"NOW=2025-05-11T00:00:00Z",
		"SEED_FILE=integration/seed.json",
		"JWT_SECRET=secret",
	}
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
		log.Fatalf("failed to wait for main.go process: %v", err)
	}

	removeBinary()

	log.Default().Printf("Test completed with code: %d", code)
	os.Exit(code)
}

func buildBinary() {
	log.Default().Println("Building binary...")
	cmdRun := exec.Command("go", "build", "-o", "bin/api", "cmd/api/main.go")
	cmdRun.Dir = "../"
	cmdRun.Stdout = os.Stdout
	cmdRun.Stderr = os.Stderr
	if err := cmdRun.Run(); err != nil {
		log.Fatalf("failed to run build: %v", err)
	}
}

func removeBinary() {
	log.Default().Println("Removing binary...")
	cmdRun := exec.Command("rm", "bin/api")
	cmdRun.Dir = "../"
	cmdRun.Stdout = os.Stdout
	cmdRun.Stderr = os.Stderr
	if err := cmdRun.Run(); err != nil {
		log.Fatalf("failed to run remove binary: %v", err)
	}
}
