all: run

build:
	@go build -o bin/api

run: build
	@./bin/api

test:
	@go test -v ./...

coverage:
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out