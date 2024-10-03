# sawitpro-recruitment

This project provides a backend service for managing plantation estates, trees, and calculating drone patrol distances. It is built using **Golang**, the **Echo web framework**, and **PostgreSQL** as the database.

## Table of Contents
- [Requirements](#requirements)
- [Running the Application](#running-the-application)
- [API Endpoints](#api-endpoints)
- [Makefile Commands](#makefile-commands)

## Requirements

To run this project locally, ensure you have the following dependencies installed:

1. **Golang**: Version 1.16 or later
2. **Docker**: Version 20.10 or later
3. **Docker Compose**: Version 1.28 or later

## Running the Application
To build,run, and test the application follow these steps (in file run.sh)
    ```bash
    #!/usr/bin/env bash
    make init                      # Initializes 
    make                           # Builds the binary
    make test                      # Runs unit tests with coverage 
    docker compose up --build -d   # Runs the docker to start the API and database.
    sleep 30                       # Wait until API runs
    make test_api                  # Runs the API testing with data.
    docker compose down --volumes  # Stops the docker containers and removes the volumes.

## API Endpoints
Here are the main API endpoints provided by this application:

1. Create Estate
Endpoint: POST /estate
Request Body:
    ```json
    {
        "width": 10,
        "length": 10
    }

Response: 201 Created with the created estate details.

2. Add Tree to Estate
Endpoint: POST /estate/:id/tree
Request Body:
    ```json
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
