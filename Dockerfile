# Dockerfile

# Use the official Golang image as the base image for building the Go binary
FROM golang:1.20.5-alpine3.17

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code into the container
COPY . .

# Build the Go binary for the main application
RUN go build -o distributed_logging

# Set the entrypoint for the main application
ENTRYPOINT ["./distributed_logging"]
