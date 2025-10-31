.PHONY: run build test

# Run the application
run:
	go run ./cmd/api

# Build the application
build:
	go build -o ./bin/gopherso ./cmd/api

# Run tests
test:
	go test -v ./...
