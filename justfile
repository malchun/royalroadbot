# royalroadbot/justfile

# List all available commands
default:
    @just --list

# Build the application
build:
    go build -o bin/royalroadbot ./app

# Rebuild the Docker image (force rebuild without cache)
rebuild-all:
    docker-compose build --no-cache

# Run the MongoDB development environment with Mongo Express client
run-dev-mongo:
    docker-compose -f docker-compose-mongo.yaml -f docker-compose-dev.yaml up -d

# Stop the MongoDB development environment
stop-dev-mongo:
    docker-compose -f docker-compose-mongo.yaml -f docker-compose-dev.yaml down

# Build and run the application locally with dev MongoDB
run-dev-local: build run-dev-mongo
    MONGODB_URI="mongodb://admin:password@127.0.0.1:27017" ./bin/royalroadbot

# Run the full container stack
run:
    docker-compose -f docker-compose.yaml -f docker-compose-mongo.yaml up

# Rebuild containers and run again the full container stack
re-run: rebuild-all run

# Run the mongo db
run-mongo:
    docker-compose -f docker-compose-mongo.yaml up

# Stop the container
stop:
    docker-compose down

# Show logs of the running container
logs:
    docker-compose logs -f

# Clean up - remove containers, images, and volumes
clean:
    docker-compose -f docker-compose.yaml -f docker-compose-mongo.yaml -f docker-compose-dev.yaml down --rmi all --volumes

# One command to build and run
restart:
    docker-compose up --build

# Run tests locally
test-local:
    go test ./app/... -v
