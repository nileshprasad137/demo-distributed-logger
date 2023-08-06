
### Overview
This project showcases a distributed logging system using Redis Streams and the Go programming language.

The main objective of this project is to explore the use of Redis Streams for real-time log processing in distributed applications.

### Usage

```

make build: Builds the Go binary (distributed_logging) for the main application.

make up: Brings up the containers (Redis and the main application) using Docker Compose.

make down: Stops and removes the containers created by Docker Compose.

make clean: Removes the generated distributed_logging binary.

```

Once the containers are up and running, the project will simulate log messages being produced by multiple services and processed by the consumer group.

