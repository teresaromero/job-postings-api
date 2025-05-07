# Go parameters
BINARY_NAME=job-posting-api
MAIN_FILE=cmd/api/main.go

.PHONY: all build run clean debug

all: build

build:
	@echo "Building..."
	go build -o bin/$(BINARY_NAME) $(MAIN_FILE)

run:
	@echo "Running..."
	go run $(MAIN_FILE)

clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)

debug:
	@echo "Debugging..."
	dlv debug $(MAIN_FILE)

test:
	@echo "Testing..."
	go test ./...

.PHONY: help
help:
	@echo "Make commands:"
	@echo "build - Build the binary"
	@echo "run   - Run the application"
	@echo "clean - Remove binary and cache"
	@echo "debug - Debug using Delve"
	@echo "test  - Run tests"
	@echo "dev   - Run with live reload (requires air)"