# Makefile for managing the application

.PHONY: build run test docker-up docker-down

build:
	go build -o main ./cmd

run: build
	./main

docker-up:
	docker-compose up --build

docker-down:
	docker-compose down
