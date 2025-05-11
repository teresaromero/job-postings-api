# Go parameters
BINARY_NAME=job-posting-api
MAIN_FILE=cmd/api/main.go

# Local environment variables
export SEED_DATA=true
export JWT_SECRET=secret

.PHONY: all build run clean debug

all: build

build:
	@echo "Building..."
	go build -o bin/$(BINARY_NAME) $(MAIN_FILE)

dev:
	@echo "Running..."
	go run $(MAIN_FILE)

run: build
	@echo "Running..."
	./bin/$(BINARY_NAME)

clean:
	@echo "Cleaning..."
	go clean
	rm -f $(BINARY_NAME)

test:
	@echo "Testing..."
	go test ./...

lint:
	@echo "Linting..."
	go vet ./...
	vacuum lint -dexq openapi.yaml
