.PHONY: run build test

# Run the application
run:
	cd gopherso && go run ./cmd/api

# Build the application
build:
	cd gopherso && go build -o ../bin/gopherso ./cmd/api

# Run tests
test:
	cd gopherso && go test -v ./...
