# Makefile

# With this Makefile, you can use the following commands:

# make or make build: Builds the Go binary (distributed_logging) for the main application.
# make up: Brings up the containers (Redis and the main application) using Docker Compose.
# make down: Stops and removes the containers created by Docker Compose.
# make clean: Removes the generated distributed_logging binary.

# Variables
DOCKER_COMPOSE = docker-compose

# Targets
.PHONY: all build up down clean

# Default target (build the Go binary)
all: build

# Build the Go binary
build:
	go build -o distributed_logging

# Bring up the containers using Docker Compose
up:
	$(DOCKER_COMPOSE) up

# Stop and remove the containers using Docker Compose
down:
	$(DOCKER_COMPOSE) down

# Clean generated files
clean:
	rm -f distributed_logging
