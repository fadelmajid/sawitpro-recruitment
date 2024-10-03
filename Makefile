# Makefile for managing the application

.PHONY: build all init docker-up docker-down generated

all: build/main

build/main: cmd/main.go generated
	@echo "Building..."
	go build -o main ./cmd

clean:
	rm -rf generated

init: clean generated
	go mod tidy
	go mod vendor

docker-up: generated
	docker-compose up --build -d

docker-down:
	docker-compose down --volumes

test:
	go clean -testcache
	go test -short -cover  -coverprofile=coverage.out ./handlers ./repositories
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

test_api:
	go clean -testcache
	go test ./tests/...

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go