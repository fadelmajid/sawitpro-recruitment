# Makefile for managing the application

.PHONY: build run test docker-up docker-down generated

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o main ./cmd

clean:
	rm -rf generated

init: clean generate
	go mod tidy
	go mod vendor

run: build
	./main

docker-up: generated
	docker-compose up --build

docker-down:
	docker-compose down

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -short -v ./...

test_cover:
	go clean -testcache
	go test -cover ./...

test_api:
	go clean -testcache
	go test ./tests/...

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go