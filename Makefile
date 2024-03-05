
# Makefile for building, running, and managing a Go application with systemd

# Application name
APP_NAME := xiaoma-test

# Service name (should match the systemd service file name)
SERVICE_NAME := backend

# Build directory
BUILD_DIR := ./bin

# Source directory
SRC_DIR := ./cmd/app

local-up:
	migrate -path ./migration -database 'postgres://postgres:test@localhost:5436/postgres?sslmode=disable' up
local-down:
	migrate -path ./migration -database 'postgres://postgres:test@localhost:5436/postgres?sslmode=disable' down
server:
	go run cmd/app/main.go
post-up:
	docker run --name=postgres -e POSTGRES_PASSWORD='test' -p 5436:5432 -d --rm postgres
post-down:
	docker stop test-db
# Default make target
all: clean build

# Build the application
build: clean
	@echo "Building the application..."
	@GO111MODULE=on go build -o $(BUILD_DIR)/$(APP_NAME) $(SRC_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Clean build files
clean:
	@echo "Cleaning up..."
	@rm -f $(BUILD_DIR)/$(APP_NAME)
	@echo "Clean complete."

# Run the application
run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(APP_NAME)

# Restart the systemd service
restart-service:
	@echo "Restarting the systemd service..."
	@sudo systemctl restart $(SERVICE_NAME)

.PHONY: all build clean run restart-service


