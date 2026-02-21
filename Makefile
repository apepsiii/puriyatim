.PHONY: help build run dev migrate clean install

# Default target
help:
	@echo "Available commands:"
	@echo "  make build    - Build the application"
	@echo "  make run      - Run the application"
	@echo "  make dev      - Run in development mode with hot reload"
	@echo "  make migrate  - Run database migrations"
	@echo "  make clean    - Clean build artifacts"
	@echo "  make install  - Install dependencies"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/puriyatim-app cmd/server/main.go

# Run the application
run: build
	@echo "Running application..."
	./bin/puriyatim-app

# Development mode with hot reload
dev:
	@echo "Running in development mode..."
	go run cmd/server/main.go

# Install dependencies
install:
	@echo "Installing dependencies..."
	go mod download
	go mod tidy

# Run database migrations
migrate:
	@echo "Running database migrations..."
	go run cmd/migrate/main.go

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f puriyatim-app

# Build CSS with Tailwind
css:
	@echo "Building CSS with Tailwind..."
	npx tailwindcss -i ./static/css/input.css -o ./static/css/app.css --watch

# Build CSS for production
css-prod:
	@echo "Building CSS for production..."
	npx tailwindcss -i ./static/css/input.css -o ./static/css/app.css --minify