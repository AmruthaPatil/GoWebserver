# Use the official Golang image to create a build artifact.
# https://hub.docker.com/_/golang
FROM golang:1.20 as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
# -o webserver is the name of the output binary.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o webserver

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/webserver /webserver

# Run the web service on container startup.
CMD ["/webserver"]
