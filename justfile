# royalroadbot/justfile

# List all available commands
default:
    @just --list

# Build the application
build:
    go build -o bin/royalroadbot ./app

# Build the Docker image
build-docker:
    sudo docker-compose build

# Rebuild the Docker image (force rebuild without cache)
rebuild-all:
    sudo docker-compose build --no-cache

# Run the MongoDB development environment with Mongo Express client
run-dev-mongo:
    sudo docker-compose -f docker-compose-mongo.yaml -f docker-compose-dev.yaml up -d

# Stop the MongoDB development environment
stop-dev-mongo:
    sudo docker-compose -f docker-compose-mongo.yaml -f docker-compose-dev.yaml down

# Build and run the application locally with dev MongoDB
run-dev-local: build run-dev-mongo
    MONGODB_URI="mongodb://admin:password@127.0.0.1:27017" ./bin/royalroadbot

# Run the full container stack
run:
    sudo docker-compose -f docker-compose.yaml -f docker-compose-mongo.yaml up

# Run the mongo db
run-mongo:
    sudo docker-compose -f docker-compose-mongo.yaml up

# Run the container in detached mode
run-detached:
    sudo docker-compose -f docker-compose.yaml -f docker-compose-mongo.yaml up -d

# Stop the container
stop:
    sudo docker-compose down

# Stop all the containers
stop-all:
    sudo docker-compose -f docker-compose.yaml -f docker-compose-mongo.yaml -f docker-compose-dev.yaml down

# Show logs of the running container
logs:
    sudo docker-compose logs -f

# Clean up - remove containers, images, and volumes
clean:
    sudo docker-compose -f docker-compose.yaml -f docker-compose-mongo.yaml -f docker-compose-dev.yaml down --rmi all --volumes

# One command to build and run
restart:
    sudo docker-compose up --build

# Run tests in Docker
test-docker:
    sudo docker-compose -f docker-compose-test.yaml -f docker-compose-mongo.yaml up --build --abort-on-container-exit test
    sudo docker-compose -f docker-compose-test.yaml -f docker-compose-mongo.yaml down

# Run tests locally
test-local:
    go test ./app/... -v
