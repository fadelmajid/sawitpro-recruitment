# Use the official Golang image as the base image for building
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

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

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
