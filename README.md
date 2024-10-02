# sawitpro-recruitment

This project provides a backend service for managing plantation estates, trees, and calculating drone patrol distances. It is built using **Golang**, the **Echo web framework**, and **PostgreSQL** as the database.

## Table of Contents
- [Requirements](#requirements)
- [Installation](#installation)
- [Docker Setup](#docker-setup)
- [Running the Application](#running-the-application)
- [Running Tests](#running-tests)
- [API Endpoints](#api-endpoints)
- [Makefile Commands](#makefile-commands)
- [Conclusion](#conclusion)

## Requirements

To run this project locally, ensure you have the following dependencies installed:

1. **Golang**: Version 1.16 or later
2. **Docker**: Version 20.10 or later
3. **Docker Compose**: Version 1.28 or later

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/your-username/your-repo-name.git
   cd your-repo-name
2. **Install Go dependencies (if you are not using Docker):**
    Use go mod to install all dependencies:
    ```bash
    Copy code
    go mod tidy

## Docker Setup
This project is dockerized for easy deployment. It includes a Dockerfile for building the application and a Docker Compose file to manage the application and its PostgreSQL database.

## Docker Compose Configuration
The docker-compose.yml file defines two services:
- app: The Go application service.
- db: The PostgreSQL database service.

## Running the Application
To build and run the application using Docker Compose, follow these steps:

1. Start the application:
    ```bash
    Copy code
    make docker-up
    
    This command will:
    - Build the Docker images.
    - Start the application and the PostgreSQL database container.

2. Access the API:
Once the application is running, you can access the API at http://localhost:8080.

## Running Tests
You can run tests locally or as part of the Docker build process.

1. Run tests locally:

    ```bash
    Copy code
    make test
    
    This command will execute all tests in the project.

2. Run tests with Docker:

When you run the application using Docker Compose, tests will run automatically before the application starts. If any tests fail, the application will not start.

## API Endpoints
Here are the main API endpoints provided by this application:

1. Create Estate
Endpoint: POST /estate
Request Body:
json
Copy code
{
  "width": 10,
  "length": 10
}
Response: 201 Created with the created estate details.

2. Add Tree to Estate
Endpoint: POST /estate/:id/tree
Request Body:
json
Copy code
{
  "x": 3,
  "y": 2,
  "height": 15
}
Response: 201 Created with the added tree details.

3. Get Estate Stats
Endpoint: GET /estate/:id/stats
Response: 200 OK with the statistics of trees in the estate.

4. Calculate Drone Patrol Distance
Endpoint: GET /estate/:id/drone-plan
Optional Query Parameter:
max_distance: Limit the total distance the drone can travel before landing.
Response: 200 OK with the total distance or the point where the drone will land.

## Makefile Commands
The project includes a Makefile for easy command execution:

Build the application:

    ```bash
    Copy code
    make build
    
Run the application locally:

    ```bash
    Copy code
    make run

Run tests:

    ```bash
    Copy code
    make test

Start the application with Docker Compose:

    ```bash
    Copy code
    make docker-up

Stop and remove containers:

    ```bash
    Copy code
    make docker-down