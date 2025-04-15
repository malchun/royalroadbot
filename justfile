# royalroadbot/justfile

# List all available commands
default:
    @just --list

# Build the Docker image
build:
    sudo docker-compose build

# Rebuild the Docker image (force rebuild without cache)
full-rebuild:
    sudo docker-compose build --no-cache

# Run the container
run:
    sudo docker-compose -f docker-compose.yaml -f docker-compose-mongo.yaml up

# Run the container in detached mode
run-detached:
    sudo docker-compose -f docker-compose.yaml -f docker-compose-mongo.yaml up -d

# Stop the container
stop:
    sudo docker-compose down

# Show logs of the running container
logs:
    sudo docker-compose logs -f

# Clean up - remove containers, images, and volumes
clean:
    sudo docker-compose -f docker-compose.yaml -f docker-compose-mongo.yaml down --rmi all --volumes

# One command to build and run
restart:
    sudo docker-compose up --build

# Run tests in Docker
test:
    sudo docker-compose -f docker-compose-test.yaml -f docker-compose-mongo.yaml up --build --abort-on-container-exit test
    sudo docker-compose -f docker-compose-test.yaml -f docker-compose-mongo.yaml down
