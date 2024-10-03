# Makefile for managing the application

.PHONY: build run test docker-up docker-down

build:
	go build -o main ./cmd

run: build
	./main

docker-up: generated
	docker-compose up --build

docker-down:
	docker-compose down

test:
	go clean -testcache
	go test -cover ./...

test_api:
	go clean -testcache
	go test ./tests/...

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go