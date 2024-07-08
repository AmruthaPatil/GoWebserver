# Go Messaging Web Server with Kafka and Redis Integration

This repository contains a Docker-based setup for a Go web server that integrates with Apache Kafka for message queuing and Redis for data storage. It's designed to demonstrate a basic real-time data processing pipeline, where the web server handles HTTP requests, processes them through Kafka, and utilizes Redis for caching results.

## System Components

- **Go Web Server**: Processes incoming HTTP requests and interacts with Kafka and Redis.
- **Kafka**: Queues messages from the web server and ensures reliable processing.
- **Redis**: Stores processed messages for quick retrieval, acting as a fast data cache.
- **Zookeeper**: Manages Kafka cluster state and configurations.

## Prerequisites

Before you begin, ensure you have the following installed on your machine:
- **Docker**: Required to containerize the Kafka, Zookeeper, Redis, and Go web server. [Install Docker](https://www.docker.com/get-started)
- **Go**: Necessary to build the Go web server if you're not using pre-built Docker images. [Install Go](https://golang.org/dl/)

## Files Overview
Here's a breakdown of each file and its purpose:

### `main.go`

- This is the main file for the Go web server. It initializes the server, sets up routing, and includes the main functions to interact with Kafka and Redis.
- **Key Functions**:
  - `initKafka()`: Configures and initializes the Kafka producer.
  - `initRedis()`: Establishes a connection with the Redis server.
  - `postDataHandler()`: Handles POST requests to send data to Kafka.
  - `getDataHandler()`: Handles GET requests to retrieve data from Redis.
  - `consumeMessages()`: Continuously consumes messages from Kafka and stores them in Redis.

### `go.mod` & `go.sum`

- These files handle the project's dependencies.
  - `go.mod`: Declares the module's dependencies and other module-related directives.
  - `go.sum`: Contains expected cryptographic checksums of the content of specific module versions.

### `Dockerfile`

- Specifies the steps to create the Docker container for the web server.
- **Key Components**:
  - Uses the official Golang image to build the web server binary.
  - Copies the built application to a Debian-based production image for a lightweight setup.

### `docker-compose.yml`

- Defines and configures all the services needed to run the application, including the web server, Kafka, Zookeeper, and Redis.
- **Services Configured**:
  - `zookeeper`: Manages Kafka state.
  - `kafka`: Handles message queuing and processing.
  - `redis`: Provides a high-performance data store.
  - `webserver`: The custom Go web server designed to interact with Kafka and Redis.


## Setup and Running

### Build and Run the Containers

Navigate to the project directory where `docker-compose.yml` is located and run the following command to build and start all services:

```bash
docker-compose up --build
```

## Usage

### Sending Data to Kafka

To send data through the Go web server to Kafka, use the following `curl` command:

```bash
curl -X POST http://localhost:8080/data -d "message=Hello Kafka"
```
This command should return a success message indicating that the data was successfully sent to Kafka.

### Retrieving Data from Redis
Retrieve data from Redis via the web server:

```bash
curl http://localhost:8080/retrieve
```
This command should return the data stored in Redis if any processing logic from Kafka to Redis is successfully implemented and operational.

## Testing Kafka and Redis
### Accessing Redis CLI
To interact directly with Redis to check entries, use the following command:

```bash
docker exec -it [redis-container-id] redis-cli
```

### Inside the Redis CLI, you can run commands like:

```redis
keys *
get [your-key]
```

## Maintenance and Updates
Ensure that Docker is kept up to date, and periodically check for updates to dependencies used in this project, such as Kafka, Redis, and the base images used in the Dockerfile.

### Conclusion
This setup provides a foundational architecture for integrating a Go web server with Kafka and Redis, demonstrating a simple yet powerful pipeline for data processing and caching.