# Use the official Golang image as the base image for building
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Install swag for generating Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the source code into the container
COPY . .

# Generate Swagger documentation
RUN swag init

# Run tests
RUN go test -v ./...

# Build the Go app with CGO disabled for static linking
RUN CGO_ENABLED=0 go build -o main .

# Second stage: Create a smaller image for production
FROM alpine:latest

# Install necessary utilities
RUN apk add --no-cache bash

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .
# Copy the Swagger documentation
COPY --from=builder /app/docs ./docs

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]